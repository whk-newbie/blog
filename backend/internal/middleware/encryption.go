package middleware

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
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
	"/api/v1/fingerprint",
	"/api/v1/visit",
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

		// Decrypt request body (AES-256-CBC, matches frontend crypto-js)
		if c.Request.Body != nil && c.Request.ContentLength > 0 {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body.Close()
			if len(bodyBytes) > 0 {
				plaintext, err := aesCBCDecrypt(aesKey, bodyBytes)
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
		encrypted, err := aesCBCEncrypt(aesKey, blw.body.Bytes())
		if err != nil {
			logger.Error("Failed to encrypt response for path %s: %v", c.Request.URL.Path, err)
			blw.ResponseWriter.Header().Set("Content-Type", "application/json")
			blw.ResponseWriter.WriteHeader(http.StatusInternalServerError)
			_, _ = blw.ResponseWriter.Write([]byte(`{"code":50001,"message":"response encryption failed"}`))
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

// aesCBCEncrypt encrypts with AES-256-CBC, PKCS7 padding.
// Returns base64(IV || ciphertext) to match frontend crypto-js format.
func aesCBCEncrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Add PKCS7 padding
	plaintext = pkcs7Pad(plaintext, aes.BlockSize)

	// Generate random IV
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	// Return base64 encoded (IV + ciphertext)
	result := make([]byte, base64.StdEncoding.EncodedLen(len(ciphertext)))
	base64.StdEncoding.Encode(result, ciphertext)
	return result, nil
}

// aesCBCDecrypt decrypts with AES-256-CBC, PKCS7 padding.
// Expects base64(IV || ciphertext) matching frontend crypto-js format.
func aesCBCDecrypt(key, ciphertext []byte) ([]byte, error) {
	// Base64 decode
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(ciphertext)))
	n, err := base64.StdEncoding.Decode(decoded, ciphertext)
	if err != nil {
		return nil, err
	}
	decoded = decoded[:n]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(decoded) < aes.BlockSize {
		return nil, io.ErrUnexpectedEOF
	}

	iv := decoded[:aes.BlockSize]
	ciphertext = decoded[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, io.ErrUnexpectedEOF
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// Remove PKCS7 padding
	plaintext, err = pkcs7Unpad(plaintext)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	padding := int(data[len(data)-1])
	if padding > len(data) || padding > aes.BlockSize {
		return nil, io.ErrUnexpectedEOF
	}
	return data[:len(data)-padding], nil
}
