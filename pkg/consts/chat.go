package consts

import (
	"errors"
	"time"
)

const ChatStreamExpire = 10 * time.Minute

var (
	ErrChatNotFound                        = errors.New("未找到该聊天")
	ErrChatStreaming                       = errors.New("聊天流正在进行中，无法添加消息")
	ErrChatStreamingPleaseWait             = errors.New("聊天流正在进行中，请稍后再试")
	ErrChatStreamNotOpen                   = errors.New("聊天流未打开，你无法添加消息，请先获取之前的聊天流")
	ErrChatStreamNotOpenAndOverrideMessage = errors.New(ErrChatStreamNotOpen.Error() + "，我们已经更正上一条消息为本消息的内容")
	ErrChatStreamNotFound                  = errors.New("未找到该聊天流，你无法获取消息。它可能已经超时了，重新生成一个试试")
	ErrChatIdNotProvided                   = errors.New("未提供聊天 ID, 这样无法确保聊天记录存在指定的聊天中")
	ErrChatCanNotDeleteBecauseNotCleared   = errors.New("该聊天还未被清理，无法删除")
	ErrNoHumanMessage                      = errors.New("没有 Human 角色的消息")
	ErrRoleCanNotBeImage                   = errors.New("角色不能是图片，请使用其他端点上传图片")
	ErrCreateChatMessage                   = errors.New("无法创建聊天消息")
)
