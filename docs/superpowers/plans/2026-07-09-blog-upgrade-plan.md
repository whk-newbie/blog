# 博客系统升级优化实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 对博客系统进行全量升级：加密通信、代码混淆、后台入口隐藏、去爬虫留指纹、AI 翻译+聊天、zyyo.net 风格全站 UI 改版。

**Architecture:** 分层推进 —— 安全基础设施（RSA+AES 加密、混淆、入口配置）→ 功能层（去爬虫、AI 翻译、AI 提供方管理、AI 聊天）→ UI 层（极简双栏布局、时间轴、暗色模式）。

**Tech Stack:** Go 1.24 + Gin + GORM + PostgreSQL + Redis, Vue 3 + Vite 5 + Element Plus + Pinia

---

## File Structure Map

### 后端新增/修改
| 文件 | 操作 | 职责 |
|------|------|------|
| `backend/internal/config/config.go` | 修改 | 新增 RSA 密钥配置项 |
| `backend/internal/pkg/crypto/rsa.go` | 新增 | RSA 密钥对生成与管理 |
| `backend/internal/pkg/crypto/session.go` | 新增 | Session→AES key 的 Redis 存取 |
| `backend/internal/middleware/encryption.go` | 新增 | 全局加解密中间件 |
| `backend/internal/handler/encryption_handler.go` | 新增 | 公钥分发 + 密钥协商接口 |
| `backend/internal/middleware/crawler_auth.go` | 删除 | 爬虫认证中间件 |
| `backend/internal/handler/crawler_handler.go` | 删除 | 爬虫处理器 |
| `backend/internal/service/crawl_service.go` | 删除 | 爬虫服务 |
| `backend/internal/repository/crawl_task_repo.go` | 删除 | 爬虫仓库 |
| `backend/internal/models/crawl_task.go` | 删除 | 爬虫模型 |
| `backend/internal/websocket/hub.go` | 修改 | 移除爬虫 WebSocket 相关 |
| `backend/internal/models/article.go` | 修改 | 新增英文翻译字段 |
| `backend/internal/models/ai_provider.go` | 新增 | AI 提供方模型 |
| `backend/internal/models/ai_chat_history.go` | 新增 | AI 聊天记录模型 |
| `backend/internal/repository/ai_provider_repo.go` | 新增 | AI 提供方仓库 |
| `backend/internal/repository/ai_chat_repo.go` | 新增 | AI 聊天记录仓库 |
| `backend/internal/service/ai_service.go` | 新增 | AI 翻译 + 聊天服务 |
| `backend/internal/handler/ai_handler.go` | 新增 | AI 接口处理器 |
| `backend/internal/router/router.go` | 修改 | 移除爬虫路由，新增加密/翻译/聊天/AI 管理路由 |
| `backend/internal/service/article_service.go` | 修改 | 新增翻译触发逻辑 |
| `backend/internal/handler/article_handler.go` | 修改 | 新增翻译触发接口 |
| `backend/internal/models/system_config.go` | 修改 | 新增 ConfigTypeSiteConfig |
| `backend/migrations/` | 新增 | 新增文章翻译字段、AI 提供方表、聊天记录表迁移 |

### 前端新增/修改
| 文件 | 操作 | 职责 |
|------|------|------|
| `frontend/src/api/http.js` | 修改 | 添加加解密拦截器 |
| `frontend/src/api/ai.js` | 新增 | AI 相关 API 封装 |
| `frontend/src/api/crawler.js` | 删除 | 爬虫 API |
| `frontend/src/views/admin/CrawlerMonitor.vue` | 删除 | 爬虫监控页面 |
| `frontend/src/views/admin/AiProviders.vue` | 新增 | AI 提供方管理 |
| `frontend/src/views/admin/AiChat.vue` | 新增 | AI 聊天页面 |
| `frontend/src/views/Timeline.vue` | 新增 | 时间轴页面 |
| `frontend/src/views/Home.vue` | 修改 | zyyo.net 风格重写 |
| `frontend/src/views/Tools.vue` | 修改 | 极简卡片风格 |
| `frontend/src/views/ArticleDetail.vue` | 修改 | 极简风格重写 |
| `frontend/src/views/Articles.vue` | 修改 | 极简风格重写 |
| `frontend/src/components/layout/MainLayout.vue` | 修改 | 双栏侧边栏布局 |
| `frontend/src/components/layout/Header.vue` | 修改 | 移除后台入口，新增语言切换/搜索 |
| `frontend/src/components/layout/Sidebar.vue` | 新增 | 侧边栏（头像/统计/标签云/社交） |
| `frontend/src/components/layout/Footer.vue` | 修改 | 极简样式 |
| `frontend/src/components/common/SearchOverlay.vue` | 新增 | 全屏搜索弹窗 |
| `frontend/src/components/chat/AiChatPanel.vue` | 新增 | AI 聊天悬浮面板 |
| `frontend/src/router/index.js` | 修改 | 动态后台路径、时间轴路由、移除爬虫路由 |
| `frontend/src/store/language.js` | 修改 | 同步 i18n locale |
| `frontend/vite.config.js` | 修改 | 新增 obfuscator 插件配置 |
| `frontend/src/assets/styles/variables.less` | 修改 | 极简配色变量 |
| `frontend/src/assets/styles/common.less` | 修改 | 全局样式调整 |
| `frontend/package.json` | 修改 | 新增依赖 |

---

## Layer 1: 安全基础设施

### Task 1.1: 后端 RSA 密钥管理 + Session 加密存储

**Files:**
- Create: `backend/internal/pkg/crypto/rsa.go`
- Create: `backend/internal/pkg/crypto/session.go`
- Modify: `backend/internal/config/config.go`

- [ ] **Step 1: 添加 RSA 配置项到 Config**

`backend/internal/config/config.go` — 在 `CryptoConfig` 中新增字段：

```go
// CryptoConfig 加密配置
type CryptoConfig struct {
	MasterKey  string `yaml:"master_key"`  // AES-256密钥(32字节)
	RSAPrivKey string `yaml:"rsa_priv_key"` // RSA 私钥 (PEM, 可选，不配置则自动生成)
	RSAPubKey  string `yaml:"rsa_pub_key"`  // RSA 公钥 (PEM, 可选)
}
```

- [ ] **Step 2: 在 `overrideFromEnv` 中新增 RSA 环境变量覆盖**

```go
if privKey := os.Getenv("RSA_PRIV_KEY"); privKey != "" {
    cfg.Crypto.RSAPrivKey = privKey
}
if pubKey := os.Getenv("RSA_PUB_KEY"); pubKey != "" {
    cfg.Crypto.RSAPubKey = pubKey
}
```

- [ ] **Step 3: 创建 RSA 密钥管理文件**

`backend/internal/pkg/crypto/rsa.go`:

```go
package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

var (
	ErrRSAKeyTooShort  = errors.New("RSA key too short")
	ErrRSADecryptFailed = errors.New("RSA decryption failed")
)

const RSAKeyBits = 2048

// RSAKeyPair RSA 密钥对
type RSAKeyPair struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
	PublicKeyPEM  string
	PrivateKeyPEM string
}

// NewRSAKeyPair 生成新的 RSA 密钥对
func NewRSAKeyPair() (*RSAKeyPair, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, RSAKeyBits)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	privKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	})

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}
	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	return &RSAKeyPair{
		PublicKey:     &privKey.PublicKey,
		PrivateKey:    privKey,
		PublicKeyPEM:  string(pubKeyPEM),
		PrivateKeyPEM: string(privKeyPEM),
	}, nil
}

// LoadRSAKeyPair 从 PEM 字符串加载 RSA 密钥对
func LoadRSAKeyPair(privKeyPEM, pubKeyPEM string) (*RSAKeyPair, error) {
	privBlock, _ := pem.Decode([]byte(privKeyPEM))
	if privBlock == nil {
		return nil, errors.New("failed to decode private key PEM")
	}
	privKey, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	pubBlock, _ := pem.Decode([]byte(pubKeyPEM))
	if pubBlock == nil {
		return nil, errors.New("failed to decode public key PEM")
	}
	pubKeyInterface, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	pubKey, ok := pubKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	return &RSAKeyPair{
		PublicKey:     pubKey,
		PrivateKey:    privKey,
		PublicKeyPEM:  pubKeyPEM,
		PrivateKeyPEM: privKeyPEM,
	}, nil
}

// EncryptWithPublicKey 使用 RSA 公钥加密数据（OAEP + SHA-256）
func EncryptWithPublicKey(pubKey *rsa.PublicKey, plaintext []byte) (string, error) {
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, plaintext, nil)
	if err != nil {
		return "", fmt.Errorf("RSA encryption failed: %w", err)
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptWithPrivateKey 使用 RSA 私钥解密数据
func DecryptWithPrivateKey(privKey *rsa.PrivateKey, ciphertextB64 string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, ciphertext, nil)
	if err != nil {
		return nil, ErrRSADecryptFailed
	}
	return plaintext, nil
}
```

- [ ] **Step 4: 创建 Session AES Key 存储**

`backend/internal/pkg/crypto/session.go`:

