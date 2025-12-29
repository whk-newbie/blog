package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// CalculateFingerprintHash 计算指纹哈希值（SHA256）
func CalculateFingerprintHash(fingerprintData interface{}) (string, error) {
	// 将指纹数据转换为JSON字符串
	jsonData, err := json.Marshal(fingerprintData)
	if err != nil {
		return "", err
	}

	// 计算SHA256哈希
	hash := sha256.Sum256(jsonData)
	return hex.EncodeToString(hash[:]), nil
}
