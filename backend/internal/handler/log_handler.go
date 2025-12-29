package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/pkg/response"
	"github.com/whk-newbie/blog/internal/service"
)

// LogHandler 日志处理器
type LogHandler struct {
	logService service.LogService
}

// NewLogHandler 创建日志处理器
func NewLogHandler(logService service.LogService) *LogHandler {
	return &LogHandler{
		logService: logService,
	}
}

// GetLogs 获取日志列表
// @Summary 获取日志列表
// @Description 获取系统日志列表（支持分页和筛选）
// @Tags 日志管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param level query string false "日志级别"
// @Param source query string false "日志来源"
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Success 200 {object} response.Response "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/logs [get]
func (h *LogHandler) GetLogs(c *gin.Context) {
	req := &service.LogQueryRequest{
		Page:     1,
		PageSize: 20,
	}

	// 解析分页参数
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			req.Page = page
		}
	}
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			req.PageSize = pageSize
		}
	}

	// 解析筛选参数
	if level := c.Query("level"); level != "" {
		req.Level = level
	}
	if source := c.Query("source"); source != "" {
		req.Source = source
	}

	// 解析日期参数
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			req.StartDate = startDate
		}
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			// 设置为当天的23:59:59
			endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())
			req.EndDate = endDate
		}
	}

	logs, err := h.logService.GetLogs(req)
	if err != nil {
		response.InternalServerError(c, "获取日志列表失败: "+err.Error())
		return
	}

	response.Success(c, logs)
}

// GetLogByID 获取日志详情
// @Summary 获取日志详情
// @Description 根据ID获取日志详情
// @Tags 日志管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "日志ID"
// @Success 200 {object} response.Response "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "日志不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/logs/:id [get]
func (h *LogHandler) GetLogByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的日志ID")
		return
	}

	log, err := h.logService.GetLogByID(uint(id))
	if err != nil {
		if err == service.ErrLogNotFound {
			response.NotFound(c, "日志不存在")
			return
		}
		response.InternalServerError(c, "获取日志失败: "+err.Error())
		return
	}

	response.Success(c, log)
}

// CleanupLogs 清理旧日志
// @Summary 清理旧日志
// @Description 手动清理旧日志
// @Tags 日志管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object true "清理请求" example({"retention_days": 90})
// @Success 200 {object} response.Response "清理完成"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/logs/cleanup [post]
func (h *LogHandler) CleanupLogs(c *gin.Context) {
	var req struct {
		RetentionDays int `json:"retention_days"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// 如果没有提供参数，使用默认值
		req.RetentionDays = 90
	}

	result, err := h.logService.CleanupOldLogs(req.RetentionDays)
	if err != nil {
		response.InternalServerError(c, "清理日志失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "清理完成", result)
}
