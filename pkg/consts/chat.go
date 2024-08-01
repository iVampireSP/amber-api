package consts

import (
	"errors"
	"time"
)

const ChatStreamExpire = 10 * time.Minute

var (
	ErrChatNotFound                        = errors.New("未找到该聊天")
	ErrChatLastMessageIsHuman              = errors.New("最后一条消息是用户发的消息，AI 还没有回复，请等待一段时间")
	ErrChatNotYours                        = errors.New("你不能访问这个聊天，因为它不是你创建的")
	ErrChatNameTooLong                     = errors.New("名称太长了")
	ErrChatNameTooShort                    = errors.New("名称太短了")
	ErrChatNameMustBeNotEmpty              = errors.New("名称不能为空")
	ErrChatAlreadyExists                   = errors.New("该聊天已存在")
	ErrChatStreamNotOpen                   = errors.New("聊天流未打开，你无法添加消息，请先获取之前的聊天流")
	ErrChatStreamNotOpenAndOverrideMessage = errors.New(ErrChatStreamNotOpen.Error() + "，我们已经更正上一条消息为本消息的内容")
)
