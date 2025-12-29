package service

import (
	"github.com/whk-newbie/blog/internal/models"
	"github.com/whk-newbie/blog/internal/repository"
)

// StatsService 统计服务接口
type StatsService interface {
	// 获取仪表盘统计数据
	GetDashboardStats() (*DashboardStatsResponse, error)
}

// DashboardStatsResponse 仪表盘统计响应
type DashboardStatsResponse struct {
	ArticleCount   int64            `json:"article_count"`   // 文章总数
	CategoryCount  int64            `json:"category_count"`  // 分类总数
	TagCount       int64            `json:"tag_count"`       // 标签总数
	RecentArticles []models.Article `json:"recent_articles"` // 最近文章列表
}

// statsService 统计服务实现
type statsService struct {
	articleRepo  repository.ArticleRepository
	categoryRepo repository.CategoryRepository
	tagRepo      repository.TagRepository
}

// NewStatsService 创建统计服务
func NewStatsService(
	articleRepo repository.ArticleRepository,
	categoryRepo repository.CategoryRepository,
	tagRepo repository.TagRepository,
) StatsService {
	return &statsService{
		articleRepo:  articleRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
	}
}

// GetDashboardStats 获取仪表盘统计数据
func (s *statsService) GetDashboardStats() (*DashboardStatsResponse, error) {
	// 获取文章总数
	_, articleCount, err := s.articleRepo.List(&repository.ArticleFilter{}, 0, 0)
	if err != nil {
		return nil, err
	}

	// 获取分类总数
	_, categoryCount, err := s.categoryRepo.List(0, 0)
	if err != nil {
		return nil, err
	}

	// 获取标签总数
	_, tagCount, err := s.tagRepo.List(0, 0)
	if err != nil {
		return nil, err
	}

	// 获取最近10篇文章
	recentArticles, _, err := s.articleRepo.List(&repository.ArticleFilter{}, 0, 10)
	if err != nil {
		return nil, err
	}

	return &DashboardStatsResponse{
		ArticleCount:   articleCount,
		CategoryCount:  categoryCount,
		TagCount:       tagCount,
		RecentArticles: recentArticles,
	}, nil
}
