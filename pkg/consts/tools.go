package consts

import "errors"

var (
	ErrToolAlreadyExists                       = errors.New("该工具已存在")
	ErrToolNotFound                            = errors.New("未找到该工具")
	ErrToolNotYours                            = errors.New("你不能访问这个工具，因为他不是你创建的。")
	ErrToolFailedDeleteBecauseHasBindAssistant = errors.New("该工具已经绑定过助理了，不能删除。请先解除绑定后再尝试删除。")
)
