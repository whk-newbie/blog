package service

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/whk-newbie/blog/internal/models"
	"github.com/whk-newbie/blog/internal/pkg/crypto"
	"github.com/whk-newbie/blog/internal/pkg/redis"
	"github.com/whk-newbie/blog/internal/repository"
)

var (
	ErrConfigNotFound     = errors.New("config not found")
	ErrConfigExists       = errors.New("config already exists")
	ErrInvalidConfigType  = errors.New("invalid config type")
	ErrCryptoNotAvailable = errors.New("crypto not available")
)

// ConfigService 配置服务接口
type ConfigService interface {
	// 获取配置列表
	GetConfigs(configType string) ([]*ConfigResponse, error)
	// 获取配置详情
	GetConfigByID(id uint) (*ConfigResponse, error)
	// 获取配置值（解密后）
	GetConfigValue(key string) (string, error)
	// 创建配置
	CreateConfig(req *CreateConfigRequest, userID uint) (*ConfigResponse, error)
	// 更新配置
	UpdateConfig(id uint, req *UpdateConfigRequest, userID uint) (*ConfigResponse, error)
	// 删除配置
	DeleteConfig(id uint) error
	// 生成爬虫Token
	GenerateCrawlerToken(name string, userID uint) (*CrawlerTokenResponse, error)
}