```go
package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/whk-newbie/blog/internal/pkg/redis"
)

const (
	SessionKeyPrefix = "session:enc:"
	SessionTTL       = 30 * time.Minute
	AESKeySize       = 32 // AES-256
)

// SessionKey 存储 session → AES key 映射
type SessionKey struct {
	SessionID string
	AESKey    []byte
}

// GenerateSessionID 生成新的 session ID
func GenerateSessionID() string {
	return uuid.New().String()
}

// GenerateAESKey 生成随机 AES-256 密钥
func GenerateAESKey() ([]byte, error) {
	key := make([]byte, AESKeySize)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate AES key: %w", err)
	}
	return key, nil
}

// StoreSessionKey 将 session→AES key 存入 Redis
func StoreSessionKey(sessionID string, aesKey []byte) error {
	key := SessionKeyPrefix + sessionID
	return redis.Set(key, base64.StdEncoding.EncodeToString(aesKey), SessionTTL)
}

// GetSessionKey 从 Redis 获取 AES key
func GetSessionKey(sessionID string) ([]byte, error) {
	key := SessionKeyPrefix + sessionID
	val, err := redis.GetValue(key)
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(val)
}

// RefreshSessionKey 刷新 session TTL
func RefreshSessionKey(sessionID string) error {
	key := SessionKeyPrefix + sessionID
	return redis.Expire(key, SessionTTL)
}
```

- [ ] **Step 5: 提交**

```bash
git add backend/internal/pkg/crypto/rsa.go backend/internal/pkg/crypto/session.go backend/internal/config/config.go
git commit -m "feat: add RSA key management and session encryption storage"
```

---

### Task 1.2: 后端加解密中间件 + 密钥协商接口

**Files:**
- Create: `backend/internal/middleware/encryption.go`
- Create: `backend/internal/handler/encryption_handler.go`
- Modify: `backend/internal/router/router.go`

- [ ] **Step 1: 创建加解密中间件**

`backend/internal/middleware/encryption.go`:

```go
package middleware

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/pkg/crypto"
	"github.com/whk-newbie/blog/internal/pkg/logger"
	"github.com/whk-newbie/blog/internal/pkg/response"
)

// EncryptionWhitelist 不需要加密的路径列表
var EncryptionWhitelist = []string{
	"/api/v1/public-key",
	"/api/v1/session/key",
	"/api/v1/health",
	"/swagger",
	"/docs",
	"/uploads",
	"/ws",
	"/health",
}

func isWhitelisted(path string) bool {
	for _, prefix := range EncryptionWhitelist {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

// Encryption 加解密中间件
func Encryption(rsaKeyPair *crypto.RSAKeyPair) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 白名单路径跳过
		if isWhitelisted(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 从 Header 获取 session ID
		sessionID := c.GetHeader("X-Session-Id")
		if sessionID == "" {
			response.Error(c, 40001, "missing session ID, please negotiate encryption key first")
			c.Abort()
			return
		}

		// 从 Redis 获取 AES key
		aesKey, err := crypto.GetSessionKey(sessionID)
		if err != nil {
			response.Error(c, 40002, "session expired or invalid, please re-negotiate encryption key")
			c.Abort()
			return
		}

		// 解密请求体
		if c.Request.Body != nil && c.Request.ContentLength > 0 {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			c.Request.Body.Close()
			if err != nil {
				response.BadRequest(c, "failed to read request body")
				c.Abort()
				return
			}

			if len(bodyBytes) > 0 {
				plaintext, err := aesGCMDecrypt(aesKey, bodyBytes)
				if err != nil {
					logger.Warn("Decryption failed for path %s: %v", c.Request.URL.Path, err)
					response.Error(c, 40003, "decryption failed, request body may be corrupted")
					c.Abort()
					return
				}
				c.Request.Body = io.NopCloser(bytes.NewReader(plaintext))
				c.Request.ContentLength = int64(len(plaintext))
			}
		}

		// 包装 ResponseWriter 以加密响应
		ew := &encryptedWriter{
			ResponseWriter: c.Writer,
			aesKey:         aesKey,
		}
		c.Writer = ew

		c.Next()

		// 刷新 session TTL
		_ = crypto.RefreshSessionKey(sessionID)
	}
}

// encryptedWriter 包装 gin.ResponseWriter，在 Write 时加密
type encryptedWriter struct {
	gin.ResponseWriter
	aesKey []byte
	buf    bytes.Buffer
}

func (w *encryptedWriter) Write(data []byte) (int, error) {
	return w.buf.Write(data)
}

func (w *encryptedWriter) WriteString(s string) (int, error) {
	return w.buf.WriteString(s)
}

// FlushEncrypted 在最终写入时加密整个缓冲区内容
func (w *encryptedWriter) FlushEncrypted() error {
	if w.buf.Len() == 0 {
		return nil
	}
	ciphertext, err := aesGCMEncrypt(w.aesKey, w.buf.Bytes())
	if err != nil {
		return err
	}
	w.ResponseWriter.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.ResponseWriter.Write(ciphertext)
	return err
}

// aesGCMEncrypt AES-256-GCM 加密
func aesGCMEncrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(bytes.NewReader(make([]byte, gcm.NonceSize()))); err != nil {
		// 使用固定 nonce（仅用于演示，生产环境需要随机 nonce）
		for i := range nonce {
			nonce[i] = byte(i)
		}
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// aesGCMDecrypt AES-256-GCM 解密
func aesGCMDecrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, io.ErrUnexpectedEOF
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
```

**注意**：由于 Gin 的 `c.Writer.Write()` 调用时机在中间件返回之后，完整加密方案需要对 `Write` 方法做延迟处理。上述 `encryptedWriter` 是一个缓冲实现，需要在 gin 框架生命周期中正确 flush。实际实现时可以考虑通过 `c.Abort()` 后的自定义写入来控制。

**简化方案（推荐）**：不包装 Writer，而是在中间件中 Hook `c.Next()` 之后通过 `c.Writer.Size()` 检查并在必要时手动处理。或者——更简单的方案——使用一个全局的 response writer wrapper。

对于实际实现，推荐创建一个更简洁的版本：

```go
// Encryption 加解密中间件（简化但完整的实现）
func Encryption(rsaKeyPair *crypto.RSAKeyPair) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isWhitelisted(c.Request.URL.Path) {
			c.Next()
			return
		}

		sessionID := c.GetHeader("X-Session-Id")
		if sessionID == "" {
			response.Error(c, 40001, "missing session ID")
			c.Abort()
			return
		}

		aesKey, err := crypto.GetSessionKey(sessionID)
		if err != nil {
			response.Error(c, 40002, "session expired")
			c.Abort()
			return
		}

		// 解密请求体
		if c.Request.Body != nil && c.Request.ContentLength > 0 {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body.Close()
			if len(bodyBytes) > 0 {
				plaintext, err := aesGCMDecrypt(aesKey, bodyBytes)
				if err != nil {
					response.Error(c, 40003, "decryption failed")
					c.Abort()
					return
				}
				c.Request.Body = io.NopCloser(bytes.NewReader(plaintext))
				c.Request.ContentLength = int64(len(plaintext))
			} else {
				c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			}
		}

		// 包装 writer
		blw := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
			aesKey:         aesKey,
		}
		c.Writer = blw

		c.Next()

		// 加密并写入响应体
		encrypted, err := aesGCMEncrypt(aesKey, blw.body.Bytes())
		if err != nil {
			return
		}
		blw.ResponseWriter.Header().Set("Content-Type", "application/octet-stream")
		blw.ResponseWriter.Write(encrypted)

		_ = crypto.RefreshSessionKey(sessionID)
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	aesKey []byte
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *bodyLogWriter) WriteString(s string) (int, error) {
	return w.body.WriteString(s)
}
```

- [ ] **Step 2: 创建密钥协商 Handler**

`backend/internal/handler/encryption_handler.go`:

```go
package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/pkg/crypto"
	"github.com/whk-newbie/blog/internal/pkg/response"
)

// EncryptionHandler 加密密钥协商处理器
type EncryptionHandler struct {
	rsaKeyPair *crypto.RSAKeyPair
}

// NewEncryptionHandler 创建加密处理器
func NewEncryptionHandler(rsaKeyPair *crypto.RSAKeyPair) *EncryptionHandler {
	return &EncryptionHandler{rsaKeyPair: rsaKeyPair}
}

// PublicKeyResponse 公钥响应
type PublicKeyResponse struct {
	PublicKey string `json:"public_key"`
	SessionID string `json:"session_id"`
}

// GetPublicKey 获取 RSA 公钥
func (h *EncryptionHandler) GetPublicKey(c *gin.Context) {
	sessionID := crypto.GenerateSessionID()
	response.Success(c, PublicKeyResponse{
		PublicKey: h.rsaKeyPair.PublicKeyPEM,
		SessionID: sessionID,
	})
}

// SessionKeyRequest 密钥协商请求
type SessionKeyRequest struct {
	EncryptedKey string `json:"encrypted_key" binding:"required"`
	SessionID    string `json:"session_id" binding:"required"`
}

// NegotiateKey 协商 AES 密钥
func (h *EncryptionHandler) NegotiateKey(c *gin.Context) {
	var req SessionKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	// RSA 解密得到 AES 密钥
	aesKey, err := crypto.DecryptWithPrivateKey(h.rsaKeyPair.PrivateKey, req.EncryptedKey)
	if err != nil {
		response.Error(c, 40004, "failed to decrypt AES key")
		return
	}

	if len(aesKey) != 32 {
		response.BadRequest(c, "AES key must be 32 bytes")
		return
	}

	// 存入 Redis
	if err := crypto.StoreSessionKey(req.SessionID, aesKey); err != nil {
		response.InternalServerError(c, "failed to store session key")
		return
	}

	response.Success(c, gin.H{"status": "ok"})
}
```

- [ ] **Step 3: 修改 router.go —— 注入 RSA key pair 和新路由**

`backend/internal/router/router.go` — 修改点：

在 `Setup` 函数开头，`jwtManager` 初始化之后添加：

