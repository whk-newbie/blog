package repository

import (
	"errors"
	"time"

	"github.com/whk-newbie/blog/internal/models"
	"gorm.io/gorm"
)

var (
	ErrArticleNotFound = errors.New("article not found")
	ErrArticleExists   = errors.New("article slug already exists")
)

// ArticleFilter 文章筛选条件
type ArticleFilter struct {
	CategoryID *uint
	TagID      *uint
	Status     *models.ArticleStatus
	IsTop      *bool
	IsFeatured *bool
	Keyword    string // 搜索关键词
}

// ArticleRepository 文章仓库接口
type ArticleRepository interface {
	// 创建文章
	Create(article *models.Article) error
	// 根据ID查找文章
	FindByID(id uint) (*models.Article, error)
	// 根据ID查找文章（包含关联）
	FindByIDWithAssociations(id uint) (*models.Article, error)
	// 根据Slug查找文章
	FindBySlug(slug string) (*models.Article, error)
	// 根据Slug查找文章（包含关联）
	FindBySlugWithAssociations(slug string) (*models.Article, error)
	// 更新文章
	Update(article *models.Article) error
	// 删除文章（软删除）
	Delete(id uint) error
	// 获取文章列表（带筛选）
	List(filter *ArticleFilter, offset, limit int) ([]models.Article, int64, error)
	// 检查Slug是否存在
	ExistsBySlug(slug string, excludeID uint) (bool, error)
	// 增加浏览量
	IncrementViewCount(id uint) error
	// 更新文章标签关联
	UpdateTags(articleID uint, tagIDs []uint) error
	// 获取已发布文章列表（公开访问）
	ListPublished(filter *ArticleFilter, offset, limit int) ([]models.Article, int64, error)
	// 全文搜索
	Search(keyword string, offset, limit int) ([]models.Article, int64, error)
	// 获取需要发布的文章（定时发布）
	GetPendingPublish() ([]models.Article, error)
	// 获取按日期统计的文章发布数量
	GetPublishStatsByDate(startDate, endDate time.Time) ([]PublishStat, error)
}

// PublishStat 文章发布统计
type PublishStat struct {
	Date  string `json:"date"`  // 日期 (YYYY-MM-DD)
	Count int64  `json:"count"` // 发布数量
}

// articleRepository 文章仓库实现
type articleRepository struct {
	db *gorm.DB
}

// NewArticleRepository 创建文章仓库
func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

// Create 创建文章
func (r *articleRepository) Create(article *models.Article) error {
	// 检查Slug是否已存在
	exists, err := r.ExistsBySlug(article.Slug, 0)
	if err != nil {
		return err
	}
	if exists {
		return ErrArticleExists
	}

	return r.db.Create(article).Error
}

// FindByID 根据ID查找文章
func (r *articleRepository) FindByID(id uint) (*models.Article, error) {
	var article models.Article
	err := r.db.First(&article, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}
	return &article, nil
}

// FindByIDWithAssociations 根据ID查找文章（包含关联）
func (r *articleRepository) FindByIDWithAssociations(id uint) (*models.Article, error) {
	var article models.Article
	err := r.db.Preload("Category").
		Preload("Tags").
		Preload("Author").
		First(&article, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}
	return &article, nil
}

// FindBySlug 根据Slug查找文章
func (r *articleRepository) FindBySlug(slug string) (*models.Article, error) {
	var article models.Article
	err := r.db.Where("slug = ?", slug).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}
	return &article, nil
}

// FindBySlugWithAssociations 根据Slug查找文章（包含关联）
func (r *articleRepository) FindBySlugWithAssociations(slug string) (*models.Article, error) {
	var article models.Article
	err := r.db.Preload("Category").
		Preload("Tags").
		Preload("Author").
		Where("slug = ?", slug).
		First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}
	return &article, nil
}

// Update 更新文章
func (r *articleRepository) Update(article *models.Article) error {
	// 检查Slug是否被其他文章使用
	exists, err := r.ExistsBySlug(article.Slug, article.ID)
	if err != nil {
		return err
	}
	if exists {
		return ErrArticleExists
	}

	return r.db.Save(article).Error
}

// Delete 删除文章（软删除）
func (r *articleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Article{}, id).Error
}