// ConfigResponse 配置响应
type ConfigResponse struct {
	ID          uint   `json:"id"`
	ConfigKey   string `json:"config_key"`
	ConfigValue string `json:"config_value"` // 敏感信息已脱敏
	ConfigType  string `json:"config_type"`
	IsEncrypted bool   `json:"is_encrypted"`
	IsActive    bool   `json:"is_active"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateConfigRequest 创建配置请求
type CreateConfigRequest struct {
	ConfigKey   string `json:"config_key" binding:"required"`
	ConfigValue string `json:"config_value" binding:"required"`
	ConfigType  string `json:"config_type" binding:"required"`
	IsEncrypted bool   `json:"is_encrypted"`
	IsActive    bool   `json:"is_active"`
	Description string `json:"description"`
}

// UpdateConfigRequest 更新配置请求
type UpdateConfigRequest struct {
	ConfigValue string `json:"config_value"`
	IsEncrypted *bool  `json:"is_encrypted"`
	IsActive    *bool  `json:"is_active"`
	Description string `json:"description"`
}

// CrawlerTokenResponse 爬虫Token响应
type CrawlerTokenResponse struct {
	Token string `json:"token"`
	Name  string `json:"name"`
}

// configService 配置服务实现
type configService struct {
	configRepo repository.ConfigRepository
	crypto     *crypto.Crypto
}

// NewConfigService 创建配置服务
func NewConfigService(configRepo repository.ConfigRepository, masterKey string) (ConfigService, error) {
	cryptoInstance, err := crypto.NewCrypto(masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize crypto: %w", err)
	}

	return &configService{
		configRepo: configRepo,
		crypto:     cryptoInstance,
	}, nil
}

// GetConfigs 获取配置列表
func (s *configService) GetConfigs(configType string) ([]*ConfigResponse, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("configs:%s", configType)
	if cached, err := s.getCachedConfigs(cacheKey); err == nil {
		return cached, nil
	}

	var configs []*models.SystemConfig
	var err error

	if configType != "" {
		configs, err = s.configRepo.FindByType(configType)
	} else {
		configs, err = s.configRepo.FindAll()
	}

	if err != nil {
		return nil, err
	}

	responses := make([]*ConfigResponse, 0, len(configs))
	for _, config := range configs {
		responses = append(responses, s.toConfigResponse(config, true))
	}

	// 缓存结果（配置信息缓存10分钟）
	_ = s.cacheConfigs(cacheKey, responses, 10*time.Minute)

	return responses, nil
}

// GetConfigByID 获取配置详情
func (s *configService) GetConfigByID(id uint) (*ConfigResponse, error) {
	config, err := s.configRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrConfigNotFound) {
			return nil, ErrConfigNotFound
		}
		return nil, err
	}

	return s.toConfigResponse(config, false), nil
}

// GetConfigValue 获取配置值（解密后）
func (s *configService) GetConfigValue(key string) (string, error) {
	config, err := s.configRepo.FindByKey(key)
	if err != nil {
		if errors.Is(err, repository.ErrConfigNotFound) {
			return "", ErrConfigNotFound
		}
		return "", err
	}

	if !config.IsActive {
		return "", ErrConfigNotFound
	}

	// 如果加密存储，需要解密
	if config.IsEncrypted {
		decrypted, err := s.crypto.Decrypt(config.ConfigValue)
		if err != nil {
			return "", fmt.Errorf("failed to decrypt config value: %w", err)
		}
		return decrypted, nil
	}

	return config.ConfigValue, nil
}

// CreateConfig 创建配置
func (s *configService) CreateConfig(req *CreateConfigRequest, userID uint) (*ConfigResponse, error) {
	// 加密配置值（如果需要）
	configValue := req.ConfigValue
	if req.IsEncrypted {
		encrypted, err := s.crypto.Encrypt(req.ConfigValue)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt config value: %w", err)
		}
		configValue = encrypted
	}

	config := &models.SystemConfig{
		ConfigKey:   req.ConfigKey,
		ConfigValue: configValue,
		ConfigType:  req.ConfigType,
		IsEncrypted: req.IsEncrypted,
		IsActive:    req.IsActive,
		Description: req.Description,
		CreatedBy:   &userID,
	}

	if err := s.configRepo.Create(config); err != nil {
		if errors.Is(err, repository.ErrConfigExists) {
			return nil, ErrConfigExists
		}
		return nil, err
	}

	// 清除配置缓存
	s.clearConfigCache(config.ConfigType)

	return s.toConfigResponse(config, true), nil
}

// UpdateConfig 更新配置
func (s *configService) UpdateConfig(id uint, req *UpdateConfigRequest, userID uint) (*ConfigResponse, error) {
	config, err := s.configRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrConfigNotFound) {
			return nil, ErrConfigNotFound
		}
		return nil, err
	}

	// 更新配置值
	if req.ConfigValue != "" {
		configValue := req.ConfigValue
		isEncrypted := config.IsEncrypted
		if req.IsEncrypted != nil {
			isEncrypted = *req.IsEncrypted
		}

		// 如果设置为加密，需要加密
		if isEncrypted {
			encrypted, err := s.crypto.Encrypt(req.ConfigValue)
			if err != nil {
				return nil, fmt.Errorf("failed to encrypt config value: %w", err)
			}
			configValue = encrypted
		}
		config.ConfigValue = configValue
		config.IsEncrypted = isEncrypted
	}

	// 更新其他字段
	if req.IsActive != nil {
		config.IsActive = *req.IsActive
	}
	if req.Description != "" {
		config.Description = req.Description
	}

	if err := s.configRepo.Update(config); err != nil {
		return nil, err
	}

	// 清除配置缓存
	s.clearConfigCache(config.ConfigType)

	return s.toConfigResponse(config, true), nil
}

// DeleteConfig 删除配置
func (s *configService) DeleteConfig(id uint) error {
	config, err := s.configRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrConfigNotFound) {
			return ErrConfigNotFound
		}
		return err
	}

	// 清除配置缓存
	s.clearConfigCache(config.ConfigType)

	return s.configRepo.Delete(id)
}

// GenerateCrawlerToken 生成爬虫Token
func (s *configService) GenerateCrawlerToken(name string, userID uint) (*CrawlerTokenResponse, error) {
	// 生成随机Token（32字节）
	tokenBytes := make([]byte, 24)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Base64编码，添加前缀
	token := "cr_" + base64.URLEncoding.EncodeToString(tokenBytes)

	// 保存到数据库
	configKey := fmt.Sprintf("crawler_token_%d", time.Now().Unix())
	config := &models.SystemConfig{
		ConfigKey:   configKey,
		ConfigValue: token,
		ConfigType:  models.ConfigTypeCrawlerToken,
		IsEncrypted: false, // Token本身不需要加密
		IsActive:    true,
		Description: name,
		CreatedBy:   &userID,
	}

	if err := s.configRepo.Create(config); err != nil {
		return nil, fmt.Errorf("failed to save token: %w", err)
	}

	return &CrawlerTokenResponse{
		Token: token,
		Name:  name,
	}, nil
}

// toConfigResponse 转换为响应格式
func (s *configService) toConfigResponse(config *models.SystemConfig, maskSensitive bool) *ConfigResponse {
	configValue := config.ConfigValue

	// 敏感信息脱敏
	if maskSensitive && configValue != "" {
		if config.IsEncrypted {
			// 加密的配置：尝试解密后脱敏
			decrypted, err := s.crypto.Decrypt(configValue)
			if err == nil {
				// 解密成功，根据配置类型进行智能脱敏
				configValue = s.maskSensitiveValue(decrypted, config.ConfigType)
			} else {
				// 解密失败，使用默认脱敏
				configValue = "***"
			}
		} else {
			// 未加密的配置：直接根据类型脱敏
			configValue = s.maskSensitiveValue(configValue, config.ConfigType)
		}
	}

	return &ConfigResponse{
		ID:          config.ID,
		ConfigKey:   config.ConfigKey,
		ConfigValue: configValue,
		ConfigType:  config.ConfigType,
		IsEncrypted: config.IsEncrypted,
		IsActive:    config.IsActive,
		Description: config.Description,
		CreatedAt:   config.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   config.UpdatedAt.Format(time.RFC3339),
	}
}

// maskSensitiveValue 根据配置类型智能脱敏
func (s *configService) maskSensitiveValue(value string, configType string) string {
	if value == "" {
		return "***"
	}

	switch configType {
	case models.ConfigTypeEmail:
		// 邮箱脱敏：***@example.com
		for i := 0; i < len(value); i++ {
			if value[i] == '@' {
				return "***" + value[i:]
			}
		}
		return "***@***"
	case models.ConfigTypeAPIToken, models.ConfigTypeCrawlerToken:
		// Token脱敏：显示前4位和后4位
		if len(value) > 8 {
			return value[:4] + "***" + value[len(value)-4:]
		}
		return "***"
	default:
		// 默认脱敏：显示后8位
		if len(value) > 8 {
			return "***" + value[len(value)-8:]
		}
		return "***"
	}
}

// cacheConfigs 缓存配置列表
func (s *configService) cacheConfigs(key string, configs []*ConfigResponse, ttl time.Duration) error {
	data, err := json.Marshal(configs)
	if err != nil {
		return err
	}
	return redis.Set(key, string(data), ttl)
}

// getCachedConfigs 获取缓存的配置列表
func (s *configService) getCachedConfigs(key string) ([]*ConfigResponse, error) {
	data, err := redis.GetValue(key)
	if err != nil {
		return nil, err
	}

	var configs []*ConfigResponse
	if err := json.Unmarshal([]byte(data), &configs); err != nil {
		return nil, err
	}

	return configs, nil
}

// clearConfigCache 清除配置缓存
func (s *configService) clearConfigCache(configType string) {
	cacheKey := fmt.Sprintf("configs:%s", configType)
	_ = redis.Del(cacheKey)
	// 同时清除所有配置的缓存
	_ = redis.Del("configs:")
}
