package security

import (
	"encoding/hex"
	"golang.org/x/crypto/argon2"
	"panelium/backend/internal/config"
	"panelium/common/random"
)

// hashPasswordInternal is a helper function that hashes the password using Argon2 provided salt and pepper.
func hashPasswordInternal(password string, salt string) string {
	pepperedPassword := password + config.SecretsInstance.GetPepper()
	passwordHashData := argon2.IDKey([]byte(pepperedPassword), []byte(salt), 1, 64*1024, 4, 32) // TODO: potentially make these parameters configurable
	passwordHash := hex.EncodeToString(passwordHashData)
	return passwordHash
}

// HashPassword is used for hashing passwords before storing them in the database and generating a salt.
// returns the password hash and the generated salt
// pepper is retrieved from app config
func HashPassword(password string) (passwordHash string, salt string, err error) {
	salt, err = random.GenerateSalt()
	if err != nil {
		return "", "", err
	}
	passwordHash = hashPasswordInternal(password, salt)
	return passwordHash, salt, nil
}

// VerifyPassword is used to verify password user input against the stored hash
// only the password is user input
// salt and passwordHash is retrieved from the database
// pepper is retrieved from app config
func VerifyPassword(password string, salt string, passwordHash string) bool {
	passwordHashNew := hashPasswordInternal(password, salt)
	return passwordHashNew == passwordHash
}
