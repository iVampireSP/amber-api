package assistant

import (
	"context"
	"rag-new/internal/entity"
	"rag-new/pkg/consts"
)

func (s *Service) BindTool(ctx context.Context, assistant *entity.Assistant, tool *entity.Tool) (*entity.AssistantTool, error) {
	// 检测是否绑定
	count, err := s.dao.WithContext(ctx).AssistantTool.Where(s.dao.AssistantTool.AssistantId.Eq(uint(assistant.Id))).Where(s.dao.AssistantTool.ToolId.Eq(uint(tool.Id))).Count()
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, consts.ErrAssistantAlreadyBindTheTool
	}

	// 检测 tool 是否存在
	toolExists, err := s.ToolExists(ctx, tool)
	if err != nil {
		return nil, err
	}
	if !toolExists {
		return nil, consts.ErrToolNotFound
	}

	// 检测 assistant 是否存在
	assistantExists, err := s.AssistantExists(ctx, assistant)
	if err != nil {
		return nil, err
	}
	if !assistantExists {
		return nil, consts.ErrAssistantNotFound
	}

	// 检测是否绑定
	var toolBind = &entity.AssistantTool{}

	toolBind.ToolId = tool.Id
	toolBind.AssistantId = assistant.Id

	// 保存
	err = s.dao.WithContext(ctx).AssistantTool.Save(toolBind)

	return toolBind, err
}

func (s *Service) UnbindTool(ctx context.Context, assistant *entity.Assistant, tool *entity.Tool) error {
	// 检测是否被 bind（记录是否存在）
	assistantTool, err := s.dao.WithContext(ctx).
		AssistantTool.Where(s.dao.AssistantTool.AssistantId.Eq(uint(assistant.Id))).
		Where(s.dao.AssistantTool.ToolId.Eq(uint(tool.Id))).
		First()
	if err != nil {
		return err
	}

	count, err := s.dao.WithContext(ctx).AssistantTool.Where(s.dao.AssistantTool.AssistantId.Eq(uint(assistant.Id))).
		Where(s.dao.AssistantTool.ToolId.Eq(uint(tool.Id))).Count()
	if err != nil {
		return err
	}

	if count == 0 {
		return consts.ErrToolNotBind
	}

	_, err = s.dao.WithContext(ctx).AssistantTool.Where(s.dao.AssistantTool.AssistantId.Eq(uint(assistant.Id))).Where(s.dao.AssistantTool.ToolId.Eq(uint(tool.Id))).Delete(assistantTool)

	return err
}

func (s *Service) ToolExists(ctx context.Context, tool *entity.Tool) (bool, error) {
	count, err := s.dao.WithContext(ctx).Tool.Where(s.dao.Tool.Id.Eq(uint(tool.Id))).Count()
	return count > 0, err
}

func (s *Service) AssistantExists(ctx context.Context, assistant *entity.Assistant) (bool, error) {
	count, err := s.dao.WithContext(ctx).Assistant.Where(s.dao.Assistant.Id.Eq(uint(assistant.Id))).Count()
	return count > 0, err
}

func (s *Service) ListAssistantTool(ctx context.Context, assistant *entity.Assistant) ([]*entity.AssistantTool, error) {
	assistantTools, err := s.dao.WithContext(ctx).AssistantTool.Where(s.dao.AssistantTool.AssistantId.Eq(uint(assistant.Id))).Preload(s.dao.AssistantTool.Tool).Find()

	return assistantTools, err
}

//func (s *Service) ListAssistantToolWithType(ctx context.Context, assistant *entity.Assistant) ([]*entity.AssistantToolType, error) {
//	var assistantToolType = make([]*entity.AssistantToolType, 0)
//	err := s.x.Context(ctx).
//		Join("INNER", "assistants", "assistants.id = assistant_tools.assistant_id").
//		Join("INNER", "tools", "tools.id = assistant_tools.tool_id").
//		Where("assistant_id = ?", assistant.Id).Find(&assistantToolType)
//	return assistantToolType, err
//}
