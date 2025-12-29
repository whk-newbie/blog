package service

import (
	"errors"
	"strings"

	"github.com/gosimple/slug"
	"github.com/iambaby/blog/internal/models"
	"github.com/iambaby/blog/internal/repository"
)

var (
	ErrTagNameRequired = errors.New("tag name is required")
	ErrTagNotFound     = repository.ErrTagNotFound
)

// TagService 标签服务接口
type TagService interface {
	// 创建标签
	Create(req *CreateTagRequest, createdBy uint) (*models.Tag, error)
	// 获取标签详情
	GetByID(id uint) (*models.Tag, error)
	// 获取标签详情（通过Slug）
	GetBySlug(slug string) (*models.Tag, error)
	// 更新标签
	Update(id uint, req *UpdateTagRequest) (*models.Tag, error)
	// 删除标签
	Delete(id uint) error
	// 获取标签列表
	List(page, pageSize int) (*TagListResponse, error)
}

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug"`
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug"`
}

// TagListResponse 标签列表响应
type TagListResponse struct {
	Items      []models.Tag `json:"items"`
	Total      int64        `json:"total"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalPages int          `json:"total_pages"`
}

// tagService 标签服务实现
type tagService struct {
	tagRepo repository.TagRepository
}

// NewTagService 创建标签服务
func NewTagService(tagRepo repository.TagRepository) TagService {
	return &tagService{
		tagRepo: tagRepo,
	}
}

// Create 创建标签
func (s *tagService) Create(req *CreateTagRequest, createdBy uint) (*models.Tag, error) {
	// 验证名称
	if strings.TrimSpace(req.Name) == "" {
		return nil, ErrTagNameRequired
	}

	// 生成Slug（如果未提供）
	tagSlug := req.Slug
	if tagSlug == "" {
		tagSlug = slug.Make(req.Name)
	} else {
		tagSlug = slug.Make(tagSlug)
	}

	tag := &models.Tag{
		Name:      req.Name,
		Slug:      tagSlug,
		CreatedBy: &createdBy,
	}

	if err := s.tagRepo.Create(tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// GetByID 获取标签详情
func (s *tagService) GetByID(id uint) (*models.Tag, error) {
	return s.tagRepo.FindByID(id)
}

// GetBySlug 获取标签详情（通过Slug）
func (s *tagService) GetBySlug(slug string) (*models.Tag, error) {
	return s.tagRepo.FindBySlug(slug)
}

// Update 更新标签
func (s *tagService) Update(id uint, req *UpdateTagRequest) (*models.Tag, error) {
	// 验证名称
	if strings.TrimSpace(req.Name) == "" {
		return nil, ErrTagNameRequired
	}

	// 查找标签
	tag, err := s.tagRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 生成Slug（如果未提供）
	tagSlug := req.Slug
	if tagSlug == "" {
		tagSlug = slug.Make(req.Name)
	} else {
		tagSlug = slug.Make(tagSlug)
	}

	// 更新字段
	tag.Name = req.Name
	tag.Slug = tagSlug

	if err := s.tagRepo.Update(tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// Delete 删除标签
func (s *tagService) Delete(id uint) error {
	// 检查标签是否存在
	_, err := s.tagRepo.FindByID(id)
	if err != nil {
		return err
	}

	return s.tagRepo.Delete(id)
}

// List 获取标签列表
func (s *tagService) List(page, pageSize int) (*TagListResponse, error) {
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
	tags, total, err := s.tagRepo.List(offset, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &TagListResponse{
		Items:      tags,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
