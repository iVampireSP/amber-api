package consts

import (
	"errors"
	"rag-new/internal/schema"
)

const (
	AuthHeader = "Authorization"
	AuthPrefix = "Bearer"

	//AnonymousUser schema.UserId = 1
	AnonymousUser schema.UserId = "anonymous"

	AuthMiddlewareKey               = "auth.user"
	AuthAssistantShareMiddlewareKey = "auth.assistant.share"

	// AccessTokenExpiration 访问令牌过期时间 (24小时)
	AccessTokenExpiration = 24 * 60 * 60
)

var (
	ErrNotValidToken  = errors.New("无效的 JWT 令牌")
	ErrJWTFormatError = errors.New("JWT 格式错误")
	ErrNotBearerType  = errors.New("不是 Bearer 类型")
	ErrEmptyResponse  = errors.New("我们的服务器返回了空请求，可能某些环节出了问题")
	ErrTokenError     = errors.New("token 类型错误")
	ErrBearerToken    = errors.New("无效的 Bearer 令牌")

	ErrNotYourResource  = errors.New("你不能修改这个资源，因为它不是你创建的。")
	ErrPermissionDenied = errors.New("没有权限访问此资源")

	// 用户认证相关错误
	ErrUserNotFound       = errors.New("用户不存在")
	ErrUserAlreadyExists  = errors.New("用户名已存在")
	ErrEmailAlreadyExists = errors.New("邮箱已存在")
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrPasswordTooShort   = errors.New("密码长度不能小于6位")
	ErrTokenExpired       = errors.New("令牌已过期")
	ErrTokenInvalid       = errors.New("无效的令牌")
)
