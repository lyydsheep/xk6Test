package util

import (
	"encoding/base64"
)

// EncodeBase64 编码字符串为 base64
func EncodeBase64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

// DecodeBase64 解码 base64 编码的字符串
func DecodeBase64(encodedString string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}
