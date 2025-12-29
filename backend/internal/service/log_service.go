package service

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/whk-newbie/blog/internal/models"
	"github.com/whk-newbie/blog/internal/repository"
)

var (
	ErrLogNotFound = errors.New("log not found")
)

// LogService 日志服务接口
type LogService interface {
	// 记录日志
	Log(level models.LogLevel, message, source string, context map[string]interface{}, userID *uint, ipAddress string) error
	// 查询日志列表
	GetLogs(req *LogQueryRequest) (*LogListResponse, error)
	// 获取日志详情
	GetLogByID(id uint) (*LogResponse, error)
	// 清理旧日志
	CleanupOldLogs(retentionDays int) (*CleanupResponse, error)
}

// LogQueryRequest 日志查询请求
type LogQueryRequest struct {
	Page      int       `json:"page"`
	PageSize  int       `json:"page_size"`
	Level     string    `json:"level"`
	Source    string    `json:"source"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// LogListResponse 日志列表响应
type LogListResponse struct {
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	Total      int64          `json:"total"`
	TotalPages int            `json:"total_pages"`
	Items      []*LogResponse `json:"items"`
}

// LogResponse 日志响应
type LogResponse struct {
	ID        uint                   `json:"id"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Context   map[string]interface{} `json:"context"`
	Source    string                 `json:"source"`
	UserID    *uint                  `json:"user_id"`
	IPAddress string                 `json:"ip_address"`
	CreatedAt string                 `json:"created_at"`
}

// CleanupResponse 清理响应
type CleanupResponse struct {
	DeletedCount int64 `json:"deleted_count"`
}

// logService 日志服务实现
type logService struct {
	logRepo repository.LogRepository
}

// NewLogService 创建日志服务
func NewLogService(logRepo repository.LogRepository) LogService {
	return &logService{
		logRepo: logRepo,
	}
}

// Log 记录日志
func (s *logService) Log(level models.LogLevel, message, source string, context map[string]interface{}, userID *uint, ipAddress string) error {
	var contextJSON []byte
	var err error

	if context != nil {
		contextJSON, err = json.Marshal(context)
		if err != nil {
			return err
		}
	}

	log := &models.SystemLog{
		Level:     level,
		Message:   message,
		Context:   contextJSON,
		Source:    source,
		UserID:    userID,
		IPAddress: ipAddress,
	}

	return s.logRepo.Create(log)
}

// GetLogs 查询日志列表
func (s *logService) GetLogs(req *LogQueryRequest) (*LogListResponse, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	var startDate, endDate *time.Time
	if !req.StartDate.IsZero() {
		startDate = &req.StartDate
	}
	if !req.EndDate.IsZero() {
		endDate = &req.EndDate
	}

	logs, total, err := s.logRepo.FindLogs(
		req.Page,
		req.PageSize,
		req.Level,
		req.Source,
		startDate,
		endDate,
	)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	items := make([]*LogResponse, 0, len(logs))
	for _, log := range logs {
		items = append(items, s.toLogResponse(log))
	}

	// 计算总页数
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize != 0 {
		totalPages++
	}

	return &LogListResponse{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      total,
		TotalPages: totalPages,
		Items:      items,
	}, nil
}

// GetLogByID 获取日志详情
func (s *logService) GetLogByID(id uint) (*LogResponse, error) {
	log, err := s.logRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrLogNotFound) {
			return nil, ErrLogNotFound
		}
		return nil, err
	}

	return s.toLogResponse(log), nil
}

// CleanupOldLogs 清理旧日志
func (s *logService) CleanupOldLogs(retentionDays int) (*CleanupResponse, error) {
	if retentionDays <= 0 {
		retentionDays = 90 // 默认90天
	}

	deletedCount, err := s.logRepo.CleanupOldLogs(retentionDays)
	if err != nil {
		return nil, err
	}

	return &CleanupResponse{
		DeletedCount: deletedCount,
	}, nil
}

// toLogResponse 转换为响应格式
func (s *logService) toLogResponse(log *models.SystemLog) *LogResponse {
	var context map[string]interface{}
	if len(log.Context) > 0 {
		json.Unmarshal(log.Context, &context)
	}

	return &LogResponse{
		ID:        log.ID,
		Level:     string(log.Level),
		Message:   log.Message,
		Context:   context,
		Source:    log.Source,
		UserID:    log.UserID,
		IPAddress: log.IPAddress,
		CreatedAt: log.CreatedAt.Format(time.RFC3339),
	}
}
