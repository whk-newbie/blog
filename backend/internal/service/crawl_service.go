package service

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/whk-newbie/blog/internal/models"
	"github.com/whk-newbie/blog/internal/repository"
	"github.com/whk-newbie/blog/internal/websocket"
	"gorm.io/datatypes"
)

var (
	ErrCrawlTaskNotFound    = repository.ErrCrawlTaskNotFound
	ErrCrawlTaskExists      = repository.ErrCrawlTaskExists
	ErrInvalidTaskStatus    = errors.New("invalid task status")
	ErrInvalidProgress      = errors.New("progress must be between 0 and 100")
	ErrTaskAlreadyCompleted = errors.New("task already completed")
	ErrTaskAlreadyFailed    = errors.New("task already failed")
)

// CrawlService 爬虫任务服务接口
type CrawlService interface {
	// 注册任务
	RegisterTask(req *RegisterTaskRequest, token string) (*models.CrawlTask, error)
	// 更新任务状态
	UpdateTaskStatus(taskID string, req *UpdateTaskStatusRequest, token string) (*models.CrawlTask, error)
	// 完成任务
	CompleteTask(taskID string, req *CompleteTaskRequest, token string) (*models.CrawlTask, error)
	// 任务失败
	FailTask(taskID string, req *FailTaskRequest, token string) (*models.CrawlTask, error)
	// 获取任务列表（管理员）
	ListTasks(req *CrawlTaskListRequest) (*CrawlTaskListResponse, error)
	// 获取任务详情（管理员）
	GetTaskByID(id uint) (*models.CrawlTask, error)
	// 获取任务详情（通过TaskID）
	GetTaskByTaskID(taskID string) (*models.CrawlTask, error)
}

// RegisterTaskRequest 注册任务请求
type RegisterTaskRequest struct {
	TaskID   string                 `json:"task_id" binding:"required"`
	TaskName string                 `json:"task_name" binding:"required"`
	Metadata map[string]interface{} `json:"metadata"`
}

// UpdateTaskStatusRequest 更新任务状态请求
type UpdateTaskStatusRequest struct {
	Status   models.CrawlTaskStatus `json:"status"`
	Progress int                    `json:"progress"`
	Message  string                 `json:"message"`
}

// CompleteTaskRequest 完成任务请求
type CompleteTaskRequest struct {
	Message  string                 `json:"message"`
	Metadata map[string]interface{} `json:"metadata"`
}

// FailTaskRequest 任务失败请求
type FailTaskRequest struct {
	Message  string                 `json:"message"`
	Error    string                 `json:"error"`
	Metadata map[string]interface{} `json:"metadata"`
}

// CrawlTaskListRequest 任务列表请求
type CrawlTaskListRequest struct {
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
	Status   *models.CrawlTaskStatus `json:"status"`
	TaskID   string                  `json:"task_id"`
}

