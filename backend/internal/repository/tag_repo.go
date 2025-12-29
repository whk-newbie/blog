package repository

import (
	"errors"

	"github.com/iambaby/blog/internal/models"
	"gorm.io/gorm"
)

var (
	ErrTagNotFound = errors.New("tag not found")
	ErrTagExists   = errors.New("tag already exists")
)

// TagRepository 标签仓库接口
type TagRepository interface {
	// 创建标签
	Create(tag *models.Tag) error
	// 根据ID查找标签
	FindByID(id uint) (*models.Tag, error)
	// 根据Slug查找标签
	FindBySlug(slug string) (*models.Tag, error)
	// 根据多个ID查找标签
	FindByIDs(ids []uint) ([]models.Tag, error)
	// 更新标签
	Update(tag *models.Tag) error
	// 删除标签（软删除）
	Delete(id uint) error
	// 获取所有标签列表
	List(offset, limit int) ([]models.Tag, int64, error)
	// 检查标签名是否存在
	ExistsByName(name string, excludeID uint) (bool, error)
	// 检查标签Slug是否存在
	ExistsBySlug(slug string, excludeID uint) (bool, error)
	// 增加文章数量
	IncrementArticleCount(id uint) error
	// 减少文章数量
	DecrementArticleCount(id uint) error
}

// tagRepository 标签仓库实现
type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository 创建标签仓库
func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

// Create 创建标签
func (r *tagRepository) Create(tag *models.Tag) error {
	// 检查名称是否已存在
	exists, err := r.ExistsByName(tag.Name, 0)
	if err != nil {
		return err
	}
	if exists {
		return ErrTagExists
	}

	// 检查Slug是否已存在
	exists, err = r.ExistsBySlug(tag.Slug, 0)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("tag slug already exists")
	}

	return r.db.Create(tag).Error
}

// FindByID 根据ID查找标签
func (r *tagRepository) FindByID(id uint) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.First(&tag, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTagNotFound
		}
		return nil, err
	}
	return &tag, nil
}

// FindBySlug 根据Slug查找标签
func (r *tagRepository) FindBySlug(slug string) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.Where("slug = ?", slug).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTagNotFound
		}
		return nil, err
	}
	return &tag, nil
}

// FindByIDs 根据多个ID查找标签
func (r *tagRepository) FindByIDs(ids []uint) ([]models.Tag, error) {
	var tags []models.Tag
	if len(ids) == 0 {
		return tags, nil
	}
	err := r.db.Where("id IN ?", ids).Find(&tags).Error
	return tags, err
}

// Update 更新标签
func (r *tagRepository) Update(tag *models.Tag) error {
	// 检查名称是否被其他标签使用
	exists, err := r.ExistsByName(tag.Name, tag.ID)
	if err != nil {
		return err
	}
	if exists {
		return ErrTagExists
	}

	// 检查Slug是否被其他标签使用
	exists, err = r.ExistsBySlug(tag.Slug, tag.ID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("tag slug already exists")
	}

	return r.db.Save(tag).Error
}

// Delete 删除标签（软删除）
func (r *tagRepository) Delete(id uint) error {
	return r.db.Delete(&models.Tag{}, id).Error
}

// List 获取所有标签列表
func (r *tagRepository) List(offset, limit int) ([]models.Tag, int64, error) {
	var tags []models.Tag
	var total int64

	// 获取总数
	if err := r.db.Model(&models.Tag{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取标签列表
	query := r.db.Order("article_count DESC, created_at DESC")
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&tags).Error; err != nil {
		return nil, 0, err
	}

	return tags, total, nil
}

// ExistsByName 检查标签名是否存在
func (r *tagRepository) ExistsByName(name string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.Tag{}).Where("name = ?", name)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// ExistsBySlug 检查标签Slug是否存在
func (r *tagRepository) ExistsBySlug(slug string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.Tag{}).Where("slug = ?", slug)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// IncrementArticleCount 增加文章数量
func (r *tagRepository) IncrementArticleCount(id uint) error {
	return r.db.Model(&models.Tag{}).
		Where("id = ?", id).
		Update("article_count", gorm.Expr("article_count + ?", 1)).Error
}

// DecrementArticleCount 减少文章数量
func (r *tagRepository) DecrementArticleCount(id uint) error {
	return r.db.Model(&models.Tag{}).
		Where("id = ?", id).
		Update("article_count", gorm.Expr("GREATEST(article_count - 1, 0)")).Error
}
