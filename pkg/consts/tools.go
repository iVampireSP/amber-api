package consts

import "errors"

var (
	ErrToolNameMustBeNotEmpty = errors.New("名称不能为空")
	ErrToolNameTooLong        = errors.New("名称太长了")
	ErrToolNameTooShort       = errors.New("名称太短了")
	ErrToolAlreadyExists      = errors.New("该工具已存在")
)
