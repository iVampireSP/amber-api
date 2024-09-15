package assistant

import (
	"context"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/random"
)

func (s *Service) CrateToken(ctx context.Context, assistant *entity.Assistant) (*entity.AssistantToken, error) {
	var assistantToken = &entity.AssistantToken{}

	assistantToken.AssistantId = assistant.Id
	// 生成 token
	assistantToken.Token = random.String(40)

	//// 检测 assistant 是否存在
	//assistant, err := s.GetAssistant(ctx, assistant.Id)
	//if err != nil {
	//	return nil, err
	//}
	//if assistant.Id == consts.NoRecord {
	//	return nil, consts.ErrAssistantNotFound
	//}

	err := s.dao.WithContext(ctx).AssistantToken.Create(assistantToken)

	return assistantToken, err
}

func (s *Service) GetToken(ctx context.Context, assistantTokenId schema.EntityId) (*entity.AssistantToken, error) {
	assistantToken, err := s.dao.WithContext(ctx).AssistantToken.Where(s.dao.AssistantToken.Id.Eq(uint(assistantTokenId))).Preload(s.dao.AssistantToken.Assistant).
		First()

	if err != nil {
		return nil, err
	}

	return assistantToken, err
}

func (s *Service) GetTokenBySecret(ctx context.Context, assistantTokenToken string) (*entity.AssistantToken, error) {
	assistantToken, err := s.dao.WithContext(ctx).AssistantToken.Where(s.dao.AssistantToken.Token.Eq(assistantTokenToken)).Preload(s.dao.AssistantToken.Assistant).
		First()

	if err != nil {
		return nil, err
	}

	return assistantToken, err
}

// ListToken 获取当前助理的所有分享
func (s *Service) ListToken(ctx context.Context, assistant *entity.Assistant) ([]*entity.AssistantToken, error) {
	assistantTokens, err := s.dao.WithContext(ctx).AssistantToken.Where(s.dao.AssistantToken.AssistantId.Eq(uint(assistant.Id))).Find()

	return assistantTokens, err
}

func (s *Service) DeleteToken(ctx context.Context, assistantToken *entity.AssistantToken) error {
	_, err := s.dao.WithContext(ctx).AssistantToken.Delete(assistantToken)

	return err
}

func (s *Service) DeleteAllToken(ctx context.Context, assistant *entity.Assistant) error {
	_, err := s.dao.WithContext(ctx).AssistantToken.Where(s.dao.AssistantToken.AssistantId.Eq(uint(assistant.Id))).Delete()

	return err
}

//func (s *Service) UpdateToken(ctx context.Context, assistant *entity.Assistant) error {
//	_, err := s.x.Context(ctx).ID(assistant.Id).AllCols().Update(assistant)
//	return err
//}
