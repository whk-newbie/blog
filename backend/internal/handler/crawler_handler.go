package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/models"
	"github.com/whk-newbie/blog/internal/pkg/response"
	"github.com/whk-newbie/blog/internal/service"
)

// CrawlerHandler 爬虫任务处理器
type CrawlerHandler struct {
	crawlService service.CrawlService
}

// NewCrawlerHandler 创建爬虫任务处理器
func NewCrawlerHandler(crawlService service.CrawlService) *CrawlerHandler {
	return &CrawlerHandler{
		crawlService: crawlService,
	}
}

// RegisterTask 注册任务
// @Summary 注册爬虫任务
// @Description 注册新的爬虫任务
// @Tags 爬虫任务
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body service.RegisterTaskRequest true "任务信息"
// @Success 201 {object} response.Response "注册成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /crawler/tasks [post]
func (h *CrawlerHandler) RegisterTask(c *gin.Context) {
	var req service.RegisterTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取Token
	token, _ := c.Get("crawlerToken")
	tokenStr := token.(string)

	task, err := h.crawlService.RegisterTask(&req, tokenStr)
	if err != nil {
		if err == service.ErrCrawlTaskExists {
			response.BadRequest(c, "任务ID已存在")
			return
		}
		response.InternalServerError(c, "注册任务失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "任务注册成功", task)
}

// UpdateTaskStatus 更新任务状态
// @Summary 更新任务状态
// @Description 更新爬虫任务的状态和进度
// @Tags 爬虫任务
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task_id path string true "任务ID"
// @Param body body service.UpdateTaskStatusRequest true "状态信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "任务不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /crawler/tasks/{task_id} [put]
func (h *CrawlerHandler) UpdateTaskStatus(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		response.BadRequest(c, "任务ID不能为空")
		return
	}

	var req service.UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取Token
	token, _ := c.Get("crawlerToken")
	tokenStr := token.(string)

	task, err := h.crawlService.UpdateTaskStatus(taskID, &req, tokenStr)
	if err != nil {
		if err == service.ErrCrawlTaskNotFound {
			response.NotFound(c, "任务不存在")
			return
		}
		if err == service.ErrTaskAlreadyCompleted || err == service.ErrTaskAlreadyFailed {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalServerError(c, "更新任务状态失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "状态更新成功", task)
}

// CompleteTask 完成任务
// @Summary 完成任务
// @Description 标记爬虫任务为已完成
// @Tags 爬虫任务
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task_id path string true "任务ID"
// @Param body body service.CompleteTaskRequest true "完成信息"
// @Success 200 {object} response.Response "任务完成"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "任务不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /crawler/tasks/{task_id}/complete [put]
func (h *CrawlerHandler) CompleteTask(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		response.BadRequest(c, "任务ID不能为空")
		return
	}

	var req service.CompleteTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取Token
	token, _ := c.Get("crawlerToken")
	tokenStr := token.(string)

	task, err := h.crawlService.CompleteTask(taskID, &req, tokenStr)
	if err != nil {
		if err == service.ErrCrawlTaskNotFound {
			response.NotFound(c, "任务不存在")
			return
		}
		if err == service.ErrTaskAlreadyCompleted || err == service.ErrTaskAlreadyFailed {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalServerError(c, "完成任务失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "任务完成", task)
}

// FailTask 任务失败
// @Summary 任务失败
// @Description 标记爬虫任务为失败
// @Tags 爬虫任务
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task_id path string true "任务ID"
// @Param body body service.FailTaskRequest true "失败信息"
// @Success 200 {object} response.Response "任务状态已更新"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "任务不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /crawler/tasks/{task_id}/fail [put]
func (h *CrawlerHandler) FailTask(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		response.BadRequest(c, "任务ID不能为空")
		return
	}

	var req service.FailTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取Token
	token, _ := c.Get("crawlerToken")
	tokenStr := token.(string)

	task, err := h.crawlService.FailTask(taskID, &req, tokenStr)
	if err != nil {
		if err == service.ErrCrawlTaskNotFound {
			response.NotFound(c, "任务不存在")
			return
		}
		if err == service.ErrTaskAlreadyCompleted || err == service.ErrTaskAlreadyFailed {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalServerError(c, "更新任务状态失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "任务状态已更新", task)
}

// ListTasks 获取任务列表（管理员）
// @Summary 获取任务列表
// @Description 获取爬虫任务列表（管理员）
// @Tags 爬虫任务
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param status query string false "状态筛选" Enums(running, completed, failed)
// @Param task_id query string false "任务ID"
// @Success 200 {object} response.Response "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/crawler/tasks [get]
func (h *CrawlerHandler) ListTasks(c *gin.Context) {
	var req service.CrawlTaskListRequest

	// 解析查询参数
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			req.Page = page
		}
	}
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			req.PageSize = pageSize
		}
	}
	if statusStr := c.Query("status"); statusStr != "" {
		status := models.CrawlTaskStatus(statusStr)
		req.Status = &status
	}
	if taskID := c.Query("task_id"); taskID != "" {
		req.TaskID = taskID
	}

	resp, err := h.crawlService.ListTasks(&req)
	if err != nil {
		response.InternalServerError(c, "获取任务列表失败: "+err.Error())
		return
	}

	response.Success(c, resp)
}

// GetTaskByID 获取任务详情（管理员）
// @Summary 获取任务详情
// @Description 根据ID获取爬虫任务详细信息（管理员）
// @Tags 爬虫任务
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task_id path string true "任务ID"
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "任务不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/crawler/tasks/{task_id} [get]
func (h *CrawlerHandler) GetTaskByID(c *gin.Context) {
	taskIDStr := c.Param("task_id")
	if taskIDStr == "" {
		response.BadRequest(c, "任务ID不能为空")
		return
	}

	// 尝试解析为数字ID
	if id, err := strconv.ParseUint(taskIDStr, 10, 32); err == nil {
		task, err := h.crawlService.GetTaskByID(uint(id))
		if err != nil {
			if err == service.ErrCrawlTaskNotFound {
				response.NotFound(c, "任务不存在")
				return
			}
			response.InternalServerError(c, "获取任务失败: "+err.Error())
			return
		}
		response.Success(c, task)
		return
	}

	// 如果不是数字，则作为TaskID处理
	task, err := h.crawlService.GetTaskByTaskID(taskIDStr)
	if err != nil {
		if err == service.ErrCrawlTaskNotFound {
			response.NotFound(c, "任务不存在")
			return
		}
		response.InternalServerError(c, "获取任务失败: "+err.Error())
		return
	}

	response.Success(c, task)
}
