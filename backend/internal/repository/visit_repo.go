package repository

import (
	"errors"
	"time"

	"github.com/whk-newbie/blog/internal/models"
	"gorm.io/gorm"
)

var (
	ErrVisitNotFound = errors.New("visit not found")
)

// VisitFilter 访问记录筛选条件
type VisitFilter struct {
	FingerprintID *uint
	ArticleID     *uint
	StartDate     *time.Time
	EndDate       *time.Time
}

// VisitRepository 访问记录仓库接口
type VisitRepository interface {
	// 创建访问记录
	Create(visit *models.Visit) error
	// 根据ID查找访问记录
	FindByID(id uint) (*models.Visit, error)
	// 获取访问记录列表（带筛选）
	List(filter *VisitFilter, offset, limit int) ([]models.Visit, int64, error)
	// 获取PV统计（按日期）
	GetPVStats(startDate, endDate time.Time, groupBy string) ([]PVStat, error)
	// 获取UV统计（按日期）
	GetUVStats(startDate, endDate time.Time, groupBy string) ([]UVStat, error)
	// 获取平均停留时间（按日期）
	GetAvgStayDuration(startDate, endDate time.Time, groupBy string) ([]AvgStayDurationStat, error)
	// 获取热门文章统计
	GetPopularArticles(limit int, days int) ([]PopularArticle, error)
	// 获取访问来源统计
	GetReferrerStats(startDate, endDate time.Time) (*ReferrerStats, error)
}

// PVStat PV统计
type PVStat struct {
	Date string `json:"date"`
	PV   int64  `json:"pv"`
}

// UVStat UV统计
type UVStat struct {
	Date string `json:"date"`
	UV   int64  `json:"uv"`
}

// AvgStayDurationStat 平均停留时间统计
type AvgStayDurationStat struct {
	Date            string  `json:"date"`
	AvgStayDuration float64 `json:"avg_stay_duration"`
}

// PopularArticle 热门文章
type PopularArticle struct {
	ArticleID       uint    `json:"article_id"`
	Title           string  `json:"title"`
	ViewCount       int64   `json:"view_count"`
	VisitCount      int64   `json:"visit_count"`
	AvgStayDuration float64 `json:"avg_stay_duration"`
}

// ReferrerStats 访问来源统计
type ReferrerStats struct {
	Direct       int64         `json:"direct"`
	SearchEngine int64         `json:"search_engine"`
	ExternalLink int64         `json:"external_link"`
	TopReferrers []TopReferrer `json:"top_referrers"`
}

// TopReferrer 热门来源
type TopReferrer struct {
	Referrer string `json:"referrer"`
	Count    int64  `json:"count"`
}

// visitRepository 访问记录仓库实现
type visitRepository struct {
	db *gorm.DB
}

// NewVisitRepository 创建访问记录仓库
func NewVisitRepository(db *gorm.DB) VisitRepository {
	return &visitRepository{db: db}
}

// Create 创建访问记录
func (r *visitRepository) Create(visit *models.Visit) error {
	return r.db.Create(visit).Error
}

// FindByID 根据ID查找访问记录
func (r *visitRepository) FindByID(id uint) (*models.Visit, error) {
	var visit models.Visit
	err := r.db.First(&visit, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVisitNotFound
		}
		return nil, err
	}
	return &visit, nil
}

// List 获取访问记录列表（带筛选）
func (r *visitRepository) List(filter *VisitFilter, offset, limit int) ([]models.Visit, int64, error) {
	var visits []models.Visit
	var total int64

	query := r.db.Model(&models.Visit{})

	// 应用筛选条件
	query = r.applyFilter(query, filter)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表
	query = query.Preload("Fingerprint").
		Preload("Article").
		Order("visit_time DESC")

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&visits).Error; err != nil {
		return nil, 0, err
	}

	return visits, total, nil
}

// applyFilter 应用筛选条件
func (r *visitRepository) applyFilter(query *gorm.DB, filter *VisitFilter) *gorm.DB {
	if filter == nil {
		return query
	}

	if filter.FingerprintID != nil {
		query = query.Where("fingerprint_id = ?", *filter.FingerprintID)
	}

	if filter.ArticleID != nil {
		query = query.Where("article_id = ?", *filter.ArticleID)
	}

	if filter.StartDate != nil {
		query = query.Where("visit_time >= ?", *filter.StartDate)
	}

	if filter.EndDate != nil {
		query = query.Where("visit_time <= ?", *filter.EndDate)
	}

	return query
}

