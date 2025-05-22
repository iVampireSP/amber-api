package entity

import (
	"rag-new/internal/schema"
	"time"
)

// User 用户实体
type User struct {
	Model
	Username     string     `json:"username"`      // 用户名
	PasswordHash string     `json:"-"`             // 密码哈希（不返回给客户端）
	Email        string     `json:"email"`         // 邮箱
	Name         string     `json:"name"`          // 昵称
	Avatar       string     `json:"avatar"`        // 头像
	LastLoginAt  *time.Time `json:"last_login_at"` // 最后登录时间
}

// TableName 表名
func (u *User) TableName() string {
	return "users"
}

// GetUserId 实现 HasUserInterface 接口
func (u *User) GetUserId() schema.UserId {
	return schema.UserId(u.Id.String())
}

// ToResponse 转换为响应对象
func (u *User) ToResponse() *schema.UserResponse {
	return &schema.UserResponse{
		Id:       u.Id,
		Username: u.Username,
		Email:    u.Email,
		Name:     u.Name,
		Avatar:   u.Avatar,
	}
}
