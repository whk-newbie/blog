#!/bin/bash

# 启动博客系统（生产环境）

set -e

echo "🚀 Starting Blog System (Production)..."

# 检查Docker和Docker Compose
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed"
    exit 1
fi

if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo "❌ Docker Compose is not installed"
    exit 1
fi

# 创建配置文件（如果不存在）
if [ ! -f "backend/config/config.yaml" ]; then
    echo "📝 Creating config file..."
    cp backend/config/config.example.yaml backend/config/config.yaml
    echo "⚠️  Please edit backend/config/config.yaml before starting"
    exit 0
fi

# 启动服务
echo "🐳 Starting Docker containers..."
docker-compose up -d

echo ""
echo "✅ Blog System started successfully!"
echo ""
echo "🌐 Access URLs:"
echo "  - Frontend: http://localhost"
echo "  - Backend API: http://localhost/api"
echo "  - Backend Health: http://localhost/health"
echo "  - API Docs (Swagger): http://localhost/swagger/index.html"
echo ""
echo "📝 View logs:"
echo "  docker-compose logs -f"
echo ""
echo "🛑 Stop services:"
echo "  ./scripts/stop.sh"

