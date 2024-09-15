package consts

import "errors"

var (
	ErrApiKeyNotFound           = errors.New("api Key 不存在")
	ErrApiKeyInvalid            = errors.New("api key 无效")
	ErrAssistantTokenNotFound   = errors.New("未找到助理令牌")
	ErrAssistantTokenIsRequired = errors.New("助理令牌是必须的")
)
