package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/models"
	"github.com/whk-newbie/blog/internal/pkg/response"
	"github.com/whk-newbie/blog/internal/service"
)

// ArticleHandler 文章处理器
type ArticleHandler struct {
	articleService service.ArticleService
}

// NewArticleHandler 创建文章处理器
func NewArticleHandler(articleService service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
	}
}

// Create 创建文章
// @Summary 创建文章
// @Description 创建新的文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body service.CreateArticleRequest true "文章信息"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/articles [post]
func (h *ArticleHandler) Create(c *gin.Context) {
	var req service.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID, _ := c.Get("userID")

	article, err := h.articleService.Create(&req, userID.(uint))
	if err != nil {
		response.InternalServerError(c, "创建文章失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "创建成功", article)
}

// GetByID 获取文章详情
// @Summary 获取文章详情
// @Description 根据ID获取文章详细信息（公开接口）
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 404 {object} response.Response "文章不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /articles/{id} [get]
func (h *ArticleHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文章ID")
		return
	}

	article, err := h.articleService.GetByID(uint(id))
	if err != nil {
		if err == service.ErrArticleNotFound {
			response.NotFound(c, "文章不存在")
			return
		}
		response.InternalServerError(c, "获取文章失败: "+err.Error())
		return
	}

	// 如果是已发布的文章，增加浏览量
	if article.Status == models.ArticleStatusPublished {
		go h.articleService.IncrementViewCount(article.ID)
	}

	response.Success(c, article)
}

// GetBySlug 根据Slug获取文章
// @Summary 根据Slug获取文章
// @Description 根据Slug获取文章详细信息（公开接口）
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param slug path string true "文章Slug"
// @Success 200 {object} response.Response "获取成功"
// @Failure 404 {object} response.Response "文章不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /articles/slug/{slug} [get]
func (h *ArticleHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.BadRequest(c, "Slug不能为空")
		return
	}

	article, err := h.articleService.GetBySlug(slug)
	if err != nil {
		if err == service.ErrArticleNotFound {
			response.NotFound(c, "文章不存在")
			return
		}
		response.InternalServerError(c, "获取文章失败: "+err.Error())
		return
	}

	// 如果是已发布的文章，增加浏览量
	if article.Status == models.ArticleStatusPublished {
		go h.articleService.IncrementViewCount(article.ID)
	}

	response.Success(c, article)
}

// Update 更新文章
// @Summary 更新文章
// @Description 更新文章信息
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Param body body service.UpdateArticleRequest true "文章信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "文章不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/articles/{id} [put]
func (h *ArticleHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文章ID")
		return
	}

	var req service.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	article, err := h.articleService.Update(uint(id), &req)
	if err != nil {
		if err == service.ErrArticleNotFound {
			response.NotFound(c, "文章不存在")
			return
		}
		response.InternalServerError(c, "更新文章失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新成功", article)
}

// Delete 删除文章
// @Summary 删除文章
// @Description 删除指定的文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "文章不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/articles/{id} [delete]
func (h *ArticleHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文章ID")
		return
	}

	err = h.articleService.Delete(uint(id))
	if err != nil {
		if err == service.ErrArticleNotFound {
			response.NotFound(c, "文章不存在")
			return
		}
		response.InternalServerError(c, "删除文章失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// List 获取文章列表（管理员）
// @Summary 获取文章列表（管理员）
// @Description 获取所有文章列表（支持分页和筛选）
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param category_id query int false "分类ID"
// @Param tag_id query int false "标签ID"
// @Param status query string false "状态" Enums(draft, published)
// @Param is_top query bool false "是否置顶"
// @Param is_featured query bool false "是否推荐"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} response.Response "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/articles [get]
func (h *ArticleHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	req := &service.ArticleListRequest{
		Page:     page,
		PageSize: pageSize,
		Keyword:  c.Query("keyword"),
	}

	// 解析分类ID
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			id := uint(categoryID)
			req.CategoryID = &id
		}
	}

	// 解析标签ID
	if tagIDStr := c.Query("tag_id"); tagIDStr != "" {
		if tagID, err := strconv.ParseUint(tagIDStr, 10, 32); err == nil {
			id := uint(tagID)
			req.TagID = &id
		}
	}

	// 解析状态
	if statusStr := c.Query("status"); statusStr != "" {
		status := models.ArticleStatus(statusStr)
		req.Status = &status
	}

	// 解析是否置顶
	if isTopStr := c.Query("is_top"); isTopStr != "" {
		if isTop, err := strconv.ParseBool(isTopStr); err == nil {
			req.IsTop = &isTop
		}
	}

	// 解析是否推荐
	if isFeaturedStr := c.Query("is_featured"); isFeaturedStr != "" {
		if isFeatured, err := strconv.ParseBool(isFeaturedStr); err == nil {
			req.IsFeatured = &isFeatured
		}
	}

	result, err := h.articleService.List(req)
	if err != nil {
		response.InternalServerError(c, "获取文章列表失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// ListPublished 获取已发布文章列表（公开）
// @Summary 获取已发布文章列表（公开）
// @Description 获取所有已发布的文章列表（支持分页和筛选）
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param category_id query int false "分类ID"
// @Param tag_id query int false "标签ID"
// @Param is_featured query bool false "是否推荐"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} response.Response "获取成功"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /articles [get]
func (h *ArticleHandler) ListPublished(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	req := &service.ArticleListRequest{
		Page:     page,
		PageSize: pageSize,
		Keyword:  c.Query("keyword"),
	}

	// 解析分类ID
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			id := uint(categoryID)
			req.CategoryID = &id
		}
	}

	// 解析标签ID
	if tagIDStr := c.Query("tag_id"); tagIDStr != "" {
		if tagID, err := strconv.ParseUint(tagIDStr, 10, 32); err == nil {
			id := uint(tagID)
			req.TagID = &id
		}
	}

	// 解析是否推荐
	if isFeaturedStr := c.Query("is_featured"); isFeaturedStr != "" {
		if isFeatured, err := strconv.ParseBool(isFeaturedStr); err == nil {
			req.IsFeatured = &isFeatured
		}
	}

	result, err := h.articleService.ListPublished(req)
	if err != nil {
		response.InternalServerError(c, "获取文章列表失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// Publish 发布文章
// @Summary 发布文章
// @Description 将草稿文章发布
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response "发布成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "文章不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/articles/{id}/publish [post]
func (h *ArticleHandler) Publish(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文章ID")
		return
	}

	err = h.articleService.Publish(uint(id))
	if err != nil {
		if err == service.ErrArticleNotFound {
			response.NotFound(c, "文章不存在")
			return
		}
		response.InternalServerError(c, "发布文章失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "发布成功", nil)
}

// Unpublish 取消发布
// @Summary 取消发布
// @Description 将已发布文章转为草稿
// @Tags 文章管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response "操作成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "文章不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/articles/{id}/unpublish [post]
func (h *ArticleHandler) Unpublish(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文章ID")
		return
	}

	err = h.articleService.Unpublish(uint(id))
	if err != nil {
		if err == service.ErrArticleNotFound {
			response.NotFound(c, "文章不存在")
			return
		}
		response.InternalServerError(c, "操作失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "操作成功", nil)
}

// Search 搜索文章
// @Summary 搜索文章
// @Description 全文搜索已发布的文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param keyword query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response "搜索成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /articles/search [get]
func (h *ArticleHandler) Search(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.BadRequest(c, "搜索关键词不能为空")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	result, err := h.articleService.Search(keyword, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "搜索失败: "+err.Error())
		return
	}

	response.Success(c, result)
}
