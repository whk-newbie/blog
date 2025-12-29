package scheduler

import (
	"github.com/whk-newbie/blog/internal/service"
)

// Manager 调度器管理器
type Manager struct {
	articleScheduler *ArticleScheduler
}

// NewManager 创建调度器管理器
func NewManager(articleService service.ArticleService) *Manager {
	return &Manager{
		articleScheduler: NewArticleScheduler(articleService),
	}
}

// Start 启动所有调度器
func (m *Manager) Start() error {
	// 启动文章调度器
	if err := m.articleScheduler.Start(); err != nil {
		return err
	}

	return nil
}

// Stop 停止所有调度器
func (m *Manager) Stop() {
	if m.articleScheduler != nil {
		m.articleScheduler.Stop()
	}
}

// GetArticleScheduler 获取文章调度器
func (m *Manager) GetArticleScheduler() *ArticleScheduler {
	return m.articleScheduler
}
