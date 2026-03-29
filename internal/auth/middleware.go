package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserIDKey   = "auth_user_id"
	ContextRoleKey     = "auth_role"
	ContextUsernameKey = "auth_username"
)

type Middleware struct {
	svc *Service
}

func NewMiddleware(svc *Service) *Middleware {
	return &Middleware{svc: svc}
}

func (m *Middleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c.GetHeader("Authorization"))
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Error: "未提供 Bearer Token"})
			return
		}

		claims, err := m.svc.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Error: "Token 无效或已过期"})
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUsernameKey, claims.Username)
		c.Set(ContextRoleKey, claims.Role)
		c.Next()
	}
}

func (m *Middleware) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get(ContextRoleKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Error: "未完成身份认证"})
			return
		}
		userRole, ok := value.(string)
		if !ok || userRole != role {
			c.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{Error: "权限不足"})
			return
		}
		c.Next()
	}
}

func (m *Middleware) RequireAdmin() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		m.RequireAuth()(c)
		if c.IsAborted() {
			return
		}
		m.RequireRole(RoleAdmin)(c)
	})
}

func extractBearerToken(header string) string {
	const prefix = "Bearer "
	if len(header) < len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		return ""
	}
	return strings.TrimSpace(header[len(prefix):])
}
