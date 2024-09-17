package consts

import "errors"

const (
	MaxToolFunctions = 100
)

var (
	ErrToolAlreadyExists                       = errors.New("该工具已存在")
	ErrToolNotFound                            = errors.New("未找到该工具")
	ErrToolNotYours                            = errors.New("你不能访问这个工具，因为他不是你创建的")
	ErrToolFailedDeleteBecauseHasBindAssistant = errors.New("该工具已经绑定过助理了，不能删除。请先解除绑定后再尝试删除")
	ErrToolSyntaxError                         = errors.New("工具的语法有错误")
	ErrToolDNSLookupError                      = errors.New("DNS 解析错误")
	ErrToolAddressIsInternal                   = errors.New("工具的 URL 或者回调地址是私有的，不能访问")
	StatusToolSyntaxMaybeOK                    = "工具的语法检查通过，可能可以正常运行。"
	ErrToolObjectParametersMissing             = errors.New("function 对象中，无论是否需要参数，都应该有 parameters 对象")
	ErrToolObjectRequiredMissing               = errors.New("parameters 对象中，required 数组是必须存在的，你可以不填入内容")
	ErrToolObjectPropertiesMissing             = errors.New("parameters 对象中，properties 对象是必须存在的")
	ErrToolObjectTypeMissing                   = errors.New("parameters 对象中，type 对象是必须存在的")
	ErrToolTooManyFunctions                    = errors.New("最多只能有 64 个函数")
	ErrToolFunctionTooMany                     = errors.New("助理所使用的工具中所有的函数加起来不能超过 100 个")
)
