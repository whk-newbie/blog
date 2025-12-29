package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/iambaby/blog/internal/pkg/logger"
)

// Migration 迁移结构
type Migration struct {
	Version string
	Name    string
	SQL     string
}

// RunMigrations 运行数据库迁移
func RunMigrations(db *sql.DB, migrationsPath string) error {
	// 创建迁移记录表
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("创建迁移表失败: %w", err)
	}

	// 获取所有迁移文件
	migrations, err := loadMigrations(migrationsPath)
	if err != nil {
		return fmt.Errorf("加载迁移文件失败: %w", err)
	}

	// 获取已执行的迁移
	executedMigrations, err := getExecutedMigrations(db)
	if err != nil {
		return fmt.Errorf("获取已执行迁移失败: %w", err)
	}

	// 执行未执行的迁移
	for _, migration := range migrations {
		if _, executed := executedMigrations[migration.Version]; executed {
			logger.Info("迁移 %s 已执行，跳过", migration.Version)
			continue
		}

		logger.Info("执行迁移 %s: %s", migration.Version, migration.Name)

		// 开启事务
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("开启事务失败: %w", err)
		}

		// 执行迁移SQL
		if _, err := tx.Exec(migration.SQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("执行迁移 %s 失败: %w", migration.Version, err)
		}

		// 记录迁移
		if _, err := tx.Exec(
			"INSERT INTO schema_migrations (version, name) VALUES ($1, $2)",
			migration.Version, migration.Name,
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("记录迁移 %s 失败: %w", migration.Version, err)
		}

		// 提交事务
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("提交迁移 %s 失败: %w", migration.Version, err)
		}

		logger.Info("迁移 %s 执行成功", migration.Version)
	}

	logger.Info("所有迁移执行完成")
	return nil
}

// createMigrationsTable 创建迁移记录表
func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			version VARCHAR(50) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			executed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	return err
}

// loadMigrations 加载所有迁移文件
func loadMigrations(migrationsPath string) ([]Migration, error) {
	var migrations []Migration

	// 读取迁移目录
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("读取迁移目录失败: %w", err)
	}

	// 遍历所有.sql文件
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		// 解析文件名：001_init_schema.sql -> version=001, name=init_schema
		parts := strings.SplitN(file.Name(), "_", 2)
		if len(parts) != 2 {
			continue
		}

		version := parts[0]
		name := strings.TrimSuffix(parts[1], ".sql")

		// 读取SQL内容
		filePath := filepath.Join(migrationsPath, file.Name())
		sqlContent, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("读取迁移文件 %s 失败: %w", file.Name(), err)
		}

		migrations = append(migrations, Migration{
			Version: version,
			Name:    name,
			SQL:     string(sqlContent),
		})
	}

	// 按版本号排序
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// getExecutedMigrations 获取已执行的迁移
func getExecutedMigrations(db *sql.DB) (map[string]bool, error) {
	executed := make(map[string]bool)

	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		executed[version] = true
	}

	return executed, rows.Err()
}

