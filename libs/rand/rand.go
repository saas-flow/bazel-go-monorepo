package rand

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes)[:length], nil
}

func GenerateVerificationCode(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	code := base64.URLEncoding.EncodeToString(b)
	return strings.TrimRight(code, "=")
}

func RandomSlugSuffix(length int) string {
	b := make([]byte, length) // 4 byte = 8 karakter hex
	_, err := rand.Read(b)
	if err != nil {
		panic(err) // Pastikan error ditangani
	}
	return hex.EncodeToString(b)
}

func EncodeBase64(data string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(data)), nil
}

func DecodeBase64(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}
