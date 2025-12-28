package models

import (
	"time"

	"gorm.io/gorm"
)

// Tag 标签模型
type Tag struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Slug         string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"slug"`
	ArticleCount int            `gorm:"default:0" json:"article_count"` // 文章数量（冗余字段）
	CreatedBy    *uint          `gorm:"index" json:"created_by"`        // 创建者ID
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联
	Articles []Article `gorm:"many2many:article_tags" json:"articles,omitempty"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}

