package consts

import "errors"

const NoRecord = 0

var (
	ErrPageNotFound = errors.New("找不到请求的内容，请检查 URL 地址是否正确。")
)
