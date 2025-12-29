package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/models"
	"github.com/whk-newbie/blog/internal/pkg/db"
	"github.com/whk-newbie/blog/internal/pkg/response"
	"gorm.io/gorm"
)

// CrawlerAuth 爬虫工具认证中间件（Bearer Token）
func CrawlerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "缺少认证Token")
			c.Abort()
			return
		}

		// 验证Token格式 (Bearer <token>)
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Token格式错误")
			c.Abort()
			return
		}

		token := parts[1]

		// 从数据库验证Token
		gormDB, err := db.GetSQLDB()
		if err != nil {
			response.InternalServerError(c, "数据库连接失败")
			c.Abort()
			return
		}

		var config models.SystemConfig
		err = gormDB.Where("config_type = ? AND config_value = ? AND is_active = ?",
			models.ConfigTypeCrawlerToken, token, true).
			First(&config).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				response.Unauthorized(c, "Token无效或已禁用")
			} else {
				response.InternalServerError(c, "Token验证失败")
			}
			c.Abort()
			return
		}

		// 将Token信息存入上下文
		c.Set("crawlerToken", token)
		c.Set("crawlerTokenID", config.ID)

		c.Next()
	}
}