```go
// 初始化 RSA 密钥对
var rsaKeyPair *crypto.RSAKeyPair
var err error
if cfg.Crypto.RSAPrivKey != "" && cfg.Crypto.RSAPubKey != "" {
	rsaKeyPair, err = crypto.LoadRSAKeyPair(cfg.Crypto.RSAPrivKey, cfg.Crypto.RSAPubKey)
	if err != nil {
		panic("Failed to load RSA key pair: " + err.Error())
	}
} else {
	rsaKeyPair, err = crypto.NewRSAKeyPair()
	if err != nil {
		panic("Failed to generate RSA key pair: " + err.Error())
	}
	logger.Info("Generated new RSA key pair (not persisted)")
}
```

在 Handler 初始化处添加：

```go
encryptionHandler := handler.NewEncryptionHandler(rsaKeyPair)
```

在路由注册处新增（`api := r.Group("/api/v1")` 之前或之后）：

```go
// 加密密钥协商接口（公开，不经过加密中间件）
r.GET("/api/v1/public-key", encryptionHandler.GetPublicKey)
r.POST("/api/v1/session/key", encryptionHandler.NegotiateKey)
```

在 `api` 路由组下，在限流中间件之后添加加密中间件：

```go
// 加密中间件（在限流之后，认证之前）
api.Use(middleware.Encryption(rsaKeyPair))
```

注意：同时需要删除所有爬虫相关代码（见 Task 2.1），包括：
- 删除 `crawlTaskRepo`、`crawlService`、`crawlerHandler` 的初始化
- 删除爬虫路由组
- 删除 `POST /admin/configs/generate-crawler-token` 路由
- 删除 WebSocket 爬虫路由 `r.GET("/ws/crawler/tasks", wsHandler.HandleCrawlerTasks)`

- [ ] **Step 4: 提交**

```bash
git add backend/internal/middleware/encryption.go backend/internal/handler/encryption_handler.go backend/internal/router/router.go
git commit -m "feat: add encryption middleware and key negotiation endpoints"
```

---

### Task 1.3: 前端 AES 加解密 + Axios 拦截器

**Files:**
- Modify: `frontend/src/api/http.js`
- Create: `frontend/src/utils/crypto.js`

- [ ] **Step 1: 创建前端加密工具**

`frontend/src/utils/crypto.js`:

```javascript
import CryptoJS from 'crypto-js'

const JSEncrypt = () => {
  // 使用内置 Web Crypto API 或 crypto-js 做 RSA
  // 注意：浏览器原生支持 SubtleCrypto，用于 RSA-OAEP
}

/**
 * 生成随机 AES-256 密钥（32 字节）
 */
export function generateAESKey() {
  const key = CryptoJS.lib.WordArray.random(32)
  return key
}

/**
 * AES-GCM 加密（crypto-js 默认 CBC 模式，使用 AES-CBC + HMAC 替代）
 * 返回 Base64 编码的密文
 */
export function aesEncrypt(plaintext, key) {
  const keyWordArray = typeof key === 'string'
    ? CryptoJS.enc.Base64.parse(key)
    : key
  const iv = CryptoJS.lib.WordArray.random(16)
  const encrypted = CryptoJS.AES.encrypt(plaintext, keyWordArray, {
    iv: iv,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7,
  })
  // 将 IV 拼接在密文前面
  const combined = iv.concat(encrypted.ciphertext)
  return CryptoJS.enc.Base64.stringify(combined)
}

/**
 * AES-CBC 解密
 */
export function aesDecrypt(ciphertextB64, key) {
  const keyWordArray = typeof key === 'string'
    ? CryptoJS.enc.Base64.parse(key)
    : key
  const combined = CryptoJS.enc.Base64.parse(ciphertextB64)
  const iv = CryptoJS.lib.WordArray.create(combined.words.slice(0, 4), 16)
  const ciphertext = CryptoJS.lib.WordArray.create(
    combined.words.slice(4),
    combined.sigBytes - 16
  )
  const decrypted = CryptoJS.AES.decrypt(
    CryptoJS.enc.Base64.stringify(ciphertext),
    keyWordArray,
    {
      iv: iv,
      mode: CryptoJS.mode.CBC,
      padding: CryptoJS.pad.Pkcs7,
    }
  )
  return decrypted.toString(CryptoJS.enc.Utf8)
}

/**
 * 使用 RSA 公钥加密（使用 Web Crypto API）
 */
export async function rsaEncrypt(publicKeyPEM, plaintext) {
  // 将 PEM 转为 ArrayBuffer
  const pemHeader = '-----BEGIN PUBLIC KEY-----'
  const pemFooter = '-----END PUBLIC KEY-----'
  const pemContents = publicKeyPEM
    .replace(pemHeader, '')
    .replace(pemFooter, '')
    .replace(/\s/g, '')
  const binaryDer = Uint8Array.from(atob(pemContents), c => c.charCodeAt(0))

  const publicKey = await crypto.subtle.importKey(
    'spki',
    binaryDer.buffer,
    { name: 'RSA-OAEP', hash: 'SHA-256' },
    false,
    ['encrypt']
  )

  const encoded = new TextEncoder().encode(plaintext)
  const encrypted = await crypto.subtle.encrypt(
    { name: 'RSA-OAEP' },
    publicKey,
    encoded
  )

  return btoa(String.fromCharCode(...new Uint8Array(encrypted)))
}
```

- [ ] **Step 2: 修改 Axios 实例 —— 添加加密拦截器**

`frontend/src/api/http.js` — 在文件顶部新增导入，在现有逻辑之前添加加密协商：

```javascript
import { generateAESKey, aesEncrypt, aesDecrypt, rsaEncrypt } from '@/utils/crypto'

// 加密会话状态
let sessionId = null
let aesKey = null
let negotiating = false
let negotiatePromise = null

// 协商加密密钥
async function negotiateKey() {
  if (aesKey && sessionId) return // 已有有效 key

  if (negotiating && negotiatePromise) {
    return negotiatePromise
  }

  negotiating = true
  negotiatePromise = (async () => {
    try {
      // 1. 获取公钥
      const pubKeyRes = await axios.get('/api/v1/public-key')
      const { public_key, session_id } = pubKeyRes.data.data
      sessionId = session_id

      // 2. 生成 AES 密钥
      const rawKey = generateAESKey()
      const keyB64 = CryptoJS.enc.Base64.stringify(rawKey)

      // 3. 用 RSA 公钥加密 AES 密钥
      const encryptedKey = await rsaEncrypt(public_key, keyB64)

      // 4. 发送加密后的密钥
      await axios.post('/api/v1/session/key', {
        encrypted_key: encryptedKey,
        session_id: sessionId,
      })

      aesKey = rawKey
      negotiating = false
    } catch (error) {
      negotiating = false
      negotiatePromise = null
      throw error
    }
  })()

  return negotiatePromise
}

// 请求拦截器 —— 在原有逻辑之前添加加密
http.interceptors.request.use(
  async (config) => {
    // 跳过白名单路径
    const whitelist = ['/api/v1/public-key', '/api/v1/session/key']
    if (whitelist.some(p => config.url?.includes(p))) {
      return config
    }

    // 确保已协商密钥
    if (!aesKey) {
      await negotiateKey()
    }

    config.headers['X-Session-Id'] = sessionId

    // 加密请求体
    if (config.data && typeof config.data === 'object') {
      const jsonStr = JSON.stringify(config.data)
      const encrypted = aesEncrypt(jsonStr, aesKey)
      config.data = encrypted
      // 标记为已加密，修改 content-type
    }

    // ... 原有的 token、request-id 逻辑保持不变 ...
    return config
  },
  // ... error handler 保持不变 ...
)

// 响应拦截器 —— 添加解密逻辑
http.interceptors.response.use(
  async (response) => {
    // 跳过白名单路径
    const whitelist = ['/api/v1/public-key', '/api/v1/session/key']
    if (whitelist.some(p => response.config.url?.includes(p))) {
      // 原有解密逻辑 ...
      // 注意：public-key 和 session/key 响应不需要解密
      return response.data?.data
    }

    // 解密响应体
    if (response.data && typeof response.data === 'string' && aesKey) {
      try {
        const decrypted = aesDecrypt(response.data, aesKey)
        response.data = JSON.parse(decrypted)
      } catch (e) {
        // 如果解密失败，可能服务端返回了明文错误
        // 尝试直接解析
      }
    }

    // ... 原有的业务响应处理逻辑 ...
  },
  async (error) => {
    // 如果是 session 过期错误（40002），重新协商密钥并重试
    if (error.response?.data?.code === 40002) {
      aesKey = null
      sessionId = null
      await negotiateKey()
      // 重试原请求（需要重新加密）
      // ...
    }
    // ... 原有的错误处理逻辑 ...
  }
)
```

**注意**：`http.js` 的修改需要谨慎整合，确保新的加密逻辑与现有的 token 注入、错误处理并存。实际实施时应该完整重写 `http.js`，将两者无缝合并。

- [ ] **Step 3: 提交**

```bash
git add frontend/src/utils/crypto.js frontend/src/api/http.js
git commit -m "feat: add frontend AES encryption and auto key negotiation in axios"
```

---

### Task 1.4: 后台入口可配置

**Files:**
- Modify: `frontend/src/router/index.js`
- Modify: `frontend/src/components/layout/Header.vue`
- Modify: `frontend/src/api/config.js`

- [ ] **Step 1: 修改路由 —— 动态后台路径**

`frontend/src/router/index.js`:

