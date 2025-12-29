-- 005_add_init_data.sql
-- 添加初始化数据

-- 初始化管理员（密码: admin@123）
-- 注意：SQL 可以创建管理员账号，但 Go 代码会检查并修正密码哈希
-- Go 代码会在应用启动时验证密码哈希是否正确，如果不正确会自动更新
-- 这样可以确保密码哈希是通过 Go 的 bcrypt 正确生成的（见 internal/pkg/db/init_admin.go）
INSERT INTO admins (username, password, email, is_default_password) 
VALUES ('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin@example.com', TRUE)
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

