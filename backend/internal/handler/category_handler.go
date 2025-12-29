package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/pkg/response"
	"github.com/whk-newbie/blog/internal/service"
)

// CategoryHandler 分类处理器
type CategoryHandler struct {
	categoryService service.CategoryService
}

// NewCategoryHandler 创建分类处理器
func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// Create 创建分类
// @Summary 创建分类
// @Description 创建新的文章分类
// @Tags 分类管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body service.CreateCategoryRequest true "分类信息"
// @Success 200 {object} response.Response{data=models.Category} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var req service.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID, _ := c.Get("userID")

	category, err := h.categoryService.Create(&req, userID.(uint))
	if err != nil {
		response.InternalServerError(c, "创建分类失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "创建成功", category)
}

// GetByID 获取分类详情
// @Summary 获取分类详情
// @Description 根据ID获取分类详细信息
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} response.Response{data=models.Category} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 404 {object} response.Response "分类不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的分类ID")
		return
	}

	category, err := h.categoryService.GetByID(uint(id))
	if err != nil {
		if err == service.ErrCategoryNotFound {
			response.NotFound(c, "分类不存在")
			return
		}
		response.InternalServerError(c, "获取分类失败: "+err.Error())
		return
	}

	response.Success(c, category)
}

// GetBySlug 根据Slug获取分类
// @Summary 根据Slug获取分类
// @Description 根据Slug获取分类详细信息
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param slug path string true "分类Slug"
// @Success 200 {object} response.Response{data=models.Category} "获取成功"
// @Failure 404 {object} response.Response "分类不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /categories/slug/{slug} [get]
func (h *CategoryHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.BadRequest(c, "Slug不能为空")
		return
	}

	category, err := h.categoryService.GetBySlug(slug)
	if err != nil {
		if err == service.ErrCategoryNotFound {
			response.NotFound(c, "分类不存在")
			return
		}
		response.InternalServerError(c, "获取分类失败: "+err.Error())
		return
	}

	response.Success(c, category)
}

// Update 更新分类
// @Summary 更新分类
// @Description 更新分类信息
// @Tags 分类管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "分类ID"
// @Param body body service.UpdateCategoryRequest true "分类信息"
// @Success 200 {object} response.Response{data=models.Category} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "分类不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的分类ID")
		return
	}

	var req service.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	category, err := h.categoryService.Update(uint(id), &req)
	if err != nil {
		if err == service.ErrCategoryNotFound {
			response.NotFound(c, "分类不存在")
			return
		}
		response.InternalServerError(c, "更新分类失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新成功", category)
}

// Delete 删除分类
// @Summary 删除分类
// @Description 删除指定的分类
// @Tags 分类管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "分类ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "分类不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的分类ID")
		return
	}

	err = h.categoryService.Delete(uint(id))
	if err != nil {
		if err == service.ErrCategoryNotFound {
			response.NotFound(c, "分类不存在")
			return
		}
		response.InternalServerError(c, "删除分类失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// List 获取分类列表
// @Summary 获取分类列表
// @Description 获取所有分类列表（支持分页）
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=service.CategoryListResponse} "获取成功"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	result, err := h.categoryService.List(page, pageSize)
	if err != nil {
		response.InternalServerError(c, "获取分类列表失败: "+err.Error())
		return
	}

	response.Success(c, result)
}
