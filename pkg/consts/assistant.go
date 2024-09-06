package consts

import "errors"

var (
	ErrAssistantAlreadyBindTheTool    = errors.New("这个助理已经绑定过此工具了")
	ErrAssistantNotFound              = errors.New("未找到该助理")
	ErrAssistantHasBindToolCantDelete = errors.New("这个助理有绑定的工具，请先移除所有的工具，然后再尝试删除该助理")
	//ErrAssistantHasBindLibraryCantDelete = errors.New("这个助理有绑定的资料库，请先移除助理绑定的资料库，然后再尝试删除该助理")
	ErrToolNotBind = errors.New("该工具没有绑定该助理")
)
