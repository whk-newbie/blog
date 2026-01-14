package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/pkg/logger"
	"github.com/whk-newbie/blog/internal/pkg/response"
)

// CORS 跨域中间件
func CORS(allowOrigins, allowMethods, allowHeaders, exposeHeaders []string, allowCredentials bool, maxAge int) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 如果没有Origin头（同源请求），直接放行
		if origin == "" {
			c.Next()
			return
		}

		// 检查是否允许该Origin
		allowed := false
		allowedOrigin := ""
		for _, ao := range allowOrigins {
			if ao == "*" {
				// 如果允许所有源，但需要凭据时不能使用*
				if !allowCredentials {
					allowed = true
					allowedOrigin = "*"
					break
				}
			} else if ao == origin {
				allowed = true
				allowedOrigin = origin
				break
			}
		}

		// 如果允许该源，设置CORS头
		if allowed {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
			c.Header("Access-Control-Allow-Methods", joinStrings(allowMethods, ", "))
			c.Header("Access-Control-Allow-Headers", joinStrings(allowHeaders, ", "))
			c.Header("Access-Control-Expose-Headers", joinStrings(exposeHeaders, ", "))

			if allowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}

			if maxAge > 0 {
				c.Header("Access-Control-Max-Age", fmt.Sprintf("%d", maxAge))
			}

			// 处理预检请求
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
		} else {
			// 如果不允许该源，记录日志并返回友好的错误信息
			logger.Warn("CORS policy violation: Origin '%s' is not allowed. Allowed origins: %v", origin, allowOrigins)
			response.Forbidden(c, fmt.Sprintf("CORS policy: Origin '%s' is not allowed", origin))
			c.Abort()
			return
		}

		c.Next()
	}
}

func joinStrings(arr []string, sep string) string {
	result := ""
	for i, s := range arr {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
