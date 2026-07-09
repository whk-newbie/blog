package middleware

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/pkg/crypto"
	"github.com/whk-newbie/blog/internal/pkg/logger"
	"github.com/whk-newbie/blog/internal/pkg/response"
)

// EncryptionWhitelist paths that skip encryption
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

// Encryption middleware: decrypts request body, encrypts response body
func Encryption(rsaKeyPair *crypto.RSAKeyPair) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isWhitelisted(c.Request.URL.Path) {
			c.Next()
			return
		}

		sessionID := c.GetHeader("X-Session-Id")
		if sessionID == "" {
			response.Error(c, 40001, "missing session ID, please negotiate encryption key first")
			c.Abort()
			return
		}

		aesKey, err := crypto.GetSessionKey(sessionID)
		if err != nil {
			response.Error(c, 40002, "session expired or invalid, please re-negotiate encryption key")
			c.Abort()
			return
		}

		// Decrypt request body
		if c.Request.Body != nil && c.Request.ContentLength > 0 {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body.Close()
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
			} else {
				c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			}
		}

		// Wrap response writer to capture and encrypt response
		blw := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
			aesKey:         aesKey,
		}
		c.Writer = blw

		c.Next()

		// Encrypt and write the actual response
		encrypted, err := aesGCMEncrypt(aesKey, blw.body.Bytes())
		if err != nil {
			return
		}
		// Set Content-Type before writing status code and body
		blw.ResponseWriter.Header().Set("Content-Type", "application/octet-stream")
		// Write the deferred status code; default to 200 if never set
		if blw.statusCode != 0 {
			blw.ResponseWriter.WriteHeader(blw.statusCode)
		} else {
			blw.ResponseWriter.WriteHeader(200)
		}
		blw.ResponseWriter.Write(encrypted)

		// Refresh session TTL
		_ = crypto.RefreshSessionKey(sessionID)
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	aesKey     []byte
	statusCode int
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *bodyLogWriter) WriteString(s string) (int, error) {
	return w.body.WriteString(s)
}

func (w *bodyLogWriter) WriteHeader(code int) {
	w.statusCode = code
}

// aesGCMEncrypt encrypts plaintext with AES-256-GCM
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
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// aesGCMDecrypt decrypts ciphertext with AES-256-GCM
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
