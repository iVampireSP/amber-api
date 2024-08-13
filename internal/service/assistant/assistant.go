package assistant

import (
	"context"
	"github.com/tmc/langchaingo/llms"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
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
	if assistant.Id == consts.NoRecord {
		return nil, consts.ErrAssistantNotFound
	}
	return assistant, err
}

func (s *Service) CreateAssistant(ctx context.Context, assistantReq *schema.AssistantCreateRequest) (*entity.Assistant, error) {
	var assistant entity.Assistant
	assistant.UserId = assistantReq.UserId
	assistant.Name = assistantReq.Name
	assistant.Prompt = assistantReq.Prompt
	assistant.Description = assistantReq.Description
	assistant.DisableDefaultPrompt = assistantReq.DisableDefaultPrompt

	_, err := s.x.Context(ctx).Insert(&assistant)
	return &assistant, err
}

func (s *Service) UpdateAssistant(ctx context.Context, assistant *entity.Assistant) error {
	_, err := s.x.Context(ctx).ID(assistant.Id).Cols("name", "description", "prompt", "disable_default_prompt").Update(assistant)
	return err
}

func (s *Service) DeleteAssistant(ctx context.Context, id int64) error {
	// assistant
	assistant, err := s.GetAssistant(ctx, id)
	if err != nil {
		return err
	}

	// 如果已经绑定过工具，则不能删除
	toolEntity, err := s.ListAssistantTool(ctx, assistant)
	if err != nil {
		return err
	}

	if len(toolEntity) > 0 {
		return consts.ErrAssistantHasBindToolCantDelete
	}

	_, err = s.x.Context(ctx).ID(id).Delete(&entity.Assistant{})

	return err
}

func (s *Service) ToLLMTool(ctx context.Context, assistant *entity.Assistant) ([]llms.Tool, error) {
	var toolList []llms.Tool

	toolEntity, err := s.ListAssistantToolWithType(ctx, assistant)
	if err != nil {
		return nil, err
	}

	// 转换格式
	for _, v := range toolEntity {
		var fnData = v.Tool.Data
		var llmTool = llms.Tool{
			Type: "function",
		}

		for _, v := range fnData.ToolFunctions {
			for _, j := range v.Function {
				llmTool.Function = &llms.FunctionDefinition{
					Name:        j.Name,
					Description: j.Description,
					Parameters:  j.Parameters,
				}
			}
			toolList = append(toolList, llmTool)
		}
	}
	return toolList, nil
}
