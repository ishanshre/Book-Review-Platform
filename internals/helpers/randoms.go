package helpers

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

const numString = "1234567890"

func RandomPhone(length int) string {
	var sb strings.Builder
	k := len(numString)
	for i := 0; i < length; i++ {
		c := numString[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}
