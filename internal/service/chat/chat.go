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

	if a.ID == consts.NoRecord || a.UserId != createChatRequest.UserId {
		return nil, consts.ErrAssistantNotFound
	}

	chat.Name = createChatRequest.Name
	chat.AssistantId = createChatRequest.AssistantId
	chat.UserId = createChatRequest.UserId
	chat.Owner = schema.OwnerUser

	_, err = s.x.Context(ctx).Insert(&chat)

	return &chat, err
}

// CreateGuestChat 用于创建访客对话，这个对话应是临时的，到时间后会删除
func (s *Service) CreateGuestChat(ctx context.Context, createGuestChatRequest *schema.ChatGuestCreateRequest) (*entity.Chat, error) {
	var chat = &entity.Chat{}

	chat.Name = createGuestChatRequest.Name
	chat.AssistantId = createGuestChatRequest.AssistantId
	chat.Owner = schema.OwnerGuest
	chat.GuestId = createGuestChatRequest.GuestID

	chat.ExpiredAt = time.Now().Add(time.Hour * 24)

	_, err := s.x.Context(ctx).Insert(chat)

	return chat, err
}

func (s *Service) GetChat(ctx context.Context, id int64) (*entity.Chat, error) {
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

	_, err = s.x.Context(ctx).ID(chat.ID).Delete(&entity.Chat{})
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

		_, err = s.x.Context(ctx).ID(v.ID).Delete(&entity.Chat{})
	}

	return nil
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

func (s *Service) ListChatFromAssistantIdWithAssistant(ctx context.Context, assistant *entity.Assistant) ([]*entity.Chat, error) {
	var chats []*entity.Chat
	err := s.x.Context(ctx).Where("assistant_id = ?", assistant.ID).Find(&chats)
	return chats, err
}

func (s *Service) ListChatFromAssistantByPage(ctx context.Context, assistant *entity.Assistant, page int) ([]*entity.Chat, error) {
	var chats []*entity.Chat
	var limit = 20
	err := s.x.Context(ctx).Where("assistant_id = ?", assistant.ID).Limit(limit, (page-1)*limit).Find(&chats)
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
