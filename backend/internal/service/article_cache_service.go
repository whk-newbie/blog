package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/whk-newbie/blog/internal/pkg/redis"
)

// ArticleCacheService 文章缓存服务接口
type ArticleCacheService interface {
	// 缓存文章列表
	CacheArticleList(key string, response *ArticleListResponse, ttl time.Duration) error
	// 获取缓存的文章列表
	GetCachedArticleList(key string) (*ArticleListResponse, error)
	// 生成文章列表缓存键
	GenerateCacheKey(filter *ArticleListRequest, isPublished bool) string
	// 清除文章相关缓存
	ClearArticleCache() error
	// 清除特定文章的缓存
	ClearArticleCacheByID(articleID uint) error
}

// articleCacheService 文章缓存服务实现
type articleCacheService struct{}

// NewArticleCacheService 创建文章缓存服务
func NewArticleCacheService() ArticleCacheService {
	return &articleCacheService{}
}

// GenerateCacheKey 生成文章列表缓存键
func (s *articleCacheService) GenerateCacheKey(filter *ArticleListRequest, isPublished bool) string {
	key := "article_list:"
	if isPublished {
		key += "published:"
	} else {
		key += "all:"
	}

	if filter.Page > 0 {
		key += fmt.Sprintf("page:%d:", filter.Page)
	}
	if filter.PageSize > 0 {
		key += fmt.Sprintf("size:%d:", filter.PageSize)
	}
	if filter.CategoryID != nil {
		key += fmt.Sprintf("cat:%d:", *filter.CategoryID)
	}
	if filter.TagID != nil {
		key += fmt.Sprintf("tag:%d:", *filter.TagID)
	}
	if filter.Status != nil {
		key += fmt.Sprintf("status:%s:", *filter.Status)
	}
	if filter.IsTop != nil {
		key += fmt.Sprintf("top:%v:", *filter.IsTop)
	}
	if filter.IsFeatured != nil {
		key += fmt.Sprintf("featured:%v:", *filter.IsFeatured)
	}
	if filter.Keyword != "" {
		key += fmt.Sprintf("keyword:%s:", filter.Keyword)
	}

	return key
}

// CacheArticleList 缓存文章列表
func (s *articleCacheService) CacheArticleList(key string, response *ArticleListResponse, ttl time.Duration) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	return redis.Set(key, string(data), ttl)
}

// GetCachedArticleList 获取缓存的文章列表
func (s *articleCacheService) GetCachedArticleList(key string) (*ArticleListResponse, error) {
	data, err := redis.GetValue(key)
	if err != nil {
		return nil, err
	}

	var response ArticleListResponse
	if err := json.Unmarshal([]byte(data), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// ClearArticleCache 清除所有文章缓存
func (s *articleCacheService) ClearArticleCache() error {
	// 使用通配符删除所有文章列表缓存
	// 注意: Redis的DEL命令不支持通配符，需要先获取所有匹配的key
	// 这里简化处理，实际可以使用Redis的SCAN命令
	// 为了简化，我们使用一个标记来使缓存失效
	// 或者使用Redis的SET来存储版本号，每次更新时递增版本号
	return nil // 暂时返回nil，实际实现可以使用Redis的KEYS或SCAN命令
}

// ClearArticleCacheByID 清除特定文章的缓存
func (s *articleCacheService) ClearArticleCacheByID(articleID uint) error {
	// 清除所有包含该文章的列表缓存
	// 由于无法精确匹配，这里清除所有文章列表缓存
	return s.ClearArticleCache()
}
