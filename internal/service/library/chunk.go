package library

import (
	"context"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"rag-new/internal/dao"
	"rag-new/internal/entity"
	"rag-new/pkg/consts"
	"strings"
)

const FileSize = 50 * 1024 * 1024
const ChunkSize = 1024
const ChunkOverlap = 128

// const mimeTypeMSWord = "application/msword"
const mimeTypePDF = "application/pdf"
const mimeTypeCSV = "text/csv"

// const mimeTypeOffice = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
const mimeTypeTXT = "text/plain"
const mimeTypeHTML = "text/html"

var allowedChunkMimeTypes = map[string]bool{
	mimeTypeCSV: true,
	//mimeTypeOffice: true,
	mimeTypePDF:  true,
	mimeTypeTXT:  true,
	mimeTypeHTML: true,
}

func (s *Service) CanChunk(file *entity.File) bool {
	if file.Size > FileSize {
		return false
	}

	fileMimeTypeString2 := strings.Split(file.MimeType, ";")[0]

	exists, ok := allowedChunkMimeTypes[fileMimeTypeString2]
	return exists && ok
}

func (s *Service) ChunkFileToDocument(ctx context.Context, file *entity.File, document *entity.Document) error {
	var documentLoader documentloaders.Loader
	size, io, err := s.fileService.GetBucketFile(ctx, file)
	if err != nil {
		return err
	}

	//if !s.canChunk(size, file.MimeType) {
	//	return nil, consts.ErrFileNotSupportChunk
	//}

	// 不过要做 Web 端 chunk，一次性上传那么多也是不现实的
	// --其实不用想这么多，后端只需要 chunk 纯文本即可，文档解析这类应该交给客户端处理--

	fileMimeTypeString2 := strings.Split(file.MimeType, ";")[0]

	switch fileMimeTypeString2 {
	case mimeTypePDF:
		documentLoader = documentloaders.NewPDF(io, size)
	//case mimeTypeMSWord:
	//	documentLoader = NewDocxLoader(io, size)
	//case mimeTypeOffice:
	//	documentLoader = NewDocxLoader(io, size)
	case mimeTypeTXT:
		documentLoader = documentloaders.NewText(io)
	case mimeTypeHTML:
		documentLoader = documentloaders.NewText(io)
	default:
		return consts.ErrFileNotSupportChunk
	}

	recursiveCharacter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(ChunkSize),
		textsplitter.WithChunkOverlap(ChunkOverlap),
	)

	var chunks []schema.Document
	chunks, err = documentLoader.LoadAndSplit(ctx, recursiveCharacter)
	if err != nil {
		return err
	}

	var documentChunks = make([]*entity.DocumentChunk, 0)

	for _, chunk := range chunks {
		documentChunks = append(documentChunks, &entity.DocumentChunk{
			Content:    chunk.PageContent,
			DocumentId: document.Id,
			LibraryId:  document.LibraryId,
		})
	}

	err = s.dao.Transaction(func(tx *dao.Query) error {
		return tx.DocumentChunk.WithContext(ctx).Create(documentChunks...)
	})

	if err != nil {
		return err
	}

	// 将  chunked 标记为 true
	_, err = s.dao.Document.WithContext(ctx).
		Where(s.dao.Document.Id.Eq(uint(document.Id))).
		UpdateSimple(s.dao.Document.Chunked.Value(true))

	return err

}

func (s *Service) ChunkTextToDocument(ctx context.Context, content string, document *entity.Document) error {
	var io = strings.NewReader(content)

	documentLoader := documentloaders.NewText(io)

	// 这种就很容易,后端只需要 chunk 纯文本即可，文档解析这类应该交给客户端处理

	recursiveCharacter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(ChunkSize),
		textsplitter.WithChunkOverlap(ChunkOverlap),
	)

	chunks, err := documentLoader.LoadAndSplit(ctx, recursiveCharacter)
	if err != nil {
		return err
	}

	var documentChunks = make([]*entity.DocumentChunk, 0)

	for _, chunk := range chunks {
		documentChunks = append(documentChunks, &entity.DocumentChunk{
			Content:    chunk.PageContent,
			DocumentId: document.Id,
			LibraryId:  document.LibraryId,
		})
	}

	err = s.dao.Transaction(func(tx *dao.Query) error {
		return tx.DocumentChunk.WithContext(ctx).Create(documentChunks...)
	})
	if err != nil {
		return err
	}

	// 将  chunked 标记为 true
	_, err = s.dao.Document.WithContext(ctx).
		Where(s.dao.Document.Id.Eq(uint(document.Id))).
		UpdateSimple(s.dao.Document.Chunked.Value(true))

	return err

}
