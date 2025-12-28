-- 001_init_schema.sql
-- 初始化数据库表结构

-- 管理员表
CREATE TABLE IF NOT EXISTS admins (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100),
    is_default_password BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

COMMENT ON TABLE admins IS '管理员表';
COMMENT ON COLUMN admins.password IS '密码（BCrypt加密）';
COMMENT ON COLUMN admins.is_default_password IS '是否使用默认密码（用于提示修改）';

-- 分类表
CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    sort_order INT DEFAULT 0,
    article_count INT DEFAULT 0,
    created_by BIGINT REFERENCES admins(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

COMMENT ON TABLE categories IS '文章分类表';
COMMENT ON COLUMN categories.slug IS 'URL友好的标识符';
COMMENT ON COLUMN categories.article_count IS '文章数量（冗余字段）';

-- 标签表
CREATE TABLE IF NOT EXISTS tags (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL,
    article_count INT DEFAULT 0,
    created_by BIGINT REFERENCES admins(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

COMMENT ON TABLE tags IS '文章标签表';
COMMENT ON COLUMN tags.article_count IS '文章数量（冗余字段）';

-- 文章表
CREATE TABLE IF NOT EXISTS articles (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    summary TEXT,
    content TEXT NOT NULL,
    cover_image VARCHAR(500),
    category_id BIGINT REFERENCES categories(id) ON DELETE SET NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    publish_at TIMESTAMP,
    view_count INT DEFAULT 0,
    like_count INT DEFAULT 0,
    is_top BOOLEAN DEFAULT FALSE,
    is_featured BOOLEAN DEFAULT FALSE,
    author_id BIGINT REFERENCES admins(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    search_vector tsvector
);

COMMENT ON TABLE articles IS '文章表';
COMMENT ON COLUMN articles.slug IS 'URL友好的标识符';
COMMENT ON COLUMN articles.status IS '文章状态：draft=草稿，published=已发布';
COMMENT ON COLUMN articles.publish_at IS '发布时间（支持定时发布）';

-- 添加约束
ALTER TABLE articles ADD CONSTRAINT check_article_status 
    CHECK (status IN ('draft', 'published'));

-- 文章标签关联表
CREATE TABLE IF NOT EXISTS article_tags (
    id BIGSERIAL PRIMARY KEY,
    article_id BIGINT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    tag_id BIGINT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(article_id, tag_id)
);

COMMENT ON TABLE article_tags IS '文章标签关联表';

-- 浏览器指纹表
CREATE TABLE IF NOT EXISTS fingerprints (
    id BIGSERIAL PRIMARY KEY,
    fingerprint_hash VARCHAR(64) UNIQUE NOT NULL,
    fingerprint_data JSONB NOT NULL,
    user_agent TEXT,
    first_seen_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_seen_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    visit_count INT DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE fingerprints IS '浏览器指纹表';
COMMENT ON COLUMN fingerprints.fingerprint_hash IS '指纹哈希值（SHA256）';
COMMENT ON COLUMN fingerprints.fingerprint_data IS '完整指纹信息（JSON格式）';

-- 访问记录表
CREATE TABLE IF NOT EXISTS visits (
    id BIGSERIAL PRIMARY KEY,
    fingerprint_id BIGINT REFERENCES fingerprints(id) ON DELETE SET NULL,
    url VARCHAR(500) NOT NULL,
    referrer VARCHAR(500),
    page_title VARCHAR(255),
    article_id BIGINT REFERENCES articles(id) ON DELETE SET NULL,
    visit_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    stay_duration INT,
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE visits IS '访问记录表';
COMMENT ON COLUMN visits.stay_duration IS '页面停留时间（秒）';
COMMENT ON COLUMN visits.referrer IS '来源URL（用于分析流量来源）';

-- 爬虫任务表
CREATE TABLE IF NOT EXISTS crawl_tasks (
    id BIGSERIAL PRIMARY KEY,
    task_id VARCHAR(100) UNIQUE NOT NULL,
    task_name VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'running',
    progress INT DEFAULT 0,
    message TEXT,
    start_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMP,
    duration INT,
    created_by_token VARCHAR(64),
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE crawl_tasks IS '爬虫任务表';
COMMENT ON COLUMN crawl_tasks.task_id IS '任务唯一标识（由Python SDK生成）';
COMMENT ON COLUMN crawl_tasks.status IS '任务状态：running=运行中，completed=已完成，failed=失败';
COMMENT ON COLUMN crawl_tasks.progress IS '任务进度（0-100）';

-- 添加约束
ALTER TABLE crawl_tasks ADD CONSTRAINT check_task_status 
    CHECK (status IN ('running', 'completed', 'failed'));

ALTER TABLE crawl_tasks ADD CONSTRAINT check_task_progress 
    CHECK (progress >= 0 AND progress <= 100);

-- 系统配置表
CREATE TABLE IF NOT EXISTS system_configs (
    id BIGSERIAL PRIMARY KEY,
    config_key VARCHAR(100) UNIQUE NOT NULL,
    config_value TEXT,
    config_type VARCHAR(50) NOT NULL,
    is_encrypted BOOLEAN DEFAULT TRUE,
    is_active BOOLEAN DEFAULT TRUE,
    description TEXT,
    created_by BIGINT REFERENCES admins(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

COMMENT ON TABLE system_configs IS '系统配置表';
COMMENT ON COLUMN system_configs.config_key IS '配置键（唯一）';
COMMENT ON COLUMN system_configs.config_value IS '配置值（敏感信息加密存储）';
COMMENT ON COLUMN system_configs.is_encrypted IS '是否加密存储';

-- 系统日志表
CREATE TABLE IF NOT EXISTS system_logs (
    id BIGSERIAL PRIMARY KEY,
    level VARCHAR(20) NOT NULL,
    message TEXT NOT NULL,
    context JSONB,
    source VARCHAR(100),
    user_id BIGINT REFERENCES admins(id),
    ip_address VARCHAR(45),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE system_logs IS '系统日志表';
COMMENT ON COLUMN system_logs.level IS '日志级别：DEBUG/INFO/WARN/ERROR';
COMMENT ON COLUMN system_logs.context IS '上下文信息（JSON格式）';

-- 添加约束
ALTER TABLE system_logs ADD CONSTRAINT check_log_level 
    CHECK (level IN ('DEBUG', 'INFO', 'WARN', 'ERROR'));

-- 留言记录表
CREATE TABLE IF NOT EXISTS messages (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100),
    fingerprint_id BIGINT REFERENCES fingerprints(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE messages IS '留言记录表（仅用于统计）';
COMMENT ON COLUMN messages.fingerprint_id IS '留言者的浏览器指纹ID';

