package assistant

import (
	"context"
	"rag-new/internal/entity"
	"rag-new/pkg/consts"
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

	_, err := s.x.Context(ctx).Insert(assistantShare)
	return assistantShare, err
}

func (s *Service) GetShare(ctx context.Context, assistantShareId int64) (*entity.AssistantShareType, error) {
	var assistantShareType = entity.AssistantShareType{}
	b, err := s.x.Context(ctx).
		Where("assistant_shares.id = ?", assistantShareId).
		Join("INNER", "assistants", "assistants.id = assistant_shares.assistant_id").
		Get(&assistantShareType)

	if err != nil {
		return nil, err
	}

	if !b {
		return nil, consts.ErrShareNotFound
	}

	return &assistantShareType, err
}

func (s *Service) GetShareByToken(ctx context.Context, assistantShareToken string) (*entity.AssistantShareType, error) {
	var assistantShare = &entity.AssistantShareType{}
	b, err := s.x.Context(ctx).Where("token = ?", assistantShareToken).
		Join("INNER", "assistants", "assistants.id = assistant_shares.assistant_id").
		Get(assistantShare)

	if err != nil {
		return nil, err
	}

	if !b {
		return nil, consts.ErrShareNotFound
	}

	return assistantShare, err
}

// ListShare 获取当前助理的所有分享
func (s *Service) ListShare(ctx context.Context, assistant *entity.Assistant) ([]*entity.AssistantShare, error) {
	var assistantShares = make([]*entity.AssistantShare, 0)
	err := s.x.Context(ctx).Where("assistant_id = ?", assistant.Id).Find(&assistantShares)
	return assistantShares, err
}

func (s *Service) DeleteShare(ctx context.Context, assistantShare *entity.AssistantShare) error {
	_, err := s.x.Context(ctx).ID(assistantShare.Id).Delete(&entity.AssistantShare{})
	return err
}

//func (s *Service) UpdateShare(ctx context.Context, assistant *entity.Assistant) error {
//	_, err := s.x.Context(ctx).ID(assistant.Id).AllCols().Update(assistant)
//	return err
//}
