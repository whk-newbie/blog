package models

import (
	"time"

	"gorm.io/gorm"
)

// ArticleStatus 文章状态
type ArticleStatus string

const (
	ArticleStatusDraft     ArticleStatus = "draft"     // 草稿
	ArticleStatusPublished ArticleStatus = "published" // 已发布
)

// Article 文章模型
type Article struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Title      string         `gorm:"type:varchar(255);not null" json:"title"`
	Slug       string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Summary    string         `gorm:"type:text" json:"summary"`                            // 文章摘要
	Content    string         `gorm:"type:text;not null" json:"content"`                   // 富文本内容（HTML）
	CoverImage string         `gorm:"type:varchar(500)" json:"cover_image"`                // 封面图片URL
	CategoryID *uint          `gorm:"index" json:"category_id"`                            // 分类ID
	Status     ArticleStatus  `gorm:"type:varchar(20);not null;default:'draft';index" json:"status"` // 文章状态
	PublishAt  *time.Time     `gorm:"index" json:"publish_at"`                             // 发布时间（可预设未来时间）
	ViewCount  int            `gorm:"default:0;index" json:"view_count"`                   // 浏览次数
	LikeCount  int            `gorm:"default:0" json:"like_count"`                         // 点赞数（预留）
	IsTop      bool           `gorm:"default:false;index" json:"is_top"`                   // 是否置顶
	IsFeatured bool           `gorm:"default:false" json:"is_featured"`                    // 是否推荐
	AuthorID   *uint          `gorm:"index" json:"author_id"`                              // 作者ID
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Author   *Admin    `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Tags     []Tag     `gorm:"many2many:article_tags" json:"tags,omitempty"`
	Visits   []Visit   `gorm:"foreignKey:ArticleID" json:"visits,omitempty"`
}

// TableName 指定表名
func (Article) TableName() string {
	return "articles"
}

// ArticleTag 文章标签关联模型
type ArticleTag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ArticleID uint      `gorm:"not null;index;uniqueIndex:idx_article_tag" json:"article_id"`
	TagID     uint      `gorm:"not null;index;uniqueIndex:idx_article_tag" json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (ArticleTag) TableName() string {
	return "article_tags"
}

