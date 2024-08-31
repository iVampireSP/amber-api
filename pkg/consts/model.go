package consts

import "errors"

const NoRecord = 0

// AutoModel 自动选择模型
const AutoModel = "auto"

// MaxTokenCount 最大的 token 数
const MaxTokenCount = 1000000

var (
	ErrPageNotFound    = errors.New("找不到请求的内容，请检查 URL 地址是否正确")
	ErrModelNotAllowed = errors.New("我们目前不支持这个模型，如果你不确定选什么，可以使用 auto")
	ErrTooManyTokens   = errors.New("token 太多了")
)
