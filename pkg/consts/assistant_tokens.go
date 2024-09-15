package consts

import "errors"

var (
	ErrTokenNotFound            = errors.New("未找到 Token")
	ErrTokenInvalid             = errors.New("无效的 Token")
	ErrAssistantTokenNotFound   = errors.New("未找到助理令牌")
	ErrAssistantTokenIsRequired = errors.New("助理令牌是必须的")
)
