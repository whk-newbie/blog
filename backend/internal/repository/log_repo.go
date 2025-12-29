package repository

import (
	"errors"
	"time"

	"github.com/whk-newbie/blog/internal/models"
	"gorm.io/gorm"
)

var (
	ErrLogNotFound = errors.New("log not found")
)

// LogRepository 日志仓库接口
type LogRepository interface {
	// 根据ID查找日志
	FindByID(id uint) (*models.SystemLog, error)
	// 查询日志列表（支持分页和筛选）
	FindLogs(page, pageSize int, level, source string, startDate, endDate *time.Time) ([]*models.SystemLog, int64, error)
	// 创建日志
	Create(log *models.SystemLog) error
	// 批量创建日志
	BatchCreate(logs []*models.SystemLog) error
	// 清理旧日志
	CleanupOldLogs(retentionDays int) (int64, error)
	// 统计日志数量
	Count(level, source string, startDate, endDate *time.Time) (int64, error)
}

// logRepository 日志仓库实现
type logRepository struct {
	db *gorm.DB
}

// NewLogRepository 创建日志仓库
func NewLogRepository(db *gorm.DB) LogRepository {
	return &logRepository{db: db}
}

// FindByID 根据ID查找日志
func (r *logRepository) FindByID(id uint) (*models.SystemLog, error) {
	var log models.SystemLog
	err := r.db.First(&log, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLogNotFound
		}
		return nil, err
	}
	return &log, nil
}

// FindLogs 查询日志列表（支持分页和筛选）
func (r *logRepository) FindLogs(page, pageSize int, level, source string, startDate, endDate *time.Time) ([]*models.SystemLog, int64, error) {
	var logs []*models.SystemLog
	var total int64

	query := r.db.Model(&models.SystemLog{})

	// 筛选条件
	if level != "" {
		query = query.Where("level = ?", level)
	}
	if source != "" {
		query = query.Where("source = ?", source)
	}
	if startDate != nil {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", endDate)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&logs).Error

	return logs, total, err
}

// Create 创建日志
func (r *logRepository) Create(log *models.SystemLog) error {
	return r.db.Create(log).Error
}

// BatchCreate 批量创建日志
func (r *logRepository) BatchCreate(logs []*models.SystemLog) error {
	if len(logs) == 0 {
		return nil
	}
	return r.db.CreateInBatches(logs, 100).Error
}

// CleanupOldLogs 清理旧日志
func (r *logRepository) CleanupOldLogs(retentionDays int) (int64, error) {
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)
	result := r.db.Where("created_at < ?", cutoffDate).Delete(&models.SystemLog{})
	return result.RowsAffected, result.Error
}

// Count 统计日志数量
func (r *logRepository) Count(level, source string, startDate, endDate *time.Time) (int64, error) {
	var count int64
	query := r.db.Model(&models.SystemLog{})

	if level != "" {
		query = query.Where("level = ?", level)
	}
	if source != "" {
		query = query.Where("source = ?", source)
	}
	if startDate != nil {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", endDate)
	}

	err := query.Count(&count).Error
	return count, err
}
