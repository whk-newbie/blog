-- 002_add_indexes.sql
-- 添加索引以优化查询性能

-- 管理员表索引
CREATE INDEX IF NOT EXISTS idx_admins_username ON admins(username);
CREATE INDEX IF NOT EXISTS idx_admins_deleted_at ON admins(deleted_at);

-- 分类表索引
CREATE INDEX IF NOT EXISTS idx_categories_slug ON categories(slug);
CREATE INDEX IF NOT EXISTS idx_categories_sort_order ON categories(sort_order);
CREATE INDEX IF NOT EXISTS idx_categories_deleted_at ON categories(deleted_at);

-- 标签表索引
CREATE INDEX IF NOT EXISTS idx_tags_slug ON tags(slug);
CREATE INDEX IF NOT EXISTS idx_tags_deleted_at ON tags(deleted_at);

-- 文章表索引
CREATE INDEX IF NOT EXISTS idx_articles_slug ON articles(slug);
CREATE INDEX IF NOT EXISTS idx_articles_category_id ON articles(category_id);
CREATE INDEX IF NOT EXISTS idx_articles_status ON articles(status);
CREATE INDEX IF NOT EXISTS idx_articles_publish_at ON articles(publish_at DESC);
CREATE INDEX IF NOT EXISTS idx_articles_view_count ON articles(view_count DESC);
CREATE INDEX IF NOT EXISTS idx_articles_is_top ON articles(is_top);
CREATE INDEX IF NOT EXISTS idx_articles_deleted_at ON articles(deleted_at);

-- 文章标签关联表索引
CREATE INDEX IF NOT EXISTS idx_article_tags_article_id ON article_tags(article_id);
CREATE INDEX IF NOT EXISTS idx_article_tags_tag_id ON article_tags(tag_id);

-- 浏览器指纹表索引
CREATE INDEX IF NOT EXISTS idx_fingerprints_hash ON fingerprints(fingerprint_hash);
CREATE INDEX IF NOT EXISTS idx_fingerprints_last_seen ON fingerprints(last_seen_at DESC);

-- JSONB字段索引（用于高效查询）
CREATE INDEX IF NOT EXISTS idx_fingerprints_data ON fingerprints USING GIN(fingerprint_data);

-- 访问记录表索引
CREATE INDEX IF NOT EXISTS idx_visits_fingerprint_id ON visits(fingerprint_id);
CREATE INDEX IF NOT EXISTS idx_visits_article_id ON visits(article_id);
CREATE INDEX IF NOT EXISTS idx_visits_visit_time ON visits(visit_time DESC);
CREATE INDEX IF NOT EXISTS idx_visits_url ON visits(url);

-- 爬虫任务表索引
CREATE INDEX IF NOT EXISTS idx_crawl_tasks_task_id ON crawl_tasks(task_id);
CREATE INDEX IF NOT EXISTS idx_crawl_tasks_status ON crawl_tasks(status);
CREATE INDEX IF NOT EXISTS idx_crawl_tasks_start_time ON crawl_tasks(start_time DESC);
CREATE INDEX IF NOT EXISTS idx_crawl_tasks_created_by_token ON crawl_tasks(created_by_token);

-- 系统配置表索引
CREATE INDEX IF NOT EXISTS idx_system_configs_key ON system_configs(config_key);
CREATE INDEX IF NOT EXISTS idx_system_configs_type ON system_configs(config_type);
CREATE INDEX IF NOT EXISTS idx_system_configs_active ON system_configs(is_active);
CREATE INDEX IF NOT EXISTS idx_system_configs_deleted_at ON system_configs(deleted_at);

-- 系统日志表索引
CREATE INDEX IF NOT EXISTS idx_system_logs_level ON system_logs(level);
CREATE INDEX IF NOT EXISTS idx_system_logs_created_at ON system_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_system_logs_source ON system_logs(source);
CREATE INDEX IF NOT EXISTS idx_system_logs_user_id ON system_logs(user_id);

-- JSONB字段索引
CREATE INDEX IF NOT EXISTS idx_system_logs_context ON system_logs USING GIN(context);

-- 留言记录表索引
CREATE INDEX IF NOT EXISTS idx_messages_email ON messages(email);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_messages_fingerprint_id ON messages(fingerprint_id);