// CrawlTaskListResponse 任务列表响应
type CrawlTaskListResponse struct {
	Items      []models.CrawlTask `json:"items"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

// crawlService 爬虫任务服务实现
type crawlService struct {
	taskRepo repository.CrawlTaskRepository
	hub      *websocket.Hub
}

// NewCrawlService 创建爬虫任务服务
func NewCrawlService(taskRepo repository.CrawlTaskRepository, hub *websocket.Hub) CrawlService {
	return &crawlService{
		taskRepo: taskRepo,
		hub:      hub,
	}
}

// RegisterTask 注册任务
func (s *crawlService) RegisterTask(req *RegisterTaskRequest, token string) (*models.CrawlTask, error) {
	// 检查TaskID是否已存在
	exists, err := s.taskRepo.ExistsByTaskID(req.TaskID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrCrawlTaskExists
	}

	now := time.Now()
	task := &models.CrawlTask{
		TaskID:         req.TaskID,
		TaskName:       req.TaskName,
		Status:         models.CrawlTaskStatusRunning,
		Progress:       0,
		Message:        "任务已注册",
		StartTime:      now,
		CreatedByToken: token,
	}

	// 处理Metadata
	if req.Metadata != nil {
		// 使用datatypes.JSON来存储
		metadataJSON, err := json.Marshal(req.Metadata)
		if err != nil {
			return nil, err
		}
		task.Metadata = datatypes.JSON(metadataJSON)
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, err
	}

	// 广播任务创建
	if s.hub != nil {
		s.hub.BroadcastTaskUpdate(task)
	}

	return task, nil
}

// UpdateTaskStatus 更新任务状态
func (s *crawlService) UpdateTaskStatus(taskID string, req *UpdateTaskStatusRequest, token string) (*models.CrawlTask, error) {
	// 获取任务
	task, err := s.taskRepo.FindByTaskID(taskID)
	if err != nil {
		return nil, err
	}

	// 验证Token权限
	if task.CreatedByToken != token {
		return nil, errors.New("无权访问此任务")
	}

	// 检查任务状态
	if task.Status == models.CrawlTaskStatusCompleted {
		return nil, ErrTaskAlreadyCompleted
	}
	if task.Status == models.CrawlTaskStatusFailed {
		return nil, ErrTaskAlreadyFailed
	}

	// 验证状态
	if req.Status != models.CrawlTaskStatusRunning &&
		req.Status != models.CrawlTaskStatusCompleted &&
		req.Status != models.CrawlTaskStatusFailed {
		return nil, ErrInvalidTaskStatus
	}

	// 验证进度
	if req.Progress < 0 || req.Progress > 100 {
		return nil, ErrInvalidProgress
	}

	// 更新任务
	task.Status = req.Status
	task.Progress = req.Progress
	if req.Message != "" {
		task.Message = req.Message
	}

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	// 广播任务更新
	if s.hub != nil {
		s.hub.BroadcastTaskUpdate(task)
	}

	return task, nil
}

// CompleteTask 完成任务
func (s *crawlService) CompleteTask(taskID string, req *CompleteTaskRequest, token string) (*models.CrawlTask, error) {
	// 获取任务
	task, err := s.taskRepo.FindByTaskID(taskID)
	if err != nil {
		return nil, err
	}

	// 验证Token权限
	if task.CreatedByToken != token {
		return nil, errors.New("无权访问此任务")
	}

	// 检查任务状态
	if task.Status == models.CrawlTaskStatusCompleted {
		return nil, ErrTaskAlreadyCompleted
	}
	if task.Status == models.CrawlTaskStatusFailed {
		return nil, ErrTaskAlreadyFailed
	}

	// 更新任务
	now := time.Now()
	task.Status = models.CrawlTaskStatusCompleted
	task.Progress = 100
	task.EndTime = &now
	if req.Message != "" {
		task.Message = req.Message
	}

	// 计算运行时长
	duration := int(now.Sub(task.StartTime).Seconds())
	task.Duration = &duration

	// 更新Metadata
	if req.Metadata != nil {
		metadataJSON, err := json.Marshal(req.Metadata)
		if err != nil {
			return nil, err
		}
		task.Metadata = datatypes.JSON(metadataJSON)
	}

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	// 广播任务更新
	if s.hub != nil {
		s.hub.BroadcastTaskUpdate(task)
	}

	return task, nil
}

// FailTask 任务失败
func (s *crawlService) FailTask(taskID string, req *FailTaskRequest, token string) (*models.CrawlTask, error) {
	// 获取任务
	task, err := s.taskRepo.FindByTaskID(taskID)
	if err != nil {
		return nil, err
	}

	// 验证Token权限
	if task.CreatedByToken != token {
		return nil, errors.New("无权访问此任务")
	}

	// 检查任务状态
	if task.Status == models.CrawlTaskStatusCompleted {
		return nil, ErrTaskAlreadyCompleted
	}
	if task.Status == models.CrawlTaskStatusFailed {
		return nil, ErrTaskAlreadyFailed
	}

	// 更新任务
	now := time.Now()
	task.Status = models.CrawlTaskStatusFailed
	task.EndTime = &now
	if req.Message != "" {
		task.Message = req.Message
	}

	// 计算运行时长
	duration := int(now.Sub(task.StartTime).Seconds())
	task.Duration = &duration

	// 更新Metadata（包含错误信息）
	metadata := make(map[string]interface{})
	if task.Metadata != nil {
		json.Unmarshal(task.Metadata, &metadata)
	}
	if req.Error != "" {
		metadata["error"] = req.Error
	}
	if req.Metadata != nil {
		for k, v := range req.Metadata {
			metadata[k] = v
		}
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}
	task.Metadata = datatypes.JSON(metadataJSON)

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	// 广播任务更新
	if s.hub != nil {
		s.hub.BroadcastTaskUpdate(task)
	}

	return task, nil
}

// ListTasks 获取任务列表（管理员）
func (s *crawlService) ListTasks(req *CrawlTaskListRequest) (*CrawlTaskListResponse, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// 构建筛选条件
	filter := &repository.CrawlTaskFilter{}
	if req.Status != nil {
		filter.Status = req.Status
	}
	if req.TaskID != "" {
		filter.TaskID = req.TaskID
	}

	// 计算偏移量
	offset := (req.Page - 1) * req.PageSize

	// 获取任务列表
	tasks, total, err := s.taskRepo.List(filter, offset, req.PageSize)
	if err != nil {
		return nil, err
	}

	// 计算总页数
	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	return &CrawlTaskListResponse{
		Items:      tasks,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetTaskByID 获取任务详情（管理员）
func (s *crawlService) GetTaskByID(id uint) (*models.CrawlTask, error) {
	return s.taskRepo.FindByID(id)
}

// GetTaskByTaskID 获取任务详情（通过TaskID）
func (s *crawlService) GetTaskByTaskID(taskID string) (*models.CrawlTask, error) {
	return s.taskRepo.FindByTaskID(taskID)
}
