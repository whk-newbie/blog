-- 006_add_performance_indexes.sql
-- 性能优化: 添加复合索引以优化常用查询

-- 文章表复合索引: 状态+发布时间+置顶 (用于公开文章列表查询)
CREATE INDEX IF NOT EXISTS idx_articles_status_publish_top 
ON articles(status, publish_at DESC, is_top DESC) 
WHERE status = 'published' AND deleted_at IS NULL;

-- 文章表复合索引: 分类+状态 (用于分类筛选)
CREATE INDEX IF NOT EXISTS idx_articles_category_status 
ON articles(category_id, status, publish_at DESC) 
WHERE deleted_at IS NULL;

-- 访问记录表复合索引: 访问时间+文章ID (用于统计查询)
CREATE INDEX IF NOT EXISTS idx_visits_time_article 
ON visits(visit_time DESC, article_id) 
WHERE article_id IS NOT NULL;

-- 访问记录表复合索引: 指纹ID+访问时间 (用于用户访问历史)
CREATE INDEX IF NOT EXISTS idx_visits_fingerprint_time 
ON visits(fingerprint_id, visit_time DESC);

-- 系统配置表复合索引: 类型+激活状态 (用于配置查询)
CREATE INDEX IF NOT EXISTS idx_configs_type_active 
ON system_configs(config_type, is_active) 
WHERE deleted_at IS NULL;