// List 获取文章列表（带筛选）
func (r *articleRepository) List(filter *ArticleFilter, offset, limit int) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	query := r.db.Model(&models.Article{})

	// 应用筛选条件
	query = r.applyFilter(query, filter)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取文章列表
	query = query.Preload("Category").
		Preload("Tags").
		Preload("Author").
		Order("is_top DESC, created_at DESC")

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// ListPublished 获取已发布文章列表（公开访问）
func (r *articleRepository) ListPublished(filter *ArticleFilter, offset, limit int) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	query := r.db.Model(&models.Article{}).
		Where("status = ?", models.ArticleStatusPublished).
		Where("publish_at IS NULL OR publish_at <= ?", time.Now())

	// 应用筛选条件
	query = r.applyFilter(query, filter)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取文章列表
	query = query.Preload("Category").
		Preload("Tags").
		Order("is_top DESC, publish_at DESC, created_at DESC")

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// applyFilter 应用筛选条件
func (r *articleRepository) applyFilter(query *gorm.DB, filter *ArticleFilter) *gorm.DB {
	if filter == nil {
		return query
	}

	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}

	if filter.TagID != nil {
		query = query.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Where("article_tags.tag_id = ?", *filter.TagID)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.IsTop != nil {
		query = query.Where("is_top = ?", *filter.IsTop)
	}

	if filter.IsFeatured != nil {
		query = query.Where("is_featured = ?", *filter.IsFeatured)
	}

	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("title LIKE ? OR summary LIKE ? OR content LIKE ?",
			keyword, keyword, keyword)
	}

	return query
}

// Search 全文搜索
func (r *articleRepository) Search(keyword string, offset, limit int) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	// 使用PostgreSQL全文搜索
	query := r.db.Model(&models.Article{}).
		Where("status = ?", models.ArticleStatusPublished).
		Where("publish_at IS NULL OR publish_at <= ?", time.Now()).
		Where("to_tsvector('simple', title || ' ' || summary || ' ' || content) @@ plainto_tsquery('simple', ?)", keyword)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取文章列表
	query = query.Preload("Category").
		Preload("Tags").
		Order("is_top DESC, created_at DESC")

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// ExistsBySlug 检查Slug是否存在
func (r *articleRepository) ExistsBySlug(slug string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.Article{}).Where("slug = ?", slug)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// IncrementViewCount 增加浏览量
func (r *articleRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&models.Article{}).
		Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// UpdateTags 更新文章标签关联
func (r *articleRepository) UpdateTags(articleID uint, tagIDs []uint) error {
	// 开启事务
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除旧的标签关联
		if err := tx.Where("article_id = ?", articleID).Delete(&models.ArticleTag{}).Error; err != nil {
			return err
		}

		// 创建新的标签关联
		if len(tagIDs) > 0 {
			var articleTags []models.ArticleTag
			for _, tagID := range tagIDs {
				articleTags = append(articleTags, models.ArticleTag{
					ArticleID: articleID,
					TagID:     tagID,
				})
			}
			if err := tx.Create(&articleTags).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetPendingPublish 获取需要发布的文章（定时发布）
func (r *articleRepository) GetPendingPublish() ([]models.Article, error) {
	var articles []models.Article
	now := time.Now()

	err := r.db.Where("status = ?", models.ArticleStatusDraft).
		Where("publish_at IS NOT NULL AND publish_at <= ?", now).
		Find(&articles).Error

	return articles, err
}

// GetPublishStatsByDate 获取按日期统计的文章发布数量
func (r *articleRepository) GetPublishStatsByDate(startDate, endDate time.Time) ([]PublishStat, error) {
	var stats []PublishStat

	// 使用 COALESCE 处理 publish_at 为 NULL 的情况，使用 created_at
	err := r.db.Model(&models.Article{}).
		Select("DATE(COALESCE(publish_at, created_at)) as date, COUNT(*) as count").
		Where("DATE(COALESCE(publish_at, created_at)) >= ?", startDate.Format("2006-01-02")).
		Where("DATE(COALESCE(publish_at, created_at)) <= ?", endDate.Format("2006-01-02")).
		Group("DATE(COALESCE(publish_at, created_at))").
		Order("date ASC").
		Scan(&stats).Error

	return stats, err
}
