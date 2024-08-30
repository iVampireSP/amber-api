package chat

import (
	"context"
	"errors"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
	"time"
)

func (s *Service) CreateChat(ctx context.Context, createChatRequest *schema.ChatCreateRequest) (*entity.Chat, error) {
	var chat entity.Chat

	// 验证 assistant Id 是否属于存在且属于用户
	a, err := s.a.GetAssistant(ctx, createChatRequest.AssistantId)
	if err != nil {
		return nil, err
	}

	if a.Id == consts.NoRecord || a.UserId != createChatRequest.UserId {
		return nil, consts.ErrAssistantNotFound
	}

	chat.Name = createChatRequest.Name
	chat.AssistantId = createChatRequest.AssistantId
	chat.UserId = createChatRequest.UserId
	chat.Owner = schema.OwnerUser

	if createChatRequest.ExpiredAt != nil {
		// 过期时间不能小于当前时间
		if createChatRequest.ExpiredAt.Before(time.Now()) {
			return nil, consts.ErrExpiredTimeCanNotBeforeNow
		}

		// 不能大于 2038 年
		if createChatRequest.ExpiredAt.After(time.Date(2038, 1, 1, 0, 0, 0, 0, time.UTC)) {
			return nil, consts.ErrExpiredTimeCanNotAfter2038
		}

		chat.ExpiredAt = &createChatRequest.ExpiredAt.Time
	}

	_, err = s.x.Context(ctx).Insert(&chat)

	return &chat, err
}

// CreateGuestChat 用于创建访客对话，这个对话应是临时的，到时间后会删除
func (s *Service) CreateGuestChat(ctx context.Context, createGuestChatRequest *schema.ChatGuestCreateRequest) (*entity.Chat, error) {
	var chat = entity.Chat{}

	chat.Name = createGuestChatRequest.Name
	chat.AssistantId = createGuestChatRequest.AssistantId
	chat.Owner = schema.OwnerGuest
	chat.GuestId = &createGuestChatRequest.GuestID

	var t = time.Now().Add(time.Hour * 24)

	chat.ExpiredAt = &t

	_, err := s.x.Context(ctx).Insert(&chat)

	return &chat, err
}

func (s *Service) GetChat(ctx context.Context, id schema.EntityId) (*entity.Chat, error) {
	var chat entity.Chat
	_, err := s.x.Context(ctx).
		ID(id).
		Get(&chat)
	return &chat, err
}

func (s *Service) GetChatWithAssistant(ctx context.Context, chatId int64) (*entity.ChatWithAssistant, error) {
	var chat entity.ChatWithAssistant
	_, err := s.x.Context(ctx).
		Join("INNER", "assistants", "assistants.id = chats.assistant_id").
		ID(chatId).
		Get(&chat)
	return &chat, err
}

func (s *Service) ListChat(ctx context.Context) ([]*entity.Chat, error) {
	var chats []*entity.Chat
	err := s.x.Context(ctx).Find(&chats)
	return chats, err
}

func (s *Service) DeleteChat(ctx context.Context, chat *entity.Chat) error {
	// 确保所有的 message 已被删除
	count, err := s.cm.CountChatMessage(ctx, chat)

	if count > 0 {
		return consts.ErrChatCanNotDeleteBecauseNotCleared
	}

	_, err = s.x.Context(ctx).ID(chat.Id).Delete(&entity.Chat{})
	return err
}

func (s *Service) DeleteChats(ctx context.Context, chat ...*entity.Chat) error {
	//  至少提供一个
	if len(chat) == 0 {
		return errors.New("no chat provided")
	}

	for _, v := range chat {
		// 确保所有的 message 已被删除
		count, err := s.cm.CountChatMessage(ctx, v)
		if err != nil {
			return err
		}

		if count > 0 {
			return consts.ErrChatCanNotDeleteBecauseNotCleared
		}

		_, err = s.x.Context(ctx).ID(v.Id).Delete(&entity.Chat{})
	}

	return nil
}

func (s *Service) DeleteChatFromUserId(ctx context.Context, id schema.EntityId, userId schema.UserId) error {
	count, err := s.x.Context(ctx).Where("id = ?", id).Where("user_id = ?", userId).Delete(&entity.Chat{})
	if err != nil {
		return err
	}
	if count == 0 {
		return consts.ErrChatNotFound
	}

	return nil
}

func (s *Service) Exists(ctx context.Context, id schema.EntityId) (bool, error) {
	count, err := s.x.Context(ctx).Where("id = ?", id).Count(&entity.Chat{})
	return count > 0, err
}

func (s *Service) ListChatFromUserId(ctx context.Context, userId schema.UserId) ([]*entity.Chat, error) {
	var chats []*entity.Chat
	err := s.x.Context(ctx).Where("user_id = ?", userId).Find(&chats)
	return chats, err
}

func (s *Service) ListChatFromAssistantIdWithAssistant(ctx context.Context, assistant *entity.Assistant) ([]*entity.Chat, error) {
	var chats []*entity.Chat
	err := s.x.Context(ctx).Where("assistant_id = ?", assistant.Id).Find(&chats)
	return chats, err
}

func (s *Service) ListChatFromAssistantByPage(ctx context.Context, assistant *entity.Assistant, page int) ([]*entity.Chat, error) {
	var chats []*entity.Chat
	var limit = 20
	err := s.x.Context(ctx).Where("assistant_id = ?", assistant.Id).Limit(limit, (page-1)*limit).Find(&chats)
	return chats, err
}

func (s *Service) ListChatFromGuestId(ctx context.Context, guestId string) ([]*entity.Chat, error) {
	var chats []*entity.Chat
	err := s.x.Context(ctx).Where("owner = ?", schema.OwnerGuest.String()).Where("guest_id = ?", guestId).Find(&chats)
	return chats, err
}

func (s *Service) ListChatFromGuestByPage(ctx context.Context, guestId string, page int) ([]*entity.Chat, error) {
	var chats []*entity.Chat
	var limit = 20
	err := s.x.Context(ctx).Where("owner", schema.OwnerGuest.String()).Where("guest_id = ?", guestId).Limit(limit, (page-1)*limit).Find(&chats)
	return chats, err
}

func (s *Service) DeleteExpiredChats(ctx context.Context, beforeTime time.Time) error {
	// 防止瞬时压力过大，一次删除固定数量
	var num = 1000
	_, err := s.x.Context(ctx).Where("expired_at < ?", beforeTime).Limit(num).Delete(&entity.Chat{})
	return err
}
