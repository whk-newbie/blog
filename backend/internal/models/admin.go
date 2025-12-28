package models

import (
	"time"

	"gorm.io/gorm"
)

// Admin 管理员模型
type Admin struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	Username          string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password          string         `gorm:"type:varchar(255);not null" json:"-"` // 不在JSON中返回密码
	Email             string         `gorm:"type:varchar(100)" json:"email"`
	IsDefaultPassword bool           `gorm:"default:true" json:"is_default_password"` // 是否使用默认密码
	LastLoginAt       *time.Time     `json:"last_login_at"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (Admin) TableName() string {
	return "admins"
}

