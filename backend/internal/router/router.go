package router

import (
	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/config"
	"github.com/whk-newbie/blog/internal/handler"
	"github.com/whk-newbie/blog/internal/middleware"
	"github.com/whk-newbie/blog/internal/pkg/db"
	"github.com/whk-newbie/blog/internal/pkg/jwt"
	"github.com/whk-newbie/blog/internal/repository"
	"github.com/whk-newbie/blog/internal/scheduler"
	"github.com/whk-newbie/blog/internal/service"
	"github.com/whk-newbie/blog/internal/websocket"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Setup 设置路由
func Setup(cfg *config.Config) (*gin.Engine, *scheduler.Manager) {
	r := gin.New()

	// 使用中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS(
		cfg.CORS.AllowOrigins,
		cfg.CORS.AllowMethods,
		cfg.CORS.AllowHeaders,
		cfg.CORS.ExposeHeaders,
		cfg.CORS.AllowCredentials,
		cfg.CORS.MaxAge,
	))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Blog API is running",
		})
	})

	// Swagger文档 - 使用相对路径指向swagger.json
	url := ginSwagger.URL("/docs/swagger.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// 直接提供docs目录的静态文件服务
	r.Static("/docs", "./docs")

	// 初始化依赖
	gormDB, _ := db.GetSQLDB()
	jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.ExpireTime, cfg.JWT.Issuer)

	// 初始化Repository
	adminRepo := repository.NewAdminRepository(gormDB)
	categoryRepo := repository.NewCategoryRepository(gormDB)
	tagRepo := repository.NewTagRepository(gormDB)
	articleRepo := repository.NewArticleRepository(gormDB)
	fingerprintRepo := repository.NewFingerprintRepository(gormDB)
	visitRepo := repository.NewVisitRepository(gormDB)
	crawlTaskRepo := repository.NewCrawlTaskRepository(gormDB)
	configRepo := repository.NewConfigRepository(gormDB)
	logRepo := repository.NewLogRepository(gormDB)

	// 初始化Service
	authService := service.NewAuthService(adminRepo, jwtManager, cfg.JWT.ExpireTime)
	categoryService := service.NewCategoryService(categoryRepo)
	tagService := service.NewTagService(tagRepo)
	articleService := service.NewArticleService(articleRepo, categoryRepo, tagRepo)

	// 访问统计相关服务
	visitCacheService := service.NewVisitCacheService()
	visitService := service.NewVisitService(visitRepo, visitCacheService)
	fingerprintService := service.NewFingerprintService(fingerprintRepo)
	statsService := service.NewStatsService(articleRepo, categoryRepo, tagRepo, visitService)

	// 初始化WebSocket Hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// 初始化爬虫任务服务（需要Hub）
	crawlService := service.NewCrawlService(crawlTaskRepo, wsHub)

	// 初始化配置和日志服务
	configService, err := service.NewConfigService(configRepo, cfg.Crypto.MasterKey)
	if err != nil {
		panic("Failed to initialize config service: " + err.Error())
	}
	logService := service.NewLogService(logRepo)

	// 初始化Handler
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	tagHandler := handler.NewTagHandler(tagService)
	articleHandler := handler.NewArticleHandler(articleService)
	uploadHandler := handler.NewUploadHandler("uploads", 10) // 10MB max size
	statsHandler := handler.NewStatsHandler(statsService)
	fingerprintHandler := handler.NewFingerprintHandler(fingerprintService)
	visitHandler := handler.NewVisitHandler(visitService)
	crawlerHandler := handler.NewCrawlerHandler(crawlService)
	configHandler := handler.NewConfigHandler(configService)
	logHandler := handler.NewLogHandler(logService)

	// 初始化WebSocket Handler
	wsHandler := handler.NewWebSocketHandler(wsHub, jwtManager)

	// API路由组
	api := r.Group("/api/v1")
	{
		// 认证相关接口（公开）
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// 认证相关接口（需要认证）
		authProtected := api.Group("/auth")
		authProtected.Use(middleware.Auth(jwtManager))
		{
			authProtected.GET("/verify", authHandler.VerifyToken)
			authProtected.PUT("/password", authHandler.ChangePassword)
		}

		// 公开接口 - 分类
		api.GET("/categories", categoryHandler.List)
		api.GET("/categories/:id", categoryHandler.GetByID)
		api.GET("/categories/slug/:slug", categoryHandler.GetBySlug)

		// 公开接口 - 标签
		api.GET("/tags", tagHandler.List)
		api.GET("/tags/:id", tagHandler.GetByID)
		api.GET("/tags/slug/:slug", tagHandler.GetBySlug)

		// 公开接口 - 文章
		api.GET("/articles", articleHandler.ListPublished)
		api.GET("/articles/:id", articleHandler.GetByID)
		api.GET("/articles/slug/:slug", articleHandler.GetBySlug)
		api.GET("/articles/search", articleHandler.Search)

		// 公开接口 - 指纹和访问统计
		api.POST("/fingerprint", fingerprintHandler.CollectFingerprint)
		api.POST("/visit", visitHandler.RecordVisit)

		// 爬虫任务接口（需要Bearer Token认证）
		crawler := api.Group("/crawler")
		crawler.Use(middleware.CrawlerAuth())
		{
			crawler.POST("/tasks", crawlerHandler.RegisterTask)
			crawler.PUT("/tasks/:id", crawlerHandler.UpdateTaskStatus)
			crawler.PUT("/tasks/:id/complete", crawlerHandler.CompleteTask)
			crawler.PUT("/tasks/:id/fail", crawlerHandler.FailTask)
		}

		// 管理接口（需要认证）
		admin := api.Group("/admin")
		admin.Use(middleware.Auth(jwtManager))
		{
			// 分类管理
			admin.POST("/categories", categoryHandler.Create)
			admin.PUT("/categories/:id", categoryHandler.Update)
			admin.DELETE("/categories/:id", categoryHandler.Delete)

			// 标签管理
			admin.POST("/tags", tagHandler.Create)
			admin.PUT("/tags/:id", tagHandler.Update)
			admin.DELETE("/tags/:id", tagHandler.Delete)

			// 文章管理
			admin.GET("/articles", articleHandler.List)
			admin.POST("/articles", articleHandler.Create)
			admin.PUT("/articles/:id", articleHandler.Update)
			admin.DELETE("/articles/:id", articleHandler.Delete)
			admin.POST("/articles/:id/publish", articleHandler.Publish)
			admin.POST("/articles/:id/unpublish", articleHandler.Unpublish)

			// 文件上传
			admin.POST("/upload/image", uploadHandler.UploadImage)
			admin.POST("/upload/article-image", uploadHandler.UploadArticleImage)

			// 统计数据
			admin.GET("/stats/dashboard", statsHandler.GetDashboardStats)
			admin.GET("/stats/visits", statsHandler.GetVisitStats)
			admin.GET("/stats/popular-articles", statsHandler.GetPopularArticles)
			admin.GET("/stats/referrers", statsHandler.GetReferrerStats)

			// 指纹管理
			admin.GET("/fingerprints", fingerprintHandler.ListFingerprints)
			admin.GET("/fingerprints/:id", fingerprintHandler.GetFingerprint)
			admin.PUT("/fingerprints/:id", fingerprintHandler.UpdateFingerprint)
			admin.DELETE("/fingerprints/:id", fingerprintHandler.DeleteFingerprint)

			// 爬虫任务管理
			admin.GET("/crawler/tasks", crawlerHandler.ListTasks)
			admin.GET("/crawler/tasks/:task_id", crawlerHandler.GetTaskByID)

			// 配置管理
			admin.GET("/configs", configHandler.GetConfigs)
			admin.GET("/configs/:id", configHandler.GetConfigByID)
			admin.POST("/configs", configHandler.CreateConfig)
			admin.PUT("/configs/:id", configHandler.UpdateConfig)
			admin.DELETE("/configs/:id", configHandler.DeleteConfig)
			admin.POST("/configs/generate-crawler-token", configHandler.GenerateCrawlerToken)

			// 日志管理
			admin.GET("/logs", logHandler.GetLogs)
			admin.GET("/logs/:id", logHandler.GetLogByID)
			admin.POST("/logs/cleanup", logHandler.CleanupLogs)
		}
	}

	// WebSocket路由
	r.GET("/ws/crawler/tasks", wsHandler.HandleCrawlerTasks)

	// 静态文件服务 - 上传的文件
	r.Static("/uploads", "./uploads")

	// 创建调度器管理器（日志保留90天）
	schedulerManager := scheduler.NewManager(articleService, logService, 90)

	return r, schedulerManager
}
