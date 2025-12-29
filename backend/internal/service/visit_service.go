package service

import (
	"time"

	"github.com/whk-newbie/blog/internal/models"
	"github.com/whk-newbie/blog/internal/repository"
)

// VisitService 访问记录服务接口
type VisitService interface {
	// 记录访问
	RecordVisit(req *RecordVisitRequest) error
	// 获取访问统计
	GetVisitStats(req *VisitStatsRequest) (*VisitStatsResponse, error)
	// 获取热门文章
	GetPopularArticles(limit, days int) ([]repository.PopularArticle, error)
	// 获取访问来源统计
	GetReferrerStats(startDate, endDate time.Time) (*ReferrerStatsResponse, error)
}

// RecordVisitRequest 记录访问请求
type RecordVisitRequest struct {
	FingerprintID *uint  `json:"fingerprint_id"`
	URL           string `json:"url" binding:"required"`
	Referrer      string `json:"referrer"`
	PageTitle     string `json:"page_title"`
	ArticleID     *uint  `json:"article_id"`
	StayDuration  *int   `json:"stay_duration"`
	UserAgent     string `json:"user_agent"`
}

// VisitStatsRequest 访问统计请求
type VisitStatsRequest struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Type      string    `json:"type"` // daily, weekly, monthly
}

// VisitStatsResponse 访问统计响应
type VisitStatsResponse struct {
	Summary struct {
		TotalPV         int64   `json:"total_pv"`
		TotalUV         int64   `json:"total_uv"`
		AvgStayDuration float64 `json:"avg_stay_duration"`
	} `json:"summary"`
	DailyStats []DailyStat `json:"daily_stats"`
}

// DailyStat 每日统计
type DailyStat struct {
	Date            string  `json:"date"`
	PV              int64   `json:"pv"`
	UV              int64   `json:"uv"`
	AvgStayDuration float64 `json:"avg_stay_duration"`
}

// ReferrerStatsResponse 访问来源统计响应
type ReferrerStatsResponse struct {
	Direct       int64                    `json:"direct"`
	SearchEngine int64                    `json:"search_engine"`
	ExternalLink int64                    `json:"external_link"`
	TopReferrers []repository.TopReferrer `json:"top_referrers"`
}

// visitService 访问记录服务实现
type visitService struct {
	visitRepo repository.VisitRepository
}

// NewVisitService 创建访问记录服务
func NewVisitService(visitRepo repository.VisitRepository) VisitService {
	return &visitService{
		visitRepo: visitRepo,
	}
}

// RecordVisit 记录访问
func (s *visitService) RecordVisit(req *RecordVisitRequest) error {
	visit := &models.Visit{
		FingerprintID: req.FingerprintID,
		URL:           req.URL,
		Referrer:      req.Referrer,
		PageTitle:     req.PageTitle,
		ArticleID:     req.ArticleID,
		StayDuration:  req.StayDuration,
		UserAgent:     req.UserAgent,
		VisitTime:     time.Now(),
	}

	return s.visitRepo.Create(visit)
}

// GetVisitStats 获取访问统计
func (s *visitService) GetVisitStats(req *VisitStatsRequest) (*VisitStatsResponse, error) {
	// 设置默认日期范围（最近30天）
	if req.StartDate.IsZero() {
		req.StartDate = time.Now().AddDate(0, 0, -30)
	}
	if req.EndDate.IsZero() {
		req.EndDate = time.Now()
	}

	// 设置默认统计类型
	if req.Type == "" {
		req.Type = "daily"
	}

	// 获取PV统计
	pvStats, err := s.visitRepo.GetPVStats(req.StartDate, req.EndDate, req.Type)
	if err != nil {
		return nil, err
	}

	// 获取UV统计
	uvStats, err := s.visitRepo.GetUVStats(req.StartDate, req.EndDate, req.Type)
	if err != nil {
		return nil, err
	}

	// 获取平均停留时间统计
	stayDurationStats, err := s.visitRepo.GetAvgStayDuration(req.StartDate, req.EndDate, req.Type)
	if err != nil {
		return nil, err
	}

	// 合并统计数据
	dailyStatsMap := make(map[string]*DailyStat)

	// 初始化PV
	for _, stat := range pvStats {
		dailyStatsMap[stat.Date] = &DailyStat{
			Date: stat.Date,
			PV:   stat.PV,
		}
	}

	// 合并UV
	for _, stat := range uvStats {
		if existing, ok := dailyStatsMap[stat.Date]; ok {
			existing.UV = stat.UV
		} else {
			dailyStatsMap[stat.Date] = &DailyStat{
				Date: stat.Date,
				UV:   stat.UV,
			}
		}
	}

	// 合并平均停留时间
	for _, stat := range stayDurationStats {
		if existing, ok := dailyStatsMap[stat.Date]; ok {
			existing.AvgStayDuration = stat.AvgStayDuration
		} else {
			dailyStatsMap[stat.Date] = &DailyStat{
				Date:            stat.Date,
				AvgStayDuration: stat.AvgStayDuration,
			}
		}
	}

	// 转换为切片
	var dailyStats []DailyStat
	for _, stat := range dailyStatsMap {
		dailyStats = append(dailyStats, *stat)
	}

	// 计算总计
	var totalPV int64
	var totalUV int64
	var totalStayDuration float64
	var stayDurationCount int

	for _, stat := range dailyStats {
		totalPV += stat.PV
		if stat.UV > totalUV {
			totalUV = stat.UV // UV取最大值（去重后的总数）
		}
		if stat.AvgStayDuration > 0 {
			totalStayDuration += stat.AvgStayDuration
			stayDurationCount++
		}
	}

	avgStayDuration := 0.0
	if stayDurationCount > 0 {
		avgStayDuration = totalStayDuration / float64(stayDurationCount)
	}

	response := &VisitStatsResponse{
		DailyStats: dailyStats,
	}
	response.Summary.TotalPV = totalPV
	response.Summary.TotalUV = totalUV
	response.Summary.AvgStayDuration = avgStayDuration

	return response, nil
}

// GetPopularArticles 获取热门文章
func (s *visitService) GetPopularArticles(limit, days int) ([]repository.PopularArticle, error) {
	if limit <= 0 {
		limit = 10
	}
	if days <= 0 {
		days = 7
	}

	return s.visitRepo.GetPopularArticles(limit, days)
}

// GetReferrerStats 获取访问来源统计
func (s *visitService) GetReferrerStats(startDate, endDate time.Time) (*ReferrerStatsResponse, error) {
	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, 0, -30)
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}

	stats, err := s.visitRepo.GetReferrerStats(startDate, endDate)
	if err != nil {
		return nil, err
	}

	return &ReferrerStatsResponse{
		Direct:       stats.Direct,
		SearchEngine: stats.SearchEngine,
		ExternalLink: stats.ExternalLink,
		TopReferrers: stats.TopReferrers,
	}, nil
}
