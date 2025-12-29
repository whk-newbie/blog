package scheduler

import (
	"time"

	"github.com/robfig/cron/v3"
	"github.com/whk-newbie/blog/internal/pkg/logger"
	"github.com/whk-newbie/blog/internal/service"
)

// LogScheduler 日志清理调度器
type LogScheduler struct {
	cron          *cron.Cron
	logService    service.LogService
	retentionDays int
}

// NewLogScheduler 创建日志清理调度器
func NewLogScheduler(logService service.LogService, retentionDays int) *LogScheduler {
	// 创建带秒级精度的cron调度器
	c := cron.New(cron.WithSeconds())

	if retentionDays <= 0 {
		retentionDays = 90 // 默认90天
	}

	return &LogScheduler{
		cron:          c,
		logService:    logService,
		retentionDays: retentionDays,
	}
}

// Start 启动调度器
func (s *LogScheduler) Start() error {
	// 每天凌晨2点执行日志清理
	_, err := s.cron.AddFunc("0 0 2 * * *", s.cleanupOldLogs)
	if err != nil {
		logger.Error("Failed to add log cleanup job: %v", err)
		return err
	}

	// 启动调度器
	s.cron.Start()
	logger.Info("Log cleanup scheduler started (retention: %d days)", s.retentionDays)

	return nil
}

// Stop 停止调度器
func (s *LogScheduler) Stop() {
	if s.cron != nil {
		ctx := s.cron.Stop()
		<-ctx.Done()
		logger.Info("Log cleanup scheduler stopped")
	}
}

// cleanupOldLogs 清理旧日志
func (s *LogScheduler) cleanupOldLogs() {
	logger.Info("Starting log cleanup at %s", time.Now().Format("2006-01-02 15:04:05"))

	result, err := s.logService.CleanupOldLogs(s.retentionDays)
	if err != nil {
		logger.Error("Failed to cleanup old logs: %v", err)
		return
	}

	logger.Info("Log cleanup completed: deleted %d logs", result.DeletedCount)
}

// SetRetentionDays 设置保留天数
func (s *LogScheduler) SetRetentionDays(days int) {
	if days > 0 {
		s.retentionDays = days
		logger.Info("Log retention days updated to %d", days)
	}
}
