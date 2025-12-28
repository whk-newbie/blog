package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS 跨域中间件
func CORS(allowOrigins, allowMethods, allowHeaders, exposeHeaders []string, allowCredentials bool, maxAge int) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// 检查是否允许该Origin
		allowed := false
		for _, allowedOrigin := range allowOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			if origin != "" {
				c.Header("Access-Control-Allow-Origin", origin)
			} else {
				c.Header("Access-Control-Allow-Origin", "*")
			}
			
			c.Header("Access-Control-Allow-Methods", joinStrings(allowMethods, ", "))
			c.Header("Access-Control-Allow-Headers", joinStrings(allowHeaders, ", "))
			c.Header("Access-Control-Expose-Headers", joinStrings(exposeHeaders, ", "))
			
			if allowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
			
			if maxAge > 0 {
				c.Header("Access-Control-Max-Age", string(rune(maxAge)))
			}
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
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

