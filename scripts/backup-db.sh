#!/bin/bash

# 备份数据库

set -e

BACKUP_DIR="./backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/backup_$TIMESTAMP.sql"

echo "💾 Starting database backup..."

# 创建备份目录
mkdir -p "$BACKUP_DIR"

# 执行备份
docker exec blog-postgres pg_dump -U blog_user blog_db > "$BACKUP_FILE"

# 压缩备份
gzip "$BACKUP_FILE"

echo "✅ Database backup completed: ${BACKUP_FILE}.gz"
echo "📦 Backup size: $(du -h ${BACKUP_FILE}.gz | cut -f1)"

