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

	err = s.dao.WithContext(ctx).Chat.Create(&chat)

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

	err := s.dao.WithContext(ctx).Chat.Create(&chat)

	return &chat, err
}

func (s *Service) GetChat(ctx context.Context, id schema.EntityId) (*entity.Chat, error) {
	chat, err := s.dao.WithContext(ctx).Chat.Preload(s.dao.Chat.Assistant).Where(s.dao.Chat.Id.Eq(uint(id))).First()

	return chat, err
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

		_, err = s.dao.WithContext(ctx).Chat.Where(s.dao.Chat.Id.Eq(uint(v.Id))).Delete()
		if err != nil {
			return err
		}

	}

	return nil
}

func (s *Service) DeleteChatFromUserId(ctx context.Context, id schema.EntityId, userId schema.UserId) error {
	_, err := s.dao.WithContext(ctx).Chat.Where(s.dao.Chat.Id.Eq(uint(id))).Where(s.dao.Chat.UserId.Eq(userId.String())).Delete()

	return err
}

func (s *Service) ListChatFromUserId(ctx context.Context, userId schema.UserId) ([]*entity.Chat, error) {
	chats, err := s.dao.WithContext(ctx).Chat.Where(s.dao.Chat.UserId.Eq(userId.String())).Find()

	return chats, err
}

func (s *Service) ListChatFromAssistantIdWithAssistant(ctx context.Context, assistant *entity.Assistant) ([]*entity.Chat, error) {
	chats, err := s.dao.WithContext(ctx).Chat.Where(s.dao.Chat.AssistantId.Eq(uint(assistant.Id))).Find()

	return chats, err
}

func (s *Service) ListChatFromAssistantByPage(ctx context.Context, assistant *entity.Assistant, page int) ([]*entity.Chat, error) {
	var limit = 20

	chats, _, err := s.dao.WithContext(ctx).Chat.Where(s.dao.Chat.AssistantId.Eq(uint(assistant.Id))).FindByPage((page-1)*limit, limit)

	return chats, err
}

func (s *Service) ListChatFromGuestId(ctx context.Context, guestId string) ([]*entity.Chat, error) {
	chats, err := s.dao.WithContext(ctx).Chat.Where(s.dao.Chat.Owner.Eq(schema.OwnerGuest.String())).Where(s.dao.Chat.GuestId.Eq(guestId)).Find()

	return chats, err
}

func (s *Service) DeleteExpiredChats(ctx context.Context, beforeTime time.Time) error {
	// 防止瞬时压力过大，一次删除固定数量
	var num = 1000

	_, err := s.dao.Chat.WithContext(ctx).Where(s.dao.Chat.ExpiredAt.Lt(beforeTime)).Limit(num).Delete()
	if err != nil {
		return err
	}

	return err
}
