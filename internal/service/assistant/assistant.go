package assistant

import (
	"context"
	"github.com/tmc/langchaingo/llms"
	"gorm.io/gen/field"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
)

var defaultTemperature = 0.7

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

	// 判断 assistantReq.Temperature 是否设置了值
	if assistantReq.Temperature == 0 {
		assistantReq.Temperature = defaultTemperature // 如果没有设置，则使用默认值 0.7
	} else if assistantReq.Temperature < 0.1 || assistantReq.Temperature > 1 {
		// 检查小数点后是否只有一位
		if float64(int(assistantReq.Temperature*10)) != assistantReq.Temperature*10 {
			// 设置默认
			assistantReq.Temperature = defaultTemperature
		}
	}

	assistant.Temperature = assistantReq.Temperature

	return &assistant, s.dao.WithContext(ctx).Assistant.Create(&assistant)
}

func (s *Service) UpdateAssistant(ctx context.Context, assistant *entity.Assistant) error {
	var assignExpr = []field.AssignExpr{
		s.dao.Assistant.Name.Value(assistant.Name),
		s.dao.Assistant.Description.Value(assistant.Description),
		s.dao.Assistant.Prompt.Value(assistant.Prompt),
		s.dao.Assistant.DisableDefaultPrompt.Value(assistant.DisableDefaultPrompt),
		s.dao.Assistant.DisableInternetSearch.Value(assistant.DisableInternetSearch),
		s.dao.Assistant.DisableWebBrowsing.Value(assistant.DisableWebBrowsing),
		s.dao.Assistant.DisableMemory.Value(assistant.DisableMemory),
		s.dao.Assistant.EnableMemoryForAssistantAPI.Value(assistant.EnableMemoryForAssistantAPI),
		s.dao.Assistant.Public.Value(assistant.Public),
	}
	// 这里不能直接设置 library_id
	if assistant.LibraryId != nil {
		assignExpr = append(assignExpr, s.dao.Assistant.LibraryId.Value(uint(*assistant.LibraryId)))
	} else {
		assignExpr = append(assignExpr, s.dao.Assistant.LibraryId.Null())
	}

	// 判断 assistantReq.Temperature 是否设置了值
	if assistant.Temperature == 0 || assistant.Temperature < 0.1 || assistant.Temperature > 1 {
		assistant.Temperature = defaultTemperature // 如果没有设置，则使用默认值 0.7
	} else {
		// 检查小数点后是否只有一位
		if float64(int(assistant.Temperature*10)) != assistant.Temperature*10 {
			// 设置默认
			assistant.Temperature = defaultTemperature
		}
	}

	assignExpr = append(assignExpr, s.dao.Assistant.Temperature.Value(assistant.Temperature))

	_, err := s.dao.WithContext(ctx).Assistant.Where(s.dao.Assistant.Id.Eq(uint(assistant.Id))).
		UpdateSimple(assignExpr...)

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

	// 如果助理在聊天消息列表中，则将他们设置为 null
	_, err = s.dao.WithContext(ctx).ChatMessage.Where(s.dao.ChatMessage.AssistantId.Eq(uint(assistant.Id))).
		UpdateSimple(s.dao.ChatMessage.AssistantId.Null())

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

func (s *Service) BindLibrary(ctx context.Context, assistant *entity.Assistant, library *entity.Library) error {
	// 重复绑定不处理
	if assistant.LibraryId != nil && *assistant.LibraryId == library.Id {
		return nil
	}
	//if library.UserId != assistant.UserId {
	//	return consts.ErrPermissionDenied
	//}

	_, err := s.dao.WithContext(ctx).Assistant.Where(s.dao.Assistant.Id.Eq(uint(assistant.Id))).
		Update(s.dao.Assistant.LibraryId, library.Id)

	return err
}

func (s *Service) UnbindLibrary(ctx context.Context, assistant *entity.Assistant, library *entity.Library) error {
	//if library.UserId != assistant.UserId {
	//	return consts.ErrPermissionDenied
	//}

	if assistant.LibraryId != nil && *assistant.LibraryId != library.Id {
		return consts.ErrPermissionDenied
	}

	_, err := s.dao.WithContext(ctx).Assistant.Where(s.dao.Assistant.Id.Eq(uint(assistant.Id))).
		Update(s.dao.Assistant.LibraryId, nil)

	return err
}

func (s *Service) IncrementTotalTokenUsage(ctx context.Context, assistant *entity.Assistant, token int64) error {
	_, err := s.dao.WithContext(ctx).Assistant.Where(s.dao.Assistant.Id.Eq(assistant.Id.Uint())).
		Update(s.dao.Assistant.TotalTokenUsage, s.dao.Assistant.TotalTokenUsage.Add(token))

	return err
}
