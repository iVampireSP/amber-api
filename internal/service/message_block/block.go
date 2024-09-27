package message_block

import (
	"context"
	"fmt"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	entity2 "github.com/milvus-io/milvus-sdk-go/v2/entity"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/md5"
	"slices"
)

var allowedRoles = []schema.ChatRole{
	schema.RoleHuman,
	schema.RoleHumanLater,
	schema.RoleHideHuman,
	schema.RoleSystem,
	schema.RoleHideSystem,
	schema.RoleSystemOverride,
	schema.RoleAssistant,
}

func (s *Service) MessageToBlock(cms []*entity.ChatMessage) ([]*entity.MessageBlock, error) {
	// 记录是否成为了一组
	var blocked = false
	var blockMessage []*entity.MessageBlock
	var tempMessageList []*entity.ChatMessage
	var tempMessageContent string
	var currentChatId schema.EntityId

	for _, v := range cms {
		// 如果不在 allowedRoles, 则跳过
		if !slices.Contains(allowedRoles, v.Role) {
			continue
		}

		if blocked {
			blocked = false
			tempMessageContent = ""
			currentChatId = 0
			tempMessageList = []*entity.ChatMessage{}
		}

		currentChatId = v.ChatId

		if v.Role != schema.RoleAssistant {
			tempMessageContent += v.Content + "\n"
			tempMessageList = append(tempMessageList, v)
		} else if v.Role == schema.RoleAssistant {
			// 如果是 Assistant，则消息块结束
			blocked = true
			// 最后收尾不需要换行
			tempMessageContent += v.Content

			contentHash, err := md5.Md5(tempMessageContent)
			if err != nil {
				return nil, err
			}

			tempMessageList = append(tempMessageList, v)

			// 收尾工作
			blockMessage = append(blockMessage, &entity.MessageBlock{
				Hash:        contentHash,
				FullContent: tempMessageContent,
				Message:     tempMessageList,
				ChatId:      currentChatId,
				Model: entity.Model{
					CreatedAt: v.CreatedAt,
					UpdatedAt: v.UpdatedAt,
				},
			})
		}
	}

	//if outputTemp {
	//	// 如果要输出没有收尾的消息
	//	contentHash, err := md5.Md5(tempMessageContent)
	//	if err != nil {
	//		return nil, err
	//	}
	//	blockMessage = append(blockMessage, &entity.MessageBlock{
	//		Hash:        contentHash,
	//		FullContent: tempMessageContent,
	//		Message:     tempMessageList,
	//		ChatId:      currentChatId,
	//		Temp:        true,
	//	})
	//}

	return blockMessage, nil

}

func (s *Service) SaveBlock(ctx context.Context, messageBlock []*entity.MessageBlock) error {
	var notExistsBlocks = make([]*entity.MessageBlock, 0)
	var pendingHashes = make([]string, 0)
	for _, bm := range messageBlock {
		// 不处理不完整的块
		if bm.Temp {
			continue
		}

		notExistsBlocks = append(notExistsBlocks, bm)
		pendingHashes = append(pendingHashes, bm.Hash)
	}

	// 批量寻找不存在的块
	var dao = s.dao.MessageBlock.WithContext(ctx)

	find, err := dao.Where(s.dao.MessageBlock.Hash.In(pendingHashes...)).Find()
	if err != nil {
		return err
	}

	for _, bm := range find {
		// 找到的块，从待处理列表中删除
		for i, v := range notExistsBlocks {
			if v.Hash == bm.Hash {
				notExistsBlocks = append(notExistsBlocks[:i], notExistsBlocks[i+1:]...)
				break
			}
		}
	}

	for _, bm := range notExistsBlocks {
		var content = bm.FullContent
		// 如果 content > 8192
		if len(content) > s.config.OpenAI.EmbeddingMaxToken {
			// 剪裁
			content = content[:s.config.OpenAI.EmbeddingMaxToken]
		}

		err = s.dao.WithContext(ctx).MessageBlock.Create(bm)
		if err != nil {
			return err
		}

		emb, err := s.embedding.TextEmbedding(ctx, []string{content})
		if err != nil {
			return err
		}

		var entityCols = []entity2.Column{
			entity2.NewColumnFloatVector("vector", s.config.OpenAI.EmbeddingDim, emb),
			entity2.NewColumnVarChar("model", []string{s.config.OpenAI.EmbeddingModel}),
			entity2.NewColumnInt64("block_id", []int64{int64(bm.Id)}),
			entity2.NewColumnInt64("chat_id", []int64{int64(bm.ChatId)}),
		}

		// insert to milvus
		_, err = s.milvus.Upsert(ctx, s.config.Milvus.MessageBlockCollection, "", entityCols...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) BlockExists(ctx context.Context, chatId schema.EntityId, hash string) (bool, error) {
	i, err := s.dao.WithContext(ctx).MessageBlock.
		Where(s.dao.MessageBlock.ChatId.Eq(chatId.Uint())).
		Where(s.dao.MessageBlock.Hash.Eq(hash)).Count()

	return i > 0, err
}

func (s *Service) SearchMessageBlock(ctx context.Context, chat *entity.Chat, content string) ([]*entity.MessageBlock, error) {
	emb, err := s.embedding.TextEmbedding(ctx, []string{content})
	if err != nil {
		return nil, err
	}
	var filter = fmt.Sprintf("chat_id == %d && model == '%s'", chat.Id, s.config.OpenAI.EmbeddingModel)
	sp, err := entity2.NewIndexAUTOINDEXSearchParam(1)
	if err != nil {
		return nil, err
	}
	vector := entity2.FloatVector(emb[0])
	existingChunks, err := s.milvus.Search(ctx, s.config.Milvus.MessageBlockCollection,
		[]string{},
		filter,
		[]string{"block_id"},
		[]entity2.Vector{vector},
		"vector",
		entity2.L2,
		3,
		sp, client.WithLimit(7))

	var ids []uint

	for _, res := range existingChunks {
		// 没找到，直接返回空的
		if res.ResultCount == 0 {
			return make([]*entity.MessageBlock, 0), nil
		}

		var blockIdColumn *entity2.ColumnInt64
		for _, field := range res.Fields {
			if field.Name() == "block_id" {
				c, ok := field.(*entity2.ColumnInt64)
				if ok {
					blockIdColumn = c
				}
			}
		}

		// 没有记录
		if blockIdColumn == nil {
			return make([]*entity.MessageBlock, 0), nil
			//return nil, fmt.Errorf("block_id column not found")
		}

		for i := 0; i < res.ResultCount; i++ {
			id, err := blockIdColumn.ValueByIdx(i)
			if err != nil {
				return nil, err
			}

			ids = append(ids, uint(id))

		}
	}

	messageBlocks, err := s.dao.MessageBlock.Where(s.dao.MessageBlock.Where(s.dao.MessageBlock.Id.In(ids...))).Find()
	return messageBlocks, err
}

func (s *Service) ClearMessageBlock(ctx context.Context, chat *entity.Chat) error {
	var filter = fmt.Sprintf(`chat_id == %d`, chat.Id)
	errDelete := s.milvus.Delete(ctx, s.config.Milvus.MessageBlockCollection, "", filter)

	if errDelete != nil {
		return errDelete
	}

	_, err := s.dao.MessageBlock.Where(s.dao.MessageBlock.ChatId.Eq(chat.Id.Uint())).Delete()

	if err != nil {
		return err
	}

	return nil
}
