package mfa

import (
	"crypto/rand"
	"math/big"
)

const mfaCodeCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// TODO: maybe config?
const DefaultMFACodeLength = 8
const DefaultRecoveryCodeLength = 8
const DefaultRecoveryCodesCount = 10

func GenerateMFACode(mfaCodeLength int) (string, error) {
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

func GenerateRecoveryCodes(recoveryCodeLength int, recoveryCodesCount int) ([]string, error) {
	codes := make([]string, recoveryCodesCount)

	for i := 0; i < recoveryCodesCount; i++ {
		code, err := GenerateMFACode(recoveryCodeLength)
		if err != nil {
			return nil, err
		}
		codes[i] = code
	}

	return codes, nil
}
