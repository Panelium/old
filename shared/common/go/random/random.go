package random

import (
	"crypto/rand"
	"encoding/base32"
)

// GenerateSecureRandomString is used for generating the salt and pepper
func GenerateSecureRandomString() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 256 bits - TODO: maybe make this configurable
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	str := base32.StdEncoding.EncodeToString(bytes)
	return str, nil
}
