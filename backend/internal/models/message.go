package models

import (
	"time"
)

// Message 留言记录模型（仅用于统计）
type Message struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"type:varchar(100)" json:"name"`
	Email         string    `gorm:"type:varchar(100);index" json:"email"`
	FingerprintID *uint     `gorm:"index" json:"fingerprint_id"` // 留言者的浏览器指纹ID
	CreatedAt     time.Time `gorm:"index" json:"created_at"`

	// 关联
	Fingerprint *Fingerprint `gorm:"foreignKey:FingerprintID" json:"fingerprint,omitempty"`
}

// TableName 指定表名
func (Message) TableName() string {
	return "messages"
}

