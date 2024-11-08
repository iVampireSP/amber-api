package assistant

import (
	"context"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
)

func (s *Service) CreateScenePrompt(ctx context.Context, label string, prompt string, assistantEntity *entity.Assistant) (*entity.ScenePrompt, error) {
	sp := &entity.ScenePrompt{
		AssistantId: assistantEntity.Id,
		Label:       label,
		Prompt:      prompt,
	}

	err := s.dao.WithContext(ctx).ScenePrompt.Create(sp)

	return sp, err
}

func (s *Service) GetScenePromptById(ctx context.Context, scenePromptId schema.EntityId) (*entity.ScenePrompt, error) {
	return s.dao.WithContext(ctx).ScenePrompt.Where(s.dao.ScenePrompt.Id.Eq(uint(scenePromptId))).First()
}

func (s *Service) GetAssistantScenePrompts(ctx context.Context, assistantEntity *entity.Assistant) ([]*entity.ScenePrompt, error) {
	return s.dao.WithContext(ctx).ScenePrompt.Where(s.dao.ScenePrompt.AssistantId.Eq(assistantEntity.Id.Uint())).Find()
}

func (s *Service) CountAssistantScenePrompts(ctx context.Context, assistantEntity *entity.Assistant) (int64, error) {
	var q = s.dao.WithContext(ctx).ScenePrompt.Where(s.dao.ScenePrompt.AssistantId.Eq(assistantEntity.Id.Uint()))

	return q.Count()
}

func (s *Service) GetScenePromptExistsByLabel(ctx context.Context, label string, assistantEntity *entity.Assistant) (bool, error) {
	var q = s.dao.WithContext(ctx).ScenePrompt.Where(s.dao.ScenePrompt.Label.Eq(label))

	if assistantEntity != nil {
		q = q.Where(s.dao.ScenePrompt.AssistantId.Eq(assistantEntity.Id.Uint()))
	}

	c, err := q.Count()

	return c > 0, err
}

func (s *Service) GetScenePromptByLabel(ctx context.Context, label string, assistantEntity *entity.Assistant) (*entity.ScenePrompt, error) {
	var q = s.dao.WithContext(ctx).ScenePrompt.Where(s.dao.ScenePrompt.Label.Eq(label))

	if assistantEntity != nil {
		q = q.Where(s.dao.ScenePrompt.AssistantId.Eq(assistantEntity.Id.Uint()))
	}

	return q.First()
}

func (s *Service) DeleteScenePrompt(ctx context.Context, scenePrompt *entity.ScenePrompt) error {
	_, err := s.dao.WithContext(ctx).ScenePrompt.Delete(scenePrompt)

	return err
}
