package models

import (
	"time"

	"gorm.io/gorm"
)

// Category 分类模型
type Category struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Slug         string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"slug"` // URL友好的标识符
	Description  string         `gorm:"type:text" json:"description"`
	SortOrder    int            `gorm:"default:0;index" json:"sort_order"`       // 排序权重
	ArticleCount int            `gorm:"default:0" json:"article_count"`          // 文章数量（冗余字段）
	CreatedBy    *uint          `gorm:"index" json:"created_by"`                 // 创建者ID
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联
	Articles []Article `gorm:"foreignKey:CategoryID" json:"articles,omitempty"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}

