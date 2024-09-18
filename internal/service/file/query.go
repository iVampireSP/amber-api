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
	count, err := s.dao.WithContext(ctx).File.Where(s.dao.File.UrlHash.Eq(urlHash)).Count()

	return count > 0, err
}

func (s *Service) GetFileByUrlHash(ctx context.Context, urlHash string) (*entity.File, error) {
	file, err := s.dao.File.WithContext(ctx).Where(s.dao.File.UrlHash.Eq(urlHash)).First()

	return file, err
}

func (s *Service) FileHashExists(ctx context.Context, fileHash string) (bool, error) {
	count, err := s.dao.WithContext(ctx).File.Where(s.dao.File.FileHash.Eq(fileHash)).Count()

	return count > 0, err
}

func (s *Service) GetFileByFileHash(ctx context.Context, fileHash string) (*entity.File, error) {
	file, err := s.dao.File.WithContext(ctx).Where(s.dao.File.FileHash.Eq(fileHash)).First()

	return file, err
}

func (s *Service) GetFileById(ctx context.Context, fileId schema.EntityId) (*entity.File, error) {
	file, err := s.dao.File.WithContext(ctx).Where(s.dao.File.Id.Eq(uint(fileId))).First()
	//if err != nil {
	//	return nil, err
	//}
	//
	//// 如果 expired_at 少于 RenewBeforeDAY 天，则延长至当前天 + ExpiredDAY
	//if file.ExpiredAt != nil && file.ExpiredAt.Before(time.Now().AddDate(0, 0, RenewBeforeDAY)) {
	//	var expired = time.Now().AddDate(0, 0, ExpiredDAY)
	//	file.ExpiredAt = &expired
	//
	//	_, err = s.dao.File.WithContext(ctx).Where(s.dao.File.Id.Eq(uint(fileId))).Update(s.dao.File.ExpiredAt, expired)
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	return file, err
}

func (s *Service) ExistsFileById(ctx context.Context, fileId schema.EntityId) (bool, error) {
	count, err := s.dao.WithContext(ctx).File.Where(s.dao.File.Id.Eq(uint(fileId))).Count()

	return count > 0, err
}

func (s *Service) GetImageUrl(file *entity.File) (string, error) {
	if file == nil {
		return "", consts.ErrFileNotExists
	}

	if s.config.Http.Url == "" {
		return "", errors.New("http url is empty")
	}

	var url = s.config.Http.Url + "/api/v1/files/download/" + file.FileHash

	return url, nil
}

// GetFilesByIds get file by ids
func (s *Service) GetFilesByIds(ctx context.Context, ids []schema.EntityId) ([]*entity.File, error) {
	// ids to uint(ids)
	ids2 := make([]uint, 0)
	for _, v := range ids {
		ids2 = append(ids2, uint(v))
	}

	files, err := s.dao.File.WithContext(ctx).Where(s.dao.File.Id.In(ids2...)).Find()

	return files, err
}

func (s *Service) Renew(ctx context.Context, files ...*entity.File) error {
	// 如果 expired_at 少于 RenewBeforeDAY 天，则延长至当前天 + ExpiredDAY
	for _, v := range files {
		if v.ExpiredAt != nil && v.ExpiredAt.Before(time.Now().AddDate(0, 0, RenewBeforeDAY)) {
			var expired = time.Now().AddDate(0, 0, ExpiredDAY)
			v.ExpiredAt = &expired
			_, err := s.dao.File.WithContext(ctx).Where(s.dao.File.Id.Eq(uint(v.Id))).Update(s.dao.File.ExpiredAt, expired)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

//func (s *Service) getCacheKey(key string) string {
//	return fmt.Sprintf("file:%s", key)
//}
