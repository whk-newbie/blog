# 个人博客系统

一个功能完整的现代化博客系统，采用前后端分离架构，支持内容加密、访客追踪、爬虫监控等高级功能。

## 🎯 项目特点

- **前后端分离**: Go后端 + Vue 3前端
- **容器化部署**: 完整的Docker支持，一键启动
- **JWT认证**: 基于Token的安全认证机制
- **内容加密**: 支持文章内容AES-256-GCM加密
- **访客追踪**: 基于浏览器指纹的访客识别
- **爬虫监控**: Python SDK支持内容爬取和监控
- **实时通信**: WebSocket支持实时消息推送
- **高性能**: Redis缓存 + PostgreSQL数据库

## 🏗️ 技术栈

### 后端
- **语言**: Go 1.21+
- **框架**: Gin
- **数据库**: PostgreSQL 15+
- **缓存**: Redis 7+
- **ORM**: GORM

### 前端
- **框架**: Vue 3
- **构建工具**: Vite 5
- **UI组件**: Element Plus
- **状态管理**: Pinia
- **路由**: Vue Router

### Python SDK
- **语言**: Python 3.10+
- **包管理**: uv
- **爬虫**: BeautifulSoup + Requests

### 基础设施
- **容器化**: Docker + Docker Compose
- **反向代理**: Nginx
- **CI/CD**: GitHub Actions（规划中）

## 📁 项目结构

```
blog/
├── backend/            # Go后端服务
├── frontend/           # Vue 3前端应用
├── python-sdk/         # Python爬虫监控SDK
├── docker/             # Docker配置文件
├── scripts/            # 运维脚本
├── docs/               # 项目文档
├── docker-compose.yml  # 生产环境配置
└── docker-compose.dev.yml  # 开发环境配置
```

## 🚀 快速开始

### 前置要求

- Docker 20.10+
- Docker Compose 2.0+

### 生产环境部署

1. 克隆项目
```bash
git clone <repository-url>
cd blog
```

2. 配置环境变量（可选）
```bash
cp backend/config/config.example.yaml backend/config/config.yaml
# 编辑配置文件，修改数据库密码、JWT密钥等敏感信息
```

3. 启动服务
```bash
chmod +x scripts/*.sh
./scripts/start.sh
```

4. 访问应用
- 前端: http://localhost
- 后端API: http://localhost/api
- 健康检查: http://localhost/health
- API文档: http://localhost/swagger/index.html

5. 默认管理员账号
- 用户名: `admin`
- 密码: `admin@123`
- **⚠️ 首次登录后请立即修改密码！**

### 开发环境

开发环境支持热重载，代码修改后自动重启服务。

```bash
./scripts/start-dev.sh
```

访问地址：
- 前端: http://localhost:5173
- 后端: http://localhost:8080
- API文档: http://localhost:8080/swagger/index.html
- PostgreSQL: localhost:5432
- Redis: localhost:6379

默认管理员账号：
- 用户名: `admin`
- 密码: `admin@123`

## 📝 常用命令

### 服务管理

```bash
# 启动服务（生产）
./scripts/start.sh

# 启动服务（开发）
./scripts/start-dev.sh

# 停止服务
./scripts/stop.sh  # 生产环境
./scripts/stop-dev.sh  # 开发环境

# 查看日志
./scripts/logs.sh          # 所有服务
./scripts/logs.sh backend  # 指定服务
```

### 数据库管理

```bash
# 备份数据库
./scripts/backup-db.sh

# 恢复数据库
./scripts/restore-db.sh backups/backup_20240101_120000.sql.gz
```

### 清理资源

```bash
# 清理所有Docker资源（谨慎使用）
./scripts/clean.sh
```

### Docker Compose 命令

```bash
# 查看服务状态
docker-compose ps

# 重启特定服务
docker-compose restart backend

# 查看服务日志
docker-compose logs -f backend

# 进入容器
docker-compose exec backend sh
docker-compose exec postgres psql -U blog_user -d blog_db

# 重新构建镜像
docker-compose build --no-cache

# 停止并删除所有容器
docker-compose down -v
```

## ✨ 功能特性

### 已实现功能

#### 认证系统 ✅
- [x] JWT Token认证
- [x] 管理员登录
- [x] 密码修改（BCrypt加密）
- [x] Token刷新
- [x] 路由守卫
- [x] 默认管理员初始化

