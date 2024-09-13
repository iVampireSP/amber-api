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
	// ErrChatCanNotDeleteBecauseNotCleared ErrChatIdNotProvided                   = errors.New("未提供聊天 ID, 这样无法确保聊天记录存在指定的聊天中")
	ErrChatCanNotDeleteBecauseNotCleared = errors.New("该聊天还未被清理，无法删除")
	ErrNoHumanMessage                    = errors.New("没有 Human 角色的消息")
	ErrRoleCanNotBeFile                  = errors.New("角色不能是文件，请使用其他端点上传文件")
	ErrCreateChatMessage                 = errors.New("无法创建聊天消息")
	ErrExpiredTimeCanNotBeforeNow        = errors.New("过期时间不能早于当前时间")
	ErrExpiredTimeCanNotAfter2038        = errors.New("过期时间不能晚于或等于 2038 年")
	ErrProvideSameImage                  = errors.New("上一条消息的文件和本条消息相同，已忽略")
	ErrImageIsRequired                   = errors.New("当 type 是 image_url 时，图片是必须的")
	ErrImageUrlCannotBeEmpty             = errors.New("当 type 是 image_url 时，图片 URL 是必须的")
	ErrTextCannotBeEmpty                 = errors.New("当 type 是 text 时，文本是必须的")
	// ErrTypeRequired ErrTextIsTooLong                       = errors.New("当 type 是 text 时，文本太长了")
	//ErrTextRequired                        = errors.New("当 type 是 text 时，文本是必须的")
	ErrTypeRequired          = errors.New("type 是必须的")
	ErrFileUrlNotURL         = errors.New("当 type 是 file_url 时，文件 URL 是必须的")
	ErrFileUrlNotValidBase64 = errors.New("文件 URL 不是有效的 base64")
)
