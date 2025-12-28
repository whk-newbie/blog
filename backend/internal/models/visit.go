package models

import (
	"time"
)

// Visit 访问记录模型
type Visit struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	FingerprintID *uint      `gorm:"index" json:"fingerprint_id"` // 浏览器指纹ID
	URL           string     `gorm:"type:varchar(500);not null;index" json:"url"` // 访问的URL
	Referrer      string     `gorm:"type:varchar(500)" json:"referrer"` // 来源URL
	PageTitle     string     `gorm:"type:varchar(255)" json:"page_title"` // 页面标题
	ArticleID     *uint      `gorm:"index" json:"article_id"` // 如果访问的是文章页
	VisitTime     time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP;index" json:"visit_time"`
	StayDuration  *int       `json:"stay_duration"` // 停留时间（秒）
	UserAgent     string     `gorm:"type:text" json:"user_agent"`
	CreatedAt     time.Time  `json:"created_at"`

	// 关联
	Fingerprint *Fingerprint `gorm:"foreignKey:FingerprintID" json:"fingerprint,omitempty"`
	Article     *Article     `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
}

// TableName 指定表名
func (Visit) TableName() string {
	return "visits"
}

