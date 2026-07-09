package repository

import (
	"errors"

	"github.com/whk-newbie/blog/internal/models"
	"gorm.io/gorm"
)

var ErrAIProviderNotFound = errors.New("AI provider not found")

// AIProviderRepository AI provider data access interface
type AIProviderRepository interface {
	FindAll() ([]*models.AIProvider, error)
	FindEnabled() ([]*models.AIProvider, error)
	FindByID(id uint) (*models.AIProvider, error)
	Create(provider *models.AIProvider) error
	Update(provider *models.AIProvider) error
	Delete(id uint) error
}

type aiProviderRepository struct {
	db *gorm.DB
}

// NewAIProviderRepository creates a new AI provider repository
func NewAIProviderRepository(db *gorm.DB) AIProviderRepository {
	return &aiProviderRepository{db: db}
}

func (r *aiProviderRepository) FindAll() ([]*models.AIProvider, error) {
	var providers []*models.AIProvider
	err := r.db.Order("sort_order ASC").Find(&providers).Error
	return providers, err
}

func (r *aiProviderRepository) FindEnabled() ([]*models.AIProvider, error) {
	var providers []*models.AIProvider
	err := r.db.Where("is_enabled = ?", true).Order("sort_order ASC").Find(&providers).Error
	return providers, err
}

func (r *aiProviderRepository) FindByID(id uint) (*models.AIProvider, error) {
	var provider models.AIProvider
	err := r.db.First(&provider, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAIProviderNotFound
	}
	return &provider, err
}

func (r *aiProviderRepository) Create(provider *models.AIProvider) error {
	return r.db.Create(provider).Error
}

func (r *aiProviderRepository) Update(provider *models.AIProvider) error {
	return r.db.Save(provider).Error
}

func (r *aiProviderRepository) Delete(id uint) error {
	return r.db.Delete(&models.AIProvider{}, id).Error
}
