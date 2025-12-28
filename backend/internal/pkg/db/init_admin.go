package db

import (
	"github.com/iambaby/blog/internal/models"
	"github.com/iambaby/blog/internal/pkg/logger"
	"github.com/iambaby/blog/internal/repository"
	"github.com/iambaby/blog/internal/service"
	"gorm.io/gorm"
)

const (
	DefaultAdminUsername = "admin"
	DefaultAdminPassword = "admin123"
	DefaultAdminEmail    = "admin@example.com"
)

// InitDefaultAdmin 初始化默认管理员
func InitDefaultAdmin(db *gorm.DB) error {
	adminRepo := repository.NewAdminRepository(db)

	// 检查是否已存在管理员
	exists, err := adminRepo.Exists(DefaultAdminUsername)
	if err != nil {
		return err
	}

	if exists {
		logger.Info("Default admin already exists, skipping initialization")
		return nil
	}

	// 加密密码
	hashedPassword, err := service.HashPassword(DefaultAdminPassword)
	if err != nil {
		return err
	}

	// 创建默认管理员
	admin := &models.Admin{
		Username:          DefaultAdminUsername,
		Password:          hashedPassword,
		Email:             DefaultAdminEmail,
		IsDefaultPassword: true,
	}

	if err := adminRepo.Create(admin); err != nil {
		return err
	}

	logger.Info("Default admin created successfully")
	logger.Info("Username: %s", DefaultAdminUsername)
	logger.Info("Password: %s", DefaultAdminPassword)
	logger.Info("⚠️  Please change the default password after first login!")

	return nil
}

