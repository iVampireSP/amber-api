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
	var assistants = make([]*entity.Assistant, 0)
	err := s.x.Context(ctx).Where("user_id = ?", userId).Find(&assistants)
	return assistants, err
}

func (s *Service) GetAssistant(ctx context.Context, id int64) (*entity.Assistant, error) {
	assistant := new(entity.Assistant)
	_, err := s.x.Context(ctx).ID(id).Get(assistant)
	return assistant, err
}

func (s *Service) CreateAssistant(ctx context.Context, assistantReq *schema.AssistantCreateRequest) (*entity.Assistant, error) {
	var assistant entity.Assistant
	assistant.UserId = assistantReq.UserId
	assistant.Name = assistantReq.Name
	assistant.Prompt = assistantReq.Prompt
	assistant.Description = assistantReq.Description

	_, err := s.x.Context(ctx).Insert(&assistant)
	return &assistant, err
}

//
//func (s *Service) UpdateAssistant(ctx context.Context, assistant *entity.Assistant) error {
//	_, err := s.x.Context(ctx).ID(assistant.Id).AllCols().Update(assistant)
//	return err
//}

func (s *Service) DeleteAssistant(ctx context.Context, id int64) error {
	_, err := s.x.Context(ctx).ID(id).Delete(&entity.Assistant{})
	return err
}