// GetPVStats 获取PV统计（按日期）
func (r *visitRepository) GetPVStats(startDate, endDate time.Time, groupBy string) ([]PVStat, error) {
	var stats []PVStat

	dateFormat := "YYYY-MM-DD"
	if groupBy == "weekly" {
		dateFormat = "YYYY-\"W\"WW"
	} else if groupBy == "monthly" {
		dateFormat = "YYYY-MM"
	}

	query := r.db.Model(&models.Visit{}).
		Select("TO_CHAR(visit_time, ?) as date, COUNT(*) as pv", dateFormat).
		Where("visit_time >= ? AND visit_time <= ?", startDate, endDate).
		Group("date").
		Order("date ASC")

	if err := query.Scan(&stats).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// GetUVStats 获取UV统计（按日期）
func (r *visitRepository) GetUVStats(startDate, endDate time.Time, groupBy string) ([]UVStat, error) {
	var stats []UVStat

	dateFormat := "YYYY-MM-DD"
	if groupBy == "weekly" {
		dateFormat = "YYYY-\"W\"WW"
	} else if groupBy == "monthly" {
		dateFormat = "YYYY-MM"
	}

	query := r.db.Model(&models.Visit{}).
		Select("TO_CHAR(visit_time, ?) as date, COUNT(DISTINCT fingerprint_id) as uv", dateFormat).
		Where("visit_time >= ? AND visit_time <= ?", startDate, endDate).
		Where("fingerprint_id IS NOT NULL").
		Group("date").
		Order("date ASC")

	if err := query.Scan(&stats).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// GetAvgStayDuration 获取平均停留时间（按日期）
func (r *visitRepository) GetAvgStayDuration(startDate, endDate time.Time, groupBy string) ([]AvgStayDurationStat, error) {
	var stats []AvgStayDurationStat

	dateFormat := "YYYY-MM-DD"
	if groupBy == "weekly" {
		dateFormat = "YYYY-\"W\"WW"
	} else if groupBy == "monthly" {
		dateFormat = "YYYY-MM"
	}

	query := r.db.Model(&models.Visit{}).
		Select("TO_CHAR(visit_time, ?) as date, AVG(stay_duration) as avg_stay_duration", dateFormat).
		Where("visit_time >= ? AND visit_time <= ?", startDate, endDate).
		Where("stay_duration IS NOT NULL").
		Group("date").
		Order("date ASC")

	if err := query.Scan(&stats).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// GetPopularArticles 获取热门文章统计
func (r *visitRepository) GetPopularArticles(limit int, days int) ([]PopularArticle, error) {
	var articles []PopularArticle

	startDate := time.Now().AddDate(0, 0, -days)

	query := r.db.Model(&models.Visit{}).
		Select(`
			article_id,
			COUNT(*) as visit_count,
			AVG(stay_duration) as avg_stay_duration
		`).
		Where("visit_time >= ?", startDate).
		Where("article_id IS NOT NULL").
		Group("article_id").
		Order("visit_count DESC").
		Limit(limit)

	if err := query.Scan(&articles).Error; err != nil {
		return nil, err
	}

	// 获取文章标题和浏览量
	for i := range articles {
		var article models.Article
		if err := r.db.Select("title, view_count").First(&article, articles[i].ArticleID).Error; err == nil {
			articles[i].Title = article.Title
			articles[i].ViewCount = int64(article.ViewCount)
		}
	}

	return articles, nil
}

// GetReferrerStats 获取访问来源统计
func (r *visitRepository) GetReferrerStats(startDate, endDate time.Time) (*ReferrerStats, error) {
	stats := &ReferrerStats{}

	// 统计直接访问（referrer为空或为本站）
	directQuery := r.db.Model(&models.Visit{}).
		Where("visit_time >= ? AND visit_time <= ?", startDate, endDate).
		Where("referrer = '' OR referrer IS NULL")

	if err := directQuery.Count(&stats.Direct).Error; err != nil {
		return nil, err
	}

	// 统计搜索引擎来源
	searchEngines := []string{"google.com", "bing.com", "baidu.com", "yahoo.com", "duckduckgo.com"}
	searchQuery := r.db.Model(&models.Visit{}).
		Where("visit_time >= ? AND visit_time <= ?", startDate, endDate).
		Where("referrer != '' AND referrer IS NOT NULL")

	for _, engine := range searchEngines {
		searchQuery = searchQuery.Or("referrer LIKE ?", "%"+engine+"%")
	}

	if err := searchQuery.Count(&stats.SearchEngine).Error; err != nil {
		return nil, err
	}

	// 统计外部链接（排除搜索引擎）
	externalQuery := r.db.Model(&models.Visit{}).
		Where("visit_time >= ? AND visit_time <= ?", startDate, endDate).
		Where("referrer != '' AND referrer IS NOT NULL")

	for _, engine := range searchEngines {
		externalQuery = externalQuery.Where("referrer NOT LIKE ?", "%"+engine+"%")
	}

	if err := externalQuery.Count(&stats.ExternalLink).Error; err != nil {
		return nil, err
	}

	// 获取热门来源
	var topReferrers []TopReferrer
	topQuery := r.db.Model(&models.Visit{}).
		Select("referrer, COUNT(*) as count").
		Where("visit_time >= ? AND visit_time <= ?", startDate, endDate).
		Where("referrer != '' AND referrer IS NOT NULL").
		Group("referrer").
		Order("count DESC").
		Limit(10)

	if err := topQuery.Scan(&topReferrers).Error; err != nil {
		return nil, err
	}

	stats.TopReferrers = topReferrers

	return stats, nil
}
