package consts

import "errors"

var (
	ErrUnableOpenFile             = errors.New("无法打开文件,你可能没有上传")
	ErrMimeTypeNotFound           = errors.New("未找到对应的文件类型")
	ErrFileSizeTooLarge           = errors.New("文件大小超过限制")
	ErrContentLengthHeaderMissing = errors.New("Content-Length 头信息缺失")
	ErrMimeTypeNotAllowed         = errors.New("不允许上传的文件类型")
	ErrFileNotExists              = errors.New("文件不存在")
)
