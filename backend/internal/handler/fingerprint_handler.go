package handler

import (
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
