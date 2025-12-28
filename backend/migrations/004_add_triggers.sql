-- 004_add_triggers.sql
-- 添加触发器

-- 自动更新updated_at字段的触发器函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 应用到所有需要的表
DROP TRIGGER IF EXISTS update_admins_updated_at ON admins;
CREATE TRIGGER update_admins_updated_at BEFORE UPDATE ON admins
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;
CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_tags_updated_at ON tags;
CREATE TRIGGER update_tags_updated_at BEFORE UPDATE ON tags
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_articles_updated_at ON articles;
CREATE TRIGGER update_articles_updated_at BEFORE UPDATE ON articles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_fingerprints_updated_at ON fingerprints;
CREATE TRIGGER update_fingerprints_updated_at BEFORE UPDATE ON fingerprints
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_crawl_tasks_updated_at ON crawl_tasks;
CREATE TRIGGER update_crawl_tasks_updated_at BEFORE UPDATE ON crawl_tasks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_system_configs_updated_at ON system_configs;
CREATE TRIGGER update_system_configs_updated_at BEFORE UPDATE ON system_configs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 更新分类/标签文章数量的函数
CREATE OR REPLACE FUNCTION update_category_article_count(
    p_category_id BIGINT
) RETURNS VOID AS $$
BEGIN
    UPDATE categories
    SET article_count = (
        SELECT COUNT(*)
        FROM articles
        WHERE category_id = p_category_id
          AND status = 'published'
          AND deleted_at IS NULL
    )
    WHERE id = p_category_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_tag_article_count(
    p_tag_id BIGINT
) RETURNS VOID AS $$
BEGIN
    UPDATE tags
    SET article_count = (
        SELECT COUNT(*)
        FROM article_tags at
        INNER JOIN articles a ON at.article_id = a.id
        WHERE at.tag_id = p_tag_id
          AND a.status = 'published'
          AND a.deleted_at IS NULL
    )
    WHERE id = p_tag_id;
END;
$$ LANGUAGE plpgsql;

-- 自动同步分类文章数量的触发器
CREATE OR REPLACE FUNCTION sync_category_article_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' OR TG_OP = 'UPDATE' THEN
        IF NEW.category_id IS NOT NULL AND NEW.status = 'published' THEN
            PERFORM update_category_article_count(NEW.category_id);
        END IF;
    END IF;
    
    IF TG_OP = 'UPDATE' THEN
        IF OLD.category_id IS NOT NULL AND OLD.category_id != COALESCE(NEW.category_id, 0) THEN
            PERFORM update_category_article_count(OLD.category_id);
        END IF;
    END IF;
    
    IF TG_OP = 'DELETE' THEN
        IF OLD.category_id IS NOT NULL THEN
            PERFORM update_category_article_count(OLD.category_id);
        END IF;
        RETURN OLD;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS sync_category_count ON articles;
CREATE TRIGGER sync_category_count AFTER INSERT OR UPDATE OR DELETE ON articles
    FOR EACH ROW EXECUTE FUNCTION sync_category_article_count();

-- 增加文章浏览次数的函数
CREATE OR REPLACE FUNCTION increment_article_view_count(
    p_article_id BIGINT
) RETURNS VOID AS $$
BEGIN
    UPDATE articles 
    SET view_count = view_count + 1,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = p_article_id;
END;
$$ LANGUAGE plpgsql;

-- 清理旧日志的函数
CREATE OR REPLACE FUNCTION cleanup_old_logs(
    p_retention_days INT DEFAULT 30
) RETURNS INT AS $$
DECLARE
    deleted_count INT;
BEGIN
    DELETE FROM system_logs
    WHERE created_at < CURRENT_TIMESTAMP - (p_retention_days || ' days')::INTERVAL;
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION cleanup_old_logs IS '清理旧日志（默认保留30天）';

