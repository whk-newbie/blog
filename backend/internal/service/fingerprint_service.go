package service

import (
	"encoding/json"
	"time"

	"github.com/whk-newbie/blog/internal/models"
	"github.com/whk-newbie/blog/internal/pkg/utils"
	"github.com/whk-newbie/blog/internal/repository"
	"gorm.io/datatypes"
)

// FingerprintService 指纹服务接口
type FingerprintService interface {
	// 收集并存储指纹
	CollectFingerprint(fingerprintData interface{}, userAgent string) (*FingerprintResponse, error)
	// 根据哈希查找指纹
	GetByHash(hash string) (*models.Fingerprint, error)
	// 获取指纹列表
	List(page, pageSize int) (*FingerprintListResponse, error)
}

// FingerprintResponse 指纹响应
type FingerprintResponse struct {
	FingerprintID   uint   `json:"fingerprint_id"`
	FingerprintHash string `json:"fingerprint_hash"`
	IsNew           bool   `json:"is_new"`
}

// FingerprintListResponse 指纹列表响应
type FingerprintListResponse struct {
	Items      []models.Fingerprint `json:"items"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalPages int                  `json:"total_pages"`
}

// fingerprintService 指纹服务实现
type fingerprintService struct {
	fingerprintRepo repository.FingerprintRepository
}

// NewFingerprintService 创建指纹服务
func NewFingerprintService(fingerprintRepo repository.FingerprintRepository) FingerprintService {
	return &fingerprintService{
		fingerprintRepo: fingerprintRepo,
	}
}

// CollectFingerprint 收集并存储指纹
func (s *fingerprintService) CollectFingerprint(fingerprintData interface{}, userAgent string) (*FingerprintResponse, error) {
	// 计算指纹哈希
	hash, err := utils.CalculateFingerprintHash(fingerprintData)
	if err != nil {
		return nil, err
	}

	// 查找是否已存在
	existing, err := s.fingerprintRepo.FindByHash(hash)
	if err != nil && err != repository.ErrFingerprintNotFound {
		return nil, err
	}

	// 如果已存在，更新最后访问时间
	if existing != nil {
		if err := s.fingerprintRepo.UpdateLastSeen(existing.ID); err != nil {
			return nil, err
		}
		return &FingerprintResponse{
			FingerprintID:   existing.ID,
			FingerprintHash: hash,
			IsNew:           false,
		}, nil
	}

	// 如果不存在，创建新指纹
	// 将指纹数据转换为JSON
	jsonBytes, err := json.Marshal(fingerprintData)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	fingerprint := &models.Fingerprint{
		FingerprintHash: hash,
		FingerprintData: datatypes.JSON(jsonBytes),
		UserAgent:       userAgent,
		FirstSeenAt:     now,
		LastSeenAt:      now,
		VisitCount:      1,
	}

	if err := s.fingerprintRepo.Create(fingerprint); err != nil {
		return nil, err
	}

	return &FingerprintResponse{
		FingerprintID:   fingerprint.ID,
		FingerprintHash: hash,
		IsNew:           true,
	}, nil
}

// GetByHash 根据哈希查找指纹
func (s *fingerprintService) GetByHash(hash string) (*models.Fingerprint, error) {
	return s.fingerprintRepo.FindByHash(hash)
}

// List 获取指纹列表
func (s *fingerprintService) List(page, pageSize int) (*FingerprintListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	fingerprints, total, err := s.fingerprintRepo.List(offset, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &FingerprintListResponse{
		Items:      fingerprints,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
