package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"panelium/common/random"
	"sync"
	"time"
)

const BasePath = "/etc/panelium/backend"

// File names
const gitignoreFileName = ".gitignore"
const configFileName = "config.json"
const secretsFileName = "secrets.json"
const jwtPrivateKeyFileName = "jwt_private_key.pem"
const jwtPublicKeyFileName = "jwt_public_key.pem"
const DatabaseFileName = "panelium.db"

// File Locations
const gitignoreLocation = BasePath + "/" + gitignoreFileName
const configLocation = BasePath + "/" + configFileName
const secretsLocation = BasePath + "/" + secretsFileName
const jwtPrivateKeyLocation = BasePath + "/" + jwtPrivateKeyFileName
const jwtPublicKeyLocation = BasePath + "/" + jwtPublicKeyFileName
const DatabaseLocation = BasePath + "/" + DatabaseFileName

const gitignoreContent = "*.db\n*.pem\nsecrets.json\n"

const jwtKeySize = 2048

var ConfigInstance *Config
var SecretsInstance *Secrets
var JWTPrivateKeyInstance *rsa.PrivateKey

func Init() error {
	if err := os.MkdirAll(BasePath, 0755); err != nil {
		return err
	}

	if _, err := os.Stat(gitignoreLocation); os.IsNotExist(err) {
		if err := os.WriteFile(gitignoreLocation, []byte(gitignoreContent), 0644); err != nil {
			return err
		}
	}

	if _, err := os.Stat(configLocation); os.IsNotExist(err) {
		config := newConfig()
		if err := config.Save(); err != nil {
			return err
		}
	}
	if _, err := os.Stat(secretsLocation); os.IsNotExist(err) {
		secrets, err := newSecrets()
		if err != nil {
			return err
		}
		if err := secrets.Save(); err != nil {
			return err
		}
	}

	if _, err := os.Stat(jwtPrivateKeyLocation); os.IsNotExist(err) {
		privateKey, err := rsa.GenerateKey(rand.Reader, jwtKeySize)
		if err != nil {
			return err
		}

		err = privateKey.Validate()
		if err != nil {
			return err
		}

		privateKey.Precompute()

		derPrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
		privatePEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: derPrivateKey})
		if err := os.WriteFile(jwtPrivateKeyLocation, privatePEM, 0600); err != nil {
			return err
		}

		derPublicKey := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
		publicPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: derPublicKey})
		if err := os.WriteFile(jwtPublicKeyLocation, publicPEM, 0644); err != nil {
			return err
		}
	}

	config, err := loadConfig()
	if err != nil {
		return err
	}
	if err := config.Migrate(); err != nil {
		return err
	}
	ConfigInstance = config

	secrets, err := loadSecrets()
	if err != nil {
		return err
	}
	if err := secrets.Migrate(); err != nil {
		return err
	}
	SecretsInstance = secrets

	jwtPrivateKey, err := loadJWTPrivateKey()
	if err != nil {
		return err
	}
	JWTPrivateKeyInstance = jwtPrivateKey

	return nil
}

// Default Config Values
const DefaultAccessTokenDuration = 5 * time.Minute         // 5 minutes
const DefaultRefreshTokenDuration = 24 * time.Hour         // 24 hours
const DefaultPasswordResetTokenDuration = 15 * time.Minute // 15 minutes
const DefaultMFATokenDuration = 15 * time.Minute           // 15 minutes

const DefaultMFACodeLength = 8
const DefaultRecoveryCodeLength = 8
const DefaultRecoveryCodesCount = 10

// Config values should never be accessed or modified directly as that could lead to race conditions.
type Config struct {
	lock sync.RWMutex
	// Durations of tokens in seconds
	JWTDurations struct {
		Access        uint `json:"access"`
		Refresh       uint `json:"refresh"`
		PasswordReset uint `json:"password_reset"`
		MFA           uint `json:"mfa"`
	} `json:"jwt_durations"`
	MFA struct {
		CodeLength         int `json:"code_length"`          // Length of the standard MFA codes (sms, email, etc.)
		RecoveryCodeLength int `json:"recovery_code_length"` // Length of the recovery codes
		RecoveryCodesCount int `json:"recovery_codes_count"` // Number of recovery codes
	} `json:"mfa"`
}

func newConfig() *Config {
	return &Config{
		lock: sync.RWMutex{},
		JWTDurations: struct {
			Access        uint `json:"access"`
			Refresh       uint `json:"refresh"`
			PasswordReset uint `json:"password_reset"`
			MFA           uint `json:"mfa"`
		}{
			Access:        uint(DefaultAccessTokenDuration.Seconds()),
			Refresh:       uint(DefaultRefreshTokenDuration.Seconds()),
			PasswordReset: uint(DefaultPasswordResetTokenDuration.Seconds()),
			MFA:           uint(DefaultMFATokenDuration.Seconds()),
		},
		MFA: struct {
			CodeLength         int `json:"code_length"`
			RecoveryCodeLength int `json:"recovery_code_length"`
			RecoveryCodesCount int `json:"recovery_codes_count"`
		}{
			CodeLength:         DefaultMFACodeLength,
			RecoveryCodeLength: DefaultRecoveryCodeLength,
			RecoveryCodesCount: DefaultRecoveryCodesCount,
		},
	}
}

