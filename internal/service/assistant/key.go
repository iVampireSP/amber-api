package assistant

import (
	"context"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/random"
)

func (s *Service) CreateKey(ctx context.Context, assistant *entity.Assistant) (*entity.AssistantKey, error) {
	var assistantKey = &entity.AssistantKey{}

	assistantKey.AssistantId = assistant.Id
	// 生成 Secret
	assistantKey.Secret = random.String(40)

	//// 检测 assistant 是否存在
	//assistant, err := s.GetAssistant(ctx, assistant.Id)
	//if err != nil {
	//	return nil, err
	//}
	//if assistant.Id == consts.NoRecord {
	//	return nil, consts.ErrAssistantNotFound
	//}

	err := s.dao.WithContext(ctx).AssistantKey.Create(assistantKey)

	return assistantKey, err
}

func (s *Service) GetKey(ctx context.Context, assistantKeyId schema.EntityId) (*entity.AssistantKey, error) {
	assistantKey, err := s.dao.WithContext(ctx).AssistantKey.Where(s.dao.AssistantKey.Id.Eq(uint(assistantKeyId))).Preload(s.dao.AssistantKey.Assistant).
		First()

	if err != nil {
		return nil, err
	}

	return assistantKey, err
}

func (s *Service) GetByKey(ctx context.Context, secret string) (*entity.AssistantKey, error) {
	assistantKey, err := s.dao.WithContext(ctx).AssistantKey.Where(s.dao.AssistantKey.Secret.Eq(secret)).Preload(s.dao.AssistantKey.Assistant).
		First()

	if err != nil {
		return nil, err
	}

	return assistantKey, err
}

// ListKey 获取当前助理的所有的密钥
func (s *Service) ListKey(ctx context.Context, assistant *entity.Assistant) ([]*entity.AssistantKey, error) {
	assistantKeys, err := s.dao.WithContext(ctx).AssistantKey.Where(s.dao.AssistantKey.AssistantId.Eq(uint(assistant.Id))).Find()

	return assistantKeys, err
}

func (s *Service) DeleteKey(ctx context.Context, assistantKey *entity.AssistantKey) error {
	_, err := s.dao.WithContext(ctx).AssistantKey.Delete(assistantKey)

	return err
}

func (s *Service) DeleteAllKey(ctx context.Context, assistant *entity.Assistant) error {
	_, err := s.dao.WithContext(ctx).AssistantKey.Where(s.dao.AssistantKey.AssistantId.Eq(uint(assistant.Id))).Delete()

	return err
}

//func (s *Service) UpdateKey(ctx context.Context, assistant *entity.Assistant) error {
//	_, err := s.x.Context(ctx).ID(assistant.Id).AllCols().Update(assistant)
//	return err
//}
