package consts

import "errors"

var (
	ErrUnableOpenFile             = errors.New("无法打开文件,你可能没有上传")
	ErrMimeTypeNotFound           = errors.New("未找到对应的文件类型")
	ErrFileSizeTooLarge           = errors.New("文件大小超过限制")
	ErrContentLengthHeaderMissing = errors.New("Content-Length 头信息缺失")
	ErrMimeTypeNotAllowed         = errors.New("不允许上传的文件类型")
	ErrFileNotExists              = errors.New("文件不存在")
	ErrFileRequired               = errors.New("如果你没有传入 URL，那么文件是必须的")
	ErrFileUrlRequired            = errors.New("如果你没有传入文件，那么文件 URL 是必须的")
	ErrFileNotSupportChunk        = errors.New("文件无法被分块")
)
