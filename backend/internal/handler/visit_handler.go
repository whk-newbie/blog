package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/pkg/response"
	"github.com/whk-newbie/blog/internal/service"
)

// VisitHandler 访问记录处理器
type VisitHandler struct {
	visitService service.VisitService
}

// NewVisitHandler 创建访问记录处理器
func NewVisitHandler(visitService service.VisitService) *VisitHandler {
	return &VisitHandler{
		visitService: visitService,
	}
}

// RecordVisit 记录访问
// @Summary 记录访问
// @Description 记录访问行为
// @Tags 访问统计
// @Accept json
// @Produce json
// @Param body body service.RecordVisitRequest true "访问信息"
// @Success 200 {object} response.Response "记录成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /visit [post]
func (h *VisitHandler) RecordVisit(c *gin.Context) {
	var req service.RecordVisitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 如果没有提供User-Agent，从请求头获取
	if req.UserAgent == "" {
		req.UserAgent = c.GetHeader("User-Agent")
	}

	if err := h.visitService.RecordVisit(&req); err != nil {
		response.InternalServerError(c, "记录访问失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}
