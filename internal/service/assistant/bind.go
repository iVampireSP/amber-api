package assistant

import (
	"context"
	"rag-new/internal/entity"
	"rag-new/pkg/consts"
)

func (s *Service) BindTool(ctx context.Context, assistantId int64, toolId int64) (*entity.AssistantTool, error) {
	var toolBind = &entity.AssistantTool{}
	_, err := s.x.Context(ctx).Where("assistant_id = ?", assistantId).Where("tool_id = ?", toolId).Get(toolBind)
	if err != nil {
		return nil, err
	}

	// 检测是否被 bind（记录是否存在）
	if toolBind.ID != consts.NoRecord {
		return nil, consts.ErrAssistantAlreadyBindTheTool
	}

	// 检测 tool 是否存在
	toolExists, err := s.ToolExists(ctx, toolId)
	if err != nil {
		return nil, err
	}
	if !toolExists {
		return nil, consts.ErrToolNotFound
	}

	// 检测 assistant 是否存在
	assistantExists, err := s.AssistantExists(ctx, assistantId)
	if err != nil {
		return nil, err
	}
	if !assistantExists {
		return nil, consts.ErrAssistantNotFound
	}

	toolBind.ToolId = toolId
	toolBind.AssistantId = assistantId

	_, err = s.x.Context(ctx).Insert(toolBind)

	return toolBind, err
}

func (s *Service) UnbindTool(ctx context.Context, assistantId int64, toolId int64) error {
	var assistantTool = &entity.AssistantTool{}

	// 检测是否被 bind（记录是否存在）
	_, err := s.x.Context(ctx).Where("assistant_id = ?", assistantId).Where("tool_id = ?", toolId).Get(assistantTool)
	if err != nil {
		return err
	}

	if assistantTool.ID == consts.NoRecord {
		return consts.ErrAssistantNotFound
	}

	_, err = s.x.Context(ctx).Where("assistant_id = ?", assistantId).Where("tool_id = ?", toolId).Delete(&entity.AssistantTool{})

	return err
}

func (s *Service) ToolExists(ctx context.Context, id int64) (bool, error) {
	count, err := s.x.Context(ctx).Where("id = ?", id).Count(&entity.Tool{})
	return count > 0, err
}

func (s *Service) AssistantExists(ctx context.Context, id int64) (bool, error) {
	count, err := s.x.Context(ctx).Where("id = ?", id).Count(&entity.Assistant{})
	return count > 0, err
}

func (s *Service) ListAssistantTool(ctx context.Context, assistantId int64) ([]*entity.AssistantTool, error) {
	var assistantTools []*entity.AssistantTool
	err := s.x.Context(ctx).Where("assistant_id = ?", assistantId).Find(&assistantTools)
	return assistantTools, err
}
