package middleware

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/service"
)

// IPBlacklist IP黑名单中间件
// 从配置服务中获取IP黑名单，支持单个IP和CIDR格式
// 如果IP在黑名单中，返回404
func IPBlacklist(configService service.ConfigService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端IP
		clientIP := c.ClientIP()
		if clientIP == "" {
			c.Next()
			return
		}

		// 获取IP黑名单配置
		blacklistConfigs, err := configService.GetConfigs("ip_blacklist")
		if err != nil {
			// 配置服务错误，记录日志但允许请求通过（降级策略）
			c.Next()
			return
		}

		// 检查IP是否在黑名单中
		for _, config := range blacklistConfigs {
			if !config.IsActive {
				continue
			}

			// 获取配置值（可能是单个IP或CIDR格式）
			// 注意：这里使用的是脱敏后的值，需要获取原始值
			// 由于中间件无法直接访问加密服务，我们需要通过配置服务获取原始值
			blacklistValue, err := configService.GetConfigValue(config.ConfigKey)
			if err != nil {
				// 获取配置值失败，跳过
				continue
			}
			if blacklistValue == "" {
				continue
			}

			// 检查是否是CIDR格式（包含/）
			if strings.Contains(blacklistValue, "/") {
				// CIDR格式
				_, ipNet, err := net.ParseCIDR(blacklistValue)
				if err != nil {
					// CIDR格式错误，跳过
					continue
				}

				// 检查IP是否在CIDR范围内
				ip := net.ParseIP(clientIP)
				if ip != nil && ipNet.Contains(ip) {
					// IP在黑名单中，返回404
					c.AbortWithStatus(404)
					return
				}
			} else {
				// 单个IP格式
				if clientIP == blacklistValue {
					// IP在黑名单中，返回404
					c.AbortWithStatus(404)
					return
				}
			}
		}

		c.Next()
	}
}
