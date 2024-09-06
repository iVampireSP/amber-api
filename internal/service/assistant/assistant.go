package assistant

import (
	"context"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"

	"github.com/tmc/langchaingo/llms"
)

func (s *Service) ListAssistant(ctx context.Context) ([]*entity.Assistant, error) {
	return s.dao.WithContext(ctx).Assistant.Find()
}

func (s *Service) ListAssistantFromUserId(ctx context.Context, userId schema.UserId) ([]*entity.Assistant, error) {
	assistants, err := s.dao.WithContext(ctx).Assistant.Where(s.dao.Assistant.UserId.Eq(userId.String())).Find()

	return assistants, err
}

func (s *Service) GetAssistant(ctx context.Context, id schema.EntityId) (*entity.Assistant, error) {
	assistant, err := s.dao.WithContext(ctx).Assistant.Where(s.dao.Assistant.Id.Eq(uint(id))).First()

	return assistant, err
}
func (s *Service) GetAssistantTool(ctx context.Context, assistant *entity.Assistant) ([]*entity.AssistantTool, error) {
	assistantTool, err := s.dao.WithContext(ctx).AssistantTool.Where(s.dao.AssistantTool.AssistantId.Eq(uint(assistant.Id))).
		Preload(s.dao.AssistantTool.Tool).Find()

	return assistantTool, err
}

func (s *Service) CreateAssistant(ctx context.Context, assistantReq *schema.AssistantCreateRequest) (*entity.Assistant, error) {
	var assistant entity.Assistant
	assistant.UserId = assistantReq.UserId
	assistant.Name = assistantReq.Name
	assistant.Prompt = assistantReq.Prompt
	assistant.Description = assistantReq.Description
	assistant.DisableDefaultPrompt = assistantReq.DisableDefaultPrompt

	// 啊，这样也行
	return &assistant, s.dao.WithContext(ctx).Assistant.Create(&assistant)
}

func (s *Service) UpdateAssistant(ctx context.Context, assistant *entity.Assistant) error {
	_, err := s.dao.WithContext(ctx).Assistant.Where(s.dao.Assistant.Id.Eq(uint(assistant.Id))).UpdateSimple(
		s.dao.Assistant.Name.Value(assistant.Name),
		s.dao.Assistant.Description.Value(assistant.Description),
		s.dao.Assistant.Prompt.Value(assistant.Prompt),
		s.dao.Assistant.DisableDefaultPrompt.Value(assistant.DisableDefaultPrompt),
		s.dao.Assistant.DisableMemory.Value(assistant.DisableMemory),
		s.dao.Assistant.EnableMemoryForAssistantShare.Value(assistant.EnableMemoryForAssistantShare),
	)
	if err != nil {
		return err
	}

	return err
}

func (s *Service) DeleteAssistant(ctx context.Context, id schema.EntityId) error {
	// assistant
	assistant, err := s.GetAssistant(ctx, id)
	if err != nil {
		return err
	}

	//// 如果绑定了资料库，则不能删除
	//if assistant.LibraryId != nil {
	//	return consts.ErrAssistantHasBindLibraryCantDelete
	//}

	// 如果已经绑定过工具，则不能删除
	toolEntity, err := s.ListAssistantTool(ctx, assistant)
	if err != nil {
		return err
	}

	if len(toolEntity) > 0 {
		return consts.ErrAssistantHasBindToolCantDelete
	}

	_, err = s.dao.WithContext(ctx).Assistant.Delete(assistant)

	return err
}

func (s *Service) ToLLMTool(ctx context.Context, assistant *entity.Assistant) ([]llms.Tool, error) {
	var toolList []llms.Tool

	toolEntity, err := s.dao.WithContext(ctx).AssistantTool.Where(s.dao.AssistantTool.AssistantId.Eq(uint(assistant.Id))).Preload(s.dao.AssistantTool.Tool).Find()
	if err != nil {
		return nil, err
	}

	//toolEntity, err := s.ListAssistantToolWithType(ctx, assistant)
	//if err != nil {
	//	return nil, err
	//}

	// 转换格式
	for _, v := range toolEntity {
		var fnData = v.Tool.Data
		var llmTool = llms.Tool{
			Type: "function",
		}

		for _, v := range fnData.ToolFunctions {
			for _, j := range v.Functions {
				llmTool.Function = &llms.FunctionDefinition{
					Name:        j.Name,
					Description: j.Description,
					Parameters:  j.Parameters,
				}
			}

			if llmTool.Function != nil {
				toolList = append(toolList, llmTool)
			}
		}
	}
	return toolList, nil
}

func (s *Service) GetAssistantFromCtx(ctx context.Context) *entity.Assistant {
	assistantEntity := ctx.Value(consts.AuthAssistantShareMiddlewareKey)

	return assistantEntity.(*entity.Assistant)
}
