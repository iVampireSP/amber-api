package assistant

import (
	"context"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/random"
)

func (s *Service) CrateApiKey(ctx context.Context, assistant *entity.Assistant) (*entity.AssistantApiKey, error) {
	var assistantApiKey = &entity.AssistantApiKey{}

	assistantApiKey.AssistantId = assistant.Id
	// 生成 Secret
	assistantApiKey.Secret = random.String(40)

	//// 检测 assistant 是否存在
	//assistant, err := s.GetAssistant(ctx, assistant.Id)
	//if err != nil {
	//	return nil, err
	//}
	//if assistant.Id == consts.NoRecord {
	//	return nil, consts.ErrAssistantNotFound
	//}

	err := s.dao.WithContext(ctx).AssistantApiKey.Create(assistantApiKey)

	return assistantApiKey, err
}

func (s *Service) GetApiKey(ctx context.Context, assistantApiKeyId schema.EntityId) (*entity.AssistantApiKey, error) {
	assistantApiKey, err := s.dao.WithContext(ctx).AssistantApiKey.Where(s.dao.AssistantApiKey.Id.Eq(uint(assistantApiKeyId))).Preload(s.dao.AssistantApiKey.Assistant).
		First()

	if err != nil {
		return nil, err
	}

	return assistantApiKey, err
}

func (s *Service) GetApiKeyBySecret(ctx context.Context, secret string) (*entity.AssistantApiKey, error) {
	assistantApiKey, err := s.dao.WithContext(ctx).AssistantApiKey.Where(s.dao.AssistantApiKey.Secret.Eq(secret)).Preload(s.dao.AssistantApiKey.Assistant).
		First()

	if err != nil {
		return nil, err
	}

	return assistantApiKey, err
}

// ListApiKey 获取当前助理的所有分享
func (s *Service) ListApiKey(ctx context.Context, assistant *entity.Assistant) ([]*entity.AssistantApiKey, error) {
	assistantApiKeys, err := s.dao.WithContext(ctx).AssistantApiKey.Where(s.dao.AssistantApiKey.AssistantId.Eq(uint(assistant.Id))).Find()

	return assistantApiKeys, err
}

func (s *Service) DeleteApiKey(ctx context.Context, assistantApiKey *entity.AssistantApiKey) error {
	_, err := s.dao.WithContext(ctx).AssistantApiKey.Delete(assistantApiKey)

	return err
}

func (s *Service) DeleteAllApiKey(ctx context.Context, assistant *entity.Assistant) error {
	_, err := s.dao.WithContext(ctx).AssistantApiKey.Where(s.dao.AssistantApiKey.AssistantId.Eq(uint(assistant.Id))).Delete()

	return err
}

//func (s *Service) UpdateApiKey(ctx context.Context, assistant *entity.Assistant) error {
//	_, err := s.x.Context(ctx).ID(assistant.Id).AllCols().Update(assistant)
//	return err
//}
