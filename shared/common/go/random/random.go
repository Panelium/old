package random

import (
	"crypto/rand"
	"errors"
)

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GeneratePepper is a copy of the Go standard library's crypto/rand.Text function, but with a 36 character alphabet a parameterized length.
func generateText(n int) (string, error) {
	if n <= 0 {
		return "", errors.New("length must be greater than 0")
	}
	src := make([]byte, n)
	_, err := rand.Read(src)
	if err != nil {
		return "", err
	}
	for i := range src {
		src[i] = alphabet[src[i]%36]
	}
	return string(src), nil
}

// GeneratePepper generates a random string with at least 256 bits of entropy.
func GeneratePepper() (string, error) {
	// ⌈log₃₆ 2²⁵⁶⌉ = 50 chars
	return generateText(50)
}

// GenerateSalt generates a random string with at least 128 bits of entropy.
func GenerateSalt() (string, error) {
	// ⌈log₃₆ 2¹²⁸⌉ = 25 chars
	return generateText(25)
}

// GenerateAESKey generates a byte slice of 32 bytes, suitable for use as an AES-256 key.
func GenerateAES256Key() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
