# 数据持久化目录

此目录用于存储数据库和缓存的持久化数据。

## 目录结构

```
data/
├── postgres/    # PostgreSQL 数据文件
└── redis/       # Redis 数据文件
```

## 说明

- **postgres/**: 存储 PostgreSQL 数据库的所有数据
- **redis/**: 存储 Redis 的持久化数据（RDB/AOF）

## 备份建议

建议定期备份此目录，或使用以下命令进行数据库备份：

```bash
# 备份数据库
./scripts/backup-db.sh

# 恢复数据库
./scripts/restore-db.sh backups/backup_YYYYMMDD_HHMMSS.sql.gz
```

## 清理数据

如果需要清空所有数据重新开始：

```bash
# 停止服务
./scripts/stop-dev.sh

# 删除数据
rm -rf data/postgres/* data/redis/*

# 重新启动（会自动初始化）
./scripts/start-dev.sh
```

## 注意事项

⚠️ 请勿手动修改此目录下的文件，以免造成数据损坏。
