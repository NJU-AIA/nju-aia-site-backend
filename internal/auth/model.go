package auth

const (
	RoleAdmin = "admin"
)

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
