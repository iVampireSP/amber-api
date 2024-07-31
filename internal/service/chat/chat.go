package chat

import (
	"context"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
)

func (s *Service) CreateChat(ctx context.Context, createChatRequest *schema.ChatCreateRequest) (*entity.Chat, error) {
	var chat entity.Chat

	// 验证 assistant Id 是否属于存在且属于用户
	a, err := s.a.GetAssistant(ctx, createChatRequest.AssistantId)
	if err != nil {
		return nil, err
	}

	if a.ID == consts.NoRecord || a.UserId != createChatRequest.UserId {
		return nil, consts.ErrAssistantNotFound
	}

	chat.Name = createChatRequest.Name
	chat.AssistantId = createChatRequest.AssistantId
	chat.UserId = createChatRequest.UserId

	_, err = s.x.Context(ctx).Insert(&chat)

	return &chat, err
}

func (s *Service) GetChat(ctx context.Context, id int64) (*entity.Chat, error) {
	var chat entity.Chat
	_, err := s.x.Context(ctx).ID(id).Get(&chat)
	return &chat, err
}

func (s *Service) ListChat(ctx context.Context) ([]*entity.Chat, error) {
	var chats []*entity.Chat
	err := s.x.Context(ctx).Find(&chats)
	return chats, err
}

func (s *Service) DeleteChat(ctx context.Context, id int64) error {
	_, err := s.x.Context(ctx).ID(id).Delete(&entity.Chat{})
	return err
}

func (s *Service) DeleteChatFromUserId(ctx context.Context, id int64, userId schema.UserId) error {
	count, err := s.x.Context(ctx).Where("id = ?", id).Where("user_id = ?", userId).Delete(&entity.Chat{})
	if err != nil {
		return err
	}
	if count == 0 {
		return consts.ErrChatNotFound
	}

	return nil
}

func (s *Service) Exists(ctx context.Context, id int64) (bool, error) {
	count, err := s.x.Context(ctx).Where("id = ?", id).Count(&entity.Chat{})
	return count > 0, err
}

func (s *Service) ListChatFromUserId(ctx context.Context, userId schema.UserId) ([]*entity.Chat, error) {
	var chats []*entity.Chat
	err := s.x.Context(ctx).Where("user_id = ?", userId).Find(&chats)
	return chats, err
}
