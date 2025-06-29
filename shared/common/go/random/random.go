package random

import (
	"crypto/rand"
)

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateSecureRandomString is a copy of the Go standard library's crypto/rand.Text function, but with a 36 character alphabet and at least 256 bits of entropy.
func GenerateSecureRandomString() (string, error) {
	// ⌈log₃₆ 2²⁵⁶⌉ = 50 chars
	src := make([]byte, 50)
	_, err := rand.Read(src)
	if err != nil {
		return "", err
	}
	for i := range src {
		src[i] = alphabet[src[i]%36]
	}
	return string(src), nil
}
