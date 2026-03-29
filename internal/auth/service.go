package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	repo      *Repository
	jwtSecret []byte
	tokenTTL  time.Duration
	eccKey    []byte
}

func NewService(repo *Repository, jwtSecret string, tokenTTL time.Duration, eccKeyB64 string) *Service {
	secret := strings.TrimSpace(jwtSecret)
	if secret == "" {
		secret = "change-me-in-production"
	}

	var eccKey []byte
	if k := strings.TrimSpace(eccKeyB64); k != "" {
		var err error
		eccKey, err = ParseECCKey(k)
		if err != nil {
			panic("TOTP_ECC_KEY 解析失败: " + err.Error())
		}
	}

	return &Service{
		repo:      repo,
		jwtSecret: []byte(secret),
		tokenTTL:  tokenTTL,
		eccKey:    eccKey,
	}
}

func (s *Service) EnsureDefaultAdmin() error {
	exists, err := s.repo.ExistsAdmin()
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	admin := &User{
		Username: "admin",
		Role:     RoleAdmin,
		Status:   "active",
	}
	return s.repo.Create(admin)
}

// Login 使用 ECC-TOTP 动态码验证身份。
func (s *Service) Login(username, code string) (string, error) {
	if len(s.eccKey) == 0 {
		return "", errors.New("TOTP 未配置，请设置 TOTP_ECC_KEY 环境变量")
	}

	if !VerifyTOTP(s.eccKey, code, time.Now()) {
		return "", errors.New("动态码无效或已过期")
	}

	var (
		user *User
		err  error
	)

	uname := strings.TrimSpace(username)
	if uname == "" {
		user, err = s.repo.FindFirstAdmin()
	} else {
		user, err = s.repo.FindByUsername(uname)
	}
	if err != nil {
		return "", err
	}
	if user.Status != "active" {
		return "", errors.New("account disabled")
	}

	return s.generateToken(user.ID, user.Username, user.Role)
}

func (s *Service) generateToken(userID uint, username, role string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"uid":      userID,
		"username": username,
		"role":     role,
		"iat":      now.Unix(),
		"exp":      now.Add(s.tokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *Service) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	uidFloat, ok := mapClaims["uid"].(float64)
	if !ok {
		return nil, errors.New("uid missing")
	}
	username, ok := mapClaims["username"].(string)
	if !ok {
		return nil, errors.New("username missing")
	}
	role, ok := mapClaims["role"].(string)
	if !ok || strings.TrimSpace(role) == "" {
		return nil, errors.New("role missing")
	}

	return &Claims{UserID: uint(uidFloat), Username: username, Role: role}, nil
}
