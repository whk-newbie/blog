package repository

import (
	"errors"

	"github.com/whk-newbie/blog/internal/models"
	"gorm.io/gorm"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryExists   = errors.New("category already exists")
)

// CategoryRepository 分类仓库接口
type CategoryRepository interface {
	// 创建分类
	Create(category *models.Category) error
	// 根据ID查找分类
	FindByID(id uint) (*models.Category, error)
	// 根据Slug查找分类
	FindBySlug(slug string) (*models.Category, error)
	// 更新分类
	Update(category *models.Category) error
	// 删除分类（软删除）
	Delete(id uint) error
	// 获取所有分类列表
	List(offset, limit int) ([]models.Category, int64, error)
	// 检查分类名是否存在
	ExistsByName(name string, excludeID uint) (bool, error)
	// 检查分类Slug是否存在
	ExistsBySlug(slug string, excludeID uint) (bool, error)
	// 增加文章数量
	IncrementArticleCount(id uint) error
	// 减少文章数量
	DecrementArticleCount(id uint) error
}

// categoryRepository 分类仓库实现
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓库
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// Create 创建分类
func (r *categoryRepository) Create(category *models.Category) error {
	// 检查名称是否已存在
	exists, err := r.ExistsByName(category.Name, 0)
	if err != nil {
		return err
	}
	if exists {
		return ErrCategoryExists
	}

	// 检查Slug是否已存在
	exists, err = r.ExistsBySlug(category.Slug, 0)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("category slug already exists")
	}

	return r.db.Create(category).Error
}

// FindByID 根据ID查找分类
func (r *categoryRepository) FindByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

// FindBySlug 根据Slug查找分类
func (r *categoryRepository) FindBySlug(slug string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

// Update 更新分类
func (r *categoryRepository) Update(category *models.Category) error {
	// 检查名称是否被其他分类使用
	exists, err := r.ExistsByName(category.Name, category.ID)
	if err != nil {
		return err
	}
	if exists {
		return ErrCategoryExists
	}

	// 检查Slug是否被其他分类使用
	exists, err = r.ExistsBySlug(category.Slug, category.ID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("category slug already exists")
	}

	return r.db.Save(category).Error
}

// Delete 删除分类（软删除）
func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}

// List 获取所有分类列表
func (r *categoryRepository) List(offset, limit int) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	// 获取总数
	if err := r.db.Model(&models.Category{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分类列表
	query := r.db.Model(&models.Category{}).Order("sort_order ASC, created_at DESC")
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

// ExistsByName 检查分类名是否存在
func (r *categoryRepository) ExistsByName(name string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.Category{}).Where("name = ?", name)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// ExistsBySlug 检查分类Slug是否存在
func (r *categoryRepository) ExistsBySlug(slug string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.Category{}).Where("slug = ?", slug)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// IncrementArticleCount 增加文章数量
func (r *categoryRepository) IncrementArticleCount(id uint) error {
	return r.db.Model(&models.Category{}).
		Where("id = ?", id).
		Update("article_count", gorm.Expr("article_count + ?", 1)).Error
}

// DecrementArticleCount 减少文章数量
func (r *categoryRepository) DecrementArticleCount(id uint) error {
	return r.db.Model(&models.Category{}).
		Where("id = ?", id).
		Update("article_count", gorm.Expr("GREATEST(article_count - 1, 0)")).Error
}
