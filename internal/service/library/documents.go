package library

import (
	"context"
	"rag-new/internal/dao"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
)

func (s *Service) ListDocuments(ctx context.Context, library *entity.Library) ([]*entity.Document, error) {
	var documentDao = s.dao.Document.WithContext(ctx)

	documents, err := documentDao.Where(s.dao.Document.LibraryId.Eq(uint(library.Id))).Find()
	return documents, err
}

func (s *Service) GetDocument(ctx context.Context, id schema.EntityId) (*entity.Document, error) {
	var documentDao = s.dao.Document.WithContext(ctx)
	return documentDao.Where(s.dao.Document.Id.Eq(uint(id))).First()
}

func (s *Service) CreateDocument(ctx context.Context, document *entity.Document) error {
	// 如果有 file_id，则寻找是否存在相同的
	//if document.FileId != 0 {
	//	count, err := s.dao.WithContext(ctx).Document.Where(s.dao.Document.FileId.Eq(uint(document.FileId))).Count()
	//	if err != nil {
	//		return err
	//	}
	//}

	var documentDao = s.dao.Document.WithContext(ctx)
	return documentDao.Create(document)
}

func (s *Service) UpdateDocument(ctx context.Context, document *entity.Document) error {
	return s.dao.Transaction(func(tx *dao.Query) error {
		_, err := s.dao.DocumentChunk.Where(s.dao.DocumentChunk.DocumentId.Eq(uint(document.Id))).Delete()
		if err != nil {
			return err
		}

		document.Chunked = false

		_, err = tx.Document.WithContext(ctx).Updates(document)

		return err
	})

}

func (s *Service) DeleteDocument(ctx context.Context, document *entity.Document) error {
	err := s.dao.Transaction(func(tx *dao.Query) error {
		_, err := s.dao.DocumentChunk.Where(s.dao.DocumentChunk.DocumentId.Eq(uint(document.Id))).Delete()
		if err != nil {
			return err
		}

		_, err = tx.Document.WithContext(ctx).Delete(document)

		return nil
	})
	return err
}

func (s *Service) AddDocumentChunk(ctx context.Context, chunk ...*entity.DocumentChunk) error {
	return s.dao.DocumentChunk.WithContext(ctx).Create(chunk...)
}

func (s *Service) DeleteDocumentChunk(ctx context.Context, document *entity.Document) error {
	_, err := s.dao.DocumentChunk.WithContext(ctx).Where(s.dao.DocumentChunk.DocumentId.Eq(uint(document.Id))).Delete()
	return err
}

func (s *Service) GetDocumentByFileId(ctx context.Context, fileId schema.EntityId) (*entity.Document, error) {
	return s.dao.Document.WithContext(ctx).Where(s.dao.Document.FileId.Eq(uint(fileId))).First()
}

func (s *Service) GetDocumentByFileAndLibrary(ctx context.Context, file *entity.File, library *entity.Library) (*entity.Document, error) {
	return s.dao.Document.WithContext(ctx).Where(s.dao.Document.FileId.Eq(uint(file.Id)), s.dao.Document.LibraryId.Eq(uint(library.Id))).First()
}
