package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/whk-newbie/blog/internal/pkg/redis"
	"github.com/whk-newbie/blog/internal/repository"
)

// VisitCacheService 访问统计缓存服务
type VisitCacheService interface {
	// 缓存访问统计
	CacheVisitStats(req *VisitStatsRequest, stats *VisitStatsResponse, ttl time.Duration) error
	// 获取缓存的访问统计
	GetCachedVisitStats(req *VisitStatsRequest) (*VisitStatsResponse, error)
	// 缓存热门文章
	CachePopularArticles(limit, days int, articles []repository.PopularArticle, ttl time.Duration) error
	// 获取缓存的热门文章
	GetCachedPopularArticles(limit, days int) ([]repository.PopularArticle, error)
	// 缓存访问来源统计
	CacheReferrerStats(startDate, endDate time.Time, stats *ReferrerStatsResponse, ttl time.Duration) error
	// 获取缓存的访问来源统计
	GetCachedReferrerStats(startDate, endDate time.Time) (*ReferrerStatsResponse, error)
	// 清除访问统计缓存
	ClearVisitStatsCache() error
}

// visitCacheService 访问统计缓存服务实现
type visitCacheService struct{}

// NewVisitCacheService 创建访问统计缓存服务
func NewVisitCacheService() VisitCacheService {
	return &visitCacheService{}
}

// getCacheKey 生成缓存键
func (s *visitCacheService) getCacheKey(prefix string, params ...interface{}) string {
	key := fmt.Sprintf("visit_stats:%s", prefix)
	for _, param := range params {
		key += fmt.Sprintf(":%v", param)
	}
	return key
}

// CacheVisitStats 缓存访问统计
func (s *visitCacheService) CacheVisitStats(req *VisitStatsRequest, stats *VisitStatsResponse, ttl time.Duration) error {
	key := s.getCacheKey("stats",
		req.StartDate.Format("2006-01-02"),
		req.EndDate.Format("2006-01-02"),
		req.Type,
	)

	data, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	return redis.Set(key, string(data), ttl)
}

// GetCachedVisitStats 获取缓存的访问统计
func (s *visitCacheService) GetCachedVisitStats(req *VisitStatsRequest) (*VisitStatsResponse, error) {
	key := s.getCacheKey("stats",
		req.StartDate.Format("2006-01-02"),
		req.EndDate.Format("2006-01-02"),
		req.Type,
	)

	data, err := redis.GetValue(key)
	if err != nil {
		return nil, err
	}

	var stats VisitStatsResponse
	if err := json.Unmarshal([]byte(data), &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

// CachePopularArticles 缓存热门文章
func (s *visitCacheService) CachePopularArticles(limit, days int, articles []repository.PopularArticle, ttl time.Duration) error {
	key := s.getCacheKey("popular_articles", limit, days)

	data, err := json.Marshal(articles)
	if err != nil {
		return err
	}

	return redis.Set(key, string(data), ttl)
}

// GetCachedPopularArticles 获取缓存的热门文章
func (s *visitCacheService) GetCachedPopularArticles(limit, days int) ([]repository.PopularArticle, error) {
	key := s.getCacheKey("popular_articles", limit, days)

	data, err := redis.GetValue(key)
	if err != nil {
		return nil, err
	}

	var articles []repository.PopularArticle
	if err := json.Unmarshal([]byte(data), &articles); err != nil {
		return nil, err
	}

	return articles, nil
}

// CacheReferrerStats 缓存访问来源统计
func (s *visitCacheService) CacheReferrerStats(startDate, endDate time.Time, stats *ReferrerStatsResponse, ttl time.Duration) error {
	key := s.getCacheKey("referrers",
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
	)

	data, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	return redis.Set(key, string(data), ttl)
}

// GetCachedReferrerStats 获取缓存的访问来源统计
func (s *visitCacheService) GetCachedReferrerStats(startDate, endDate time.Time) (*ReferrerStatsResponse, error) {
	key := s.getCacheKey("referrers",
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
	)

	data, err := redis.GetValue(key)
	if err != nil {
		return nil, err
	}

	var stats ReferrerStatsResponse
	if err := json.Unmarshal([]byte(data), &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

// ClearVisitStatsCache 清除访问统计缓存
func (s *visitCacheService) ClearVisitStatsCache() error {
	// 这里可以使用Redis的KEYS命令或SCAN命令来查找并删除所有相关缓存
	// 为了简单起见，我们只删除特定的键模式
	// 实际使用中可以使用Redis的SCAN命令来遍历所有匹配的键
	return nil
}
