package helpers

import (
	"crypto/rand"
	"encoding/base64"
)

// Generate a token with random number if string
func GenerateRandomToken(length int) (string, error) {
	token := make([]byte, length)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(token), nil
}
