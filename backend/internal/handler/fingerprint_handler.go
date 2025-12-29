package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/pkg/response"
	"github.com/whk-newbie/blog/internal/service"
)

// FingerprintHandler 指纹处理器
type FingerprintHandler struct {
	fingerprintService service.FingerprintService
}

// NewFingerprintHandler 创建指纹处理器
func NewFingerprintHandler(fingerprintService service.FingerprintService) *FingerprintHandler {
	return &FingerprintHandler{
		fingerprintService: fingerprintService,
	}
}

// CollectFingerprint 收集指纹
// @Summary 收集浏览器指纹
// @Description 收集并存储浏览器指纹信息
// @Tags 指纹
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "指纹数据"
// @Success 200 {object} response.Response "收集成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /fingerprint [post]
func (h *FingerprintHandler) CollectFingerprint(c *gin.Context) {
	var fingerprintData map[string]interface{}
	if err := c.ShouldBindJSON(&fingerprintData); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 提取指纹数据和User-Agent
	userAgent := c.GetHeader("User-Agent")

	// 如果请求体中有user_agent字段，优先使用
	if ua, ok := fingerprintData["user_agent"].(string); ok && ua != "" {
		userAgent = ua
	}

	// 提取指纹数据（排除user_agent）
	fpData := make(map[string]interface{})
	for k, v := range fingerprintData {
		if k != "user_agent" {
			fpData[k] = v
		}
	}

	result, err := h.fingerprintService.CollectFingerprint(fpData, userAgent)
	if err != nil {
		response.InternalServerError(c, "收集指纹失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// ListFingerprints 获取指纹列表
// @Summary 获取指纹列表
// @Description 获取浏览器指纹列表（管理员）
// @Tags 指纹
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/fingerprints [get]
func (h *FingerprintHandler) ListFingerprints(c *gin.Context) {
	page := 1
	pageSize := 10

	// 解析分页参数
	if pageStr := c.Query("page"); pageStr != "" {
		if parsed, err := parseInt(pageStr); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if parsed, err := parseInt(pageSizeStr); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	result, err := h.fingerprintService.List(page, pageSize)
	if err != nil {
		response.InternalServerError(c, "获取指纹列表失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// GetFingerprint 获取指纹详情
// @Summary 获取指纹详情
// @Description 根据ID获取指纹详细信息（管理员）
// @Tags 指纹
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "指纹ID"
// @Success 200 {object} response.Response "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "指纹不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/fingerprints/{id} [get]
func (h *FingerprintHandler) GetFingerprint(c *gin.Context) {
	idStr := c.Param("id")
	id, err := parseInt(idStr)
	if err != nil {
		response.BadRequest(c, "无效的指纹ID")
		return
	}

	fingerprint, err := h.fingerprintService.GetByID(uint(id))
	if err != nil {
		if err.Error() == "fingerprint not found" {
			response.NotFound(c, "指纹不存在")
			return
		}
		response.InternalServerError(c, "获取指纹失败: "+err.Error())
		return
	}

	response.Success(c, fingerprint)
}

// UpdateFingerprint 更新指纹
// @Summary 更新指纹
// @Description 更新指纹信息（管理员）
// @Tags 指纹
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "指纹ID"
// @Param body body map[string]string true "更新数据"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "指纹不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/fingerprints/{id} [put]
func (h *FingerprintHandler) UpdateFingerprint(c *gin.Context) {
	idStr := c.Param("id")
	id, err := parseInt(idStr)
	if err != nil {
		response.BadRequest(c, "无效的指纹ID")
		return
	}

	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	userAgent := req["user_agent"]
	if userAgent == "" {
		response.BadRequest(c, "user_agent不能为空")
		return
	}

	if err := h.fingerprintService.Update(uint(id), userAgent); err != nil {
		if err.Error() == "fingerprint not found" {
			response.NotFound(c, "指纹不存在")
			return
		}
		response.InternalServerError(c, "更新指纹失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteFingerprint 删除指纹
// @Summary 删除指纹
// @Description 删除指纹（管理员）
// @Tags 指纹
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "指纹ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "指纹不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/fingerprints/{id} [delete]
func (h *FingerprintHandler) DeleteFingerprint(c *gin.Context) {
	idStr := c.Param("id")
	id, err := parseInt(idStr)
	if err != nil {
		response.BadRequest(c, "无效的指纹ID")
		return
	}

	if err := h.fingerprintService.Delete(uint(id)); err != nil {
		if err.Error() == "fingerprint not found" {
			response.NotFound(c, "指纹不存在")
			return
		}
		response.InternalServerError(c, "删除指纹失败: "+err.Error())
		return
	}

	response.NoContent(c, "删除成功")
}

// parseInt 解析整数
func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}
