package scheduler

import (
	"time"

	"github.com/robfig/cron/v3"
	"github.com/whk-newbie/blog/internal/pkg/logger"
	"github.com/whk-newbie/blog/internal/service"
)

// BackupScheduler 备份调度器
type BackupScheduler struct {
	cron           *cron.Cron
	backupService  service.BackupService
	schedule       string
	retentionCount int
}

// NewBackupScheduler 创建备份调度器
func NewBackupScheduler(backupService service.BackupService, schedule string, retentionCount int) *BackupScheduler {
	// 创建带秒级精度的cron调度器
	c := cron.New(cron.WithSeconds())

	// 默认每天凌晨3点执行备份
	if schedule == "" {
		schedule = "0 0 3 * * *" // 每天凌晨3点
	}

	// 默认保留10个备份
	if retentionCount <= 0 {
		retentionCount = 10
	}

	return &BackupScheduler{
		cron:           c,
		backupService:  backupService,
		schedule:       schedule,
		retentionCount: retentionCount,
	}
}

// Start 启动调度器
func (s *BackupScheduler) Start() error {
	// 添加备份任务
	_, err := s.cron.AddFunc(s.schedule, s.createBackup)
	if err != nil {
		logger.Error("Failed to add backup job: %v", err)
		return err
	}

	// 启动调度器
	s.cron.Start()
	logger.Info("Backup scheduler started (schedule: %s, retention: %d)", s.schedule, s.retentionCount)

	return nil
}

// Stop 停止调度器
func (s *BackupScheduler) Stop() {
	if s.cron != nil {
		ctx := s.cron.Stop()
		<-ctx.Done()
		logger.Info("Backup scheduler stopped")
	}
}

// createBackup 创建备份
func (s *BackupScheduler) createBackup() {
	logger.Info("Starting automatic backup at %s", time.Now().Format("2006-01-02 15:04:05"))

	// 创建备份
	backupInfo, err := s.backupService.CreateBackup()
	if err != nil {
		logger.Error("Failed to create automatic backup: %v", err)
		return
	}

	logger.Info("Automatic backup created successfully: %s (size: %d bytes)", backupInfo.Filename, backupInfo.Size)

	// 清理旧备份
	if err := s.backupService.CleanupOldBackups(s.retentionCount); err != nil {
		logger.Error("Failed to cleanup old backups: %v", err)
		return
	}

	logger.Info("Automatic backup completed")
}

// SetSchedule 设置备份计划
func (s *BackupScheduler) SetSchedule(schedule string) error {
	// 停止当前调度器
	s.Stop()

	// 更新计划
	s.schedule = schedule

	// 重新启动
	return s.Start()
}

// SetRetentionCount 设置保留数量
func (s *BackupScheduler) SetRetentionCount(count int) {
	if count > 0 {
		s.retentionCount = count
		logger.Info("Backup retention count updated to %d", count)
	}
}
