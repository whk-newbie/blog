package models

import (
	"time"

	"gorm.io/datatypes"
)

// CrawlTaskStatus 爬虫任务状态
type CrawlTaskStatus string

const (
	CrawlTaskStatusRunning   CrawlTaskStatus = "running"   // 运行中
	CrawlTaskStatusCompleted CrawlTaskStatus = "completed" // 已完成
	CrawlTaskStatusFailed    CrawlTaskStatus = "failed"    // 失败
)

// CrawlTask 爬虫任务模型
type CrawlTask struct {
	ID             uint            `gorm:"primaryKey" json:"id"`
	TaskID         string          `gorm:"type:varchar(100);uniqueIndex;not null" json:"task_id"` // 任务唯一标识
	TaskName       string          `gorm:"type:varchar(255);not null" json:"task_name"`           // 任务名称
	Status         CrawlTaskStatus `gorm:"type:varchar(20);not null;default:'running';index" json:"status"` // 任务状态
	Progress       int             `gorm:"default:0" json:"progress"`                             // 进度（0-100）
	Message        string          `gorm:"type:text" json:"message"`                              // 状态消息
	StartTime      time.Time       `gorm:"not null;default:CURRENT_TIMESTAMP;index" json:"start_time"`
	EndTime        *time.Time      `json:"end_time"`
	Duration       *int            `json:"duration"`                            // 运行时长（秒）
	CreatedByToken string          `gorm:"type:varchar(64);index" json:"-"`     // 创建任务的Token（不返回给客户端）
	Metadata       datatypes.JSON  `gorm:"type:jsonb" json:"metadata"`          // 额外元数据
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// TableName 指定表名
func (CrawlTask) TableName() string {
	return "crawl_tasks"
}

