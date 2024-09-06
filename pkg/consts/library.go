package consts

import "errors"

var (
	ErrLibraryNotFound         = errors.New("未找到指定的资料库")
	ErrLibraryHasDocuments     = errors.New("资料库内有文档，请先清空资料库内的文档后再尝试删除该资料库")
	ErrLibraryUsedByAssistants = errors.New("资料库被助理绑定了，请先解绑后再尝试删除该资料库")
)
