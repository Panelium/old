package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"os"
	"sync"
)

const BasePath = "/etc/panelium/daemon"

// File names
const gitignoreFileName = ".gitignore"
const configFileName = "config.json"
const backendJWTPublicKeyFileName = "backend_jwt_public_key.pem"
const jwtPrivateKeyFileName = "jwt_private_key.pem"
const jwtPublicKeyFileName = "jwt_public_key.pem"
const DatabaseFileName = "daemon.db"

// File Locations
const gitignoreLocation = BasePath + "/" + gitignoreFileName
const configLocation = BasePath + "/" + configFileName
const backendJWTPublicKeyLocation = BasePath + "/" + backendJWTPublicKeyFileName
const jwtPrivateKeyLocation = BasePath + "/" + jwtPrivateKeyFileName
const jwtPublicKeyLocation = BasePath + "/" + jwtPublicKeyFileName
const DatabaseLocation = BasePath + "/" + DatabaseFileName

const gitignoreContent = "*.db\n*.pem\nsecrets.json\n"

const jwtKeySize = 2048

var ConfigInstance *Config
var JWTPrivateKeyInstance *rsa.PrivateKey
var BackendJWTPublicKeyInstance *rsa.PublicKey

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

	if _, err := os.Stat(backendJWTPublicKeyLocation); os.IsNotExist(err) {
		return err
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

	jwtPrivateKey, err := loadJWTPrivateKey()
	if err != nil {
		return err
	}
	JWTPrivateKeyInstance = jwtPrivateKey

	backendJWTPublicKey, err := loadBackendJWTPublicKey()
	if err != nil {
		return err
	}
	BackendJWTPublicKeyInstance = backendJWTPublicKey

	return nil
}

// Default Config Values

// Config values should never be accessed or modified directly as that could lead to race conditions.
type Config struct {
	lock sync.RWMutex
}

func newConfig() *Config {
	return &Config{
		lock: sync.RWMutex{},
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

func loadBackendJWTPublicKey() (*rsa.PublicKey, error) {
	file, err := os.ReadFile(backendJWTPublicKeyLocation)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(file)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("invalid PEM block type or format")
	}
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}
