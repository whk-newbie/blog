-- 003_add_fulltext_search.sql
-- 添加全文搜索功能

-- 全文搜索索引（GIN索引，支持中文和英文）
CREATE INDEX IF NOT EXISTS idx_articles_search ON articles USING GIN(search_vector);

-- 触发器函数：自动更新search_vector
CREATE OR REPLACE FUNCTION articles_search_trigger() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('simple', COALESCE(NEW.title, '')), 'A') ||
        setweight(to_tsvector('simple', COALESCE(NEW.summary, '')), 'B') ||
        setweight(to_tsvector('simple', COALESCE(NEW.content, '')), 'C');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 创建触发器
DROP TRIGGER IF EXISTS tsvector_update_trigger ON articles;
CREATE TRIGGER tsvector_update_trigger
BEFORE INSERT OR UPDATE ON articles
FOR EACH ROW EXECUTE FUNCTION articles_search_trigger();

-- 更新现有数据的search_vector
UPDATE articles SET search_vector = 
    setweight(to_tsvector('simple', COALESCE(title, '')), 'A') ||
    setweight(to_tsvector('simple', COALESCE(summary, '')), 'B') ||
    setweight(to_tsvector('simple', COALESCE(content, '')), 'C')
WHERE search_vector IS NULL;

