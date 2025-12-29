package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iambaby/blog/internal/pkg/response"
	"github.com/iambaby/blog/internal/service"
)

// TagHandler 标签处理器
type TagHandler struct {
	tagService service.TagService
}

// NewTagHandler 创建标签处理器
func NewTagHandler(tagService service.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// Create 创建标签
// @Summary 创建标签
// @Description 创建新的文章标签
// @Tags 标签管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body service.CreateTagRequest true "标签信息"
// @Success 200 {object} response.Response{data=models.Tag} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/tags [post]
func (h *TagHandler) Create(c *gin.Context) {
	var req service.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID, _ := c.Get("userID")

	tag, err := h.tagService.Create(&req, userID.(uint))
	if err != nil {
		response.InternalServerError(c, "创建标签失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "创建成功", tag)
}

// GetByID 获取标签详情
// @Summary 获取标签详情
// @Description 根据ID获取标签详细信息
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Success 200 {object} response.Response{data=models.Tag} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 404 {object} response.Response "标签不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /tags/{id} [get]
func (h *TagHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的标签ID")
		return
	}

	tag, err := h.tagService.GetByID(uint(id))
	if err != nil {
		if err == service.ErrTagNotFound {
			response.NotFound(c, "标签不存在")
			return
		}
		response.InternalServerError(c, "获取标签失败: "+err.Error())
		return
	}

	response.Success(c, tag)
}

// GetBySlug 根据Slug获取标签
// @Summary 根据Slug获取标签
// @Description 根据Slug获取标签详细信息
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param slug path string true "标签Slug"
// @Success 200 {object} response.Response{data=models.Tag} "获取成功"
// @Failure 404 {object} response.Response "标签不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /tags/slug/{slug} [get]
func (h *TagHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.BadRequest(c, "Slug不能为空")
		return
	}

	tag, err := h.tagService.GetBySlug(slug)
	if err != nil {
		if err == service.ErrTagNotFound {
			response.NotFound(c, "标签不存在")
			return
		}
		response.InternalServerError(c, "获取标签失败: "+err.Error())
		return
	}

	response.Success(c, tag)
}

// Update 更新标签
// @Summary 更新标签
// @Description 更新标签信息
// @Tags 标签管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "标签ID"
// @Param body body service.UpdateTagRequest true "标签信息"
// @Success 200 {object} response.Response{data=models.Tag} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "标签不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/tags/{id} [put]
func (h *TagHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的标签ID")
		return
	}

	var req service.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	tag, err := h.tagService.Update(uint(id), &req)
	if err != nil {
		if err == service.ErrTagNotFound {
			response.NotFound(c, "标签不存在")
			return
		}
		response.InternalServerError(c, "更新标签失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新成功", tag)
}

// Delete 删除标签
// @Summary 删除标签
// @Description 删除指定的标签
// @Tags 标签管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "标签ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "标签不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/tags/{id} [delete]
func (h *TagHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的标签ID")
		return
	}

	err = h.tagService.Delete(uint(id))
	if err != nil {
		if err == service.ErrTagNotFound {
			response.NotFound(c, "标签不存在")
			return
		}
		response.InternalServerError(c, "删除标签失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// List 获取标签列表
// @Summary 获取标签列表
// @Description 获取所有标签列表（支持分页）
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=service.TagListResponse} "获取成功"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /tags [get]
func (h *TagHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	result, err := h.tagService.List(page, pageSize)
	if err != nil {
		response.InternalServerError(c, "获取标签列表失败: "+err.Error())
		return
	}

	response.Success(c, result)
}
