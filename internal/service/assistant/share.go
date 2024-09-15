package assistant

import (
	"context"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/random"
)

func (s *Service) CrateShare(ctx context.Context, assistant *entity.Assistant) (*entity.AssistantShare, error) {
	var assistantShare = &entity.AssistantShare{}

	assistantShare.AssistantId = assistant.Id
	// 生成 token
	assistantShare.Token = random.String(40)

	//// 检测 assistant 是否存在
	//assistant, err := s.GetAssistant(ctx, assistant.Id)
	//if err != nil {
	//	return nil, err
	//}
	//if assistant.Id == consts.NoRecord {
	//	return nil, consts.ErrAssistantNotFound
	//}

	err := s.dao.WithContext(ctx).AssistantShare.Create(assistantShare)

	return assistantShare, err
}

func (s *Service) GetShare(ctx context.Context, assistantShareId schema.EntityId) (*entity.AssistantShare, error) {
	assistantShare, err := s.dao.WithContext(ctx).AssistantShare.Where(s.dao.AssistantShare.Id.Eq(uint(assistantShareId))).Preload(s.dao.AssistantShare.Assistant).
		First()

	if err != nil {
		return nil, err
	}

	return assistantShare, err
}

func (s *Service) GetShareByToken(ctx context.Context, assistantShareToken string) (*entity.AssistantShare, error) {
	assistantShare, err := s.dao.WithContext(ctx).AssistantShare.Where(s.dao.AssistantShare.Token.Eq(assistantShareToken)).Preload(s.dao.AssistantShare.Assistant).
		First()

	if err != nil {
		return nil, err
	}

	return assistantShare, err
}

// ListShare 获取当前助理的所有分享
func (s *Service) ListShare(ctx context.Context, assistant *entity.Assistant) ([]*entity.AssistantShare, error) {
	assistantShares, err := s.dao.WithContext(ctx).AssistantShare.Where(s.dao.AssistantShare.AssistantId.Eq(uint(assistant.Id))).Find()

	return assistantShares, err
}

func (s *Service) DeleteShare(ctx context.Context, assistantShare *entity.AssistantShare) error {
	_, err := s.dao.WithContext(ctx).AssistantShare.Delete(assistantShare)

	return err
}

func (s *Service) DeleteAllShare(ctx context.Context, assistant *entity.Assistant) error {
	_, err := s.dao.WithContext(ctx).AssistantShare.Where(s.dao.AssistantShare.AssistantId.Eq(uint(assistant.Id))).Delete()

	return err
}

//func (s *Service) UpdateShare(ctx context.Context, assistant *entity.Assistant) error {
//	_, err := s.x.Context(ctx).ID(assistant.Id).AllCols().Update(assistant)
//	return err
//}
