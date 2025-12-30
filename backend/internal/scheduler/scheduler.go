package scheduler

import (
	"github.com/whk-newbie/blog/internal/service"
)

// Manager 调度器管理器
type Manager struct {
	articleScheduler *ArticleScheduler
	logScheduler     *LogScheduler
	backupScheduler  *BackupScheduler
}

// NewManager 创建调度器管理器
func NewManager(articleService service.ArticleService, logService service.LogService, backupService service.BackupService, logRetentionDays int, backupSchedule string, backupRetentionCount int) *Manager {
	return &Manager{
		articleScheduler: NewArticleScheduler(articleService),
		logScheduler:     NewLogScheduler(logService, logRetentionDays),
		backupScheduler:  NewBackupScheduler(backupService, backupSchedule, backupRetentionCount),
	}
}

// Start 启动所有调度器
func (m *Manager) Start() error {
	// 启动文章调度器
	if err := m.articleScheduler.Start(); err != nil {
		return err
	}

	// 启动日志清理调度器
	if err := m.logScheduler.Start(); err != nil {
		return err
	}

	// 启动备份调度器
	if m.backupScheduler != nil {
		if err := m.backupScheduler.Start(); err != nil {
			return err
		}
	}

	return nil
}

// Stop 停止所有调度器
func (m *Manager) Stop() {
	if m.articleScheduler != nil {
		m.articleScheduler.Stop()
	}
	if m.logScheduler != nil {
		m.logScheduler.Stop()
	}
	if m.backupScheduler != nil {
		m.backupScheduler.Stop()
	}
}

// GetArticleScheduler 获取文章调度器
func (m *Manager) GetArticleScheduler() *ArticleScheduler {
	return m.articleScheduler
}

// GetLogScheduler 获取日志调度器
func (m *Manager) GetLogScheduler() *LogScheduler {
	return m.logScheduler
}

// GetBackupScheduler 获取备份调度器
func (m *Manager) GetBackupScheduler() *BackupScheduler {
	return m.backupScheduler
}
