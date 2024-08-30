package consts

import "errors"

var (
	ErrAssistantAlreadyBindTheTool    = errors.New("这个助理已经绑定过此工具了")
	ErrAssistantNotFound              = errors.New("未找到该助理")
	ErrAssistantHasBindToolCantDelete = errors.New("这个助理有绑定的工具，请先移除所有的工具，然后再尝试删除该助理")
	ErrToolNotBind                    = errors.New("该工具没有绑定该助理")

	MessageBatchDeleting = "正在批量删除与该助理关联的数据"
)
