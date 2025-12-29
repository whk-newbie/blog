package service

import (
	"errors"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/whk-newbie/blog/internal/models"
	"github.com/whk-newbie/blog/internal/repository"
)

var (
	ErrArticleTitleRequired = errors.New("article title is required")
	ErrArticleNotFound      = repository.ErrArticleNotFound
)

// ArticleService 文章服务接口
type ArticleService interface {
	// 创建文章
	Create(req *CreateArticleRequest, authorID uint) (*models.Article, error)
	// 获取文章详情
	GetByID(id uint) (*models.Article, error)
	// 获取文章详情（通过Slug）
	GetBySlug(slug string) (*models.Article, error)
	// 更新文章
	Update(id uint, req *UpdateArticleRequest) (*models.Article, error)
	// 删除文章
	Delete(id uint) error
	// 获取文章列表（管理员）
	List(req *ArticleListRequest) (*ArticleListResponse, error)
	// 获取已发布文章列表（公开）
	ListPublished(req *ArticleListRequest) (*ArticleListResponse, error)
	// 发布文章
	Publish(id uint) error
	// 取消发布（转为草稿）
	Unpublish(id uint) error
	// 增加浏览量
	IncrementViewCount(id uint) error
	// 搜索文章
	Search(keyword string, page, pageSize int) (*ArticleListResponse, error)
	// 处理定时发布
	ProcessScheduledPublish() error
}

// CreateArticleRequest 创建文章请求
type CreateArticleRequest struct {
	Title      string               `json:"title" binding:"required"`
	Slug       string               `json:"slug"`
	Summary    string               `json:"summary"`
	Content    string               `json:"content" binding:"required"`
	CoverImage string               `json:"cover_image"`
	CategoryID *uint                `json:"category_id"`
	TagIDs     []uint               `json:"tag_ids"`
	Status     models.ArticleStatus `json:"status"`
	PublishAt  *time.Time           `json:"publish_at"`
	IsTop      bool                 `json:"is_top"`
	IsFeatured bool                 `json:"is_featured"`
}

// UpdateArticleRequest 更新文章请求
type UpdateArticleRequest struct {
	Title      string               `json:"title" binding:"required"`
	Slug       string               `json:"slug"`
	Summary    string               `json:"summary"`
	Content    string               `json:"content" binding:"required"`
	CoverImage string               `json:"cover_image"`
	CategoryID *uint                `json:"category_id"`
	TagIDs     []uint               `json:"tag_ids"`
	Status     models.ArticleStatus `json:"status"`
	PublishAt  *time.Time           `json:"publish_at"`
	IsTop      bool                 `json:"is_top"`
	IsFeatured bool                 `json:"is_featured"`
}

// ArticleListRequest 文章列表请求
type ArticleListRequest struct {
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
	CategoryID *uint                 `json:"category_id"`
	TagID      *uint                 `json:"tag_id"`
	Status     *models.ArticleStatus `json:"status"`
	IsTop      *bool                 `json:"is_top"`
	IsFeatured *bool                 `json:"is_featured"`
	Keyword    string                `json:"keyword"`
}

