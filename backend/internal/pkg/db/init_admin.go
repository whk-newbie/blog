package db

import (
	"errors"

	"github.com/iambaby/blog/internal/models"
	"github.com/iambaby/blog/internal/pkg/logger"
	"github.com/iambaby/blog/internal/repository"
	"github.com/iambaby/blog/internal/service"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	DefaultAdminUsername = "admin"
	DefaultAdminPassword = "admin@123"
	DefaultAdminEmail    = "admin@example.com"
)

// InitDefaultAdmin 初始化默认管理员
// 如果管理员已存在，会检查并修正密码哈希，确保是通过 Go 的 bcrypt 正确生成的
func InitDefaultAdmin(db *gorm.DB) error {
	adminRepo := repository.NewAdminRepository(db)

	// 检查是否已存在管理员
	admin, err := adminRepo.FindByUsername(DefaultAdminUsername)
	if err != nil {
		if errors.Is(err, repository.ErrAdminNotFound) {
			// 管理员不存在，创建新管理员
			return createDefaultAdmin(adminRepo)
		}
		return err
	}

	// 管理员已存在，检查并修正密码哈希
	return verifyAndFixAdminPassword(adminRepo, admin)
}

// createDefaultAdmin 创建默认管理员
func createDefaultAdmin(adminRepo repository.AdminRepository) error {
	// 使用 Go 的 bcrypt 生成密码哈希
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

// verifyAndFixAdminPassword 验证并修正管理员密码哈希
func verifyAndFixAdminPassword(adminRepo repository.AdminRepository, admin *models.Admin) error {
	// 如果密码已被修改（不再是默认密码），不需要检查和更新
	if !admin.IsDefaultPassword {
		logger.Info("Default admin exists and password has been changed, skipping verification")
		return nil
	}

	// 尝试用默认密码验证现有的密码哈希
	err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(DefaultAdminPassword))
	if err == nil {
		// 密码验证成功，说明哈希是正确的
		logger.Info("Default admin already exists with correct password hash")
		// 确保邮箱正确
		if admin.Email == "" {
			admin.Email = DefaultAdminEmail
			if err := adminRepo.Update(admin); err != nil {
				return err
			}
		}
		return nil
	}

	// 密码验证失败，且 is_default_password 为 true
	// 说明密码哈希可能不正确（比如是从 SQL 迁移脚本创建的旧哈希）
	// 使用 Go 的 bcrypt 重新生成正确的密码哈希

	logger.Info("Default admin exists but password hash may be incorrect, updating...")

	// 使用 Go 的 bcrypt 重新生成密码哈希
	hashedPassword, err := service.HashPassword(DefaultAdminPassword)
	if err != nil {
		return err
	}

	// 更新密码哈希
	admin.Password = hashedPassword
	admin.IsDefaultPassword = true
	// 确保邮箱正确
	if admin.Email == "" {
		admin.Email = DefaultAdminEmail
	}

	if err := adminRepo.Update(admin); err != nil {
		return err
	}

	logger.Info("Default admin password hash has been corrected")
	logger.Info("Username: %s", DefaultAdminUsername)
	logger.Info("Password: %s", DefaultAdminPassword)
	logger.Info("⚠️  Please change the default password after first login!")

	return nil
}
