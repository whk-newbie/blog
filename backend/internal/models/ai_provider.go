package models

import (
	"time"

	"gorm.io/gorm"
)

// AIProvider represents an AI service provider configuration
type AIProvider struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(100);not null" json:"name"`
	ProviderType string         `gorm:"type:varchar(20);not null" json:"provider_type"` // claude, openai, deepseek, custom
	APIKey       string         `gorm:"type:text;not null" json:"api_key"`               // AES encrypted
	BaseURL      string         `gorm:"type:varchar(500)" json:"base_url"`
	Model        string         `gorm:"type:varchar(100);not null" json:"model"`
	IsEnabled    bool           `gorm:"default:true;index" json:"is_enabled"`
	SortOrder    int            `gorm:"default:0" json:"sort_order"`
	Balance      float64        `gorm:"type:decimal(10,4);default:0" json:"balance"`
	LastCheckAt  *time.Time     `json:"last_check_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (AIProvider) TableName() string {
	return "ai_providers"
}

// AIChatHistory stores chat conversation history
type AIChatHistory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ProviderID uint      `gorm:"index;not null" json:"provider_id"`
	Role       string    `gorm:"type:varchar(20);not null" json:"role"` // user / assistant
	Content    string    `gorm:"type:text;not null" json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

func (AIChatHistory) TableName() string {
	return "ai_chat_history"
}
