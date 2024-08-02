package assistant

import (
	"context"
	"rag-new/internal/entity"
	"rag-new/pkg/consts"
)

func (s *Service) BindTool(ctx context.Context, assistant *entity.Assistant, tool *entity.Tool) (*entity.AssistantTool, error) {
	var toolBind = &entity.AssistantTool{}
	_, err := s.x.Context(ctx).Where("assistant_id = ?", assistant.ID).Where("tool_id = ?", tool.ID).Get(toolBind)
	if err != nil {
		return nil, err
	}

	// 检测是否被 bind（记录是否存在）
	if toolBind.ID != consts.NoRecord {
		return nil, consts.ErrAssistantAlreadyBindTheTool
	}

	// 检测 tool 是否存在
	toolExists, err := s.ToolExists(ctx, assistant)
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

	toolBind.ToolId = tool.ID
	toolBind.AssistantId = assistant.ID

	_, err = s.x.Context(ctx).Cols("tool_id", "assistant_id").Insert(toolBind)
	//_, err = s.x.Context(ctx).Cols("tool_id", "assistant_id").Insert(toolBind)

	return toolBind, err
}

func (s *Service) UnbindTool(ctx context.Context, assistant *entity.Assistant, tool *entity.Tool) error {
	var assistantTool = &entity.AssistantTool{}

	// 检测是否被 bind（记录是否存在）
	_, err := s.x.Context(ctx).Where("assistant_id = ?", assistant.ID).Where("tool_id = ?", tool.ID).Get(assistantTool)
	if err != nil {
		return err
	}

	if assistantTool.ID == consts.NoRecord {
		return consts.ErrAssistantNotFound
	}

	_, err = s.x.Context(ctx).Where("assistant_id = ?", assistant.ID).Where("tool_id = ?", tool.ID).Delete(&entity.AssistantTool{})

	return err
}

func (s *Service) ToolExists(ctx context.Context, assistant *entity.Assistant) (bool, error) {
	count, err := s.x.Context(ctx).Where("id = ?", assistant.ID).Count(&entity.Tool{})
	return count > 0, err
}

func (s *Service) AssistantExists(ctx context.Context, assistant *entity.Assistant) (bool, error) {
	count, err := s.x.Context(ctx).Where("id = ?", assistant.ID).Count(&entity.Assistant{})
	return count > 0, err
}

func (s *Service) ListAssistantTool(ctx context.Context, assistant *entity.Assistant) ([]*entity.AssistantTool, error) {
	var assistantTools []*entity.AssistantTool
	err := s.x.Context(ctx).
		Where("assistant_id = ?", assistant.ID).Find(&assistantTools)
	return assistantTools, err
}

func (s *Service) ListAssistantToolWithType(ctx context.Context, assistant *entity.Assistant) ([]*entity.AssistantToolType, error) {
	var assistantToolType = make([]*entity.AssistantToolType, 0)
	err := s.x.Context(ctx).
		Join("INNER", "tools", "tools.id = assistant_tools.tool_id").
		Join("INNER", "assistants", "assistants.id = assistant_tools.assistant_id").
		Where("assistant_id = ?", assistant.ID).Find(&assistantToolType)
	return assistantToolType, err
}
