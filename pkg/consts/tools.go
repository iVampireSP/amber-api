package consts

import "errors"

var (
	ErrToolNameMustBeNotEmpty = errors.New("名称不能为空")
	ErrToolNameTooLong        = errors.New("名称太长了")
	ErrToolNameTooShort       = errors.New("名称太短了")
	ErrToolAlreadyExists      = errors.New("该工具已存在")
	ErrToolNotFound           = errors.New("未找到该工具")
	ErrToolNotYours           = errors.New("你不能访问这个工具，因为他不是你创建的。")
)
