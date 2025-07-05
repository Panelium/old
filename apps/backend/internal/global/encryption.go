package global

import (
	"crypto/aes"
	"crypto/cipher"
	"panelium/backend/internal/config"
	"sync"
)

var (
	initOnceEncryption sync.Once
	encryptionBlock    cipher.Block
)

func InitEncryption() error {
	var err error

	initOnceEncryption.Do(func() {
		var encryptionKey []byte
		encryptionKey, err = config.SecretsInstance.GetEncryptionKey()

		encryptionBlock, err = aes.NewCipher(encryptionKey)
	})

	if err != nil {
		return err
	}
	return nil
}

func EncryptionInstance() *cipher.Block {
	if encryptionBlock == nil {
		panic("encryption not initialized, call InitEncryption() first")
	}
	return &encryptionBlock
}
