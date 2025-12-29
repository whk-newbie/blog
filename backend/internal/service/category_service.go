package service

import (
	"errors"
	"strings"

	"github.com/gosimple/slug"
	"github.com/whk-newbie/blog/internal/models"
	"github.com/whk-newbie/blog/internal/repository"
)

var (
	ErrCategoryNameRequired = errors.New("category name is required")
	ErrCategoryNotFound     = repository.ErrCategoryNotFound
)

// CategoryService 分类服务接口
type CategoryService interface {
	// 创建分类
	Create(req *CreateCategoryRequest, createdBy uint) (*models.Category, error)
	// 获取分类详情
	GetByID(id uint) (*models.Category, error)
	// 获取分类详情（通过Slug）
	GetBySlug(slug string) (*models.Category, error)
	// 更新分类
	Update(id uint, req *UpdateCategoryRequest) (*models.Category, error)
	// 删除分类
	Delete(id uint) error
	// 获取分类列表
	List(page, pageSize int) (*CategoryListResponse, error)
}

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// CategoryListResponse 分类列表响应
type CategoryListResponse struct {
	Items      []models.Category `json:"items"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

// categoryService 分类服务实现
type categoryService struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryService 创建分类服务
func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

// Create 创建分类
func (s *categoryService) Create(req *CreateCategoryRequest, createdBy uint) (*models.Category, error) {
	// 验证名称
	if strings.TrimSpace(req.Name) == "" {
		return nil, ErrCategoryNameRequired
	}

	// 生成Slug（如果未提供）
	categorySlug := req.Slug
	if categorySlug == "" {
		categorySlug = slug.Make(req.Name)
	} else {
		categorySlug = slug.Make(categorySlug)
	}

	category := &models.Category{
		Name:        req.Name,
		Slug:        categorySlug,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		CreatedBy:   &createdBy,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

// GetByID 获取分类详情
func (s *categoryService) GetByID(id uint) (*models.Category, error) {
	return s.categoryRepo.FindByID(id)
}

// GetBySlug 获取分类详情（通过Slug）
func (s *categoryService) GetBySlug(slug string) (*models.Category, error) {
	return s.categoryRepo.FindBySlug(slug)
}

// Update 更新分类
func (s *categoryService) Update(id uint, req *UpdateCategoryRequest) (*models.Category, error) {
	// 验证名称
	if strings.TrimSpace(req.Name) == "" {
		return nil, ErrCategoryNameRequired
	}

	// 查找分类
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 生成Slug（如果未提供）
	categorySlug := req.Slug
	if categorySlug == "" {
		categorySlug = slug.Make(req.Name)
	} else {
		categorySlug = slug.Make(categorySlug)
	}

	// 更新字段
	category.Name = req.Name
	category.Slug = categorySlug
	category.Description = req.Description
	category.SortOrder = req.SortOrder

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

// Delete 删除分类
func (s *categoryService) Delete(id uint) error {
	// 检查分类是否存在
	_, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return err
	}

	return s.categoryRepo.Delete(id)
}

// List 获取分类列表
func (s *categoryService) List(page, pageSize int) (*CategoryListResponse, error) {
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
	categories, total, err := s.categoryRepo.List(offset, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &CategoryListResponse{
		Items:      categories,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
