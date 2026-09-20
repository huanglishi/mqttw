// Package grand provides high performance random bytes/number/string generation functionality.
package grand

import (
	"math/rand"
	"sync"
	"time"
)

var (
	reader = rand.New(newRandSource())
	mutex  sync.Mutex
)

// newRandSource 初始化随机种子
func newRandSource() rand.Source {
	return rand.NewSource(time.Now().UnixNano())
}

// Intn 返回 [0, max) 随机整数
func Intn(max int) int {
	if max <= 0 {
		return 0
	}
	mutex.Lock()
	defer mutex.Unlock()
	return reader.Intn(max)
}

// Int 返回随机整数
func Int() int {
	mutex.Lock()
	defer mutex.Unlock()
	return reader.Int()
}

// Uint64 返回随机 uint64
func Uint64() uint64 {
	mutex.Lock()
	defer mutex.Unlock()
	return reader.Uint64()
}

// Bool 随机 true / false
func Bool() bool {
	return Intn(2) == 1
}

// String 生成指定长度的随机字母数字字符串
func String(length int) string {
	const letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	mutex.Lock()
	defer mutex.Unlock()
	for i := range b {
		b[i] = letters[reader.Intn(len(letters))]
	}
	return string(b)
}

// Bytes 生成随机字节数组
func Bytes(n int) []byte {
	b := make([]byte, n)
	mutex.Lock()
	defer mutex.Unlock()
	_, _ = reader.Read(b)
	return b
}

// Digits 生成指定位数的 纯数字随机字符串
// 例如：Digits(6) → "135792"
func Digits(length int) string {
	const digits = "0123456789"
	b := make([]byte, length)
	mutex.Lock()
	defer mutex.Unlock()
	for i := range b {
		b[i] = digits[reader.Intn(10)]
	}
	return string(b)
}
