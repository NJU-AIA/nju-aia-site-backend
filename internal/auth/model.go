package auth

import "time"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// User 使用通用用户模型，当前只开放管理员登录。
type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"size:64;uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"size:255" json:"-"`
	Role         string    `gorm:"size:16;index;not null" json:"role"`
	Status       string    `gorm:"size:16;default:'active';not null" json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// LoginRequest 为兼容现有契约保留 password 必填。
type LoginRequest struct {
	Password string `json:"password" binding:"required" example:"strong_pwd_123"`
	Username string `json:"username,omitempty" example:"admin"`
}

type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type ErrorResponse struct {
	Error string `json:"error" example:"密码错误"`
}

type Claims struct {
	UserID   uint   `json:"uid"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
