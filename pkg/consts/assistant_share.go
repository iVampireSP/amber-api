package consts

import "errors"

var (
	ErrShareNotFound = errors.New("未找到该共享内容")
	ErrShareInvalid  = errors.New("无效的共享内容")
)
