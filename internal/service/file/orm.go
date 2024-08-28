package file

import (
	"context"
	"errors"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
	"time"
)

func (s *Service) URLExists(ctx context.Context, urlHash string) (bool, error) {
	i, err := s.x.Context(ctx).Where("url_hash = ?", urlHash).Count(&entity.File{})
	return i > 0, err
}

func (s *Service) GetFileByUrlHash(ctx context.Context, urlHash string) (*entity.File, error) {
	var file entity.File
	_, err := s.x.Context(ctx).Where("url_hash = ?", urlHash).Get(&file)
	return &file, err
}

func (s *Service) FileHashExists(ctx context.Context, fileHash string) (bool, error) {
	i, err := s.x.Context(ctx).Where("file_hash = ?", fileHash).Count(&entity.File{})
	return i > 0, err
}

func (s *Service) GetFileByFileHash(ctx context.Context, fileHash string) (*entity.File, error) {
	var file entity.File
	_, err := s.x.Context(ctx).Where("file_hash = ?", fileHash).Get(&file)
	return &file, err
}

func (s *Service) GetFileById(ctx context.Context, fileId schema.EntityId) (*entity.File, error) {
	var file = &entity.File{}
	_, err := s.x.Context(ctx).ID(fileId).Get(file)

	// 如果 expired_at 少于 RenewBeforeDAY 天，则延长至当前天 + ExpiredDAY
	if file.ExpiredAt != nil && file.ExpiredAt.Before(time.Now().AddDate(0, 0, RenewBeforeDAY)) {
		var expired = time.Now().AddDate(0, 0, ExpiredDAY)
		file.ExpiredAt = &expired
		_, _ = s.x.Context(ctx).ID(fileId).Cols("expired_at").Update(file)
	}

	return file, err
}

func (s *Service) ExistsFileById(ctx context.Context, fileId schema.EntityId) (bool, error) {
	i, err := s.x.Context(ctx).ID(fileId).Count(&entity.File{})
	return i > 0, err
}

func (s *Service) GetImageUrl(file *entity.File) (string, error) {
	if file == nil {
		return "", consts.ErrFileNotExists
	}

	if s.config.Http.Url == "" {
		return "", errors.New("http url is empty")
	}

	var url = s.config.Http.Url + "/api/v1/files/" + file.Id.String() + "/download"

	return url, nil
}

// GetFilesByIds get file by ids
func (s *Service) GetFilesByIds(ctx context.Context, ids []schema.EntityId) ([]*entity.File, error) {
	var files = make([]*entity.File, 0)
	err := s.x.Context(ctx).In("id", ids).Find(&files)
	if err != nil {
		return nil, err
	}

	// 如果 expired_at 少于 RenewBeforeDAY 天，则延长至当前天 + ExpiredDAY
	for _, v := range files {
		if v.ExpiredAt != nil && v.ExpiredAt.Before(time.Now().AddDate(0, 0, RenewBeforeDAY)) {
			var expired = time.Now().AddDate(0, 0, ExpiredDAY)
			v.ExpiredAt = &expired
			_, _ = s.x.Context(ctx).ID(v.Id).Cols("expired_at").Update(v)
		}
	}

	return files, err
}
