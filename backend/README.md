# Blog 后端服务

基于 Go + Gin + PostgreSQL + Redis 的博客系统后端服务。

## 技术栈

- **框架**: Gin Web Framework
- **数据库**: PostgreSQL 15+
- **缓存**: Redis 7+
- **ORM**: GORM
- **日志**: Logrus
- **配置**: Viper

## 目录结构

```
backend/
├── cmd/server/          # 应用入口
├── internal/            # 内部代码
│   ├── config/         # 配置管理
│   ├── models/         # 数据模型
│   ├── repository/     # 数据访问层
│   ├── service/        # 业务逻辑层
│   ├── handler/        # HTTP处理器
│   ├── middleware/     # 中间件
│   ├── router/         # 路由配置
│   └── pkg/            # 工具包
├── migrations/         # 数据库迁移
├── config/             # 配置文件
└── Dockerfile          # Docker镜像
```

## 快速开始

### 使用 Docker Compose (推荐)

```bash
# 在项目根目录执行
docker-compose up -d
```

### 本地开发

1. 安装依赖
```bash
go mod download
```

2. 配置文件
```bash
cp config/config.example.yaml config/config.yaml
# 修改 config/config.yaml 中的配置
```

3. 运行服务
```bash
go run cmd/server/main.go
```

## 环境变量

- `CONFIG_PATH`: 配置文件路径（默认: config/config.yaml）
- `DB_HOST`: 数据库主机
- `DB_PORT`: 数据库端口
- `DB_USER`: 数据库用户
- `DB_PASSWORD`: 数据库密码
- `DB_NAME`: 数据库名称
- `REDIS_HOST`: Redis主机
- `REDIS_PORT`: Redis端口
- `REDIS_PASSWORD`: Redis密码
- `JWT_SECRET`: JWT密钥
- `CRYPTO_MASTER_KEY`: 加密主密钥（32字节）

## API 文档

服务启动后访问：
- 健康检查: http://localhost:8080/health
- Swagger文档: http://localhost:8080/swagger/index.html

### 生成Swagger文档

```bash
# 安装swag工具
make install-swag

# 生成文档
make swag
```

## 开发工具

### 热重载 (Air)

```bash
# 安装 Air
go install github.com/cosmtrek/air@latest

# 运行
air
```

## 测试

```bash
# 运行所有测试
go test ./...

# 运行测试并显示覆盖率
go test -cover ./...
```

## 构建

```bash
# 构建二进制文件
go build -o blog-server ./cmd/server

# 构建 Docker 镜像
docker build -t blog-backend .
```

## 许可证

MIT

