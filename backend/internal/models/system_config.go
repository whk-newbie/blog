package models

import (
	"time"

	"gorm.io/gorm"
)

// SystemConfig 系统配置模型
type SystemConfig struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ConfigKey   string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"config_key"` // 配置键
	ConfigValue string         `gorm:"type:text" json:"config_value"`                            // 配置值（加密存储）
	ConfigType  string         `gorm:"type:varchar(50);not null;index" json:"config_type"`       // 配置类型
	IsEncrypted bool           `gorm:"default:true" json:"is_encrypted"`                         // 是否加密存储
	IsActive    bool           `gorm:"default:true;index" json:"is_active"`                      // 是否启用
	Description string         `gorm:"type:text" json:"description"`                             // 配置描述
	CreatedBy   *uint          `gorm:"index" json:"created_by"`                                  // 创建者ID
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (SystemConfig) TableName() string {
	return "system_configs"
}

// 配置类型常量
const (
	ConfigTypeEmail        = "email"         // 邮箱地址
	ConfigTypeAPIToken     = "api_token"     // API Token
	ConfigTypeCrawlerToken = "crawler_token" // 爬虫认证Token
	ConfigTypeAppKey       = "application_key" // 应用密钥
	ConfigTypeSalt         = "salt"          // 加密盐
	ConfigTypeIPBlacklist  = "ip_blacklist"  // IP黑名单
	ConfigTypeSiteInfo     = "site_info"     // 站点信息(博客标题、备案信息等)
)

