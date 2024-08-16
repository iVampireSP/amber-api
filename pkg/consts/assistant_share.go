package consts

import "errors"

var (
	ErrShareNotFound            = errors.New("未找到该共享内容")
	ErrShareInvalid             = errors.New("无效的共享内容")
	ErrAssistantTokenNotFound   = errors.New("未找到助理令牌")
	ErrAssistantTokenIsRequired = errors.New("助理令牌是必须的")
)
