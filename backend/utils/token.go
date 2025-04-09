package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomToken 生成指定长度的随机令牌
// 参数：
//   - length: 令牌长度（字节数）
//
// 返回：
//   - string: 生成的随机令牌字符串
func GenerateRandomToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
