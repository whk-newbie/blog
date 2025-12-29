package repository

import (
	"errors"

	"github.com/whk-newbie/blog/internal/models"
	"gorm.io/gorm"
)

var (
	ErrCrawlTaskNotFound = errors.New("crawl task not found")
	ErrCrawlTaskExists   = errors.New("crawl task already exists")
)

// CrawlTaskFilter 爬虫任务筛选条件
type CrawlTaskFilter struct {
	Status         *models.CrawlTaskStatus
	TaskID         string
	CreatedByToken string
}

// CrawlTaskRepository 爬虫任务仓库接口
type CrawlTaskRepository interface {
	// 创建任务
	Create(task *models.CrawlTask) error
	// 根据ID查找任务
	FindByID(id uint) (*models.CrawlTask, error)
	// 根据TaskID查找任务
	FindByTaskID(taskID string) (*models.CrawlTask, error)
	// 更新任务
	Update(task *models.CrawlTask) error
	// 获取任务列表（带筛选）
	List(filter *CrawlTaskFilter, offset, limit int) ([]models.CrawlTask, int64, error)
	// 检查TaskID是否存在
	ExistsByTaskID(taskID string) (bool, error)
}

// crawlTaskRepository 爬虫任务仓库实现
type crawlTaskRepository struct {
	db *gorm.DB
}

// NewCrawlTaskRepository 创建爬虫任务仓库
func NewCrawlTaskRepository(db *gorm.DB) CrawlTaskRepository {
	return &crawlTaskRepository{db: db}
}

// Create 创建任务
func (r *crawlTaskRepository) Create(task *models.CrawlTask) error {
	// 检查TaskID是否已存在
	exists, err := r.ExistsByTaskID(task.TaskID)
	if err != nil {
		return err
	}
	if exists {
		return ErrCrawlTaskExists
	}

	return r.db.Create(task).Error
}

// FindByID 根据ID查找任务
func (r *crawlTaskRepository) FindByID(id uint) (*models.CrawlTask, error) {
	var task models.CrawlTask
	err := r.db.First(&task, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCrawlTaskNotFound
		}
		return nil, err
	}
	return &task, nil
}

// FindByTaskID 根据TaskID查找任务
func (r *crawlTaskRepository) FindByTaskID(taskID string) (*models.CrawlTask, error) {
	var task models.CrawlTask
	err := r.db.Where("task_id = ?", taskID).First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCrawlTaskNotFound
		}
		return nil, err
	}
	return &task, nil
}

// Update 更新任务
func (r *crawlTaskRepository) Update(task *models.CrawlTask) error {
	return r.db.Save(task).Error
}

// List 获取任务列表（带筛选）
func (r *crawlTaskRepository) List(filter *CrawlTaskFilter, offset, limit int) ([]models.CrawlTask, int64, error) {
	var tasks []models.CrawlTask
	var total int64

	query := r.db.Model(&models.CrawlTask{})

	// 应用筛选条件
	if filter != nil {
		if filter.Status != nil {
			query = query.Where("status = ?", *filter.Status)
		}
		if filter.TaskID != "" {
			query = query.Where("task_id = ?", filter.TaskID)
		}
		if filter.CreatedByToken != "" {
			query = query.Where("created_by_token = ?", filter.CreatedByToken)
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取任务列表
	query = query.Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// ExistsByTaskID 检查TaskID是否存在
func (r *crawlTaskRepository) ExistsByTaskID(taskID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.CrawlTask{}).Where("task_id = ?", taskID).Count(&count).Error
	return count > 0, err
}
