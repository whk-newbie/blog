package repository

import (
	"github.com/whk-newbie/blog/internal/models"
	"gorm.io/gorm"
)

type AIChatRepository interface {
	Save(msg *models.AIChatHistory) error
	FindByProviderID(providerID uint, limit int) ([]*models.AIChatHistory, error)
	DeleteByProviderID(providerID uint) error
}

type aiChatRepository struct {
	db *gorm.DB
}

func NewAIChatRepository(db *gorm.DB) AIChatRepository {
	return &aiChatRepository{db: db}
}

func (r *aiChatRepository) Save(msg *models.AIChatHistory) error {
	return r.db.Create(msg).Error
}

func (r *aiChatRepository) FindByProviderID(providerID uint, limit int) ([]*models.AIChatHistory, error) {
	var history []*models.AIChatHistory
	if limit <= 0 {
		limit = 50
	}
	err := r.db.Where("provider_id = ?", providerID).Order("created_at DESC").Limit(limit).Find(&history).Error
	// Reverse to chronological order
	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}
	return history, err
}

func (r *aiChatRepository) DeleteByProviderID(providerID uint) error {
	return r.db.Where("provider_id = ?", providerID).Delete(&models.AIChatHistory{}).Error
}