```javascript
import { createRouter, createWebHistory } from 'vue-router'

// 动态获取 admin 路径
function getAdminPath() {
  // 从站点配置中读取，默认为 'admin'
  const stored = localStorage.getItem('admin_path')
  return stored || 'admin'
}

const routes = [
  // 公开页面（保持不变）
  {
    path: '/',
    component: () => import('@/components/layout/MainLayout.vue'),
    children: [
      { path: '', name: 'Home', component: () => import('@/views/Home.vue') },
      { path: 'articles', name: 'Articles', component: () => import('@/views/Articles.vue') },
      { path: 'article/:slug', name: 'ArticleDetail', component: () => import('@/views/ArticleDetail.vue') },
      { path: 'timeline', name: 'Timeline', component: () => import('@/views/Timeline.vue') },
      { path: 'tools', name: 'Tools', component: () => import('@/views/Tools.vue') },
    ]
  },
  // 管理后台 —— 动态路径
  {
    path: `/${getAdminPath()}`,
    component: () => import('@/components/layout/AdminLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Admin',
        component: () => import('@/views/admin/Dashboard.vue'),
        meta: { requiresAuth: true }
      },
      // ... 其他管理路由（crawler 相关已删除，详见 Task 2.1）...
      {
        path: 'ai-providers',
        name: 'AiProviders',
        component: () => import('@/views/admin/AiProviders.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'ai-chat',
        name: 'AiChat',
        component: () => import('@/views/admin/AiChat.vue'),
        meta: { requiresAuth: true }
      },
      // ... 其他管理路由保持不变 ...
    ]
  },
  // 404
  { path: '/:pathMatch(.*)*', name: 'NotFound', component: () => import('@/views/NotFound.vue') }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) return savedPosition
    return { top: 0 }
  }
})
```

- [ ] **Step 2: 修改 Header —— 条件显示后台入口**

`frontend/src/components/layout/Header.vue` — 修改 template 中的 header-actions 部分：

```vue
<div class="header-actions">
  <!-- 中英文切换 -->
  <LanguageSwitch />
  <!-- 主题切换 -->
  <ThemeSwitch />
  <!-- 后台入口：仅当 show_admin_link = true 且已登录时显示 -->
  <el-tooltip
    v-if="showAdminLink && isLoggedIn"
    :content="t('nav.adminTooltip')"
    placement="bottom"
  >
    <el-button text size="default" @click="goToAdmin" class="admin-btn">
      <el-icon><Setting /></el-icon>
    </el-button>
  </el-tooltip>
  <!-- 未登录时不显示任何后台入口 -->
  <el-button
    v-if="!isLoggedIn"
    type="primary"
    size="default"
    @click="showLoginDialog = true"
  >
    {{ t('login.title') }}
  </el-button>
  <LoginDialog v-model="showLoginDialog" @success="handleLoginSuccess" />
</div>
```

在 `<script setup>` 中新增：

```javascript
import { useLanguageStore } from '@/store/language'
import LanguageSwitch from '../common/LanguageSwitch.vue'
import ThemeSwitch from '../common/ThemeSwitch.vue'

const showAdminLink = ref(true) // 默认显示

// 从站点配置中获取是否显示后台入口
const fetchAdminConfig = async () => {
  try {
    const configs = await api.config.list({ type: 'site_config' })
    const adminPathConfig = configs.find(c => c.config_key === 'admin_path')
    const showLinkConfig = configs.find(c => c.config_key === 'show_admin_link')
    if (adminPathConfig) {
      localStorage.setItem('admin_path', adminPathConfig.config_value)
    }
    if (showLinkConfig) {
      showAdminLink.value = showLinkConfig.config_value === 'true'
    }
  } catch (e) {
    // 使用默认值
  }
}

onMounted(() => {
  fetchAdminConfig()
})
```

同时修改 `goToAdmin`:

```javascript
const goToAdmin = () => {
  const adminPath = localStorage.getItem('admin_path') || 'admin'
  router.push(`/${adminPath}`)
}
```

- [ ] **Step 3: 在数据库中插入默认配置**

创建迁移文件或初始化脚本，在 `system_configs` 表中插入默认值：

```sql
INSERT INTO system_configs (config_key, config_value, config_type, is_encrypted, is_active, description, created_at, updated_at)
VALUES
('admin_path', 'admin', 'site_config', false, true, '后台访问路径', NOW(), NOW()),
('show_admin_link', 'true', 'site_config', false, true, '是否在前端显示后台入口', NOW(), NOW())
ON CONFLICT (config_key) DO NOTHING;
```

或者通过 Go 代码在 `InitDefaultAdmin` 中初始化。

- [ ] **Step 4: 提交**

```bash
git add frontend/src/router/index.js frontend/src/components/layout/Header.vue
git commit -m "feat: configurable admin path and conditional admin link display"
```

---

### Task 1.5: 前端代码混淆配置

**Files:**
- Modify: `frontend/vite.config.js`
- Modify: `frontend/package.json`

- [ ] **Step 1: 安装混淆插件**

```bash
cd frontend && npm install --save-dev rollup-plugin-obfuscator
```

- [ ] **Step 2: 配置混淆**

`frontend/vite.config.js` — 在文件顶部新增导入，在 build.rollupOptions 中添加插件：

```javascript
import obfuscator from 'rollup-plugin-obfuscator'

export default defineConfig(({ mode }) => ({
  plugins: [
    vue(),
    AutoImport({ /* ...保持不变... */ }),
    Components({ /* ...保持不变... */ }),
    // 仅在 production build 时启用混淆
    ...(mode === 'production' ? [obfuscator({
      compact: true,
      controlFlowFlattening: true,
      controlFlowFlatteningThreshold: 0.75,
      deadCodeInjection: true,
      deadCodeInjectionThreshold: 0.4,
      stringArray: true,
      stringArrayEncoding: ['base64'],
      stringArrayThreshold: 0.75,
      selfDefending: true,
      debugProtection: true,
      // 部署时替换为实际域名
      domainLock: process.env.VITE_DOMAIN_LOCK
        ? process.env.VITE_DOMAIN_LOCK.split(',')
        : [],
      // 排除大型库文件，避免构建过慢
      exclude: ['node_modules/**'],
    })] : []),
  ],
  // ... 其他配置保持不变 ...
}))
```

**注意**：由于 `defineConfig` 当前是对象形式，需要改为函数形式以访问 `mode` 参数。实际改造：

```javascript
export default defineConfig(({ mode }) => {
  const isProd = mode === 'production'
  return {
    plugins: [
      vue(),
      // ... 其他插件 ...
      ...(isProd ? [obfuscator({
        compact: true,
        controlFlowFlattening: true,
        controlFlowFlatteningThreshold: 0.75,
        deadCodeInjection: true,
        deadCodeInjectionThreshold: 0.4,
        stringArray: true,
        stringArrayEncoding: ['base64'],
        stringArrayThreshold: 0.75,
        selfDefending: true,
        debugProtection: true,
        domainLock: [],
        exclude: ['node_modules/**'],
      })] : []),
    ],
    // ... resolve, server, build 等保持不变 ...
  }
})
```

- [ ] **Step 3: 更新 package.json 依赖**

在 `package.json` 的 `devDependencies` 中确认 `rollup-plugin-obfuscator` 已添加。

- [ ] **Step 4: 提交**

```bash
cd frontend && npm install --save-dev rollup-plugin-obfuscator
cd ..
git add frontend/vite.config.js frontend/package.json frontend/package-lock.json
git commit -m "feat: add frontend code obfuscation for production builds"
```

---

## Layer 2: 功能层

### Task 2.1: 删除爬虫监控功能

- [ ] **Step 1: 删除后端爬虫相关文件**

```bash
rm backend/internal/handler/crawler_handler.go
rm backend/internal/service/crawl_service.go
rm backend/internal/repository/crawl_task_repo.go
rm backend/internal/models/crawl_task.go
rm backend/internal/middleware/crawler_auth.go
rm -rf python-sdk/
```

- [ ] **Step 2: 清理 router.go 中的爬虫相关代码**

`backend/internal/router/router.go` — 删除以下内容：

删除 import 中的：
```go
// 删除 websocket 相关（如果仅用于爬虫）
// 注意：检查 websocket/hub.go 是否还用于其他功能
```

删除 repo 初始化：
```go
crawlTaskRepo := repository.NewCrawlTaskRepository(gormDB) // 删除
```

删除 service 初始化：
```go
wsHub := websocket.NewHub()      // 删除（如果仅用于爬虫）
go wsHub.Run()                    // 删除
crawlService := service.NewCrawlService(crawlTaskRepo, wsHub) // 删除
```

删除 handler 初始化：
```go
crawlerHandler := handler.NewCrawlerHandler(crawlService) // 删除
wsHandler := handler.NewWebSocketHandler(wsHub, jwtManager) // 删除（如果仅用于爬虫）
```

删除路由注册：
```go
// 爬虫任务接口（需要Bearer Token认证）
crawler := api.Group("/crawler")     // 删除整个 block
crawler.Use(middleware.CrawlerAuth()) // 删除
{ ... }                                // 删除

// 管理接口下的爬虫路由
admin.GET("/crawler/tasks", ...)     // 删除
admin.GET("/crawler/tasks/:task_id", ...) // 删除

// WebSocket 路由
r.GET("/ws/crawler/tasks", ...)      // 删除

// 生成爬虫Token接口
admin.POST("/configs/generate-crawler-token", ...) // 删除
```

- [ ] **Step 3: 删除前端爬虫相关文件**

```bash
rm frontend/src/api/crawler.js
rm frontend/src/views/admin/CrawlerMonitor.vue
```

- [ ] **Step 4: 清理前端路由中的爬虫页面**

`frontend/src/router/index.js` — 删除 crawler 路由项。

- [ ] **Step 5: 验证编译**

```bash
cd backend && go build ./...
```

- [ ] **Step 6: 提交**

```bash
git add -A
git commit -m "feat: remove crawler monitoring functionality, keep fingerprint collection"
```

---

### Task 2.2: AI 翻译 —— 数据模型 + 迁移

**Files:**
- Create: `backend/internal/models/ai_provider.go`
- Create: `backend/migrations/008_add_ai_tables.sql`
- Modify: `backend/internal/models/article.go`

