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
	ErrRSAKeyTooShort   = errors.New("RSA key too short")
	ErrRSADecryptFailed = errors.New("RSA decryption failed")
)

const RSAKeyBits = 2048

type RSAKeyPair struct {
	PublicKey     *rsa.PublicKey
	PrivateKey    *rsa.PrivateKey
	PublicKeyPEM  string
	PrivateKeyPEM string
}

// NewRSAKeyPair generates a new RSA key pair
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

// LoadRSAKeyPair loads RSA key pair from PEM strings
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

// EncryptWithPublicKey encrypts data with RSA public key using OAEP + SHA-256
func EncryptWithPublicKey(pubKey *rsa.PublicKey, plaintext []byte) (string, error) {
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, plaintext, nil)
	if err != nil {
		return "", fmt.Errorf("RSA encryption failed: %w", err)
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptWithPrivateKey decrypts data with RSA private key
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
