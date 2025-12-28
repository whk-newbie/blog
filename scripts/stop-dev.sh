#!/bin/bash

# 停止博客系统（开发环境）

set -e

echo "🛑 Stopping Blog System (Development)..."

docker-compose -f docker-compose.dev.yml down

echo "✅ Blog System (Dev) stopped successfully!"

