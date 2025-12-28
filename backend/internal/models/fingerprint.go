package models

import (
	"time"

	"gorm.io/datatypes"
)

// Fingerprint 浏览器指纹模型
type Fingerprint struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	FingerprintHash string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"fingerprint_hash"` // 指纹哈希值（SHA256）
	FingerprintData datatypes.JSON `gorm:"type:jsonb;not null" json:"fingerprint_data"`                   // 完整指纹信息
	UserAgent       string         `gorm:"type:text" json:"user_agent"`
	FirstSeenAt     time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"first_seen_at"`
	LastSeenAt      time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP;index" json:"last_seen_at"`
	VisitCount      int            `gorm:"default:1" json:"visit_count"` // 访问次数
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`

	// 关联
	Visits []Visit `gorm:"foreignKey:FingerprintID" json:"visits,omitempty"`
}

// TableName 指定表名
func (Fingerprint) TableName() string {
	return "fingerprints"
}

