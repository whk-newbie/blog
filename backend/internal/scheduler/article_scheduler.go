package scheduler

import (
	"time"

	"github.com/iambaby/blog/internal/pkg/logger"
	"github.com/iambaby/blog/internal/service"
	"github.com/robfig/cron/v3"
)

// ArticleScheduler 文章调度器
type ArticleScheduler struct {
	cron           *cron.Cron
	articleService service.ArticleService
}

// NewArticleScheduler 创建文章调度器
func NewArticleScheduler(articleService service.ArticleService) *ArticleScheduler {
	// 创建带秒级精度的cron调度器
	c := cron.New(cron.WithSeconds())

	return &ArticleScheduler{
		cron:           c,
		articleService: articleService,
	}
}

// Start 启动调度器
func (s *ArticleScheduler) Start() error {
	// 每分钟执行一次定时发布检查
	_, err := s.cron.AddFunc("0 * * * * *", s.processScheduledPublish)
	if err != nil {
		logger.Error("Failed to add scheduled publish job: %v", err)
		return err
	}

	// 启动调度器
	s.cron.Start()
	logger.Info("Article scheduler started")

	return nil
}

// Stop 停止调度器
func (s *ArticleScheduler) Stop() {
	if s.cron != nil {
		ctx := s.cron.Stop()
		<-ctx.Done()
		logger.Info("Article scheduler stopped")
	}
}

// processScheduledPublish 处理定时发布
func (s *ArticleScheduler) processScheduledPublish() {
	logger.Debug("Processing scheduled publish at %s", time.Now().Format("2006-01-02 15:04:05"))

	err := s.articleService.ProcessScheduledPublish()
	if err != nil {
		logger.Error("Failed to process scheduled publish: %v", err)
		return
	}

	logger.Debug("Scheduled publish processed successfully")
}

// AddJob 添加自定义任务
func (s *ArticleScheduler) AddJob(spec string, cmd func()) (cron.EntryID, error) {
	return s.cron.AddFunc(spec, cmd)
}

// RemoveJob 移除任务
func (s *ArticleScheduler) RemoveJob(id cron.EntryID) {
	s.cron.Remove(id)
}

// GetEntries 获取所有任务
func (s *ArticleScheduler) GetEntries() []cron.Entry {
	return s.cron.Entries()
}
