package schema

import "time"

// UserResponse 用户基本信息响应
type UserResponse struct {
	Id       EntityId `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Avatar   string   `json:"avatar"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=30"`
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required,min=2,max=30"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse 令牌响应
type TokenResponse struct {
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
	ExpiresAt   time.Time     `json:"expires_at"`
	User        *UserResponse `json:"user"`
}