#### 文章管理 ✅
- [x] 文章CRUD（创建、读取、更新、删除）
- [x] 分类管理（CRUD）
- [x] 标签管理（CRUD）
- [x] 图片上传（支持压缩和WebP格式）
- [x] 富文本编辑器（Quill编辑器）
- [x] 文章搜索（全文搜索）
- [x] 文章发布/取消发布
- [x] 定时发布调度器

#### 访问统计 ✅
- [x] 浏览器指纹识别（Canvas、WebGL、屏幕、时区、语言等）
- [x] 访问记录（PV/UV统计）
- [x] 访问趋势分析（按日/周/月统计）
- [x] 来源分析（Referrer统计）
- [x] 热门文章统计
- [x] 访问统计可视化（ECharts图表）

#### 爬虫监控 ✅
- [x] 爬虫任务注册和管理
- [x] 任务状态实时更新（WebSocket推送）
- [x] 任务进度跟踪
- [x] Python SDK支持（TaskReporter）
- [x] Bearer Token认证

#### 系统管理 ✅
- [x] 系统配置管理（CRUD、加密存储）
- [x] 爬虫Token生成
- [x] 系统日志管理（查询、清理）
- [x] 数据备份（手动/自动备份）
- [x] 备份文件管理（下载、删除）

#### 工具集合 ✅
- [x] JSON格式化工具
- [x] Header格式化工具
- [x] Cookie格式化工具
- [x] Dict格式化工具
- [x] cURL转Python工具
- [x] Python转cURL工具
- [x] URL格式化工具
- [x] 加密解密工具集（URL编码、Base64、MD5）

#### 仪表盘 ✅
- [x] 数据概览（文章、分类、标签统计）
- [x] 访问趋势图表
- [x] 文章发布趋势图表
- [x] 最近文章列表

#### 安全功能 ✅
- [x] 数据加密传输（AES-256-GCM）
- [x] IP黑名单
- [x] 请求限流
- [x] XSS防护
- [x] SQL注入防护
- [x] CORS配置

#### 性能优化 ✅
- [x] Redis缓存（文章列表、访问统计、配置信息）
- [x] 数据库索引优化
- [x] 图片压缩（自动压缩超过1MB的图片）
- [x] 代码分割和懒加载
- [x] 前端代码混淆

## 🔧 配置说明

### 后端配置

配置文件位于 `backend/config/config.yaml`，主要配置项：

- **server**: 服务器配置（端口、模式等）
- **database**: PostgreSQL数据库配置
- **redis**: Redis缓存配置
- **jwt**: JWT认证配置（Token过期时间、签名密钥等）
- **crypto**: 加密配置（AES-256密钥）
- **upload**: 文件上传配置
- **log**: 日志配置
- **cors**: 跨域配置

### 环境变量

可通过环境变量覆盖配置文件：

- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`
- `JWT_SECRET`
- `CRYPTO_MASTER_KEY`

## 📚 开发文档

- [需求文档](./PRD.md)
- [架构设计](./ARCHITECTURE.md)
- [数据库设计](./DATABASE_DESIGN.md)
- [API设计](./API_DESIGN.md)
- [前端设计](./FRONTEND_DESIGN.md)
- [项目结构](./PROJECT_STRUCTURE.md)
- [开发计划](./DEVELOPMENT_PLAN.md)

## 🧪 测试

### 后端测试

```bash
cd backend
go test ./...
```

### 前端测试

```bash
cd frontend
npm run test
```

## 📦 构建

### 后端构建

```bash
cd backend
go build -o blog-server ./cmd/server
```

### 前端构建

```bash
cd frontend
npm run build
```

## 🔐 安全建议

1. **生产环境务必修改默认密钥**
   - JWT_SECRET
   - CRYPTO_MASTER_KEY
   - 数据库密码
   - Redis密码

2. **使用HTTPS**
   - 配置SSL证书
   - 强制HTTPS重定向

3. **定期备份**
   - 定期备份数据库
   - 备份上传文件

4. **限制访问**
   - 配置防火墙规则
   - 使用IP白名单

## 🤝 贡献

欢迎提交Issue和Pull Request！

## 📄 许可证

MIT License

## 👤 作者

iambaby

## 🙏 致谢

感谢所有开源项目的贡献者！

