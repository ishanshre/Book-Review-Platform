package helpers

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomToken(length int) (string, error) {
	token := make([]byte, length)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(token), nil
}