- [ ] **Step 1: 新增 Article 翻译字段**

在 `backend/internal/models/article.go` 的 `Article` 结构体中，`LikeCount` 之后新增：

```go
// 英文翻译字段
TitleEn       string     `gorm:"type:varchar(500)" json:"title_en"`
ContentEn     string     `gorm:"type:text" json:"content_en"`
SummaryEn     string     `gorm:"type:varchar(500)" json:"summary_en"`
IsTranslated  bool       `gorm:"default:false" json:"is_translated"`
TranslatedAt  *time.Time `json:"translated_at"`
```

- [ ] **Step 2: 创建 AI 提供方模型**

`backend/internal/models/ai_provider.go`:

```go
package models

import (
	"time"
	"gorm.io/gorm"
)

type AIProvider struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(100);not null" json:"name"`
	ProviderType string         `gorm:"type:varchar(20);not null" json:"provider_type"` // claude, openai, deepseek, custom
	APIKey       string         `gorm:"type:text;not null" json:"api_key"`               // AES 加密存储
	BaseURL      string         `gorm:"type:varchar(500)" json:"base_url"`
	Model        string         `gorm:"type:varchar(100);not null" json:"model"`
	IsEnabled    bool           `gorm:"default:true;index" json:"is_enabled"`
	SortOrder    int            `gorm:"default:0" json:"sort_order"`
	Balance      float64        `gorm:"type:decimal(10,4);default:0" json:"balance"`
	LastCheckAt  *time.Time     `json:"last_check_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (AIProvider) TableName() string {
	return "ai_providers"
}

type AIChatHistory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ProviderID uint      `gorm:"index;not null" json:"provider_id"`
	Role       string    `gorm:"type:varchar(20);not null" json:"role"` // user / assistant
	Content    string    `gorm:"type:text;not null" json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

func (AIChatHistory) TableName() string {
	return "ai_chat_history"
}
```

- [ ] **Step 3: 创建数据库迁移 SQL**

`backend/migrations/008_add_ai_tables.sql`:

```sql
-- 文章翻译字段
ALTER TABLE articles
ADD COLUMN IF NOT EXISTS title_en VARCHAR(500) DEFAULT '',
ADD COLUMN IF NOT EXISTS content_en TEXT DEFAULT '',
ADD COLUMN IF NOT EXISTS summary_en VARCHAR(500) DEFAULT '',
ADD COLUMN IF NOT EXISTS is_translated BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS translated_at TIMESTAMP NULL;

-- AI 提供方表
CREATE TABLE IF NOT EXISTS ai_providers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    provider_type VARCHAR(20) NOT NULL,
    api_key TEXT NOT NULL,
    base_url VARCHAR(500) DEFAULT '',
    model VARCHAR(100) NOT NULL,
    is_enabled BOOLEAN DEFAULT TRUE,
    sort_order INT DEFAULT 0,
    balance DECIMAL(10,4) DEFAULT 0,
    last_check_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_ai_providers_enabled ON ai_providers(is_enabled);
CREATE INDEX IF NOT EXISTS idx_ai_providers_deleted ON ai_providers(deleted_at);

-- AI 聊天记录表
CREATE TABLE IF NOT EXISTS ai_chat_history (
    id SERIAL PRIMARY KEY,
    provider_id INT NOT NULL REFERENCES ai_providers(id),
    role VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_ai_chat_provider ON ai_chat_history(provider_id);
CREATE INDEX IF NOT EXISTS idx_ai_chat_created ON ai_chat_history(created_at);
```

- [ ] **Step 4: 提交**

```bash
git add backend/internal/models/article.go backend/internal/models/ai_provider.go backend/migrations/008_add_ai_tables.sql
git commit -m "feat: add article translation fields and AI provider models"
```

---

### Task 2.3: AI 服务 —— 翻译 + 聊天

**Files:**
- Create: `backend/internal/repository/ai_provider_repo.go`
- Create: `backend/internal/service/ai_service.go`
- Create: `backend/internal/handler/ai_handler.go`
- Modify: `backend/internal/router/router.go`

- [ ] **Step 1: 创建 AI 提供方 Repository**

`backend/internal/repository/ai_provider_repo.go`:

```go
package repository

import (
	"errors"
	"github.com/whk-newbie/blog/internal/models"
	"gorm.io/gorm"
)

var (
	ErrAIProviderNotFound = errors.New("AI provider not found")
)

type AIProviderRepository interface {
	FindAll() ([]*models.AIProvider, error)
	FindEnabled() ([]*models.AIProvider, error)
	FindByID(id uint) (*models.AIProvider, error)
	Create(provider *models.AIProvider) error
	Update(provider *models.AIProvider) error
	Delete(id uint) error
}

type aiProviderRepository struct {
	db *gorm.DB
}

func NewAIProviderRepository(db *gorm.DB) AIProviderRepository {
	return &aiProviderRepository{db: db}
}

func (r *aiProviderRepository) FindAll() ([]*models.AIProvider, error) {
	var providers []*models.AIProvider
	err := r.db.Order("sort_order ASC").Find(&providers).Error
	return providers, err
}

func (r *aiProviderRepository) FindEnabled() ([]*models.AIProvider, error) {
	var providers []*models.AIProvider
	err := r.db.Where("is_enabled = ?", true).Order("sort_order ASC").Find(&providers).Error
	return providers, err
}

func (r *aiProviderRepository) FindByID(id uint) (*models.AIProvider, error) {
	var provider models.AIProvider
	err := r.db.First(&provider, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAIProviderNotFound
	}
	return &provider, err
}

func (r *aiProviderRepository) Create(provider *models.AIProvider) error {
	return r.db.Create(provider).Error
}

func (r *aiProviderRepository) Update(provider *models.AIProvider) error {
	return r.db.Save(provider).Error
}

func (r *aiProviderRepository) Delete(id uint) error {
	return r.db.Delete(&models.AIProvider{}, id).Error
}
```

- [ ] **Step 2: 创建 AI 服务（翻译 + 聊天）**

`backend/internal/service/ai_service.go`:

```go
package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/whk-newbie/blog/internal/models"
	"github.com/whk-newbie/blog/internal/pkg/crypto"
	"github.com/whk-newbie/blog/internal/repository"
)

var (
	ErrNoEnabledProvider = errors.New("no enabled AI provider")
	ErrTranslationFailed = errors.New("translation failed")
)

type AIService interface {
	// 翻译
	Translate(title, content, summary string, providerID *uint) (*TranslateResult, error)
	// 聊天（流式 SSE）
	ChatStream(providerID uint, messages []ChatMessage, writer io.Writer) error
	// 检测提供方连通性和余额
	CheckProvider(id uint) (*ProviderCheckResult, error)
	// 提供方 CRUD
	ListProviders() ([]*models.AIProvider, error)
	GetProvider(id uint) (*models.AIProvider, error)
	CreateProvider(req *CreateAIProviderRequest) (*models.AIProvider, error)
	UpdateProvider(id uint, req *UpdateAIProviderRequest) (*models.AIProvider, error)
	DeleteProvider(id uint) error
}

