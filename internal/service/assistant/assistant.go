package assistant

import (
	"context"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
)

func (s *Service) ListAssistant(ctx context.Context) ([]*entity.Assistant, error) {
	var assistants []*entity.Assistant
	err := s.x.Context(ctx).Find(&assistants)
	return assistants, err
}

func (s *Service) ListAssistantFromUserId(ctx context.Context, userId schema.UserId) ([]*entity.Assistant, error) {
	var assistants []*entity.Assistant
	err := s.x.Context(ctx).Where("user_id = ?", userId).Find(&assistants)
	return assistants, err
}

func (s *Service) GetAssistant(ctx context.Context, id int64) (*entity.Assistant, error) {
	assistant := new(entity.Assistant)
	_, err := s.x.Context(ctx).ID(id).Get(assistant)
	return assistant, err
}
