package chat_message

import (
	"context"
	"errors"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	page2 "rag-new/pkg/page"
	"sort"
)

func (s *Service) GetChatMessage(ctx context.Context, chat *entity.Chat) ([]*entity.ChatMessage, error) {
	chatMessage, err := s.dao.WithContext(ctx).ChatMessage.
		//Preload(s.dao.ChatMessage.File).
		Preload(s.dao.ChatMessage.File).
		//Preload(s.dao.ChatMessage.UserFile.File).
		Preload(s.dao.ChatMessage.Assistant).
		Where(s.dao.ChatMessage.ChatId.Eq(uint(chat.Id))).
		Where(s.dao.ChatMessage.Role.Neq(schema.RoleHideSystem.String())).
		Where(s.dao.ChatMessage.Role.Neq(schema.RoleHideHuman.String())).
		Order(s.dao.ChatMessage.CreatedAt.Asc()).
		Find()

	return chatMessage, err
}

func (s *Service) GetChatMessageWithHide(ctx context.Context, chat *entity.Chat) ([]*entity.ChatMessage, error) {
	chatMessage, err := s.dao.WithContext(ctx).ChatMessage.
		Where(s.dao.ChatMessage.ChatId.Eq(uint(chat.Id))).
		Preload(s.dao.ChatMessage.Assistant).
		Preload(s.dao.ChatMessage.File).
		//Preload(s.dao.ChatMessage.UserFile).
		//Preload(s.dao.ChatMessage.UserFile.File).
		Order(s.dao.ChatMessage.CreatedAt.Asc()).
		Find()

	return chatMessage, err
}

// GetLatestChatMessage 获取最新的消息，并保证消息的完整性
func (s *Service) GetLatestChatMessage(ctx context.Context, chat *entity.Chat, pageSize int) ([]*entity.ChatMessage, int64, error) {
	page := 1
	var r []*entity.ChatMessage
	var count int64

	for {
		newHistory, newCount, err := s.dao.WithContext(ctx).ChatMessage.
			Where(s.dao.ChatMessage.ChatId.Eq(chat.Id.Uint())).
			Preload(s.dao.ChatMessage.File).
			Order(s.dao.ChatMessage.CreatedAt.Desc()).
			FindByPage(page2.OffsetCustom(page, pageSize), pageSize)

		if err != nil {
			return nil, 0, err
		}

		if newCount == 0 {
			break
		}

		// 放进去，最后重新排序
		r = append(r, newHistory...)
		//r = append(newHistory, r...)

		count += newCount // 更新总数量

		var lastOne = newHistory[len(newHistory)-1]
		// 检查条件
		if lastOne.Role == schema.RoleToolCall || lastOne.Role == schema.RoleTool {
			page++ // 继续获取之前的消息
		} else {
			break
		}
	}

	// 根据 CreatedAt 从小到大排序
	sort.Slice(r, func(i, j int) bool {
		return r[i].CreatedAt.Before(r[j].CreatedAt)
	})

	return r, count, nil
}
func (s *Service) GetChatMessagePageAsc(ctx context.Context, chat *entity.Chat, page int, pageSize int) (cms []*entity.ChatMessage, count int64, totalCount int64, err error) {
	cms, totalCount, err = s.dao.WithContext(ctx).ChatMessage.
		Where(s.dao.ChatMessage.ChatId.Eq(chat.Id.Uint())).
		Preload(s.dao.ChatMessage.File).
		Order(s.dao.ChatMessage.CreatedAt.Asc()).
		FindByPage(page2.OffsetCustom(page, pageSize), pageSize)

	if err != nil {
		return nil, 0, 0, err
	}

	var dataLength = len(cms)

	if dataLength > 0 {
		for {
			if cms[0].Role == schema.RoleToolCall {
				// 由于是 ASC，向后翻 1 页面
				page = page + 1
				newHistory, _, err := s.dao.WithContext(ctx).ChatMessage.
					Where(s.dao.ChatMessage.ChatId.Eq(chat.Id.Uint())).
					Preload(s.dao.ChatMessage.File).
					Order(s.dao.ChatMessage.CreatedAt.Asc()).
					FindByPage(page2.OffsetCustom(page, pageSize), pageSize)

				if err != nil {
					return nil, 0, 0, err
				}

				var newCount = len(newHistory)

				if newCount == 0 {
					break
				}

				// 将新的内容放到之前内容的最前面
				cms = append(newHistory, cms...) // 把新内容放到前面
				dataLength += newCount           // 更新总数量
			} else {
				break
			}
		}
	}

	return cms, int64(dataLength), totalCount, nil
}

func (s *Service) CreateChatMessage(ctx context.Context, chatMessage *entity.ChatMessage) error {
	err := s.dao.WithContext(ctx).ChatMessage.
		Preload(s.dao.ChatMessage.File).
		//Preload(s.dao.ChatMessage.UserFile).
		//Preload(s.dao.ChatMessage.UserFile.File).
		Create(chatMessage)

	return err
}

