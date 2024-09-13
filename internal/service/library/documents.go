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

func (s *Service) CreateDocument(ctx context.Context, library *entity.Library, name string, content string) (*entity.Document, error) {
	// 如果有 file_id，则寻找是否存在相同的
	//if document.FileId != 0 {
	//	count, err := s.dao.WithContext(ctx).Document.Where(s.dao.Document.FileId.Eq(uint(document.FileId))).Count()
	//	if err != nil {
	//		return err
	//	}
	//}

	document := &entity.Document{
		LibraryId: library.Id,
		Name:      name,
		Chunked:   false,
	}

	// chunk 文档
	var documentDao = s.dao.Document.WithContext(ctx)
	err := documentDao.Create(document)
	if err != nil {
		return nil, err
	}

	if content != "" {
		go func() {
			err = s.ChunkTextToDocument(ctx, content, document)
			if err != nil {
				s.logger.Sugar.Error(err)

				// 删除文档
				_, err = documentDao.Delete(document)
				if err != nil {
					s.logger.Sugar.Error(err)
				}
			}
		}()
	}

	return document, nil
}

func (s *Service) UpdateDocument(ctx context.Context, document *entity.Document) error {
	return s.dao.Transaction(func(tx *dao.Query) error {
		_, err := s.dao.DocumentChunk.Where(s.dao.DocumentChunk.DocumentId.Eq(uint(document.Id))).Delete()
		if err != nil {
			return err
		}

		document.Chunked = false

		_, err = tx.Document.WithContext(ctx).Updates(document)
		if err != nil {
			return err
		}

		err = s.deleteMilvusChunk(ctx, document)
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
		if err != nil {
			return err
		}

		err = s.deleteMilvusChunk(ctx, document)

		return err
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

func (s *Service) GetDocumentFromLibrary(ctx context.Context, library *entity.Library, documentId schema.EntityId) (*entity.Document, error) {
	return s.dao.Document.WithContext(ctx).Where(s.dao.Document.Id.Eq(uint(documentId))).
		Where(s.dao.Document.LibraryId.Eq(uint(library.Id))).First()
}
