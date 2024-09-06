package library

import (
	"context"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
)

func (s *Service) DefaultLibrary(ctx context.Context, userId schema.UserId) (*entity.Library, error) {
	var libraryDao = s.dao.WithContext(ctx).Library

	hasDefault, err := s.HasDefaultLibrary(ctx, userId)

	var library *entity.Library

	if !hasDefault {
		library = &entity.Library{
			Name:        "Default",
			Default:     true,
			Description: nil,
			UserId:      userId,
		}
		err = s.CreateLibrary(ctx, library)
		if err != nil {
			return nil, err
		}
	} else {
		library, err = libraryDao.
			Where(s.dao.Library.UserId.Eq(int64(userId))).
			Where(s.dao.Library.Default.Is(true)).First()
	}

	return library, err
}

func (s *Service) CreateLibrary(ctx context.Context, library *entity.Library) error {
	var libraryDao = s.dao.Library.WithContext(ctx)

	hasDefault, err := s.HasDefaultLibrary(ctx, library.UserId)

	if hasDefault {
		library.Default = false
	} else {
		library.Default = true
	}

	err = libraryDao.Create(library)
	return err
}

func (s *Service) UpdateLibrary(ctx context.Context, library *entity.Library) error {
	var libraryDao = s.dao.Library.WithContext(ctx)

	hasDefault, err := s.HasDefaultLibrary(ctx, library.UserId)
	if hasDefault && library.Default {
		library.Default = false
	} else {
		library.Default = true
	}

	_, err = libraryDao.Updates(library)
	return err
}

func (s *Service) GetLibrary(ctx context.Context, id schema.EntityId) (*entity.Library, error) {
	var libraryDao = s.dao.WithContext(ctx).Library

	library, err := libraryDao.Where(s.dao.Library.Id.Eq(uint(id))).First()

	return library, err
}

func (s *Service) GetLibraryByUserId(ctx context.Context, userId schema.UserId) (*entity.Library, error) {
	var libraryDao = s.dao.WithContext(ctx).Library
	library, err := libraryDao.Where(s.dao.Library.UserId.Eq(int64(userId))).First()
	return library, err
}

// GetLibraryDocuments returns a library with its documents.
//
// The returned library object is loaded with its documents.
// This is useful for scenarios where you need to access the documents of a library.
//
// The returned error is non-nil if the library doesn't exist.
//
// Example:
//
//	library, err := service.GetLibraryDocuments(ctx, library)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	for _, document := range library.Documents {
//	    // Do something with the document.
//	}
func (s *Service) GetLibraryDocuments(ctx context.Context, library *entity.Library) (*entity.Library, error) {
	var libraryDao = s.dao.WithContext(ctx).Library
	library, err := libraryDao.Where(s.dao.Library.Id.Eq(uint(library.Id))).
		Preload(s.dao.Library.Document).First()
	return library, err
}

func (s *Service) GetLibraryDocumentsById(ctx context.Context, libraryId schema.EntityId) (*entity.Library, error) {
	var libraryDao = s.dao.WithContext(ctx).Library
	library, err := libraryDao.Where(s.dao.Library.Id.Eq(uint(libraryId))).
		Preload(s.dao.Library.Document).First()
	return library, err
}

func (s *Service) DeleteLibrary(ctx context.Context, library *entity.Library) error {
	// 检测资料库是否有文档
	count, err := s.dao.Document.WithContext(ctx).Where(s.dao.Document.LibraryId.Eq(uint(library.Id))).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return consts.ErrLibraryHasDocuments
	}

	// 如果资料库内绑定了助理
	count, err = s.dao.WithContext(ctx).Assistant.Where(s.dao.Assistant.LibraryId.Eq(uint(library.Id))).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return consts.ErrLibraryUsedByAssistants
	}

	var libraryDao = s.dao.Library.WithContext(ctx)
	_, err = libraryDao.Delete(library)
	return err
}

func (s *Service) ListLibrary(ctx context.Context, userId schema.UserId) ([]*entity.Library, error) {
	var libraryDao = s.dao.WithContext(ctx).Library
	libraries, err := libraryDao.Where(s.dao.Library.UserId.Eq(int64(userId))).Find()
	return libraries, err
}

func (s *Service) ListLibraryByUserId(ctx context.Context, userId schema.UserId) ([]*entity.Library, error) {
	var libraryDao = s.dao.WithContext(ctx).Library
	libraries, err := libraryDao.Where(s.dao.Library.UserId.Eq(int64(userId))).Find()
	return libraries, err
}

func (s *Service) HasDefaultLibrary(ctx context.Context, userId schema.UserId) (bool, error) {
	var libraryDao = s.dao.WithContext(ctx).Library
	count, err := libraryDao.
		Where(s.dao.Library.UserId.Eq(int64(userId))).
		Where(s.dao.Library.Default.Is(true)).
		Count()

	return count > 0, err
}
