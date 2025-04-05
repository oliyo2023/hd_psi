package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// QRCodeData 二维码数据结构
type QRCodeData struct {
	SKU         string `json:"sku"`
	BatchNumber string `json:"batch"`
	Timestamp   int64  `json:"ts"`
}

// GenerateQRCode 生成商品二维码数据
// 参数:
// - sku: 商品SKU
// - batchNumber: 采购批次号
// - timestamp: 时间戳
// - secretKey: 用于生成校验码的密钥
// 返回:
// - 编码后的二维码数据
// - 错误信息
func GenerateQRCode(sku, batchNumber string, timestamp int64, secretKey string) (string, error) {
	// 创建二维码数据结构
	data := QRCodeData{
		SKU:         sku,
		BatchNumber: batchNumber,
		Timestamp:   timestamp,
	}
	
	// 将数据转换为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("JSON编码失败: %v", err)
	}
	
	// 生成HMAC校验码
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(jsonData)
	checksum := h.Sum(nil)
	
	// 将JSON数据和校验码拼接
	finalData := append(jsonData, checksum...)
	
	// Base64编码
	encodedData := base64.StdEncoding.EncodeToString(finalData)
	
	return encodedData, nil
}

// VerifyQRCode 验证二维码数据
// 参数:
// - encodedData: 编码后的二维码数据
// - secretKey: 用于验证校验码的密钥
// 返回:
// - 解码后的二维码数据
// - 是否有效
// - 错误信息
func VerifyQRCode(encodedData, secretKey string) (*QRCodeData, bool, error) {
	// Base64解码
	data, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, false, fmt.Errorf("Base64解码失败: %v", err)
	}
	
	// 数据长度必须大于SHA256校验码长度(32字节)
	if len(data) <= 32 {
		return nil, false, fmt.Errorf("数据格式无效")
	}
	
	// 分离JSON数据和校验码
	jsonData := data[:len(data)-32]
	checksum := data[len(data)-32:]
	
	// 验证校验码
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(jsonData)
	expectedChecksum := h.Sum(nil)
	
	// 比较校验码
	valid := hmac.Equal(checksum, expectedChecksum)
	
	// 如果校验码有效，解析JSON数据
	if valid {
		var qrData QRCodeData
		if err := json.Unmarshal(jsonData, &qrData); err != nil {
			return nil, false, fmt.Errorf("JSON解码失败: %v", err)
		}
		return &qrData, true, nil
	}
	
	return nil, false, nil
}
