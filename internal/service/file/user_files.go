package file

import (
	"context"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
)

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

func (s *Service) GetUserFiles(ctx context.Context, user schema.UserId) ([]*entity.UserFile, error) {
	userFiles, err := s.dao.UserFile.WithContext(ctx).Where(s.dao.UserFile.UserId.Eq(int64(user))).Find()

	return userFiles, err
}

func (s *Service) GetUserFile(ctx context.Context, fileId schema.EntityId) (*entity.UserFile, error) {
	userFile, err := s.dao.UserFile.WithContext(ctx).
		Preload(s.dao.UserFile.File).
		Where(s.dao.UserFile.Id.Eq(uint(fileId))).First()

	return userFile, err
}

func (s *Service) DeleteUserFile(ctx context.Context, userFile *entity.UserFile) error {
	_, err := s.dao.UserFile.WithContext(ctx).Delete(userFile)

	return err
}

func (s *Service) ExistsUserFileById(ctx context.Context, userFileId schema.EntityId) (bool, error) {
	count, err := s.dao.UserFile.WithContext(ctx).Where(s.dao.UserFile.Id.Eq(uint(userFileId))).Count()

	return count > 0, err
}
