package consts

import "errors"

var (
	ErrWordRepeatedDetected = errors.New("检测到大量重复输出，会话已被终止")
)
