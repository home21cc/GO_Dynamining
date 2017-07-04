package dtools

import (
	"time"
	"math/rand"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // 52 possibilities
	letterIdxBits = 6                    // 6 bits to represent 64 possibilities / indexes
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax = 63 / letterIdxBits
)

var src = rand.NewSource(time.Now().UnixNano())

// UniqueId 생성용
func GenerateUniqueId(length int) string {
	uniqueName := make([]byte, length)
	for i, cache, remain := length-1, src.Int63(), letterIdxMax; i>= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			uniqueName[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(uniqueName)
}