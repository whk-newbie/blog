package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/pkg/crypto"
	"github.com/whk-newbie/blog/internal/pkg/response"
)

// EncryptionConfig 加密配置
type EncryptionConfig struct {
	// 加密工具实例
	Crypto *crypto.Crypto
	// 是否启用加密（可以通过配置控制）
	Enabled bool
}

// EncryptedRequest 加密请求格式
type EncryptedRequest struct {
	EncryptedData string `json:"encrypted_data"`
	IV            string `json:"iv,omitempty"` // 可选，用于兼容其他格式
	Timestamp     int64  `json:"timestamp,omitempty"`
}

// EncryptedResponse 加密响应格式
type EncryptedResponse struct {
	EncryptedData string `json:"encrypted_data"`
	IV            string `json:"iv,omitempty"`
	Timestamp     int64  `json:"timestamp"`
}

// 不需要加密的路径列表（公开接口和Python SDK接口）
var noEncryptionPaths = []string{
	"/auth/login",
	"/auth/refresh",
	"/fingerprint",
	"/visit",
	"/crawler",    // Python SDK 爬虫接口
	"/ws/crawler", // WebSocket 爬虫接口
}

// Encryption 数据加密中间件
// 支持请求解密和响应加密
// 请求格式: {"encrypted_data": "base64密文", "timestamp": 1234567890}
// 响应格式: {"encrypted_data": "base64密文", "timestamp": 1234567890}
func Encryption(config EncryptionConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果未启用加密，直接放行
		if !config.Enabled || config.Crypto == nil {
			c.Next()
			return
		}

		// 检查是否是排除路径（公开接口不需要加密）
		requestPath := c.Request.URL.Path
		for _, path := range noEncryptionPaths {
			if strings.Contains(requestPath, path) {
				// 公开接口，跳过加密处理
				c.Next()
				return
			}
		}

		// 处理请求解密
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			// 读取请求体
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				response.BadRequest(c, "读取请求体失败: "+err.Error())
				c.Abort()
				return
			}
			c.Request.Body.Close()

			// 如果请求体为空，直接放行
			if len(body) == 0 {
				c.Next()
				return
			}

			// 尝试解析为加密格式
			var encryptedReq EncryptedRequest
			if err := json.Unmarshal(body, &encryptedReq); err != nil {
				// 如果不是加密格式，可能是普通JSON，直接放行
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
				c.Next()
				return
			}

			// 解密数据
			plaintext, err := config.Crypto.Decrypt(encryptedReq.EncryptedData)
			if err != nil {
				response.BadRequest(c, "解密失败: "+err.Error())
				c.Abort()
				return
			}

			// 将解密后的数据替换请求体
			c.Request.Body = io.NopCloser(bytes.NewBufferString(plaintext))
			c.Request.ContentLength = int64(len(plaintext))
		}

		// 保存原始响应写入器
		originalWriter := c.Writer

		// 创建自定义响应写入器（只捕获，不写入）
		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: originalWriter,
			statusCode:     200, // 默认状态码
		}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 处理响应加密
		// 只加密JSON响应
		contentType := blw.Header().Get("Content-Type")
		if blw.statusCode == 200 && (contentType == "application/json; charset=utf-8" || contentType == "application/json") {
			responseBody := blw.body.String()
			if responseBody != "" {
				// 加密响应数据
				encryptedData, err := config.Crypto.Encrypt(responseBody)
				if err != nil {
					// 加密失败，返回原始响应
					originalWriter.Header().Set("Content-Type", contentType)
					originalWriter.Header().Set("Content-Length", strconv.Itoa(len(responseBody)))
					originalWriter.WriteHeader(blw.statusCode)
					originalWriter.WriteString(responseBody)
					return
				}

				// 构建加密响应
				encryptedResp := EncryptedResponse{
					EncryptedData: encryptedData,
					Timestamp:     time.Now().Unix(),
				}

				// 序列化为JSON
				encryptedJSON, err := json.Marshal(encryptedResp)
				if err != nil {
					// 序列化失败，返回原始响应
					originalWriter.Header().Set("Content-Type", contentType)
					originalWriter.Header().Set("Content-Length", strconv.Itoa(len(responseBody)))
					originalWriter.WriteHeader(blw.statusCode)
					originalWriter.WriteString(responseBody)
					return
				}

				// 写入加密响应（清除之前可能写入的内容）
				// 复制响应头（除了Content-Length）
				for key, values := range blw.Header() {
					if key != "Content-Length" {
						for _, value := range values {
							originalWriter.Header().Set(key, value)
						}
					}
				}
				originalWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
				originalWriter.Header().Set("Content-Length", strconv.Itoa(len(encryptedJSON)))
				originalWriter.WriteHeader(blw.statusCode)
				originalWriter.Write(encryptedJSON)
			} else {
				// 响应体为空，直接复制响应头
				for key, values := range blw.Header() {
					for _, value := range values {
						originalWriter.Header().Set(key, value)
					}
				}
				originalWriter.WriteHeader(blw.statusCode)
			}
		} else {
			// 非JSON响应或非200状态码，直接写入原始响应
			// 复制响应头
			for key, values := range blw.Header() {
				for _, value := range values {
					originalWriter.Header().Set(key, value)
				}
			}
			originalWriter.WriteHeader(blw.statusCode)
			originalWriter.Write(blw.body.Bytes())
		}
	}
}

// bodyLogWriter 自定义响应写入器，用于捕获响应体
// 注意：只捕获，不写入到原始ResponseWriter，避免重复写入
type bodyLogWriter struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	statusCode int
	written    bool
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	// 只写入到缓冲区，不写入到原始ResponseWriter
	w.body.Write(b)
	return len(b), nil
}

func (w *bodyLogWriter) WriteString(s string) (int, error) {
	// 只写入到缓冲区，不写入到原始ResponseWriter
	w.body.WriteString(s)
	return len(s), nil
}

func (w *bodyLogWriter) WriteHeader(statusCode int) {
	// 保存状态码
	w.statusCode = statusCode
	// 不调用原始ResponseWriter的WriteHeader，避免提前写入
}
