package repository

import (
	"errors"

	"github.com/whk-newbie/blog/internal/models"
	"gorm.io/gorm"
)

var (
	ErrAdminNotFound = errors.New("admin not found")
	ErrAdminExists   = errors.New("admin already exists")
)

// AdminRepository 管理员仓库接口
type AdminRepository interface {
	// 根据用户名查找管理员
	FindByUsername(username string) (*models.Admin, error)
	// 根据ID查找管理员
	FindByID(id uint) (*models.Admin, error)
	// 创建管理员
	Create(admin *models.Admin) error
	// 更新管理员
	Update(admin *models.Admin) error
	// 更新密码
	UpdatePassword(id uint, password string) error
	// 更新最后登录时间
	UpdateLastLogin(id uint) error
	// 标记密码已修改
	MarkPasswordChanged(id uint) error
	// 检查管理员是否存在
	Exists(username string) (bool, error)
}

// adminRepository 管理员仓库实现
type adminRepository struct {
	db *gorm.DB
}

// NewAdminRepository 创建管理员仓库
func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

// FindByUsername 根据用户名查找管理员
func (r *adminRepository) FindByUsername(username string) (*models.Admin, error) {
	var admin models.Admin
	err := r.db.Where("username = ?", username).First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAdminNotFound
		}
		return nil, err
	}
	return &admin, nil
}

// FindByID 根据ID查找管理员
func (r *adminRepository) FindByID(id uint) (*models.Admin, error) {
	var admin models.Admin
	err := r.db.First(&admin, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAdminNotFound
		}
		return nil, err
	}
	return &admin, nil
}

// Create 创建管理员
func (r *adminRepository) Create(admin *models.Admin) error {
	// 检查用户名是否已存在
	exists, err := r.Exists(admin.Username)
	if err != nil {
		return err
	}
	if exists {
		return ErrAdminExists
	}

	return r.db.Create(admin).Error
}

// Update 更新管理员
func (r *adminRepository) Update(admin *models.Admin) error {
	return r.db.Save(admin).Error
}

// UpdatePassword 更新密码
func (r *adminRepository) UpdatePassword(id uint, password string) error {
	return r.db.Model(&models.Admin{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"password":            password,
			"is_default_password": false,
		}).Error
}

// UpdateLastLogin 更新最后登录时间
func (r *adminRepository) UpdateLastLogin(id uint) error {
	return r.db.Model(&models.Admin{}).
		Where("id = ?", id).
		Update("last_login_at", gorm.Expr("NOW()")).Error
}

// MarkPasswordChanged 标记密码已修改
func (r *adminRepository) MarkPasswordChanged(id uint) error {
	return r.db.Model(&models.Admin{}).
		Where("id = ?", id).
		Update("is_default_password", false).Error
}

// Exists 检查管理员是否存在
func (r *adminRepository) Exists(username string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Admin{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

