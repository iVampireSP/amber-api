package file

import (
	"context"
	"errors"
	"fmt"
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
	if err != nil {
		return nil, err
	}

	// 如果 expired_at 少于 RenewBeforeDAY 天，则延长至当前天 + ExpiredDAY
	if file.ExpiredAt != nil && file.ExpiredAt.Before(time.Now().AddDate(0, 0, RenewBeforeDAY)) {
		var expired = time.Now().AddDate(0, 0, ExpiredDAY)
		file.ExpiredAt = &expired

		_, err = s.dao.File.WithContext(ctx).Where(s.dao.File.Id.Eq(uint(fileId))).Update(s.dao.File.ExpiredAt, expired)
		if err != nil {
			return nil, err
		}
	}

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

	var url = s.config.Http.Url + "/api/v1/files/" + file.Id.String() + "/download"

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

	if err != nil {
		return nil, err
	}

	// 如果 expired_at 少于 RenewBeforeDAY 天，则延长至当前天 + ExpiredDAY
	for _, v := range files {
		if v.ExpiredAt != nil && v.ExpiredAt.Before(time.Now().AddDate(0, 0, RenewBeforeDAY)) {
			var expired = time.Now().AddDate(0, 0, ExpiredDAY)
			v.ExpiredAt = &expired
			_, err = s.dao.File.WithContext(ctx).Where(s.dao.File.Id.Eq(uint(v.Id))).Update(s.dao.File.ExpiredAt, expired)
			if err != nil {
				return nil, err
			}
		}
	}

	return files, err
}

func (s *Service) getURL(ctx context.Context, file *entity.File) (string, error) {
	// TODO 生成 1 分钟的 key。不用管 key 是否存在或已生成，可直接生成新的 key
	//cmd := s.redis.Get(ctx, s.getCacheKey("temp_key_file:"+file.Id.String()))
	//result, err := cmd.Result()
	//if err != nil {
	//	if !errors.Is(err, redis.Nil) {
	//		return "", err
	//	}
	//} else {
	//
	//}

	return "", nil
}

func (s *Service) getCacheKey(key string) string {
	return fmt.Sprintf("file:%s", key)
}

func (s *Service) BindFileToUser(ctx context.Context, file *entity.File, user schema.UserId) (*entity.UserFile, error) {
	var userFile = &entity.UserFile{
		FileId: file.Id,
		UserId: user,
	}

	// 检测是否绑定过
	// count
	count, err := s.dao.UserFile.WithContext(ctx).
		Where(s.dao.UserFile.FileId.Eq(uint(file.Id)), s.dao.UserFile.UserId.Eq(int64(user))).
		Count()

	if count > 0 {
		// 获取并返回
		return s.dao.UserFile.WithContext(ctx).
			Where(s.dao.UserFile.FileId.Eq(uint(file.Id))).First()
	}

	err = s.dao.UserFile.WithContext(ctx).Create(userFile)

	return userFile, err
}

func (s *Service) UnbindFileFromUser(ctx context.Context, fileId schema.EntityId, user schema.UserId) error {
	_, err := s.dao.UserFile.
		WithContext(ctx).
		Where(s.dao.UserFile.FileId.Eq(uint(fileId)), s.dao.UserFile.UserId.Eq(int64(user))).
		Delete()

	return err
}