// ArticleListResponse 文章列表响应
type ArticleListResponse struct {
	Items      []models.Article `json:"items"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

// articleService 文章服务实现
type articleService struct {
	articleRepo  repository.ArticleRepository
	categoryRepo repository.CategoryRepository
	tagRepo      repository.TagRepository
}

// NewArticleService 创建文章服务
func NewArticleService(
	articleRepo repository.ArticleRepository,
	categoryRepo repository.CategoryRepository,
	tagRepo repository.TagRepository,
) ArticleService {
	return &articleService{
		articleRepo:  articleRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
	}
}

// Create 创建文章
func (s *articleService) Create(req *CreateArticleRequest, authorID uint) (*models.Article, error) {
	// 验证标题
	if strings.TrimSpace(req.Title) == "" {
		return nil, ErrArticleTitleRequired
	}

	// 生成Slug（如果未提供）
	articleSlug := req.Slug
	if articleSlug == "" {
		articleSlug = slug.Make(req.Title)
	} else {
		articleSlug = slug.Make(articleSlug)
	}

	// 验证分类是否存在
	if req.CategoryID != nil {
		if _, err := s.categoryRepo.FindByID(*req.CategoryID); err != nil {
			return nil, errors.New("category not found")
		}
	}

	// 验证标签是否存在
	if len(req.TagIDs) > 0 {
		tags, err := s.tagRepo.FindByIDs(req.TagIDs)
		if err != nil {
			return nil, err
		}
		if len(tags) != len(req.TagIDs) {
			return nil, errors.New("some tags not found")
		}
	}

	// 创建文章
	article := &models.Article{
		Title:      req.Title,
		Slug:       articleSlug,
		Summary:    req.Summary,
		Content:    req.Content,
		CoverImage: req.CoverImage,
		CategoryID: req.CategoryID,
		Status:     req.Status,
		PublishAt:  req.PublishAt,
		IsTop:      req.IsTop,
		IsFeatured: req.IsFeatured,
		AuthorID:   &authorID,
	}

	// 如果状态为空，默认为草稿
	if article.Status == "" {
		article.Status = models.ArticleStatusDraft
	}

	// 如果是发布状态且没有发布时间，设置为当前时间
	if article.Status == models.ArticleStatusPublished && article.PublishAt == nil {
		now := time.Now()
		article.PublishAt = &now
	}

	if err := s.articleRepo.Create(article); err != nil {
		return nil, err
	}

	// 更新标签关联
	if len(req.TagIDs) > 0 {
		if err := s.articleRepo.UpdateTags(article.ID, req.TagIDs); err != nil {
			return nil, err
		}
	}

	// 如果已发布，更新分类和标签的文章数
	if article.Status == models.ArticleStatusPublished {
		if article.CategoryID != nil {
			s.categoryRepo.IncrementArticleCount(*article.CategoryID)
		}
		for _, tagID := range req.TagIDs {
			s.tagRepo.IncrementArticleCount(tagID)
		}
	}

	// 重新加载文章（包含关联）
	return s.articleRepo.FindByIDWithAssociations(article.ID)
}

// GetByID 获取文章详情
func (s *articleService) GetByID(id uint) (*models.Article, error) {
	return s.articleRepo.FindByIDWithAssociations(id)
}

// GetBySlug 获取文章详情（通过Slug）
func (s *articleService) GetBySlug(slug string) (*models.Article, error) {
	return s.articleRepo.FindBySlugWithAssociations(slug)
}

// Update 更新文章
func (s *articleService) Update(id uint, req *UpdateArticleRequest) (*models.Article, error) {
	// 验证标题
	if strings.TrimSpace(req.Title) == "" {
		return nil, ErrArticleTitleRequired
	}

	// 查找文章
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	oldStatus := article.Status
	oldCategoryID := article.CategoryID

	// 生成Slug（如果未提供）
	articleSlug := req.Slug
	if articleSlug == "" {
		articleSlug = slug.Make(req.Title)
	} else {
		articleSlug = slug.Make(articleSlug)
	}

	// 验证分类是否存在
	if req.CategoryID != nil {
		if _, err := s.categoryRepo.FindByID(*req.CategoryID); err != nil {
			return nil, errors.New("category not found")
		}
	}

	// 验证标签是否存在
	if len(req.TagIDs) > 0 {
		tags, err := s.tagRepo.FindByIDs(req.TagIDs)
		if err != nil {
			return nil, err
		}
		if len(tags) != len(req.TagIDs) {
			return nil, errors.New("some tags not found")
		}
	}

	// 更新字段
	article.Title = req.Title
	article.Slug = articleSlug
	article.Summary = req.Summary
	article.Content = req.Content
	article.CoverImage = req.CoverImage
	article.CategoryID = req.CategoryID
	article.Status = req.Status
	article.PublishAt = req.PublishAt
	article.IsTop = req.IsTop
	article.IsFeatured = req.IsFeatured

	// 如果状态改为发布且没有发布时间，设置为当前时间
	if article.Status == models.ArticleStatusPublished && article.PublishAt == nil {
		now := time.Now()
		article.PublishAt = &now
	}

	if err := s.articleRepo.Update(article); err != nil {
		return nil, err
	}

	// 获取旧的标签
	oldArticle, _ := s.articleRepo.FindByIDWithAssociations(id)
	var oldTagIDs []uint
	if oldArticle != nil {
		for _, tag := range oldArticle.Tags {
			oldTagIDs = append(oldTagIDs, tag.ID)
		}
	}

	// 更新标签关联
	if err := s.articleRepo.UpdateTags(article.ID, req.TagIDs); err != nil {
		return nil, err
	}

	// 更新分类和标签的文章数
	// 如果状态从草稿变为发布
	if oldStatus == models.ArticleStatusDraft && article.Status == models.ArticleStatusPublished {
		if article.CategoryID != nil {
			s.categoryRepo.IncrementArticleCount(*article.CategoryID)
		}
		for _, tagID := range req.TagIDs {
			s.tagRepo.IncrementArticleCount(tagID)
		}
	}
	// 如果状态从发布变为草稿
	if oldStatus == models.ArticleStatusPublished && article.Status == models.ArticleStatusDraft {
		if oldCategoryID != nil {
			s.categoryRepo.DecrementArticleCount(*oldCategoryID)
		}
		for _, tagID := range oldTagIDs {
			s.tagRepo.DecrementArticleCount(tagID)
		}
	}
	// 如果状态保持发布，但分类改变
	if oldStatus == models.ArticleStatusPublished && article.Status == models.ArticleStatusPublished {
		if oldCategoryID != nil && article.CategoryID != nil && *oldCategoryID != *article.CategoryID {
			s.categoryRepo.DecrementArticleCount(*oldCategoryID)
			s.categoryRepo.IncrementArticleCount(*article.CategoryID)
		}
		// 处理标签变化
		for _, oldTagID := range oldTagIDs {
			found := false
			for _, newTagID := range req.TagIDs {
				if oldTagID == newTagID {
					found = true
					break
				}
			}
			if !found {
				s.tagRepo.DecrementArticleCount(oldTagID)
			}
		}
		for _, newTagID := range req.TagIDs {
			found := false
			for _, oldTagID := range oldTagIDs {
				if newTagID == oldTagID {
					found = true
					break
				}
			}
			if !found {
				s.tagRepo.IncrementArticleCount(newTagID)
			}
		}
	}

	// 重新加载文章（包含关联）
	return s.articleRepo.FindByIDWithAssociations(article.ID)
}

// Delete 删除文章
func (s *articleService) Delete(id uint) error {
	// 查找文章
	article, err := s.articleRepo.FindByIDWithAssociations(id)
	if err != nil {
		return err
	}

	// 删除文章
	if err := s.articleRepo.Delete(id); err != nil {
		return err
	}

	// 如果是已发布状态，更新分类和标签的文章数
	if article.Status == models.ArticleStatusPublished {
		if article.CategoryID != nil {
			s.categoryRepo.DecrementArticleCount(*article.CategoryID)
		}
		for _, tag := range article.Tags {
			s.tagRepo.DecrementArticleCount(tag.ID)
		}
	}

	return nil
}

// List 获取文章列表（管理员）
func (s *articleService) List(req *ArticleListRequest) (*ArticleListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	filter := &repository.ArticleFilter{
		CategoryID: req.CategoryID,
		TagID:      req.TagID,
		Status:     req.Status,
		IsTop:      req.IsTop,
		IsFeatured: req.IsFeatured,
		Keyword:    req.Keyword,
	}

	offset := (req.Page - 1) * req.PageSize
	articles, total, err := s.articleRepo.List(filter, offset, req.PageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &ArticleListResponse{
		Items:      articles,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// ListPublished 获取已发布文章列表（公开）
func (s *articleService) ListPublished(req *ArticleListRequest) (*ArticleListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	filter := &repository.ArticleFilter{
		CategoryID: req.CategoryID,
		TagID:      req.TagID,
		IsFeatured: req.IsFeatured,
		Keyword:    req.Keyword,
	}

	offset := (req.Page - 1) * req.PageSize
	articles, total, err := s.articleRepo.ListPublished(filter, offset, req.PageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &ArticleListResponse{
		Items:      articles,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// Publish 发布文章
func (s *articleService) Publish(id uint) error {
	article, err := s.articleRepo.FindByIDWithAssociations(id)
	if err != nil {
		return err
	}

	if article.Status == models.ArticleStatusPublished {
		return errors.New("article is already published")
	}

	oldStatus := article.Status
	article.Status = models.ArticleStatusPublished
	if article.PublishAt == nil {
		now := time.Now()
		article.PublishAt = &now
	}

	if err := s.articleRepo.Update(article); err != nil {
		return err
	}

	// 如果从草稿变为发布，更新分类和标签的文章数
	if oldStatus == models.ArticleStatusDraft {
		if article.CategoryID != nil {
			s.categoryRepo.IncrementArticleCount(*article.CategoryID)
		}
		for _, tag := range article.Tags {
			s.tagRepo.IncrementArticleCount(tag.ID)
		}
	}

	return nil
}

// Unpublish 取消发布（转为草稿）
func (s *articleService) Unpublish(id uint) error {
	article, err := s.articleRepo.FindByIDWithAssociations(id)
	if err != nil {
		return err
	}

	if article.Status == models.ArticleStatusDraft {
		return errors.New("article is already draft")
	}

	article.Status = models.ArticleStatusDraft

	if err := s.articleRepo.Update(article); err != nil {
		return err
	}

	// 更新分类和标签的文章数
	if article.CategoryID != nil {
		s.categoryRepo.DecrementArticleCount(*article.CategoryID)
	}
	for _, tag := range article.Tags {
		s.tagRepo.DecrementArticleCount(tag.ID)
	}

	return nil
}

// IncrementViewCount 增加浏览量
func (s *articleService) IncrementViewCount(id uint) error {
	return s.articleRepo.IncrementViewCount(id)
}

// Search 搜索文章
func (s *articleService) Search(keyword string, page, pageSize int) (*ArticleListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	articles, total, err := s.articleRepo.Search(keyword, offset, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &ArticleListResponse{
		Items:      articles,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// ProcessScheduledPublish 处理定时发布
func (s *articleService) ProcessScheduledPublish() error {
	articles, err := s.articleRepo.GetPendingPublish()
	if err != nil {
		return err
	}

	for _, article := range articles {
		if err := s.Publish(article.ID); err != nil {
			// 记录错误但继续处理其他文章
			continue
		}
	}

	return nil
}
