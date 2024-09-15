package consts

import "errors"

var (
	ErrApiKeyNotFound         = errors.New("api Key 不存在")
	ErrApiKeyInvalid          = errors.New("api key 无效")
	ErrAssistantKeyNotFound   = errors.New("未找到助理令牌")
	ErrAssistantKeyIsRequired = errors.New("助理令牌是必须的")
)