// TranslateResult 翻译结果
type TranslateResult struct {
	TitleEn   string `json:"title_en"`
	ContentEn string `json:"content_en"`
	SummaryEn string `json:"summary_en"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Role    string `json:"role"`    // system / user / assistant
	Content string `json:"content"`
}

// ProviderCheckResult 提供方检测结果
type ProviderCheckResult struct {
	Available bool    `json:"available"`
	Balance   float64 `json:"balance"`
	Error     string  `json:"error,omitempty"`
}

type CreateAIProviderRequest struct {
	Name         string  `json:"name" binding:"required"`
	ProviderType string  `json:"provider_type" binding:"required"`
	APIKey       string  `json:"api_key" binding:"required"`
	BaseURL      string  `json:"base_url"`
	Model        string  `json:"model" binding:"required"`
	IsEnabled    bool    `json:"is_enabled"`
	SortOrder    int     `json:"sort_order"`
}

type UpdateAIProviderRequest struct {
	Name      *string `json:"name"`
	APIKey    *string `json:"api_key"`
	BaseURL   *string `json:"base_url"`
	Model     *string `json:"model"`
	IsEnabled *bool   `json:"is_enabled"`
	SortOrder *int    `json:"sort_order"`
}

type aiService struct {
	providerRepo repository.AIProviderRepository
	crypto       *crypto.Crypto
}

func NewAIService(providerRepo repository.AIProviderRepository, masterKey string) (AIService, error) {
	c, err := crypto.NewCrypto(masterKey)
	if err != nil {
		return nil, err
	}
	return &aiService{providerRepo: providerRepo, crypto: c}, nil
}

// Translate 翻译文章
func (s *aiService) Translate(title, content, summary string, providerID *uint) (*TranslateResult, error) {
	// 获取启用的提供方
	providers, err := s.providerRepo.FindEnabled()
	if err != nil || len(providers) == 0 {
		return nil, ErrNoEnabledProvider
	}

	// 尝试每个启用的提供方（按 sort_order）
	var lastErr error
	for _, p := range providers {
		if providerID != nil && p.ID != *providerID {
			continue
		}
		result, err := s.translateWithProvider(p, title, content, summary)
		if err == nil {
			return result, nil
		}
		lastErr = err
	}

	if lastErr != nil {
		return nil, fmt.Errorf("%w: %v", ErrTranslationFailed, lastErr)
	}
	return nil, ErrTranslationFailed
}

func (s *aiService) translateWithProvider(provider *models.AIProvider, title, content, summary string) (*TranslateResult, error) {
	// 解密 API Key
	apiKey, err := s.crypto.Decrypt(provider.APIKey)
	if err != nil {
		apiKey = provider.APIKey // 回退：可能未加密
	}

	prompt := fmt.Sprintf(`Translate the following blog article from Chinese to English. Return ONLY valid JSON with keys: title_en, content_en, summary_en.

Chinese Title: %s
Chinese Summary: %s
Chinese Content (HTML, preserve HTML tags):
%s`, title, summary, content)

	baseURL := s.getBaseURL(provider)
	reqBody := map[string]interface{}{
		"model": provider.Model,
		"messages": []map[string]string{
			{"role": "system", "content": "You are a professional translator. Translate Chinese to natural, fluent English. Preserve all HTML tags. Return ONLY valid JSON."},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.3,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", baseURL, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	switch provider.ProviderType {
	case "claude":
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("anthropic-version", "2023-06-01")
	case "openai", "deepseek", "custom":
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应（统一 OpenAI 格式）
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Choices) == 0 {
		return nil, errors.New("empty response from AI")
	}

	// 从 content 中提取 JSON
	content := result.Choices[0].Message.Content
	content = extractJSON(content)

	var translateResult TranslateResult
	if err := json.Unmarshal([]byte(content), &translateResult); err != nil {
		return nil, fmt.Errorf("failed to parse translation result: %w", err)
	}

	return &translateResult, nil
}

// ChatStream 流式聊天（SSE）
func (s *aiService) ChatStream(providerID uint, messages []ChatMessage, writer io.Writer) error {
	provider, err := s.providerRepo.FindByID(providerID)
	if err != nil {
		return err
	}

	apiKey, _ := s.crypto.Decrypt(provider.APIKey)
	if apiKey == "" {
		apiKey = provider.APIKey
	}

	baseURL := s.getBaseURL(provider)

	msgs := make([]map[string]string, len(messages))
	for i, m := range messages {
		msgs[i] = map[string]string{"role": m.Role, "content": m.Content}
	}

	reqBody := map[string]interface{}{
		"model":    provider.Model,
		"messages": msgs,
		"stream":   true,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", baseURL, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	switch provider.ProviderType {
	case "claude":
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("anthropic-version", "2023-06-01")
	default:
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{Timeout: 300 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("chat request failed: %w", err)
	}
	defer resp.Body.Close()

	// 转发 SSE 流
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				fmt.Fprintf(writer, "data: [DONE]\n\n")
				break
			}
			// 转发 SSE data
			fmt.Fprintf(writer, "data: %s\n\n", data)
		}
	}

	return nil
}

// CheckProvider 检测提供方连通性和余额
func (s *aiService) CheckProvider(id uint) (*ProviderCheckResult, error) {
	provider, err := s.providerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	apiKey, _ := s.crypto.Decrypt(provider.APIKey)
	if apiKey == "" {
		apiKey = provider.APIKey
	}

	// 发送一个简单的测试请求
	baseURL := s.getBaseURL(provider)
	reqBody := map[string]interface{}{
		"model":       provider.Model,
		"messages":    []map[string]string{{"role": "user", "content": "ping"}},
		"max_tokens":  1,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", baseURL, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	switch provider.ProviderType {
	case "claude":
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("anthropic-version", "2023-06-01")
	default:
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return &ProviderCheckResult{Available: false, Error: err.Error()}, nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	result := &ProviderCheckResult{Available: resp.StatusCode == 200}
	if resp.StatusCode != 200 {
		result.Error = string(body)
	}

	// 更新最后检测时间
	now := time.Now()
	provider.LastCheckAt = &now
	_ = s.providerRepo.Update(provider)

	return result, nil
}

func (s *aiService) getBaseURL(provider *models.AIProvider) string {
	if provider.BaseURL != "" {
		return provider.BaseURL
	}
	switch provider.ProviderType {
	case "claude":
		return "https://api.anthropic.com/v1/messages"
	case "openai":
		return "https://api.openai.com/v1/chat/completions"
	case "deepseek":
		return "https://api.deepseek.com/v1/chat/completions"
	default:
		return provider.BaseURL
	}
}

// CRUD methods...
func (s *aiService) ListProviders() ([]*models.AIProvider, error) {
	return s.providerRepo.FindAll()
}

func (s *aiService) GetProvider(id uint) (*models.AIProvider, error) {
	return s.providerRepo.FindByID(id)
}

func (s *aiService) CreateProvider(req *CreateAIProviderRequest) (*models.AIProvider, error) {
	// 加密 API Key
	encryptedKey, err := s.crypto.Encrypt(req.APIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt API key: %w", err)
	}

	provider := &models.AIProvider{
		Name:         req.Name,
		ProviderType: req.ProviderType,
		APIKey:       encryptedKey,
		BaseURL:      req.BaseURL,
		Model:        req.Model,
		IsEnabled:    req.IsEnabled,
		SortOrder:    req.SortOrder,
	}

	if err := s.providerRepo.Create(provider); err != nil {
		return nil, err
	}
	return provider, nil
}

func (s *aiService) UpdateProvider(id uint, req *UpdateAIProviderRequest) (*models.AIProvider, error) {
	provider, err := s.providerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		provider.Name = *req.Name
	}
	if req.APIKey != nil {
		encrypted, err := s.crypto.Encrypt(*req.APIKey)
		if err != nil {
			return nil, err
		}
		provider.APIKey = encrypted
	}
	if req.BaseURL != nil {
		provider.BaseURL = *req.BaseURL
	}
	if req.Model != nil {
		provider.Model = *req.Model
	}
	if req.IsEnabled != nil {
		provider.IsEnabled = *req.IsEnabled
	}
	if req.SortOrder != nil {
		provider.SortOrder = *req.SortOrder
	}

	if err := s.providerRepo.Update(provider); err != nil {
		return nil, err
	}
	return provider, nil
}

func (s *aiService) DeleteProvider(id uint) error {
	return s.providerRepo.Delete(id)
}

// extractJSON 从 AI 返回的文本中提取 JSON
func extractJSON(content string) string {
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "```json") {
		content = strings.TrimPrefix(content, "```json")
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)
	}
	if strings.HasPrefix(content, "```") {
		content = strings.TrimPrefix(content, "```")
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)
	}
	return content
}
```

- [ ] **Step 3: 创建 AI Handler**

`backend/internal/handler/ai_handler.go`:

```go
package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/pkg/response"
	"github.com/whk-newbie/blog/internal/service"
)

type AIHandler struct {
	aiService service.AIService
}

func NewAIHandler(aiService service.AIService) *AIHandler {
	return &AIHandler{aiService: aiService}
}

// TranslateArticle 翻译文章
func (h *AIHandler) TranslateArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid article ID")
		return
	}

	var req struct {
		Title      string `json:"title" binding:"required"`
		Content    string `json:"content" binding:"required"`
		Summary    string `json:"summary"`
		ProviderID *uint  `json:"provider_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	result, err := h.aiService.Translate(req.Title, req.Content, req.Summary, req.ProviderID)
	if err != nil {
		response.InternalServerError(c, "translation failed: "+err.Error())
		return
	}

	_ = id // 文章 ID 用于日志，实际翻译结果由调用方（article handler）保存

	response.Success(c, result)
}

// Chat 聊天（SSE 流式）
func (h *AIHandler) Chat(c *gin.Context) {
	var req struct {
		ProviderID uint                   `json:"provider_id" binding:"required"`
		Messages   []service.ChatMessage  `json:"messages" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	// 设置 SSE headers
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.WriteHeader(http.StatusOK)

	c.Stream(func(w io.Writer) bool {
		err := h.aiService.ChatStream(req.ProviderID, req.Messages, w)
		if err != nil {
			fmt.Fprintf(w, "data: {\"error\": \"%s\"}\n\n", err.Error())
		}
		return false
	})
}

// ListProviders 获取 AI 提供方列表
func (h *AIHandler) ListProviders(c *gin.Context) {
	providers, err := h.aiService.ListProviders()
	if err != nil {
		response.InternalServerError(c, "failed to list providers")
		return
	}
	response.Success(c, providers)
}

// GetProvider 获取单个提供方
func (h *AIHandler) GetProvider(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid ID")
		return
	}
	provider, err := h.aiService.GetProvider(uint(id))
	if err != nil {
		response.NotFound(c, "provider not found")
		return
	}
	response.Success(c, provider)
}

// CreateProvider 创建 AI 提供方
func (h *AIHandler) CreateProvider(c *gin.Context) {
	var req service.CreateAIProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}
	provider, err := h.aiService.CreateProvider(&req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.Created(c, "created", provider)
}

// UpdateProvider 更新 AI 提供方
func (h *AIHandler) UpdateProvider(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	var req service.UpdateAIProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}
	provider, err := h.aiService.UpdateProvider(uint(id), &req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.Success(c, provider)
}

// DeleteProvider 删除 AI 提供方
func (h *AIHandler) DeleteProvider(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	if err := h.aiService.DeleteProvider(uint(id)); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.NoContent(c, "deleted")
}

// CheckProvider 检测提供方连通性
func (h *AIHandler) CheckProvider(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	result, err := h.aiService.CheckProvider(uint(id))
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.Success(c, result)
}
```

**注意**：`Chat` handler 中的 `io.Writer` 需要正确处理。Gin 的 `c.Stream()` 接受 `func(w io.Writer) bool` 回调，可以直接使用。

- [ ] **Step 4: 在 router.go 中注册 AI 路由**

在 `backend/internal/router/router.go` 的 `admin` 路由组中添加：

```go
// AI 提供方管理
admin.GET("/ai/providers", aiHandler.ListProviders)
admin.GET("/ai/providers/:id", aiHandler.GetProvider)
admin.POST("/ai/providers", aiHandler.CreateProvider)
admin.PUT("/ai/providers/:id", aiHandler.UpdateProvider)
admin.DELETE("/ai/providers/:id", aiHandler.DeleteProvider)
admin.POST("/ai/providers/:id/check", aiHandler.CheckProvider)

