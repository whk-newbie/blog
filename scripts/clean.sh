#!/bin/bash

# 清理Docker资源

set -e

echo "🧹 Cleaning Docker resources..."

read -p "⚠️  This will remove all containers, volumes and images. Continue? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Cancelled"
    exit 1
fi

# 停止所有容器
echo "Stopping containers..."
docker-compose down
docker-compose -f docker-compose.dev.yml down

# 删除数据卷
echo "Removing volumes..."
docker volume rm blog_postgres_data blog_redis_data blog_backend_uploads blog_backend_backups blog_backend_logs 2>/dev/null || true
docker volume rm blog_postgres_dev_data blog_redis_dev_data blog_go_mod_cache 2>/dev/null || true

# 清理未使用的资源
echo "Pruning unused resources..."
docker system prune -f

echo "✅ Cleanup completed!"

