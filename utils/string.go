package utils

import (
	"math/rand"
	"time"
)

const (
	intChars = "0123456789"
	alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func init() {
	rand.Seed(int64(time.Now().UnixNano()))
}

// GenRandomString 生成长度为 n 的随机字符串
func GenRandomString(n uint64, alphabets ...byte) string {
	letters := alphanum
	if len(alphabets) > 0 {
		letters = string(alphabets)
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}
