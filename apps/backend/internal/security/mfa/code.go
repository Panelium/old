package mfa

import (
	"crypto/rand"
	"math/big"
	"panelium/backend/internal/config"
)

const mfaCodeCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateCode(length int) (string, error) {
	code := make([]byte, length)

	for i := range code {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(mfaCodeCharset))))
		if err != nil {
			return "", err
		}
		code[i] = mfaCodeCharset[index.Int64()]
	}

	return string(code), nil
}

func GenerateMFACode() (string, error) {
	code, err := generateCode(config.ConfigInstance.MFA.CodeLength)
	if err != nil {
		return "", err
	}
	return code, nil
}

func generateRecoveryCode() (string, error) {
	code, err := generateCode(config.ConfigInstance.MFA.RecoveryCodeLength)
	if err != nil {
		return "", err
	}
	return code, nil
}

func GenerateRecoveryCodes() ([]string, error) {
	count := config.ConfigInstance.MFA.RecoveryCodesCount
	codes := make([]string, count)

	for i := 0; i < count; i++ {
		code, err := generateRecoveryCode()
		if err != nil {
			return nil, err
		}
		codes[i] = code
	}

	return codes, nil
}
