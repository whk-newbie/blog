#!/bin/bash

# 停止博客系统（生产环境）

set -e

echo "🛑 Stopping Blog System (Production)..."

docker-compose down

echo "✅ Blog System stopped successfully!"

