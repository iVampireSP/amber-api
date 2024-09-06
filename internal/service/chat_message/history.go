package chat_message

import (
	"context"
	"errors"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
)

func (s *Service) GetChatMessage(ctx context.Context, chat *entity.Chat) ([]*entity.ChatMessage, error) {
	chatMessage, err := s.dao.WithContext(ctx).ChatMessage.
		//Preload(s.dao.ChatMessage.File).
		//Preload(s.dao.ChatMessage.File).
		Preload(s.dao.ChatMessage.UserFile.File).
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
		Preload(s.dao.ChatMessage.File).
		//Preload(s.dao.ChatMessage.UserFile).
		Preload(s.dao.ChatMessage.UserFile.File).
		Order(s.dao.ChatMessage.CreatedAt.Asc()).
		Find()

	return chatMessage, err
}

func (s *Service) CreateChatMessage(ctx context.Context, chatMessage *entity.ChatMessage) error {
	err := s.dao.WithContext(ctx).ChatMessage.
		Preload(s.dao.ChatMessage.File).
		Preload(s.dao.ChatMessage.UserFile).
		Preload(s.dao.ChatMessage.UserFile.File).
		Create(chatMessage)

	return err
}

func (s *Service) DeleteChatMessage(ctx context.Context, ChatMessage *entity.ChatMessage) error {
	_, err := s.dao.WithContext(ctx).ChatMessage.Delete(ChatMessage)

	return err
}

func (s *Service) DeleteChatMessageByChats(ctx context.Context, chat ...*entity.Chat) error {
	if len(chat) == 0 {
		return errors.New("no chat provided")
	}

	_, err := s.dao.WithContext(ctx).Chat.Delete(chat...)

	return err
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
		Preload(s.dao.ChatMessage.UserFile).
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
	_, err := s.dao.WithContext(ctx).ChatMessage.Where(s.dao.ChatMessage.ChatId.Eq(uint(chat.Id))).Delete()

	return err
}

func (s *Service) ClearToolCall(ctx context.Context, cm *entity.ChatMessage) error {
	// 将此条 cm 的 tool_call 设置为 null
	_, err := s.dao.WithContext(ctx).ChatMessage.Where(s.dao.ChatMessage.Id.Eq(uint(cm.Id))).Update(s.dao.ChatMessage.ToolCall, nil)

	return err
}
