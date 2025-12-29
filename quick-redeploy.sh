#!/bin/bash

# 快速重新部署脚本（不清理缓存，更快）
# 用途：停止服务，重新构建并启动

set -e

echo "======================================"
echo "  快速重新部署"
echo "======================================"

# 停止服务
echo "[1/4] 停止服务..."
docker compose -f docker-compose.dev.yml down

# 清理旧容器
echo "[2/4] 清理旧容器..."
docker rm -f blog-backend-dev blog-frontend-dev blog-postgres-dev blog-redis-dev 2>/dev/null || true

# 重新构建（使用缓存）
echo "[3/4] 重新构建服务..."
docker compose -f docker-compose.dev.yml build

# 启动服务
echo "[4/4] 启动服务..."
docker compose -f docker-compose.dev.yml up -d

echo ""
echo "等待服务启动..."
sleep 5

# 检查服务状态
echo ""
echo "服务状态："
docker compose -f docker-compose.dev.yml ps

echo ""
echo "======================================"
echo "  部署完成！"
echo "======================================"
echo ""
echo "访问地址："
echo "  前端：http://localhost:5173"
echo "  后端：http://localhost:8080"
echo "  Swagger：http://localhost:8080/swagger/index.html"
echo ""
echo "查看日志："
echo "  docker logs -f blog-backend-dev"
echo "  docker logs -f blog-frontend-dev"
echo ""

