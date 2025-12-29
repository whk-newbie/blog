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

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 初始化依赖
	gormDB, _ := db.GetSQLDB()
	jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.ExpireTime, cfg.JWT.Issuer)

	// 初始化Repository
	adminRepo := repository.NewAdminRepository(gormDB)
	categoryRepo := repository.NewCategoryRepository(gormDB)
	tagRepo := repository.NewTagRepository(gormDB)
	articleRepo := repository.NewArticleRepository(gormDB)

	// 初始化Service
	authService := service.NewAuthService(adminRepo, jwtManager, cfg.JWT.ExpireTime)
	categoryService := service.NewCategoryService(categoryRepo)
	tagService := service.NewTagService(tagRepo)
	articleService := service.NewArticleService(articleRepo, categoryRepo, tagRepo)

	// 初始化Handler
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	tagHandler := handler.NewTagHandler(tagService)
	articleHandler := handler.NewArticleHandler(articleService)
	uploadHandler := handler.NewUploadHandler("uploads", 10) // 10MB max size

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
		}
	}

	// 静态文件服务 - 上传的文件
	r.Static("/uploads", "./uploads")

	// 创建调度器管理器
	schedulerManager := scheduler.NewManager(articleService)

	return r, schedulerManager
}
