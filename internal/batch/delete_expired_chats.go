package batch

import (
	"context"
	"rag-new/internal/service/chat"
	"time"
)

type ChatDeleteExpired struct {
	BeforeTime  time.Time
	ChatService *chat.Service
}

func (b *Batch) DeleteExpiredChats(ctx context.Context, ce *ChatDeleteExpired) error {
	b.logger.Sugar.Info("Batch DeleteExpiredChats: %v", ce.BeforeTime)
	go func() {
		err := ce.ChatService.DeleteExpiredChats(ctx, ce.BeforeTime)
		if err != nil {
			b.logger.Sugar.Error("delete expired chats error", err)
			return
		}

		b.logger.Sugar.Info("Batch DeleteExpiredChats: success")
	}()
	return nil
}
