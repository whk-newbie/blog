package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/config"
	"github.com/whk-newbie/blog/internal/pkg/db"
	"github.com/whk-newbie/blog/internal/pkg/logger"
	"github.com/whk-newbie/blog/internal/pkg/redis"
	"github.com/whk-newbie/blog/internal/router"

	_ "github.com/whk-newbie/blog/docs" // Swagger文档
)

// @title Blog API
// @version 1.0
// @description 个人博客系统API文档
// @termsOfService https://github.com/whk-newbie/blog

// @contact.name API Support
// @contact.url https://github.com/whk-newbie/blog/issues
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	logger.Init(logger.LogConfig{
		Level:      cfg.Log.Level,
		Format:     cfg.Log.Format,
		Output:     cfg.Log.Output,
		FilePath:   cfg.Log.FilePath,
		MaxSize:    cfg.Log.MaxSize,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAge,
		Compress:   cfg.Log.Compress,
	})

	// 初始化数据库
	if err := db.Init(db.DatabaseConfig{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		DBName:          cfg.Database.DBName,
		SSLMode:         cfg.Database.SSLMode,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	}); err != nil {
		logger.Fatal("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// 运行数据库迁移
	gormDB, err := db.GetSQLDB()
	if err != nil {
		logger.Fatal("Failed to get database instance: %v", err)
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		logger.Fatal("Failed to get sql.DB: %v", err)
	}
	if err := db.RunMigrations(sqlDB, "./migrations"); err != nil {
		logger.Fatal("Failed to run migrations: %v", err)
	}

	// 初始化默认管理员
	if err := db.InitDefaultAdmin(gormDB); err != nil {
		logger.Fatal("Failed to initialize default admin: %v", err)
	}

	// 初始化Redis
	if err := redis.Init(redis.RedisConfig{
		Host:         cfg.Redis.Host,
		Port:         cfg.Redis.Port,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		MaxRetries:   cfg.Redis.MaxRetries,
		DialTimeout:  cfg.Redis.DialTimeout,
		ReadTimeout:  cfg.Redis.ReadTimeout,
		WriteTimeout: cfg.Redis.WriteTimeout,
	}); err != nil {
		logger.Fatal("Failed to initialize redis: %v", err)
	}
	defer redis.Close()

	// 设置Gin模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化路由
	r, schedulerManager := router.Setup(cfg)

	// 启动调度器
	if err := schedulerManager.Start(); err != nil {
		logger.Fatal("Failed to start scheduler: %v", err)
	}
	logger.Info("Scheduler started successfully")

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("Starting server on %s", addr)

	// 优雅关闭
	go func() {
		if err := r.Run(addr); err != nil {
			logger.Fatal("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 停止调度器
	schedulerManager.Stop()
	logger.Info("Scheduler stopped")
}
