-- Article translation fields
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='articles' AND column_name='title_en') THEN
        ALTER TABLE articles ADD COLUMN title_en VARCHAR(500) DEFAULT '';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='articles' AND column_name='content_en') THEN
        ALTER TABLE articles ADD COLUMN content_en TEXT DEFAULT '';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='articles' AND column_name='summary_en') THEN
        ALTER TABLE articles ADD COLUMN summary_en VARCHAR(500) DEFAULT '';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='articles' AND column_name='is_translated') THEN
        ALTER TABLE articles ADD COLUMN is_translated BOOLEAN DEFAULT FALSE;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='articles' AND column_name='translated_at') THEN
        ALTER TABLE articles ADD COLUMN translated_at TIMESTAMP NULL;
    END IF;
END $$;

-- AI provider table
CREATE TABLE IF NOT EXISTS ai_providers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    provider_type VARCHAR(20) NOT NULL,
    api_key TEXT NOT NULL,
    base_url VARCHAR(500) DEFAULT '',
    model VARCHAR(100) NOT NULL,
    is_enabled BOOLEAN DEFAULT TRUE,
    sort_order INT DEFAULT 0,
    balance DECIMAL(10,4) DEFAULT 0,
    last_check_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_ai_providers_enabled ON ai_providers(is_enabled);
CREATE INDEX IF NOT EXISTS idx_ai_providers_deleted ON ai_providers(deleted_at);

-- AI chat history table
CREATE TABLE IF NOT EXISTS ai_chat_history (
    id SERIAL PRIMARY KEY,
    provider_id INT NOT NULL REFERENCES ai_providers(id),
    role VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_ai_chat_provider ON ai_chat_history(provider_id);
CREATE INDEX IF NOT EXISTS idx_ai_chat_created ON ai_chat_history(created_at);
