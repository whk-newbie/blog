-- Insert default admin path and visibility configs
INSERT INTO system_configs (config_key, config_value, config_type, is_encrypted, is_active, description, created_at, updated_at)
VALUES
('admin_path', 'admin', 'site_config', false, true, '后台访问路径', NOW(), NOW()),
('show_admin_link', 'true', 'site_config', false, true, '是否在前端显示后台入口链接', NOW(), NOW())
ON CONFLICT (config_key) DO NOTHING;
