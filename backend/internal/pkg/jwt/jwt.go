package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenExpired     = errors.New("token has expired")
	ErrTokenNotValidYet = errors.New("token not active yet")
	ErrTokenMalformed   = errors.New("token is malformed")
	ErrTokenInvalid     = errors.New("token is invalid")
)

// Claims JWT声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Manager JWT管理器
type Manager struct {
	secret     []byte
	expireTime time.Duration
	issuer     string
}

// NewManager 创建JWT管理器
func NewManager(secret string, expireTime time.Duration, issuer string) *Manager {
	return &Manager{
		secret:     []byte(secret),
		expireTime: expireTime,
		issuer:     issuer,
	}
}

// GenerateToken 生成JWT Token
func (m *Manager) GenerateToken(userID uint, username string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(m.expireTime)

	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    m.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ParseToken 解析JWT Token
func (m *Manager) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return m.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotValidYet
		} else if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// RefreshToken 刷新Token
func (m *Manager) RefreshToken(tokenString string) (string, error) {
	claims, err := m.ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	// 生成新的token
	return m.GenerateToken(claims.UserID, claims.Username)
}

// VerifyToken 验证Token是否有效
func (m *Manager) VerifyToken(tokenString string) bool {
	_, err := m.ParseToken(tokenString)
	return err == nil
}