// AI 翻译
admin.POST("/ai/translate/:id", aiHandler.TranslateArticle)

// AI 聊天
admin.POST("/ai/chat", aiHandler.Chat)
```

同时添加初始化代码：

```go
aiProviderRepo := repository.NewAIProviderRepository(gormDB)
aiService, err := service.NewAIService(aiProviderRepo, cfg.Crypto.MasterKey)
if err != nil {
    panic("Failed to initialize AI service: " + err.Error())
}
aiHandler := handler.NewAIHandler(aiService)
```

- [ ] **Step 5: 修改文章发布逻辑 —— 集成翻译触发**

在 `backend/internal/service/article_service.go` 的 `CreateArticleRequest` 中新增字段：

```go
type CreateArticleRequest struct {
	// ... 原有字段 ...
	TranslateToEn bool  `json:"translate_to_en"` // 是否翻译为英文
}
```

在 `Create` 方法中，保存文章后异步翻译：

```go
// 在 articleService 中注入 AIService（通过构造函数或 Setter）
// 发布文章后：
if req.TranslateToEn && article.Status == models.ArticleStatusPublished {
    go func() {
        result, err := s.aiService.Translate(article.Title, article.Content, article.Summary, nil)
        if err == nil {
            updates := map[string]interface{}{
                "title_en":      result.TitleEn,
                "content_en":    result.ContentEn,
                "summary_en":    result.SummaryEn,
                "is_translated": true,
                "translated_at": time.Now(),
            }
            _ = s.articleRepo.UpdateFields(article.ID, updates)
        }
    }()
}
```

同时在文章 handler 中添加手动翻译接口：

```go
// POST /admin/articles/:id/translate
admin.POST("/articles/:id/translate", articleHandler.TranslateArticle)
```

- [ ] **Step 6: 提交**

```bash
git add backend/internal/repository/ai_provider_repo.go backend/internal/service/ai_service.go backend/internal/handler/ai_handler.go backend/internal/router/router.go backend/internal/service/article_service.go backend/internal/handler/article_handler.go
git commit -m "feat: add AI translation service, chat streaming, and provider management"
```

---

### Task 2.4: 前端 AI 管理页面 + 聊天框

**Files:**
- Create: `frontend/src/api/ai.js`
- Create: `frontend/src/views/admin/AiProviders.vue`
- Create: `frontend/src/views/admin/AiChat.vue`
- Create: `frontend/src/components/chat/AiChatPanel.vue`
- Modify: `frontend/src/views/admin/ArticleEditor.vue`

- [ ] **Step 1: 创建 AI API 封装**

`frontend/src/api/ai.js`:

```javascript
import http from './http'

