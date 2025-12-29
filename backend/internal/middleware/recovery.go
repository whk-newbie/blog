package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/pkg/logger"
	"github.com/whk-newbie/blog/internal/pkg/response"
)

// Recovery 错误恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录panic信息
				logger.Error("Panic recovered: %v\n%s", err, debug.Stack())
				
				// 返回错误响应
				response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err))
				c.Abort()
			}
		}()
		c.Next()
	}
}

