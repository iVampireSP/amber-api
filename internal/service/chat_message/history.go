package chat_message

import (
	"context"
	"errors"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
)

func (s *Service) GetChatMessage(ctx context.Context, chat *entity.Chat) ([]*entity.ChatMessage, error) {
	var chatMessage []*entity.ChatMessage
	err := s.x.Context(ctx).Where("chat_id = ?", chat.ID).And("role != ?", entity.RoleHideSystem.String()).OrderBy("id asc").Find(&chatMessage)
	return chatMessage, err
}

func (s *Service) GetChatMessageWithHide(ctx context.Context, chat *entity.Chat) ([]*entity.ChatMessage, error) {
	var chatMessage []*entity.ChatMessage
	err := s.x.Context(ctx).Where("chat_id = ?", chat.ID).OrderBy("id asc").Find(&chatMessage)
	return chatMessage, err
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

func (s *Service) DeleteChatMessageByChats(ctx context.Context, chat ...*entity.Chat) error {
	// build wherein
	var ids = make([]int64, 0)
	for _, v := range chat {
		ids = append(ids, v.ID)
	}

	if len(ids) == 0 {
		return errors.New("no chat provided")
	}

	_, err := s.x.Context(ctx).In("chat_id", ids).Delete(&entity.ChatMessage{})
	return err
}

// CountChatMessage count messages
func (s *Service) CountChatMessage(ctx context.Context, chat *entity.Chat) (int64, error) {
	count, err := s.x.Context(ctx).Where("chat_id = ?", chat.ID).Count(&entity.ChatMessage{})
	return count, err
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
	_, err := s.x.Context(ctx).Where("chat_id = ?", chat.ID).Limit(1).OrderBy("id desc").Get(&chatMessage)
	return &chatMessage, err
}

// UpdateMessageContent update message content
func (s *Service) UpdateMessageContent(ctx context.Context, chatMessage *entity.ChatMessage) error {
	_, err := s.x.Context(ctx).ID(chatMessage.ID).Cols("content").Update(chatMessage)
	return err
}

func (s *Service) ClearChatMessage(ctx context.Context, chat *entity.Chat) error {
	_, err := s.x.Context(ctx).Where("chat_id = ?", chat.ID).Delete(&entity.ChatMessage{})
	return err
}
