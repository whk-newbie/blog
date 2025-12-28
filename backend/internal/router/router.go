package router

import (
	"github.com/gin-gonic/gin"
	"github.com/iambaby/blog/internal/config"
	"github.com/iambaby/blog/internal/middleware"
	
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

	// API路由组
	api := r.Group("/api/v1")
	{
		// 公开接口
		public := api.Group("/public")
		{
			public.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "pong"})
			})
		}

		// 管理接口（需要认证）
		admin := api.Group("/admin")
		{
			admin.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "admin pong"})
			})
		}
	}

	return r
}

