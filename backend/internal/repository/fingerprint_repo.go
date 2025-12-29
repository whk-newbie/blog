package repository

import (
	"errors"
	"time"

	"github.com/whk-newbie/blog/internal/models"
	"gorm.io/gorm"
)

var (
	ErrFingerprintNotFound = errors.New("fingerprint not found")
)

// FingerprintRepository 指纹仓库接口
type FingerprintRepository interface {
	// 创建指纹
	Create(fingerprint *models.Fingerprint) error
	// 根据哈希查找指纹
	FindByHash(hash string) (*models.Fingerprint, error)
	// 更新指纹（更新最后访问时间和访问次数）
	UpdateLastSeen(id uint) error
	// 增加访问次数
	IncrementVisitCount(id uint) error
	// 获取指纹列表
	List(offset, limit int) ([]models.Fingerprint, int64, error)
}

// fingerprintRepository 指纹仓库实现
type fingerprintRepository struct {
	db *gorm.DB
}

// NewFingerprintRepository 创建指纹仓库
func NewFingerprintRepository(db *gorm.DB) FingerprintRepository {
	return &fingerprintRepository{db: db}
}

// Create 创建指纹
func (r *fingerprintRepository) Create(fingerprint *models.Fingerprint) error {
	return r.db.Create(fingerprint).Error
}

// FindByHash 根据哈希查找指纹
func (r *fingerprintRepository) FindByHash(hash string) (*models.Fingerprint, error) {
	var fingerprint models.Fingerprint
	err := r.db.Where("fingerprint_hash = ?", hash).First(&fingerprint).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFingerprintNotFound
		}
		return nil, err
	}
	return &fingerprint, nil
}

// UpdateLastSeen 更新指纹（更新最后访问时间和访问次数）
func (r *fingerprintRepository) UpdateLastSeen(id uint) error {
	return r.db.Model(&models.Fingerprint{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_seen_at": time.Now(),
			"visit_count":  gorm.Expr("visit_count + ?", 1),
		}).Error
}

// IncrementVisitCount 增加访问次数
func (r *fingerprintRepository) IncrementVisitCount(id uint) error {
	return r.db.Model(&models.Fingerprint{}).
		Where("id = ?", id).
		Update("visit_count", gorm.Expr("visit_count + ?", 1)).Error
}

// List 获取指纹列表
func (r *fingerprintRepository) List(offset, limit int) ([]models.Fingerprint, int64, error) {
	var fingerprints []models.Fingerprint
	var total int64

	// 获取总数
	if err := r.db.Model(&models.Fingerprint{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表
	query := r.db.Order("last_seen_at DESC")
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&fingerprints).Error; err != nil {
		return nil, 0, err
	}

	return fingerprints, total, nil
}
