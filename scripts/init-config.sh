#!/bin/bash

# 配置初始化脚本
# 交互式配置向导，生成配置文件

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

CONFIG_FILE="backend/config/config.yaml"
EXAMPLE_FILE="backend/config/config.example.yaml"

echo -e "${BLUE}⚙️  配置初始化向导${NC}"
echo ""

# 检查示例文件
if [ ! -f "$EXAMPLE_FILE" ]; then
    echo -e "${RED}❌ 示例配置文件不存在: $EXAMPLE_FILE${NC}"
    exit 1
fi

# 如果配置文件已存在，询问是否覆盖
if [ -f "$CONFIG_FILE" ]; then
    echo -e "${YELLOW}⚠️  配置文件已存在: $CONFIG_FILE${NC}"
    read -p "是否覆盖现有配置? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}❌ 已取消${NC}"
        exit 0
    fi
fi

echo -e "${YELLOW}请填写以下配置信息（直接回车使用默认值）:${NC}"
echo ""

# 服务器配置
echo -e "${BLUE}📡 服务器配置${NC}"
read -p "服务器Host [0.0.0.0]: " SERVER_HOST
SERVER_HOST=${SERVER_HOST:-0.0.0.0}

read -p "服务器Port [8080]: " SERVER_PORT
SERVER_PORT=${SERVER_PORT:-8080}

read -p "运行模式 (debug/release) [debug]: " SERVER_MODE
SERVER_MODE=${SERVER_MODE:-debug}

# 数据库配置
echo ""
echo -e "${BLUE}🗄️  数据库配置${NC}"
read -p "数据库Host [postgres]: " DB_HOST
DB_HOST=${DB_HOST:-postgres}

read -p "数据库Port [5432]: " DB_PORT
DB_PORT=${DB_PORT:-5432}

read -p "数据库用户 [blog_user]: " DB_USER
DB_USER=${DB_USER:-blog_user}

read -p "数据库密码 [blog_password]: " DB_PASSWORD
DB_PASSWORD=${DB_PASSWORD:-blog_password}

read -p "数据库名称 [blog_db]: " DB_NAME
DB_NAME=${DB_NAME:-blog_db}

# Redis配置
echo ""
echo -e "${BLUE}📦 Redis配置${NC}"
read -p "Redis Host [redis]: " REDIS_HOST
REDIS_HOST=${REDIS_HOST:-redis}

read -p "Redis Port [6379]: " REDIS_PORT
REDIS_PORT=${REDIS_PORT:-6379}

read -p "Redis密码 (留空表示无密码): " REDIS_PASSWORD

# JWT配置
echo ""
echo -e "${BLUE}🔐 JWT配置${NC}"
read -p "JWT密钥 (留空自动生成32字符): " JWT_SECRET
if [ -z "$JWT_SECRET" ]; then
    JWT_SECRET=$(openssl rand -hex 16)
    echo -e "${GREEN}✅ 已自动生成JWT密钥${NC}"
fi

# 加密配置
echo ""
echo -e "${BLUE}🔒 加密配置${NC}"
read -p "主密钥 (留空自动生成32字节): " CRYPTO_MASTER_KEY
if [ -z "$CRYPTO_MASTER_KEY" ]; then
    CRYPTO_MASTER_KEY=$(openssl rand -hex 16)
    echo -e "${GREEN}✅ 已自动生成主密钥${NC}"
fi

# 确认配置
echo ""
echo -e "${YELLOW}📋 配置摘要:${NC}"
echo "  服务器: ${SERVER_HOST}:${SERVER_PORT} (${SERVER_MODE})"
echo "  数据库: ${DB_USER}@${DB_HOST}:${DB_PORT}/${DB_NAME}"
echo "  Redis: ${REDIS_HOST}:${REDIS_PORT}"
echo "  JWT密钥: ${JWT_SECRET:0:8}..."
echo "  主密钥: ${CRYPTO_MASTER_KEY:0:8}..."
echo ""

read -p "确认创建配置文件? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}❌ 已取消${NC}"
    exit 0
fi

# 生成配置文件
cat > "$CONFIG_FILE" <<EOF
server:
  host: ${SERVER_HOST}
  port: ${SERVER_PORT}
  mode: ${SERVER_MODE}
  read_timeout: 10s
  write_timeout: 10s

database:
  host: ${DB_HOST}
  port: ${DB_PORT}
  user: ${DB_USER}
  password: ${DB_PASSWORD}
  dbname: ${DB_NAME}
  sslmode: disable
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 1h

redis:
  host: ${REDIS_HOST}
  port: ${REDIS_PORT}
  password: "${REDIS_PASSWORD}"
  db: 0
  pool_size: 100
  min_idle_conns: 10
  max_retries: 3
  dial_timeout: 5s
  read_timeout: 3s
  write_timeout: 3s

jwt:
  secret: "${JWT_SECRET}"
  expire_time: 24h
  issuer: "blog-system"

crypto:
  master_key: "${CRYPTO_MASTER_KEY}"

upload:
  path: "./uploads"
  max_size: 10485760
  allowed_types:
    - "image/jpeg"
    - "image/png"
    - "image/gif"
    - "image/webp"
  compress_quality: 85

log:
  level: ${SERVER_MODE}
  format: json
  output: stdout
  file_path: "./logs/app.log"
  max_size: 100
  max_backups: 10
  max_age: 30
  compress: true

cors:
  allow_origins:
    - "http://localhost:5173"
    - "http://127.0.0.1:5173"
    - "http://localhost:3000"
    - "http://127.0.0.1:3000"
    - "http://localhost"
    - "http://127.0.0.1"
  allow_methods:
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
    - "PATCH"
  allow_headers:
    - "Origin"
    - "Content-Type"
    - "Accept"
    - "Authorization"
    - "X-Request-ID"
    - "X-Requested-With"
  expose_headers:
    - "Content-Length"
    - "Content-Type"
  allow_credentials: true
  max_age: 43200
EOF

echo -e "${GREEN}✅ 配置文件已创建: $CONFIG_FILE${NC}"
echo ""
echo -e "${YELLOW}⚠️  重要提示:${NC}"
echo "  1. 请妥善保管JWT密钥和主密钥"
echo "  2. 生产环境请修改默认密码"
echo "  3. 建议将配置文件添加到.gitignore"
echo ""

