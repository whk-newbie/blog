package router

import (
	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/config"
	"github.com/whk-newbie/blog/internal/handler"
	"github.com/whk-newbie/blog/internal/middleware"
	"github.com/whk-newbie/blog/internal/pkg/db"
	"github.com/whk-newbie/blog/internal/pkg/crypto"
	"github.com/whk-newbie/blog/internal/pkg/jwt"
	"github.com/whk-newbie/blog/internal/pkg/logger"
	"github.com/whk-newbie/blog/internal/repository"
	"github.com/whk-newbie/blog/internal/scheduler"
	"github.com/whk-newbie/blog/internal/service"

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

	// Initialize RSA key pair
	var rsaKeyPair *crypto.RSAKeyPair
	var err error
	if cfg.Crypto.RSAPrivKey != "" && cfg.Crypto.RSAPubKey != "" {
		rsaKeyPair, err = crypto.LoadRSAKeyPair(cfg.Crypto.RSAPrivKey, cfg.Crypto.RSAPubKey)
		if err != nil {
			panic("Failed to load RSA key pair: " + err.Error())
		}
	} else {
		rsaKeyPair, err = crypto.NewRSAKeyPair()
		if err != nil {
			panic("Failed to generate RSA key pair: " + err.Error())
		}
		logger.Info("Generated new RSA key pair (not persisted to config)")
	}

	// 初始化Repository
	adminRepo := repository.NewAdminRepository(gormDB)
	categoryRepo := repository.NewCategoryRepository(gormDB)
	tagRepo := repository.NewTagRepository(gormDB)
	articleRepo := repository.NewArticleRepository(gormDB)
	fingerprintRepo := repository.NewFingerprintRepository(gormDB)
	visitRepo := repository.NewVisitRepository(gormDB)
	configRepo := repository.NewConfigRepository(gormDB)
	logRepo := repository.NewLogRepository(gormDB)

	// 初始化Service
	authService := service.NewAuthService(adminRepo, jwtManager, cfg.JWT.ExpireTime)
	categoryService := service.NewCategoryService(categoryRepo)
	tagService := service.NewTagService(tagRepo)
	articleCacheSvc := service.NewArticleCacheService()
	articleService := service.NewArticleService(articleRepo, categoryRepo, tagRepo, articleCacheSvc)

	// 访问统计相关服务
	visitCacheService := service.NewVisitCacheService()
	visitService := service.NewVisitService(visitRepo, visitCacheService)
	fingerprintService := service.NewFingerprintService(fingerprintRepo)
	statsService := service.NewStatsService(articleRepo, categoryRepo, tagRepo, visitService)

	// 初始化配置和日志服务
	configService, err := service.NewConfigService(configRepo, cfg.Crypto.MasterKey)
	if err != nil {
		panic("Failed to initialize config service: " + err.Error())
	}
	logService := service.NewLogService(logRepo)

	// 初始化 AI 服务
	aiProviderRepo := repository.NewAIProviderRepository(gormDB)
	aiService, err := service.NewAIService(aiProviderRepo, cfg.Crypto.MasterKey)
	if err != nil {
		panic("Failed to initialize AI service: " + err.Error())
	}

	// 初始化备份服务
	backupService := service.NewBackupService(cfg)

	// 注册数据库日志钩子，自动将WARN和ERROR级别日志写入数据库
	dbHook := logger.NewDatabaseHook(logService)
	logger.AddHook(dbHook)

	// 初始化Handler
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	tagHandler := handler.NewTagHandler(tagService)
	articleHandler := handler.NewArticleHandler(articleService)
	uploadHandler := handler.NewUploadHandler("uploads", 10) // 10MB max size
	statsHandler := handler.NewStatsHandler(statsService)
	fingerprintHandler := handler.NewFingerprintHandler(fingerprintService)
	visitHandler := handler.NewVisitHandler(visitService)
	configHandler := handler.NewConfigHandler(configService)
	logHandler := handler.NewLogHandler(logService)
	backupHandler := handler.NewBackupHandler(backupService)
	encryptionHandler := handler.NewEncryptionHandler(rsaKeyPair)
	aiHandler := handler.NewAIHandler(aiService)

	// Encryption key negotiation (public, no encryption middleware)
	r.GET("/api/v1/public-key", encryptionHandler.GetPublicKey)
	r.POST("/api/v1/session/key", encryptionHandler.NegotiateKey)

	// API路由组
	api := r.Group("/api/v1")
	{
		// 安全中间件：IP黑名单（在所有公开接口之前）
		api.Use(middleware.IPBlacklist(configService))

		// 安全中间件：限流（针对非登录用户，每分钟60次）
		api.Use(middleware.RateLimit(middleware.RateLimitConfig{
			RequestsPerMinute: 60,
			SkipAuthenticated: true,
		}))

		// Encryption middleware (after rate limiting, before auth)
		api.Use(middleware.Encryption(rsaKeyPair))

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

		// 公开接口 - 站点配置
		api.GET("/site/config", configHandler.GetPublicSiteConfig)

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

			// 配置管理
			admin.GET("/configs", configHandler.GetConfigs)
			admin.GET("/configs/:id", configHandler.GetConfigByID)
			admin.POST("/configs", configHandler.CreateConfig)
			admin.PUT("/configs/:id", configHandler.UpdateConfig)
			admin.DELETE("/configs/:id", configHandler.DeleteConfig)

			// 日志管理
			admin.GET("/logs", logHandler.GetLogs)
			admin.GET("/logs/:id", logHandler.GetLogByID)
			admin.POST("/logs/cleanup", logHandler.CleanupLogs)

			// 数据备份
			admin.GET("/backups", backupHandler.GetBackups)
			admin.POST("/backups", backupHandler.CreateBackup)
			admin.GET("/backups/download/:filename", backupHandler.DownloadBackup)
			admin.DELETE("/backups/:filename", backupHandler.DeleteBackup)
			admin.POST("/backups/cleanup", backupHandler.CleanupBackups)

			// AI 提供方管理
			admin.GET("/ai/providers", aiHandler.ListProviders)
			admin.GET("/ai/providers/:id", aiHandler.GetProvider)
			admin.POST("/ai/providers", aiHandler.CreateProvider)
			admin.PUT("/ai/providers/:id", aiHandler.UpdateProvider)
			admin.DELETE("/ai/providers/:id", aiHandler.DeleteProvider)
			admin.POST("/ai/providers/:id/check", aiHandler.CheckProvider)

			// AI 翻译
			admin.POST("/ai/translate/:id", aiHandler.TranslateArticle)

			// AI 聊天
			admin.POST("/ai/chat", aiHandler.Chat)
		}
	}

	// 静态文件服务 - 上传的文件
	r.Static("/uploads", "./uploads")

	// 创建调度器管理器（日志保留90天，备份每天凌晨3点，保留10个备份）
	backupSchedule := "0 0 3 * * *" // 每天凌晨3点
	backupRetentionCount := 10      // 保留10个备份
	schedulerManager := scheduler.NewManager(articleService, logService, backupService, 90, backupSchedule, backupRetentionCount)

	return r, schedulerManager
}
