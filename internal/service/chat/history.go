package chat

import (
	"context"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
)

func (s *Service) GetChatMessage(ctx context.Context, chat *entity.Chat) ([]*entity.ChatMessage, error) {
	var ChatMessage []*entity.ChatMessage
	err := s.x.Context(ctx).Where("chat_id = ?", chat.ID).OrderBy("id asc").Find(&ChatMessage)
	return ChatMessage, err
}

func (s *Service) CreateChatMessage(ctx context.Context, ChatMessage *entity.ChatMessage) error {
	if ChatMessage.ChatId == consts.NoRecord {
		return consts.ErrChatIdNotProvided
	}

	_, err := s.x.Context(ctx).Insert(ChatMessage)
	return err
}

func (s *Service) DeleteChatMessage(ctx context.Context, ChatMessage *entity.ChatMessage) error {
	_, err := s.x.Context(ctx).ID(ChatMessage.ID).Delete(ChatMessage)
	return err
}

func (s *Service) DeleteChatMessageByChatId(ctx context.Context, chatId int64) error {
	_, err := s.x.Context(ctx).Where("chat_id = ?", chatId).Delete(&entity.ChatMessage{})
	return err
}

func (s *Service) DeleteChatMessageByAssistantId(ctx context.Context, assistantId int64) error {
	_, err := s.x.Context(ctx).Where("assistant_id = ?", assistantId).Delete(&entity.ChatMessage{})
	return err
}

func (s *Service) DeleteChatMessageByUserId(ctx context.Context, userId schema.UserId) error {
	_, err := s.x.Context(ctx).Where("user_id = ?", userId).Delete(&entity.ChatMessage{})
	return err
}

// GetLatestMessage get latest chat message
func (s *Service) GetLatestMessage(ctx context.Context, chat *entity.Chat) (*entity.ChatMessage, error) {
	var chatMessage entity.ChatMessage
	_, err := s.x.Context(ctx).Where("chat_id = ?", chat.ID).Limit(1).OrderBy("id asc").Get(&chatMessage)
	return &chatMessage, err
}

// UpdateMessageContent update message content
func (s *Service) UpdateMessageContent(ctx context.Context, chatMessage *entity.ChatMessage) error {
	_, err := s.x.Context(ctx).ID(chatMessage.ID).Cols("content").Update(chatMessage)
	return err
}
