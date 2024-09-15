package batch

import (
	"context"
	"rag-new/internal/entity"
	"rag-new/internal/service/assistant"
	"rag-new/internal/service/chat"
	"rag-new/internal/service/chat_message"
)

type AssistantDeleteBatch struct {
	AssistantEntity    *entity.Assistant
	AssistantService   *assistant.Service
	ChatService        *chat.Service
	ChatMessageService *chat_message.Service
}

func (b *Batch) AssistantDelete(ctx context.Context, adb *AssistantDeleteBatch) error {
	go func() {
		var chatPage = 1
		for {
			chatEntity, err := adb.ChatService.ListChatFromAssistantByPage(ctx, adb.AssistantEntity, chatPage)
			if err != nil {
				b.logger.Sugar.Errorf("Batch AssistantDelete: %v", err)
				return
			}

			var chatEntities = len(chatEntity)
			if chatEntities == 0 {
				b.logger.Sugar.Debugf("Batch AssistantDelete: success")
				break
			}

			// 删除所有的聊天
			err = adb.ChatMessageService.DeleteChatMessageByChats(ctx, chatEntity...)
			if err != nil {
				b.logger.Sugar.Errorf("Batch AssistantDelete: %v", err)
				return
			}

			// 删除 chat
			err = adb.ChatService.DeleteChats(ctx, chatEntity...)
			if err != nil {
				return
			}

			chatPage++
		}

		// 删除 shares
		err := adb.AssistantService.DeleteAllShare(ctx, adb.AssistantEntity)
		if err != nil {
			b.logger.Sugar.Errorf("Batch AssistantShareDelete: %v", err)
			return
		}

		err = adb.AssistantService.DeleteAssistant(ctx, adb.AssistantEntity.Id)
		if err != nil {
			b.logger.Sugar.Errorf("Batch AssistantDelete: %v", err)
			return
		}
	}()
	return nil
}
