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
		var err error

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

			// 将聊天的 Assistant ID 设置为空
			err = adb.ChatMessageService.ClearChatAssistant(ctx, chatEntity...)
			if err != nil {
				b.logger.Sugar.Errorf("Batch AssistantDelete: %v", err)
				return
			}

			chatPage++
		}

		err = adb.AssistantService.DeleteAllKey(ctx, adb.AssistantEntity)
		if err != nil {
			b.logger.Sugar.Errorf("Batch AssistantAllSecretDelete: %v", err)
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
