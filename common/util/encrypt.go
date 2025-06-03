package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// Pad 加填充
func Pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

// Unpad 去除填充
func Unpad(src []byte) (string, error) {
	if len(src) == 0 {
		return "", fmt.Errorf("input to unpad is zero length")
	}
	padding := src[len(src)-1]
	if int(padding) > len(src) {
		return "", fmt.Errorf("padding size is invalid")
	}
	return string(src[:len(src)-int(padding)]), nil
}

// Encrypt AES 加密
func Encrypt(str string, key []byte) (string, error) {
	plainText := []byte(str)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plainText = Pad(plainText, block.BlockSize())
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainText)

	// 使用 Base64 编码密文
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt AES 解密
func Decrypt(cipherTextBase64 string, key []byte) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	return Unpad(cipherText)
}
