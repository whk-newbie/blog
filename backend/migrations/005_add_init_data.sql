-- 005_add_init_data.sql
-- 添加初始化数据

-- 初始化管理员（密码: admin@123）
-- BCrypt hash of 'admin@123' with cost 10
INSERT INTO admins (username, password, is_default_password) 
VALUES ('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', TRUE)
ON CONFLICT (username) DO NOTHING;

-- 初始化默认分类
INSERT INTO categories (name, slug, description, sort_order, created_by) VALUES
('技术', 'tech', '技术相关文章', 100, 1),
('生活', 'life', '生活随笔', 90, 1),
('读书', 'reading', '读书笔记', 80, 1)
ON CONFLICT (slug) DO NOTHING;

-- 初始化默认标签
INSERT INTO tags (name, slug, created_by) VALUES
('Go', 'go', 1),
('Vue', 'vue', 1),
('PostgreSQL', 'postgresql', 1),
('Docker', 'docker', 1),
('Python', 'python', 1),
('Redis', 'redis', 1)
ON CONFLICT (slug) DO NOTHING;

-- 初始化默认配置
INSERT INTO system_configs (config_key, config_value, config_type, is_encrypted, description) VALUES
('encryption_salt', 'change_this_salt_in_production', 'salt', FALSE, '加密盐（建议修改）')
ON CONFLICT (config_key) DO NOTHING;

