package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"
)

// TODO: should make this thread safe with a mutex

const BasePath = "/etc/panelium"

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

var gitignoreContent = fmt.Sprintf("%s\n%s\n%s\n", secretsFileName, jwtPrivateKeyFileName, jwtPublicKeyFileName)

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
		secrets := newSecrets()
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
const DefaultAccessTokenDuration = 5 * time.Minute // 5 minutes
const DefaultRefreshTokenDuration = 24 * time.Hour // 24 hours
const DefaultMFATokenDuration = 15 * time.Minute   // 15 minutes

const DefaultMFACodeLength = 8
const DefaultRecoveryCodeLength = 8
const DefaultRecoveryCodesCount = 10

type Config struct {
	// Durations of tokens in seconds
	JWTDurations struct {
		Access  uint `json:"access"`
		Refresh uint `json:"refresh"`
		MFA     uint `json:"mfa"`
	} `json:"jwt_durations"`
	MFA struct {
		CodeLength         int `json:"code_length"`          // Length of the standard MFA codes (sms, email, etc.)
		RecoveryCodeLength int `json:"recovery_code_length"` // Length of the recovery codes
		RecoveryCodesCount int `json:"recovery_codes_count"` // Number of recovery codes
	} `json:"mfa"`
}

func newConfig() *Config {
	return &Config{
		JWTDurations: struct {
			Access  uint `json:"access"`
			Refresh uint `json:"refresh"`
			MFA     uint `json:"mfa"`
		}{
			Access:  uint(DefaultAccessTokenDuration.Seconds()),
			Refresh: uint(DefaultRefreshTokenDuration.Seconds()),
			MFA:     uint(DefaultMFATokenDuration.Seconds()),
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
	if c.JWTDurations.Access == 0 {
		c.JWTDurations.Access = uint(DefaultAccessTokenDuration.Seconds())
	}
	if c.JWTDurations.Refresh == 0 {
		c.JWTDurations.Refresh = uint(DefaultRefreshTokenDuration.Seconds())
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

	if err := c.Save(); err != nil {
		return err
	}

	return nil
}

func (c *Config) Save() error {
	file, err := os.Create(configLocation)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c)
}

func (c *Config) GetAccessTokenDuration() time.Duration {
	return time.Duration(c.JWTDurations.Access) * time.Second
}

func (c *Config) GetRefreshTokenDuration() time.Duration {
	return time.Duration(c.JWTDurations.Refresh) * time.Second
}

func (c *Config) GetMFATokenDuration() time.Duration {
	return time.Duration(c.JWTDurations.MFA) * time.Second
}

type Secrets struct {
	pepper string `json:"pepper"`
}

func newSecrets() *Secrets {
	return &Secrets{
		pepper: rand.Text(),
	}
}

func loadSecrets() (*Secrets, error) {
	file, err := os.Open(secretsLocation)
	if err != nil {
		if os.IsNotExist(err) {
			secrets := newSecrets()
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
	if s.pepper == "" {
		s.pepper = rand.Text()
	}

	if err := s.Save(); err != nil {
		return err
	}

	return nil
}

func (s *Secrets) Save() error {
	file, err := os.Create(secretsLocation)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(s)
}

func (s *Secrets) GetPepper() string {
	return s.pepper
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
