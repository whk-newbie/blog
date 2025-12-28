package models

import (
	"time"

	"gorm.io/datatypes"
)

// LogLevel 日志级别
type LogLevel string

const (
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
)

// SystemLog 系统日志模型
type SystemLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Level     LogLevel       `gorm:"type:varchar(20);not null;index" json:"level"` // 日志级别
	Message   string         `gorm:"type:text;not null" json:"message"`            // 日志消息
	Context   datatypes.JSON `gorm:"type:jsonb" json:"context"`                    // 上下文信息（请求信息、错误堆栈等）
	Source    string         `gorm:"type:varchar(100);index" json:"source"`        // 日志来源（模块名）
	UserID    *uint          `gorm:"index" json:"user_id"`                         // 如果是用户操作日志
	IPAddress string         `gorm:"type:varchar(45)" json:"ip_address"`           // IP地址
	CreatedAt time.Time      `gorm:"index" json:"created_at"`
}

// TableName 指定表名
func (SystemLog) TableName() string {
	return "system_logs"
}

