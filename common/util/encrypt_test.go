package util

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGen(t *testing.T) {
	// 填入密钥
	key, _ := hex.DecodeString("13c48c6c2cb3238180cd59c25a74b1cc7f2955c105ae086bb8fcdb0a3ca1535a")
	// 待加密的数据
	data := "Z7aV3tdfAD7VMZcZYLfn"
	// 加密后的数据
	mysql, _ := Encrypt(data, key)
	fmt.Println(mysql)
	de, _ := Decrypt(mysql, key)
	assert.Equal(t, data, de)
}