func loadConfig() (*Config, error) {
	file, err := os.Open(configLocation)
	if err != nil {
		if os.IsNotExist(err) {
			config := newConfig()
			if err := config.Save(); err != nil {
				return nil, err
			}
			return config, nil
		}
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) Migrate() error {
	c.lock.Lock()

	if c.JWTDurations.Access == 0 {
		c.JWTDurations.Access = uint(DefaultAccessTokenDuration.Seconds())
	}
	if c.JWTDurations.Refresh == 0 {
		c.JWTDurations.Refresh = uint(DefaultRefreshTokenDuration.Seconds())
	}
	if c.JWTDurations.PasswordReset == 0 {
		c.JWTDurations.PasswordReset = uint(DefaultPasswordResetTokenDuration.Seconds())
	}
	if c.JWTDurations.MFA == 0 {
		c.JWTDurations.MFA = uint(DefaultMFATokenDuration.Seconds())
	}

	if c.MFA.CodeLength == 0 {
		c.MFA.CodeLength = DefaultMFACodeLength
	}
	if c.MFA.RecoveryCodeLength == 0 {
		c.MFA.RecoveryCodeLength = DefaultRecoveryCodeLength
	}
	if c.MFA.RecoveryCodesCount == 0 {
		c.MFA.RecoveryCodesCount = DefaultRecoveryCodesCount
	}

	c.lock.Unlock()

	if err := c.Save(); err != nil {
		return err
	}

	return nil
}

func (c *Config) Save() error {
	c.lock.Lock()
	data, err := json.MarshalIndent(c, "", "  ")
	c.lock.Unlock()

	if err != nil {
		return err
	}

	return os.WriteFile(configLocation, data, 0644)
}

func (c *Config) GetAccessTokenDuration() time.Duration {
	c.lock.RLock()
	defer c.lock.RUnlock()

	val := time.Duration(c.JWTDurations.Access) * time.Second

	return val
}

func (c *Config) GetRefreshTokenDuration() time.Duration {
	c.lock.RLock()
	defer c.lock.RUnlock()

	val := time.Duration(c.JWTDurations.Refresh) * time.Second

	return val
}

func (c *Config) GetPasswordResetTokenDuration() time.Duration {
	c.lock.RLock()
	defer c.lock.RUnlock()

	val := time.Duration(c.JWTDurations.PasswordReset) * time.Second

	return val
}

func (c *Config) GetMFATokenDuration() time.Duration {
	c.lock.RLock()
	defer c.lock.RUnlock()

	val := time.Duration(c.JWTDurations.MFA) * time.Second

	return val
}

// TODO: Secrets should be stored in HSM when possible, or at least encrypted with the encryption key being in HSM or similar secure storage.

// Secrets values should never be accessed or modified directly as that could lead to race conditions.
type Secrets struct {
	lock          sync.RWMutex
	Pepper        string `json:"pepper"`
	EncryptionKey string `json:"encryption_key"` // AES-256 key for encrypting sensitive data
}

func newSecrets() (*Secrets, error) {
	pepper, err := random.GeneratePepper()
	if err != nil {
		return nil, errors.New("failed to generate pepper: " + err.Error())
	}

	encryptionKey, err := random.GenerateAES256Key()
	if err != nil {
		return nil, errors.New("failed to generate encryption key: " + err.Error())
	}

	keyString := hex.EncodeToString(encryptionKey)

	return &Secrets{
		lock:          sync.RWMutex{},
		Pepper:        pepper,
		EncryptionKey: keyString,
	}, nil
}

func loadSecrets() (*Secrets, error) {
	file, err := os.Open(secretsLocation)
	if err != nil {
		if os.IsNotExist(err) {
			secrets, err := newSecrets()
			if err != nil {
				return nil, err
			}
			if err := secrets.Save(); err != nil {
				return nil, err
			}
			return secrets, nil
		}
		return nil, err
	}
	defer file.Close()

	var secrets Secrets
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&secrets); err != nil {
		return nil, err
	}

	return &secrets, nil
}

func (s *Secrets) Migrate() error {
	s.lock.Lock()

	if s.Pepper == "" {
		pepper, err := random.GeneratePepper()
		if err != nil {
			s.lock.Unlock()
			return errors.New("failed to generate pepper: " + err.Error())
		}
		s.Pepper = pepper
	}

	if s.EncryptionKey == "" {
		encryptionKey, err := random.GenerateAES256Key()
		if err != nil {
			s.lock.Unlock()
			return errors.New("failed to generate encryption key: " + err.Error())
		}
		s.EncryptionKey = hex.EncodeToString(encryptionKey)
	}

	s.lock.Unlock()

	if err := s.Save(); err != nil {
		return err
	}

	return nil
}

func (s *Secrets) Save() error {
	s.lock.Lock()
	data, err := json.MarshalIndent(s, "", "  ")
	s.lock.Unlock()

	if err != nil {
		return err
	}

	return os.WriteFile(secretsLocation, data, 0600)
}

func (s *Secrets) GetPepper() string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	val := s.Pepper

	return val
}

func (s *Secrets) GetEncryptionKey() ([]byte, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.EncryptionKey == "" {
		return nil, errors.New("encryption key is not set")
	}

	val, err := hex.DecodeString(s.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encryption key: %w", err)
	}

	if len(val) != 32 {
		return nil, fmt.Errorf("invalid encryption key length: expected 32 bytes, got %d bytes", len(val))
	}

	return val, nil
}

func loadJWTPrivateKey() (*rsa.PrivateKey, error) {
	file, err := os.ReadFile(jwtPrivateKeyLocation)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(file)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("invalid PEM block type or format")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	if err := privateKey.Validate(); err != nil {
		return nil, err
	}

	return privateKey, nil
}
