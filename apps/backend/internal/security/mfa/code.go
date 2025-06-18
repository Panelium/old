package mfa

import (
	"crypto/rand"
	"math/big"
)

var mfaCodeCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var mfaCodeLength = 8

func generateMFACode() (string, error) {
	code := make([]byte, mfaCodeLength)

	for i := range code {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(mfaCodeCharset))))
		if err != nil {
			return "", err
		}
		code[i] = mfaCodeCharset[index.Int64()]
	}

	return string(code), nil
}

func SendSMS(phoneNumber string, code string) error {
	// Placeholder for SMS sending logic
	// This function should integrate with an SMS gateway to send the code
	return nil
}

func SendEmail(email string, code string) error {
	// Placeholder for email sending logic
	// This function should integrate with an email service to send the code
	return nil
}