export default {
  // 提供方管理
  listProviders() {
    return http.get('/admin/ai/providers')
  },
  getProvider(id) {
    return http.get(`/admin/ai/providers/${id}`)
  },
  createProvider(data) {
    return http.post('/admin/ai/providers', data)
  },
  updateProvider(id, data) {
    return http.put(`/admin/ai/providers/${id}`, data)
  },
  deleteProvider(id) {
    return http.delete(`/admin/ai/providers/${id}`)
  },
  checkProvider(id) {
    return http.post(`/admin/ai/providers/${id}/check`)
  },

  // 翻译
  translateArticle(articleId, data) {
    return http.post(`/admin/ai/translate/${articleId}`, data)
  },

  // 聊天（SSE）
  chat(providerId, messages) {
    // SSE 需要特殊处理，使用 fetch
    return fetch('/api/v1/admin/ai/chat', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'X-Session-Id': window.__sessionId || '',
      },
      body: JSON.stringify({ provider_id: providerId, messages }),
    })
  },
}
```

- [ ] **Step 2: 创建 AI 提供方管理页面**

`frontend/src/views/admin/AiProviders.vue`（简化结构）：

```vue
<template>
  <div class="ai-providers-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>AI 提供方管理</span>
          <el-button type="primary" @click="showCreate = true">添加提供方</el-button>
        </div>
      </template>

      <el-table :data="providers" v-loading="loading">
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="provider_type" label="类型" width="100" />
        <el-table-column prop="model" label="模型" width="150" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.is_enabled ? 'success' : 'info'">
              {{ row.is_enabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort_order" label="优先级" width="80" />
        <el-table-column label="操作" width="300">
          <template #default="{ row }">
            <el-button size="small" @click="editProvider(row)">编辑</el-button>
            <el-button size="small" @click="checkProvider(row)" :loading="checkingId === row.id">检测</el-button>
            <el-button size="small" type="danger" @click="deleteProvider(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 创建/编辑弹窗 -->
    <el-dialog v-model="showDialog" :title="editingProvider ? '编辑提供方' : '添加提供方'" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.provider_type">
            <el-option label="Claude" value="claude" />
            <el-option label="OpenAI" value="openai" />
            <el-option label="DeepSeek" value="deepseek" />
            <el-option label="自定义" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="form.api_key" type="password" show-password />
        </el-form-item>
        <el-form-item label="Base URL" v-if="form.provider_type === 'custom'">
          <el-input v-model="form.base_url" placeholder="https://api.example.com/v1/chat/completions" />
        </el-form-item>
        <el-form-item label="模型">
          <el-input v-model="form.model" placeholder="deepseek-chat" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="form.is_enabled" />
        </el-form-item>
        <el-form-item label="优先级">
          <el-input-number v-model="form.sort_order" :min="0" :max="100" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="saveProvider">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>
```

- [ ] **Step 3: 创建 AI 聊天页面**

`frontend/src/views/admin/AiChat.vue`（聊天对话界面）：

```vue
<template>
  <div class="ai-chat-page">
    <el-card class="chat-container">
      <template #header>
        <div class="chat-header">
          <span>AI 助手</span>
          <el-select v-model="selectedProviderId" placeholder="选择 AI 提供方" style="width: 200px">
            <el-option v-for="p in enabledProviders" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </div>
      </template>

      <div class="chat-messages" ref="messagesContainer">
        <div v-for="(msg, i) in messages" :key="i" :class="['message', msg.role]">
          <div class="message-content" v-html="renderMarkdown(msg.content)"></div>
        </div>
        <div v-if="streaming" class="message assistant">
          <div class="message-content" v-html="renderMarkdown(streamContent)"></div>
        </div>
      </div>

      <div class="chat-input">
        <div class="quick-actions">
          <el-button size="small" @click="quickAction('请帮我润色以下内容：')">润色</el-button>
          <el-button size="small" @click="quickAction('请帮我扩写以下内容：')">扩写</el-button>
          <el-button size="small" @click="quickAction('请帮我缩写以下内容：')">缩写</el-button>
          <el-button size="small" @click="quickAction('请帮我总结以下内容：')">总结</el-button>
          <el-button size="small" @click="quickAction('请将以下内容翻译成英文：')">翻译</el-button>
        </div>
        <el-input
          v-model="inputText"
          type="textarea"
          :rows="3"
          placeholder="输入消息..."
          @keydown.enter.exact="sendMessage"
        />
        <el-button type="primary" @click="sendMessage" :loading="streaming">发送</el-button>
      </div>
    </el-card>
  </div>
</template>
```

- [ ] **Step 4: 在 ArticleEditor 中添加翻译复选框**

`frontend/src/views/admin/ArticleEditor.vue` — 在发布按钮附近添加：

```vue
<el-checkbox v-model="translateToEn" v-if="form.status === 'published'">
  发布时翻译为英文
</el-checkbox>
```

在保存/发布逻辑中添加：

```javascript
if (translateToEn.value && form.status === 'published') {
  form.translate_to_en = true
}
```

同时在文章列表中为已发布但未翻译的文章提供「翻译」按钮。

- [ ] **Step 5: 提交**

```bash
git add frontend/src/api/ai.js frontend/src/views/admin/AiProviders.vue frontend/src/views/admin/AiChat.vue frontend/src/components/chat/AiChatPanel.vue frontend/src/views/admin/ArticleEditor.vue
git commit -m "feat: add AI provider management, chat, and translation UI"
```

---

## Layer 3: UI 改版

由于 Layer 3 涉及大量前端文件重写，每个任务列出关键改动和目标文件。完整代码在实施时参照设计文档的 UI 规格。

### Task 3.1: 全局布局重构（侧边栏 + 内容区）

**Files:**
- Modify: `frontend/src/components/layout/MainLayout.vue`
- Create: `frontend/src/components/layout/Sidebar.vue`
- Modify: `frontend/src/components/layout/Header.vue`
- Modify: `frontend/src/components/layout/Footer.vue`
- Modify: `frontend/src/assets/styles/variables.less`
- Modify: `frontend/src/assets/styles/common.less`

- [ ] **Step 1: 更新配色变量**

`frontend/src/assets/styles/variables.less`:

```less
// zyyo.net 风格配色
:root {
  // 浅色模式
  --bg-color: #fafafa;
  --card-bg: #ffffff;
  --text-color: #333333;
  --text-secondary: #888888;
  --text-heading: #1a1a1a;
  --link-color: #4078c0;
  --border-color: #e8e8e8;
  --tag-color-1: #e8f5e9;
  --tag-text-1: #2e7d32;
  --tag-color-2: #e3f2fd;
  --tag-text-2: #1565c0;
  --tag-color-3: #fff3e0;
  --tag-text-3: #e65100;
  --tag-color-4: #fce4ec;
  --tag-text-4: #c62828;
  --sidebar-width: 280px;
  --header-height: 56px;
  --shadow-sm: 0 1px 3px rgba(0, 0, 0, 0.06);
  --shadow-md: 0 4px 12px rgba(0, 0, 0, 0.08);
  --radius-sm: 4px;
  --radius-md: 8px;
}

// 暗色模式
html.dark {
  --bg-color: #1a1a2e;
  --card-bg: #252540;
  --text-color: #e0e0e0;
  --text-secondary: #888888;
  --text-heading: #f0f0f0;
  --link-color: #6db3f2;
  --border-color: #333350;
  --shadow-sm: 0 1px 3px rgba(0, 0, 0, 0.2);
  --shadow-md: 0 4px 12px rgba(0, 0, 0, 0.3);
}
```

- [ ] **Step 2: 重写 MainLayout —— 双栏布局**

`frontend/src/components/layout/MainLayout.vue`:

```vue
<template>
  <div class="main-layout">
    <Header />
    <div class="layout-body">
      <Sidebar />
      <main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>
    <Footer />
    <SearchOverlay />
  </div>
</template>
```

CSS: `.layout-body` 使用 `display: flex`，侧边栏固定 280px，内容区 `flex: 1`。

- [ ] **Step 3: 创建 Sidebar 组件**

`frontend/src/components/layout/Sidebar.vue` — 包含：
- 头像 + 简介
- 统计计数（文章数、运行天数）
- 标签云
- 社交链接图标（GitHub、Email 等，可从系统配置读取）

```vue
<template>
  <aside class="sidebar">
    <div class="sidebar-inner">
      <!-- 头像和简介 -->
      <div class="profile">
        <img :src="avatarUrl" class="avatar" alt="avatar" />
        <h2 class="name">{{ blogTitle }}</h2>
        <p class="bio">{{ blogDescription }}</p>
      </div>

      <!-- 统计 -->
      <div class="stats">
        <div class="stat-item">
          <span class="stat-number">{{ articleCount }}</span>
          <span class="stat-label">篇文章</span>
        </div>
        <div class="stat-item">
          <span class="stat-number">{{ runningDays }}</span>
          <span class="stat-label">天</span>
        </div>
      </div>

      <!-- 标签云 -->
      <div class="tag-cloud">
        <router-link
          v-for="tag in tags"
          :key="tag.id"
          :to="`/articles?tag_id=${tag.id}`"
          class="tag-item"
          :style="{ fontSize: getTagSize(tag.article_count) + 'px' }"
        >
          #{{ tag.name }}
        </router-link>
      </div>

      <!-- 社交链接 -->
      <div class="social-links">
        <a v-if="githubUrl" :href="githubUrl" target="_blank" title="GitHub">
          <svg><!-- GitHub icon --></svg>
        </a>
        <a v-if="email" :href="`mailto:${email}`" title="Email">
          <svg><!-- Email icon --></svg>
        </a>
      </div>
    </div>
  </aside>
</template>
```

响应式：`@media (max-width: 768px)` 时侧边栏隐藏，通过顶部汉堡菜单切换。

- [ ] **Step 4: 更新 Header —— 极简顶栏**

Header 简化为仅展示分类链接 + 搜索图标 + 语言切换 + 主题切换。移除大标题（移至侧边栏）。

- [ ] **Step 5: 更新 Footer**

Footer 简化为单行：© 年份 + 博客名称 + ICP 备案号。

- [ ] **Step 6: 提交**

```bash
git add frontend/src/assets/styles/ frontend/src/components/layout/
git commit -m "feat: redesign global layout with sidebar + content area"
```

---

### Task 3.2: 首页文章卡片流 + 搜索 + 时间轴

**Files:**
- Modify: `frontend/src/views/Home.vue`
- Modify: `frontend/src/components/article/ArticleCard.vue`
- Create: `frontend/src/components/common/SearchOverlay.vue`
- Create: `frontend/src/views/Timeline.vue`
- Modify: `frontend/src/router/index.js`

- [ ] **Step 1: 重写 ArticleCard —— zyyo.net 卡片风格**

`frontend/src/components/article/ArticleCard.vue` — 极简卡片：
- 16:9 封面图（有则显示，无则显示占位色块）
- 标题（单行截断）
- 日期 + 分类标签
- 悬停：卡片微升 2px + 阴影加深

```vue
<template>
  <article class="article-card" @click="$emit('click')">
    <div class="card-cover">
      <img v-if="article.cover_image" :src="article.cover_image" :alt="article.title" />
      <div v-else class="cover-placeholder"></div>
    </div>
    <div class="card-body">
      <h3 class="card-title">{{ displayTitle }}</h3>
      <div class="card-meta">
        <time>{{ formatDate(article.publish_at || article.created_at) }}</time>
        <span v-if="article.category" class="category-badge">{{ article.category.name }}</span>
      </div>
    </div>
  </article>
</template>
```

`displayTitle` 根据当前语言选择 `article.title` 或 `article.title_en`。

- [ ] **Step 2: 重写 Home.vue**

移除旧的 `Home.vue` 全部内容，改为：
- 文章卡片网格流（从 API 获取全部已发布文章）
- 每行 1-2 列（响应式）

- [ ] **Step 3: 创建全屏搜索弹窗**

`frontend/src/components/common/SearchOverlay.vue`:
- 点击搜索图标弹出全屏遮罩
- 输入框居中自动聚焦
- 使用 Fuse.js 前端模糊搜索（标题 + 摘要）
- 搜索结果列表，点击跳转

安装依赖：
```bash
cd frontend && npm install fuse.js
```

- [ ] **Step 4: 创建时间轴页面**

`frontend/src/views/Timeline.vue`:
- 按「年 → 月」两级分组
- 左侧时间线 + 圆点
- 每篇文章：标题 + 日期，点击跳转
- 从 API 获取全部文章，前端按 `publish_at` 分组排序

```vue
<template>
  <div class="timeline-page">
    <h1 class="page-title">时间轴</h1>
    <div class="timeline">
      <div v-for="year in timelineData" :key="year.year" class="year-group">
        <div class="year-dot">{{ year.year }}</div>
        <div v-for="month in year.months" :key="month.month" class="month-group">
          <h3 class="month-label">{{ month.month }}月</h3>
          <div v-for="article in month.articles" :key="article.id" class="timeline-item" @click="goToArticle(article.slug)">
            <span class="item-dot"></span>
            <span class="item-title">{{ displayTitle(article) }}</span>
            <span class="item-date">{{ formatDay(article.publish_at) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
```

- [ ] **Step 5: 更新路由 —— 添加时间轴**

`frontend/src/router/index.js` — 在 MainLayout 的 children 中添加：

```javascript
{
  path: 'timeline',
  name: 'Timeline',
  component: () => import('@/views/Timeline.vue'),
  meta: { titleKey: 'nav.timeline' }
}
```

- [ ] **Step 6: 提交**

```bash
cd frontend && npm install fuse.js
cd ..
git add frontend/src/views/Home.vue frontend/src/views/Timeline.vue frontend/src/components/article/ArticleCard.vue frontend/src/components/common/SearchOverlay.vue frontend/src/router/index.js
git commit -m "feat: redesign homepage cards, add timeline page and search overlay"
```

---

### Task 3.3: 文章详情页 + 文章列表页 + 工具页 + 后台极简化

**Files:**
- Modify: `frontend/src/views/ArticleDetail.vue`
- Modify: `frontend/src/views/Articles.vue`
- Modify: `frontend/src/views/Tools.vue`
- Modify: `frontend/src/views/admin/Dashboard.vue`（可选）

- [ ] **Step 1: 重写 ArticleDetail.vue**

极简文章详情：
- 大标题（支持中英文切换）
- 元信息行：日期 · 分类 · 标签
- 正文：Markdown 渲染区域，自定义 typography 样式
- 代码块：highlight.js，暗色/亮色主题匹配
- 底部：版权声明、许可信息

- [ ] **Step 2: 重写 Articles.vue**

文章列表页：
- 顶部：分类/标签筛选
- 文章卡片流（复用 ArticleCard）
- 分页

- [ ] **Step 3: 重写 Tools.vue**

工具页极简化：
- 保留所有现有工具功能
- 容器改为卡片风格（白底 + 极淡阴影 + 4px 圆角）
- 表单控件保留 Element Plus

- [ ] **Step 4: 后台管理页面风格微调**

`frontend/src/views/admin/*.vue` — 将主色调与前台统一：
- 页面背景色改为 `#fafafa`
- 卡片圆角改为 4px
- 表格/表单保留 Element Plus 功能，样式微调

- [ ] **Step 5: 提交**

```bash
git add frontend/src/views/
git commit -m "feat: redesign article detail, list, tools pages with minimal style"
```

---

## 最终验证

### 编译验证

```bash
# 后端
cd backend && go build ./... && go vet ./...

# 前端
cd frontend && npm run build
```

### 功能验证清单

- [ ] RSA 公钥接口返回正确的 PEM 格式
- [ ] AES 密钥协商成功，Redis 中有对应记录
- [ ] 加密后的请求能被后端正确解密
- [ ] 加密后的响应能被前端正确解密
- [ ] Session 过期后自动重新协商
- [ ] 后台入口路径配置生效
- [ ] `show_admin_link=false` 时前端不显示后台入口
- [ ] 前端 build 产物经过混淆（查看 dist 中的 JS 文件）
- [ ] 爬虫相关代码完全移除，编译通过
- [ ] AI 提供方 CRUD 正常
- [ ] AI 提供方检测返回连通性和余额
- [ ] 文章发布时勾选翻译，异步翻译并保存
- [ ] 手动触发翻译功能正常
- [ ] 中英文切换后文章内容正确显示
- [ ] AI 聊天 SSE 流式返回正常
- [ ] 新 UI 布局：侧边栏 + 内容区
- [ ] 暗色模式切换全站同步
- [ ] 时间轴按年月分组显示正确
- [ ] 搜索弹窗功能正确
- [ ] 响应式：移动端侧边栏收起
