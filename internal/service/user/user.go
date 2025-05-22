package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"rag-new/internal/base/conf"
	"rag-new/internal/base/logger"
	"rag-new/internal/dao"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Service 用户服务
type Service struct {
	dao    *dao.Query
	config *conf.Config
	logger *logger.Logger
	// 用于JWT签名的密钥
	secretKey []byte
}

// NewService 创建用户服务
func NewService(dao *dao.Query, config *conf.Config, logger *logger.Logger) *Service {
	// 生成随机密钥，或者从配置中读取
	secretKey := make([]byte, 32)
	_, err := rand.Read(secretKey)
	if err != nil {
		panic("无法生成随机密钥：" + err.Error())
	}

	return &Service{
		dao:       dao,
		config:    config,
		logger:    logger,
		secretKey: secretKey,
	}
}

// Register 注册新用户
func (s *Service) Register(ctx context.Context, req *schema.RegisterRequest) (*entity.User, error) {
	// 检查用户名是否已存在
	existingUser, err := s.dao.User.WithContext(ctx).Where(s.dao.User.Username.Eq(req.Username)).First()
	if err == nil && existingUser != nil {
		return nil, consts.ErrUserAlreadyExists
	}

	// 检查邮箱是否已存在
	existingEmail, err := s.dao.User.WithContext(ctx).Where(s.dao.User.Email.Eq(req.Email)).First()
	if err == nil && existingEmail != nil {
		return nil, consts.ErrEmailAlreadyExists
	}

	// 生成密码哈希
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &entity.User{
		Username:     req.Username,
		PasswordHash: string(passwordHash),
		Email:        req.Email,
		Name:         req.Name,
	}

	// 保存到数据库
	err = s.dao.User.WithContext(ctx).Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *Service) Login(ctx context.Context, req *schema.LoginRequest) (*schema.TokenResponse, error) {
	// 查找用户
	user, err := s.dao.User.WithContext(ctx).Where(s.dao.User.Username.Eq(req.Username)).First()
	if err != nil {
		return nil, consts.ErrInvalidCredentials
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, consts.ErrInvalidCredentials
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	_, err = s.dao.User.WithContext(ctx).Where(s.dao.User.Id.Eq(uint(user.Id))).Updates(user)
	if err != nil {
		s.logger.Sugar.Warnf("更新用户最后登录时间失败: %v", err)
	}

	// 生成JWT令牌
	token, expiresAt, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return &schema.TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt,
		User:        user.ToResponse(),
	}, nil
}

// VerifyToken 验证令牌并返回用户信息
func (s *Service) VerifyToken(ctx context.Context, tokenString string) (*entity.User, error) {
	// 解析令牌
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 确保签名方法是我们期望的
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, consts.ErrTokenInvalid
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, consts.ErrTokenInvalid
	}

	// 验证令牌有效性
	if !token.Valid {
		return nil, consts.ErrTokenInvalid
	}

	// 获取令牌中的claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, consts.ErrTokenInvalid
	}

	// 验证过期时间
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, consts.ErrTokenExpired
		}
	}

	// 获取用户ID
	userIdStr, ok := claims["sub"].(string)
	if !ok {
		return nil, consts.ErrTokenInvalid
	}

	// 将用户ID转换为实体ID
	userId, err := schema.EntityIdFromString(userIdStr)
	if err != nil {
		return nil, consts.ErrTokenInvalid
	}

	// 查找用户
	user, err := s.dao.User.WithContext(ctx).Where(s.dao.User.Id.Eq(uint(userId))).First()
	if err != nil {
		return nil, consts.ErrUserNotFound
	}

	return user, nil
}

// GetUserById 根据ID获取用户
func (s *Service) GetUserById(ctx context.Context, userId schema.EntityId) (*entity.User, error) {
	user, err := s.dao.User.WithContext(ctx).Where(s.dao.User.Id.Eq(uint(userId))).First()
	if err != nil {
		return nil, consts.ErrUserNotFound
	}
	return user, nil
}

// 生成JWT令牌
func (s *Service) generateToken(user *entity.User) (string, time.Time, error) {
	// 设置过期时间
	expiresAt := time.Now().Add(time.Duration(consts.AccessTokenExpiration) * time.Second)

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.Id.String(), // 用户ID
		"iat":   time.Now().Unix(),
		"exp":   expiresAt.Unix(),
		"name":  user.Name,
		"email": user.Email,
	})

	// 签名令牌
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// GenerateRandomKey 生成随机密钥
func GenerateRandomKey(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

// GenerateToken 为指定用户生成新令牌
func (s *Service) GenerateToken(user *entity.User) (*schema.TokenResponse, error) {
	// 生成JWT令牌
	token, expiresAt, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return &schema.TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt,
		User:        user.ToResponse(),
	}, nil
}

// ConvertToSchemaUser 将实体用户转换为架构用户
func (s *Service) ConvertToSchemaUser(user *entity.User) *schema.User {
	return &schema.User{
		Token: schema.UserTokenInfo{
			Sub:    schema.UserId(user.Id.String()),
			Name:   user.Name,
			Email:  user.Email,
			Avatar: user.Avatar,
		},
		Valid: true,
	}
}
