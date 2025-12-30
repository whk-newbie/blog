package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
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

		// 创建自定义响应写入器
		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: originalWriter,
		}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 处理响应加密
		// 只加密JSON响应
		if c.Writer.Status() == 200 && c.Writer.Header().Get("Content-Type") == "application/json; charset=utf-8" {
			responseBody := blw.body.String()
			if responseBody != "" {
				// 加密响应数据
				encryptedData, err := config.Crypto.Encrypt(responseBody)
				if err != nil {
					// 加密失败，返回原始响应
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
					originalWriter.WriteString(responseBody)
					return
				}

				// 写入加密响应
				originalWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
				originalWriter.Header().Set("Content-Length", strconv.Itoa(len(encryptedJSON)))
				originalWriter.Write(encryptedJSON)
			}
		} else {
			// 非JSON响应，直接写入
			originalWriter.Write(blw.body.Bytes())
		}
	}
}

// bodyLogWriter 自定义响应写入器，用于捕获响应体
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
