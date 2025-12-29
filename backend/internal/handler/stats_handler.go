package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/pkg/response"
	"github.com/whk-newbie/blog/internal/service"
)

// StatsHandler 统计处理器
type StatsHandler struct {
	statsService service.StatsService
}

// NewStatsHandler 创建统计处理器
func NewStatsHandler(statsService service.StatsService) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
	}
}

// GetDashboardStats 获取仪表盘统计数据
// @Summary 获取仪表盘统计数据
// @Description 获取仪表盘统计数据（文章数量、分类数量、标签数量、最近文章）
// @Tags 统计
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/stats/dashboard [get]
func (h *StatsHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.statsService.GetDashboardStats()
	if err != nil {
		response.InternalServerError(c, "获取统计数据失败: "+err.Error())
		return
	}

	response.Success(c, stats)
}

// GetVisitStats 获取访问统计
// @Summary 获取访问统计
// @Description 获取访问统计数据（PV/UV、平均停留时间等）
// @Tags 统计
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Param type query string false "统计类型 (daily/weekly/monthly)" Enums(daily, weekly, monthly) default(daily)
// @Success 200 {object} response.Response "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/stats/visits [get]
func (h *StatsHandler) GetVisitStats(c *gin.Context) {
	req := &service.VisitStatsRequest{}

	// 解析开始日期
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			req.StartDate = startDate
		}
	}

	// 解析结束日期
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			req.EndDate = endDate
		}
	}

	// 解析统计类型
	if typeStr := c.Query("type"); typeStr != "" {
		req.Type = typeStr
	}

	stats, err := h.statsService.GetVisitStats(req)
	if err != nil {
		response.InternalServerError(c, "获取访问统计失败: "+err.Error())
		return
	}

	response.Success(c, stats)
}

// GetPopularArticles 获取热门文章
// @Summary 获取热门文章
// @Description 获取热门文章统计
// @Tags 统计
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "返回数量" default(10)
// @Param days query int false "统计天数" default(7)
// @Success 200 {object} response.Response "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/stats/popular-articles [get]
func (h *StatsHandler) GetPopularArticles(c *gin.Context) {
	limit := 10
	days := 7

	// 解析limit参数
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// 解析days参数
	if daysStr := c.Query("days"); daysStr != "" {
		if parsedDays, err := strconv.Atoi(daysStr); err == nil && parsedDays > 0 {
			days = parsedDays
		}
	}

	articles, err := h.statsService.GetPopularArticles(limit, days)
	if err != nil {
		response.InternalServerError(c, "获取热门文章失败: "+err.Error())
		return
	}

	response.Success(c, articles)
}

// GetReferrerStats 获取访问来源统计
// @Summary 获取访问来源统计
// @Description 获取访问来源统计数据
// @Tags 统计
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Success 200 {object} response.Response "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/stats/referrers [get]
func (h *StatsHandler) GetReferrerStats(c *gin.Context) {
	var startDate, endDate time.Time

	// 解析开始日期
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		parsed, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = parsed
		}
	}

	// 解析结束日期
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		parsed, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			endDate = parsed
		}
	}

	stats, err := h.statsService.GetReferrerStats(startDate, endDate)
	if err != nil {
		response.InternalServerError(c, "获取访问来源统计失败: "+err.Error())
		return
	}

	response.Success(c, stats)
}