// DeleteChatMessage 删除指定的 chat message，但是现在不推荐了，因为有 Message Block。Chat Message 是不可变的。
func (s *Service) DeleteChatMessage(ctx context.Context, ChatMessage *entity.ChatMessage) error {
	_, err := s.dao.WithContext(ctx).ChatMessage.Delete(ChatMessage)

	return err
}

func (s *Service) DeleteChatMessageByChats(ctx context.Context, chat ...*entity.Chat) error {
	if len(chat) == 0 {
		return errors.New("no chat provided")
	}

	for _, c := range chat {
		err := s.messageBlock.ClearMessageBlock(ctx, c)
		if err != nil {
			return err
		}
		_, err = s.dao.WithContext(ctx).Chat.Delete(c)
	}

	return nil
}

// ClearChatAssistant 重置对话的 Assistant ID 为 nil
func (s *Service) ClearChatAssistant(ctx context.Context, chats ...*entity.Chat) error {
	var ids = make([]uint, len(chats))

	for i, c := range chats {
		ids[i] = c.Id.Uint()
	}

	_, err := s.dao.WithContext(ctx).Chat.Where(s.dao.Chat.Id.In(ids...)).UpdateSimple(s.dao.Chat.AssistantId.Null())
	if err != nil {
		return err
	}

	// 将 Chat 的 Chat Message 的 Assistant ID 设置为 nil
	_, err = s.dao.WithContext(ctx).ChatMessage.Where(s.dao.ChatMessage.ChatId.In(ids...)).UpdateSimple(s.dao.ChatMessage.AssistantId.Null())
	if err != nil {
		return err
	}

	return nil
}

// ClearAllChatAssistant 重置对话的 Assistant ID 为 nil
func (s *Service) ClearAllChatAssistant(ctx context.Context, assistantId *schema.EntityId) error {
	_, err := s.dao.WithContext(ctx).Chat.
		Where(s.dao.Chat.AssistantId.Eq(assistantId.Uint())).
		UpdateSimple(s.dao.Chat.AssistantId.Null())

	if err != nil {
		return err
	}

	// 清理 ChatMessage 的 Assistant ID
	_, err = s.dao.WithContext(ctx).ChatMessage.
		Where(s.dao.ChatMessage.AssistantId.Eq(assistantId.Uint())).
		UpdateSimple(s.dao.ChatMessage.AssistantId.Null())
	if err != nil {
		return err
	}

	return nil
}

// CountChatMessage count messages
func (s *Service) CountChatMessage(ctx context.Context, chat *entity.Chat) (int64, error) {
	count, err := s.dao.WithContext(ctx).ChatMessage.Where(s.dao.ChatMessage.ChatId.Eq(uint(chat.Id))).Count()

	return count, err
}

// GetLatestMessage get latest chat message
func (s *Service) GetLatestMessage(ctx context.Context, chat *entity.Chat) (*entity.ChatMessage, error) {
	count, err := s.dao.WithContext(ctx).ChatMessage.Where(s.dao.ChatMessage.ChatId.Eq(uint(chat.Id))).Count()

	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, nil
	}

	chatMessage, err := s.dao.WithContext(ctx).ChatMessage.
		Preload(s.dao.ChatMessage.File).
		//Preload(s.dao.ChatMessage.UserFile).
		Where(s.dao.ChatMessage.ChatId.Eq(uint(chat.Id))).
		Order(s.dao.ChatMessage.CreatedAt.Desc()).
		Order(s.dao.ChatMessage.Id.Desc()).First()

	return chatMessage, err
}

// UpdateMessageContent update message content
func (s *Service) UpdateMessageContent(ctx context.Context, chatMessage *entity.ChatMessage) error {
	_, err := s.dao.WithContext(ctx).ChatMessage.
		Where(s.dao.ChatMessage.Id.Eq(uint(chatMessage.ChatId))).
		Update(s.dao.ChatMessage.Content, chatMessage.Content)

	return err
}

func (s *Service) ClearChatMessage(ctx context.Context, chat *entity.Chat) error {
	// 先删除 Message Block
	err := s.messageBlock.ClearMessageBlock(ctx, chat)
	if err != nil {
		return err
	}

	_, err = s.dao.WithContext(ctx).ChatMessage.Where(s.dao.ChatMessage.ChatId.Eq(uint(chat.Id))).Delete()

	return err
}

//func (s *Service) ClearToolCall(ctx context.Context, cm *entity.ChatMessage) error {
//	// 将此条 cm 的 tool_call 设置为 null
//	_, err := s.dao.WithContext(ctx).ChatMessage.Where(s.dao.ChatMessage.Id.Eq(uint(cm.Id))).Update(s.dao.ChatMessage.ToolCall, nil)
//
//	return err
//}
