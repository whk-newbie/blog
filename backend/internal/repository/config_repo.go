package repository

import (
	"errors"

	"github.com/whk-newbie/blog/internal/models"
	"gorm.io/gorm"
)

var (
	ErrConfigNotFound = errors.New("config not found")
	ErrConfigExists   = errors.New("config already exists")
)

// ConfigRepository 配置仓库接口
type ConfigRepository interface {
	// 根据ID查找配置
	FindByID(id uint) (*models.SystemConfig, error)
	// 根据配置键查找配置
	FindByKey(key string) (*models.SystemConfig, error)
	// 根据配置类型查找配置列表
	FindByType(configType string) ([]*models.SystemConfig, error)
	// 获取所有配置
	FindAll() ([]*models.SystemConfig, error)
	// 创建配置
	Create(config *models.SystemConfig) error
	// 更新配置
	Update(config *models.SystemConfig) error
	// 删除配置（软删除）
	Delete(id uint) error
	// 检查配置键是否存在
	Exists(key string) (bool, error)
}

// configRepository 配置仓库实现
type configRepository struct {
	db *gorm.DB
}

// NewConfigRepository 创建配置仓库
func NewConfigRepository(db *gorm.DB) ConfigRepository {
	return &configRepository{db: db}
}

// FindByID 根据ID查找配置
func (r *configRepository) FindByID(id uint) (*models.SystemConfig, error) {
	var config models.SystemConfig
	err := r.db.First(&config, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConfigNotFound
		}
		return nil, err
	}
	return &config, nil
}

// FindByKey 根据配置键查找配置
func (r *configRepository) FindByKey(key string) (*models.SystemConfig, error) {
	var config models.SystemConfig
	err := r.db.Where("config_key = ?", key).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConfigNotFound
		}
		return nil, err
	}
	return &config, nil
}

// FindByType 根据配置类型查找配置列表
func (r *configRepository) FindByType(configType string) ([]*models.SystemConfig, error) {
	var configs []*models.SystemConfig
	err := r.db.Where("config_type = ?", configType).Find(&configs).Error
	return configs, err
}

// FindAll 获取所有配置
func (r *configRepository) FindAll() ([]*models.SystemConfig, error) {
	var configs []*models.SystemConfig
	err := r.db.Find(&configs).Error
	return configs, err
}

// Create 创建配置
func (r *configRepository) Create(config *models.SystemConfig) error {
	// 检查配置键是否已存在
	exists, err := r.Exists(config.ConfigKey)
	if err != nil {
		return err
	}
	if exists {
		return ErrConfigExists
	}

	return r.db.Create(config).Error
}

// Update 更新配置
func (r *configRepository) Update(config *models.SystemConfig) error {
	return r.db.Save(config).Error
}

// Delete 删除配置（软删除）
func (r *configRepository) Delete(id uint) error {
	return r.db.Delete(&models.SystemConfig{}, id).Error
}

// Exists 检查配置键是否存在
func (r *configRepository) Exists(key string) (bool, error) {
	var count int64
	err := r.db.Model(&models.SystemConfig{}).Where("config_key = ?", key).Count(&count).Error
	return count > 0, err
}
