# 部署文档

本文档介绍如何部署博客系统到生产环境。

## 📋 目录

- [前置要求](#前置要求)
- [快速部署](#快速部署)
- [配置说明](#配置说明)
- [SSL证书配置](#ssl证书配置)
- [故障排除](#故障排除)
- [维护操作](#维护操作)

## 前置要求

### 系统要求

- **操作系统**: Linux (推荐 Ubuntu 20.04+ 或 CentOS 7+)
- **Docker**: 20.10+
- **Docker Compose**: 2.0+
- **磁盘空间**: 至少 10GB 可用空间
- **内存**: 至少 2GB RAM
- **CPU**: 至少 2 核心

### 端口要求

确保以下端口未被占用：

- `80` - HTTP (前端)
- `443` - HTTPS (前端)
- `8080` - 后端API (开发环境)
- `5432` - PostgreSQL (开发环境)
- `6379` - Redis (开发环境)

## 快速部署

### 1. 克隆项目

```bash
git clone <repository-url>
cd blog
```

### 2. 初始化配置

运行配置初始化脚本：

```bash
chmod +x scripts/*.sh
./scripts/init-config.sh
```

脚本会引导你完成以下配置：
- 数据库密码
- Redis密码
- JWT密钥
- 加密主密钥
- 域名配置（用于SSL）

### 3. 启动服务

```bash
./scripts/start.sh
```

### 4. 验证部署

访问以下地址验证服务是否正常：

- 前端: http://your-domain
- 健康检查: http://your-domain/health
- API文档: http://your-domain/swagger/index.html

### 5. 首次登录

- 用户名: `admin`
- 密码: `admin@123`
- **⚠️ 首次登录后请立即修改密码！**

## 配置说明

### 后端配置

配置文件位置: `backend/config/config.yaml`

主要配置项：

```yaml
server:
  port: 8080
  mode: release  # release 或 debug

database:
  host: postgres
  port: 5432
  user: blog_user
  password: your_password
  dbname: blog_db

redis:
  host: redis
  port: 6379
  password: your_redis_password

jwt:
  secret: your_jwt_secret
  expire_time: 24h
  issuer: blog

crypto:
  master_key: your_master_key
```

### 环境变量

可以通过环境变量覆盖配置：

```bash
export DB_PASSWORD=your_password
export REDIS_PASSWORD=your_redis_password
export JWT_SECRET=your_jwt_secret
export CRYPTO_MASTER_KEY=your_master_key
```

### Nginx配置

Nginx配置文件位置: `docker/nginx/nginx.conf`

主要配置项：
- 反向代理到后端API
- 静态文件服务
- SSL配置
- 安全头设置

## SSL证书配置

### 使用Let's Encrypt（推荐）

1. 运行SSL配置脚本：

```bash
./scripts/setup-ssl.sh
```

2. 选择证书类型：
   - `new` - 申请新证书
   - `renew` - 续期现有证书

3. 输入域名和邮箱：

```bash
Domain: example.com
Email: admin@example.com
```

4. 脚本会自动：
   - 安装Certbot
   - 申请证书
   - 配置Nginx
   - 设置自动续期

### 手动配置SSL

1. 将证书文件放到 `docker/nginx/ssl/` 目录：
   - `cert.pem` - 证书文件
   - `key.pem` - 私钥文件

2. 更新 `docker/nginx/nginx.conf` 中的SSL配置

3. 重启Nginx服务：

```bash
docker-compose restart nginx
```

### SSL证书续期

Let's Encrypt证书有效期为90天，需要定期续期：

```bash
./scripts/setup-ssl.sh renew
```

建议设置定时任务自动续期：

```bash
# 添加到 crontab
0 0 1 * * /path/to/blog/scripts/setup-ssl.sh renew
```

## 故障排除

### 服务无法启动

1. 检查端口占用：

```bash
netstat -tulpn | grep -E '80|443|8080|5432|6379'
```

2. 检查Docker服务：

```bash
docker ps
docker-compose ps
```

3. 查看日志：

```bash
./scripts/logs.sh
./scripts/logs.sh backend  # 查看特定服务日志
```

### 数据库连接失败

1. 检查PostgreSQL容器状态：

```bash
docker-compose ps postgres
docker-compose logs postgres
```

2. 验证数据库配置：

```bash
docker-compose exec postgres psql -U blog_user -d blog_db
```

3. 检查网络连接：

```bash
docker-compose exec backend ping postgres
```

### Redis连接失败

1. 检查Redis容器状态：

```bash
docker-compose ps redis
docker-compose logs redis
```

2. 测试Redis连接：

```bash
docker-compose exec redis redis-cli -a your_password ping
```

### 前端无法访问后端API

1. 检查Nginx配置：

```bash
docker-compose exec nginx nginx -t
```

2. 检查后端服务：

```bash
curl http://localhost:8080/health
```

3. 查看Nginx日志：

```bash
docker-compose logs nginx
```

### 文件上传失败

1. 检查上传目录权限：

```bash
ls -la backend/uploads/
```

2. 确保目录可写：

```bash
chmod -R 755 backend/uploads/
```

### 备份失败

1. 检查备份目录权限：

```bash
ls -la backend/backups/
```

2. 检查磁盘空间：

```bash
df -h
```

3. 手动测试备份：

```bash
./scripts/backup-db.sh
```

## 维护操作

### 数据备份

#### 手动备份

```bash
./scripts/backup-db.sh
```

备份文件保存在 `backend/backups/` 目录，格式：`backup_YYYYMMDD_HHMMSS.sql.gz`

#### 自动备份

系统已配置自动备份调度器，每天凌晨3点自动备份，保留最近10个备份。

#### 恢复备份

```bash
./scripts/restore-db.sh backend/backups/backup_20240101_120000.sql.gz
```

### 日志管理

#### 查看日志

```bash
# 查看所有服务日志
./scripts/logs.sh

# 查看特定服务日志
./scripts/logs.sh backend
./scripts/logs.sh frontend
./scripts/logs.sh postgres
./scripts/logs.sh redis
```

#### 清理日志

系统日志会自动清理90天前的记录。也可以手动清理：

1. 登录管理后台
2. 进入"日志管理"页面
3. 点击"清理旧日志"按钮

### 更新服务

1. 拉取最新代码：

```bash
git pull origin main
```

2. 重新构建镜像：

```bash
docker-compose build --no-cache
```

3. 重启服务：

```bash
docker-compose down
docker-compose up -d
```

### 健康检查

运行健康检查脚本：

```bash
./scripts/health-check.sh
```

脚本会检查：
- 所有服务是否运行
- 数据库连接
- Redis连接
- API健康状态

### 清理资源

⚠️ **谨慎使用** - 会删除所有容器和数据卷

```bash
./scripts/clean.sh
```

## 性能优化建议

### 数据库优化

1. 定期执行VACUUM：

```bash
docker-compose exec postgres psql -U blog_user -d blog_db -c "VACUUM ANALYZE;"
```

2. 监控慢查询：

```bash
docker-compose exec postgres psql -U blog_user -d blog_db -c "SELECT * FROM pg_stat_statements ORDER BY total_time DESC LIMIT 10;"
```

### Redis优化

1. 监控内存使用：

```bash
docker-compose exec redis redis-cli INFO memory
```

2. 清理过期键：

```bash
docker-compose exec redis redis-cli --scan --pattern "*" | xargs redis-cli DEL
```

### 系统监控

建议使用监控工具（如Prometheus + Grafana）监控：
- CPU使用率
- 内存使用率
- 磁盘I/O
- 网络流量
- 数据库连接数
- Redis连接数

## 安全建议

1. **修改默认密码**
   - 数据库密码
   - Redis密码
   - JWT密钥
   - 加密主密钥

2. **使用HTTPS**
   - 配置SSL证书
   - 强制HTTPS重定向

3. **防火墙配置**
   - 只开放必要端口（80, 443）
   - 限制管理端口访问

4. **定期更新**
   - 定期更新Docker镜像
   - 定期更新系统依赖

5. **备份策略**
   - 定期备份数据库
   - 备份上传的文件
   - 测试备份恢复流程

## 常见问题

### Q: 如何修改管理员密码？

A: 登录管理后台，进入"修改密码"页面修改。

### Q: 如何添加新的爬虫Token？

A: 登录管理后台，进入"配置管理"页面，点击"生成爬虫Token"。

### Q: 如何查看系统日志？

A: 登录管理后台，进入"日志管理"页面查看。

### Q: 如何配置邮件通知？

A: 在"配置管理"页面添加邮箱配置（类型：email）。

### Q: 如何限制访问频率？

A: 系统已内置限流功能，非登录用户每分钟限制60次请求。

## 获取帮助

如果遇到问题，可以：

1. 查看日志文件
2. 检查健康状态
3. 查看GitHub Issues
4. 联系技术支持

