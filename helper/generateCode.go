package helper

import (
	"math/rand"
	"time"
)

func GenerateCodeHelper() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 6)
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for i := 0; i < 6; i++ {
		result[i] = chars[r.Intn(len(chars))]
	}

	return string(result)
}
