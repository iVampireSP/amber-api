package auth

import (
	"context"
	"rag-new/internal/base/conf"
	"rag-new/internal/base/logger"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/internal/service/user"
	"rag-new/pkg/consts"
	"strings"

	"github.com/gin-gonic/gin"
)

type Service struct {
	config      *conf.Config
	logger      *logger.Logger
	userService *user.Service
}

func NewAuthService(config *conf.Config, logger *logger.Logger, userService *user.Service) *Service {
	return &Service{
		config:      config,
		logger:      logger,
		userService: userService,
	}
}

func (a *Service) GinMiddlewareAuth(c *gin.Context) (*schema.User, error) {
	if a.config.Debug.Enabled {
		// 调试模式下，返回测试用户
		return &schema.User{
			Token: schema.UserTokenInfo{
				Sub:   consts.AnonymousUser,
				Name:  "Debug User",
				Email: "debug@example.com",
			},
			Valid: true,
		}, nil
	}

	authorization := c.Request.Header.Get(consts.AuthHeader)

	if authorization == "" {
		return nil, consts.ErrJWTFormatError
	}

	authSplit := strings.Split(authorization, " ")
	if len(authSplit) != 2 {
		return nil, consts.ErrJWTFormatError
	}

	if authSplit[0] != consts.AuthPrefix {
		return nil, consts.ErrNotBearerType
	}

	// 验证令牌
	user, err := a.userService.VerifyToken(c, authSplit[1])
	if err != nil {
		return nil, err
	}

	// 转换为架构用户
	return a.userService.ConvertToSchemaUser(user), nil
}

func (a *Service) GinUser(c *gin.Context) *schema.User {
	user, _ := c.Get(consts.AuthMiddlewareKey)
	return user.(*schema.User)
}

func (a *Service) GetUserId(ctx context.Context) schema.UserId {
	user := a.GetUser(ctx)
	return user.Token.Sub
}

func (a *Service) GetUser(ctx context.Context) *schema.User {
	user := ctx.Value(consts.AuthMiddlewareKey)
	return user.(*schema.User)
}

func (a *Service) SetUser(ctx context.Context, user *schema.User) context.Context {
	return context.WithValue(ctx, consts.AuthMiddlewareKey, user)
}

// Register 注册新用户
func (a *Service) Register(ctx context.Context, req *schema.RegisterRequest) (*schema.TokenResponse, error) {
	// 调用用户服务注册
	user, err := a.userService.Register(ctx, req)
	if err != nil {
		return nil, err
	}

	// 生成令牌
	return a.userService.GenerateToken(user)
}

// Login 用户登录
func (a *Service) Login(ctx context.Context, req *schema.LoginRequest) (*schema.TokenResponse, error) {
	// 调用用户服务登录
	return a.userService.Login(ctx, req)
}

// GetCurrentUser 获取当前登录用户
func (a *Service) GetCurrentUser(ctx context.Context) (*entity.User, error) {
	userId := a.GetUserId(ctx)
	// 尝试将userId转换为字符串后再转为EntityId
	entityId, err := schema.EntityIdFromString(userId.String())
	if err != nil {
		return nil, err
	}
	return a.userService.GetUserById(ctx, entityId)
}
