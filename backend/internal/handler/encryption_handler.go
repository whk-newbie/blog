package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/pkg/crypto"
	"github.com/whk-newbie/blog/internal/pkg/response"
)

type EncryptionHandler struct {
	rsaKeyPair *crypto.RSAKeyPair
}

func NewEncryptionHandler(rsaKeyPair *crypto.RSAKeyPair) *EncryptionHandler {
	return &EncryptionHandler{rsaKeyPair: rsaKeyPair}
}

type PublicKeyResponse struct {
	PublicKey string `json:"public_key"`
	SessionID string `json:"session_id"`
}

// GetPublicKey returns the RSA public key and a new session ID
func (h *EncryptionHandler) GetPublicKey(c *gin.Context) {
	sessionID := crypto.GenerateSessionID()
	response.Success(c, PublicKeyResponse{
		PublicKey: h.rsaKeyPair.PublicKeyPEM,
		SessionID: sessionID,
	})
}

type SessionKeyRequest struct {
	EncryptedKey string `json:"encrypted_key" binding:"required"`
	SessionID    string `json:"session_id" binding:"required"`
}

// NegotiateKey decrypts the AES key with RSA private key and stores it
func (h *EncryptionHandler) NegotiateKey(c *gin.Context) {
	var req SessionKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	aesKey, err := crypto.DecryptWithPrivateKey(h.rsaKeyPair.PrivateKey, req.EncryptedKey)
	if err != nil {
		response.Error(c, 40004, "failed to decrypt AES key")
		return
	}

	if len(aesKey) != 32 {
		response.BadRequest(c, "AES key must be 32 bytes")
		return
	}

	if err := crypto.StoreSessionKey(req.SessionID, aesKey); err != nil {
		response.InternalServerError(c, "failed to store session key")
		return
	}

	response.Success(c, gin.H{"status": "ok"})
}
