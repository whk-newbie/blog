#!/bin/bash

# 启动博客系统（开发环境）

set -e

echo "🚀 Starting Blog System (Development)..."

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
fi

# 启动开发环境
echo "🐳 Starting Docker containers (development mode)..."
docker-compose -f docker-compose.dev.yml up -d

echo ""
echo "✅ Blog System (Dev) started successfully!"
echo ""
echo "📊 Service Status:"
docker-compose -f docker-compose.dev.yml ps
echo ""
echo "🌐 Access URLs:"
echo "  - Frontend (Dev): http://localhost:5173"
echo "  - Backend API (Dev): http://localhost:8080/api/v1"
echo "  - Backend Health: http://localhost:8080/health"
echo "  - API Docs (Swagger): http://localhost:8080/swagger/index.html"
echo "  - PostgreSQL: localhost:5432"
echo "  - Redis: localhost:6379"
echo ""
echo "📝 View logs:"
echo "  docker-compose -f docker-compose.dev.yml logs -f"
echo ""
echo "🛑 Stop services:"
echo "  ./scripts/stop-dev.sh"
echo ""
echo "💡 Hot reload enabled for both frontend and backend!"

