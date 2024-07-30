package consts

import "errors"

var (
	ErrAssistantAlreadyBindTheTool = errors.New("这个助理已经绑定过此工具了")
	ErrAssistantNotFound           = errors.New("未找到该助理")
)
