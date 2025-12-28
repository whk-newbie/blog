package router

import (
	"github.com/gin-gonic/gin"
	"github.com/iambaby/blog/internal/config"
	"github.com/iambaby/blog/internal/handler"
	"github.com/iambaby/blog/internal/middleware"
	"github.com/iambaby/blog/internal/pkg/db"
	"github.com/iambaby/blog/internal/pkg/jwt"
	"github.com/iambaby/blog/internal/repository"
	"github.com/iambaby/blog/internal/service"
	
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Setup 设置路由
func Setup(cfg *config.Config) *gin.Engine {
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
			"status": "ok",
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
	
	// 初始化Service
	authService := service.NewAuthService(adminRepo, jwtManager, cfg.JWT.ExpireTime)
	
	// 初始化Handler
	authHandler := handler.NewAuthHandler(authService)

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

		// 公开接口
		public := api.Group("/public")
		{
			public.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "pong"})
			})
		}

		// 管理接口（需要认证）
		admin := api.Group("/admin")
		admin.Use(middleware.Auth(jwtManager))
		{
			admin.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "admin pong"})
			})
		}
	}

	return r
}

