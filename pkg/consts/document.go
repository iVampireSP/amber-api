package consts

import "errors"

var (
	ErrDocumentNotFound = errors.New("找不到该文档。")
	ErrDocumentEmpty    = errors.New("文档内容为空。")
	ErrDocumentInvalid  = errors.New("无效的文档。")
)
