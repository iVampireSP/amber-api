package assistant

import (
	"context"
	"errors"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"

	page2 "github.com/iVampireSP/pkg/page"
	"gorm.io/gorm"
)

// ListPublicAssistant 获取公开的助理(paginate)
func (s *Service) ListPublicAssistant(ctx context.Context, page int) (*page2.PagedResult[*schema.AssistantPublic], error) {
	var pagedResult = page2.NewPagedResult[*schema.AssistantPublic]()

	pagedResult.Page = page

	var err error

	var assistants []*entity.Assistant

	pagedResult.TotalCount, err = s.dao.WithContext(ctx).Assistant.
		Where(s.dao.Assistant.Public.Is(true)).
		ScanByPage(&assistants, pagedResult.Offset(), pagedResult.PageSize)
	if err != nil {
		return nil, err
	}

	for _, v := range assistants {
		pagedResult.Data = append(pagedResult.Data, v.ToPublic())
	}

	return pagedResult.Output(), err
}

// ListUserFavoriteAssistants 获取用户收藏的助理
func (s *Service) ListUserFavoriteAssistants(ctx context.Context, userId schema.UserId, page int) (*page2.PagedResult[*schema.AssistantPublic], error) {
	var pagedResult = page2.NewPagedResult[*schema.AssistantPublic]()

	pagedResult.Page = page

	var err error

	var favoriteAssistants []*entity.FavoriteAssistants

	favoriteAssistants, pagedResult.TotalCount, err = s.dao.WithContext(ctx).FavoriteAssistants.
		Preload(s.dao.FavoriteAssistants.Assistant).
		Where(s.dao.FavoriteAssistants.UserId.Eq(userId.String())).
		FindByPage(pagedResult.Offset(), pagedResult.PageSize)
	if err != nil {
		return nil, err
	}

	for _, v := range favoriteAssistants {
		if v.Assistant == nil {
			continue
		}

		pagedResult.Data = append(pagedResult.Data, v.Assistant.ToPublic())
	}

	return pagedResult.Output(), err
}

func (s *Service) FavoriteAssistant(ctx context.Context, userId schema.UserId, assistant *entity.Assistant) error {
	favorite, err := s.HasFavoriteAssistant(ctx, userId, assistant)
	if err != nil {
		return err
	}

	if favorite {
		return consts.ErrAlreadyFavorite
	}

	err = s.dao.WithContext(ctx).FavoriteAssistants.Create(&entity.FavoriteAssistants{
		AssistantId: assistant.Id,
		UserId:      userId,
	})

	return err
}

func (s *Service) UnFavoriteAssistant(ctx context.Context, userId schema.UserId, assistant *entity.Assistant) error {
	// 检测是否 favorite
	_, err := s.dao.WithContext(ctx).FavoriteAssistants.Where(
		s.dao.FavoriteAssistants.AssistantId.Eq(assistant.Id.Uint()),
		s.dao.FavoriteAssistants.UserId.Eq(userId.String()),
	).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return consts.ErrNotFavorite
		}

		return err
	}

	_, err = s.dao.WithContext(ctx).FavoriteAssistants.Where(
		s.dao.FavoriteAssistants.AssistantId.Eq(assistant.Id.Uint()),
		s.dao.FavoriteAssistants.UserId.Eq(userId.String()),
	).Delete()

	return err
}

func (s *Service) HasFavoriteAssistant(ctx context.Context, userId schema.UserId, assistant *entity.Assistant) (bool, error) {
	count, err := s.dao.WithContext(ctx).FavoriteAssistants.Where(
		s.dao.FavoriteAssistants.Deleted.Is(false),
		s.dao.FavoriteAssistants.AssistantId.Eq(assistant.Id.Uint()),
		s.dao.FavoriteAssistants.UserId.Eq(userId.String()),
	).Count()

	return count > 0, err
}

func (s *Service) ClearAssistantFavorite(ctx context.Context, assistantId schema.EntityId) error {
	_, err := s.dao.WithContext(ctx).FavoriteAssistants.
		Where(s.dao.FavoriteAssistants.AssistantId.Eq(assistantId.Uint())).
		UpdateSimple(s.dao.FavoriteAssistants.Deleted.Value(true))

	return err
}

func (s *Service) CanUse(ctx context.Context, userId schema.UserId, assistantId schema.EntityId) (bool, error) {
	assistantEntity, err := s.GetAssistant(ctx, assistantId)
	if err != nil {
		return false, err
	}

	// 检测是否公开
	if !assistantEntity.Public && assistantEntity.UserId != userId {
		return false, err
	}

	// 检测是不是收藏的
	hasFavorite, err := s.HasFavoriteAssistant(ctx, userId, assistantEntity)
	if err != nil {
		return false, err
	}

	if !hasFavorite {
		if assistantEntity.UserId != userId {
			return false, nil
		}
	}

	return true, nil
}
