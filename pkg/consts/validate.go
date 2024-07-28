package consts

import "errors"

var (
	ErrValidateNameMustBeNotEmpty = errors.New("名称不能为空")
	ErrValidateNameTooLong        = errors.New("名称太长了")
	ErrValidateNameTooShort       = errors.New("名称太短了")
)
