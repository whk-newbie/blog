package service

import (
	"errors"
	"time"

	"github.com/whk-newbie/blog/internal/pkg/jwt"
	"github.com/whk-newbie/blog/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrPasswordTooShort   = errors.New("password must be at least 6 characters")
	ErrSamePassword       = errors.New("new password cannot be the same as old password")
)

// AuthService 认证服务接口
type AuthService interface {
	// 登录
	Login(username, password string) (*LoginResponse, error)
	// 修改密码
	ChangePassword(userID uint, oldPassword, newPassword string) error
	// 验证Token
	VerifyToken(token string) (*jwt.Claims, error)
	// 刷新Token
	RefreshToken(token string) (string, error)
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token             string        `json:"token"`
	ExpiresIn         int64         `json:"expires_in"` // 过期时间（秒）
	User              *AdminInfo    `json:"user"`
	IsDefaultPassword bool          `json:"is_default_password"` // 是否使用默认密码
}

// AdminInfo 管理员信息
type AdminInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// authService 认证服务实现
type authService struct {
	adminRepo  repository.AdminRepository
	jwtManager *jwt.Manager
	expireTime time.Duration
}

// NewAuthService 创建认证服务
func NewAuthService(adminRepo repository.AdminRepository, jwtManager *jwt.Manager, expireTime time.Duration) AuthService {
	return &authService{
		adminRepo:  adminRepo,
		jwtManager: jwtManager,
		expireTime: expireTime,
	}
}

// Login 登录
func (s *authService) Login(username, password string) (*LoginResponse, error) {
	// 查找管理员
	admin, err := s.adminRepo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, repository.ErrAdminNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// 生成JWT Token
	token, err := s.jwtManager.GenerateToken(admin.ID, admin.Username)
	if err != nil {
		return nil, err
	}

	// 更新最后登录时间
	if err := s.adminRepo.UpdateLastLogin(admin.ID); err != nil {
		// 记录错误但不影响登录
		// logger可以在这里记录
	}

	return &LoginResponse{
		Token:     token,
		ExpiresIn: int64(s.expireTime.Seconds()),
		User: &AdminInfo{
			ID:       admin.ID,
			Username: admin.Username,
			Email:    admin.Email,
		},
		IsDefaultPassword: admin.IsDefaultPassword,
	}, nil
}

// ChangePassword 修改密码
func (s *authService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	// 验证新密码长度
	if len(newPassword) < 6 {
		return ErrPasswordTooShort
	}

	// 查找管理员
	admin, err := s.adminRepo.FindByID(userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(oldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	// 检查新旧密码是否相同
	if oldPassword == newPassword {
		return ErrSamePassword
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	return s.adminRepo.UpdatePassword(userID, string(hashedPassword))
}

// VerifyToken 验证Token
func (s *authService) VerifyToken(token string) (*jwt.Claims, error) {
	return s.jwtManager.ParseToken(token)
}

// RefreshToken 刷新Token
func (s *authService) RefreshToken(token string) (string, error) {
	return s.jwtManager.RefreshToken(token)
}

// HashPassword 加密密码（工具方法）
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

