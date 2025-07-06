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
const secretsFileName = "secrets.json"
const backendJWTPublicKeyFileName = "backend_jwt_public_key.pem"
const jwtPrivateKeyFileName = "jwt_private_key.pem"
const jwtPublicKeyFileName = "jwt_public_key.pem"
const DatabaseFileName = "daemon.db"

// File Locations
const gitignoreLocation = BasePath + "/" + gitignoreFileName
const configLocation = BasePath + "/" + configFileName
const secretsLocation = BasePath + "/" + secretsFileName
const backendJWTPublicKeyLocation = BasePath + "/" + backendJWTPublicKeyFileName
const jwtPrivateKeyLocation = BasePath + "/" + jwtPrivateKeyFileName
const jwtPublicKeyLocation = BasePath + "/" + jwtPublicKeyFileName
const DatabaseLocation = BasePath + "/" + DatabaseFileName

const gitignoreContent = "*.db\n*.pem\nsecrets.json\n"

const jwtKeySize = 2048

var ConfigInstance *Config
var SecretsInstance *Secrets
var JWTPrivateKeyInstance *rsa.PrivateKey
var BackendJWTPublicKeyInstance *rsa.PublicKey // Note: can be nil if not set yet

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

	_, err = os.Stat(backendJWTPublicKeyLocation)
	if !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		backendJWTPublicKey, err := loadBackendJWTPublicKey()
		if err != nil {
			return err
		}
		BackendJWTPublicKeyInstance = backendJWTPublicKey
	}

	return nil
}

// Default Config Values
const DefaultDashboardHost = "https://example.com"
const DefaultBackendHost = "https://example.com:9090"
const DefaultDaemonHost = "https://example.com:9000"

// Config values should never be accessed or modified directly as that could lead to race conditions.
type Config struct {
	lock  sync.RWMutex
	Hosts struct {
		Dashboard string `json:"dashboard"`
		Backend   string `json:"backend"`
		Daemon    string `json:"daemon"` // host for this daemon instance
	}
}

func newConfig() *Config {
	return &Config{
		lock: sync.RWMutex{},
		Hosts: struct {
			Dashboard string `json:"dashboard"`
			Backend   string `json:"backend"`
			Daemon    string `json:"daemon"`
		}{
			Dashboard: DefaultDashboardHost,
			Backend:   DefaultBackendHost,
			Daemon:    DefaultDaemonHost,
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

	if c.Hosts.Dashboard == "" {
		c.Hosts.Dashboard = DefaultDashboardHost
	}
	if c.Hosts.Backend == "" {
		c.Hosts.Backend = DefaultBackendHost
	}
	if c.Hosts.Daemon == "" {
		c.Hosts.Daemon = DefaultDaemonHost
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

func (c *Config) GetDashboardHost() string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.Hosts.Dashboard
}

func (c *Config) GetBackendHost() string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.Hosts.Backend
}

func (c *Config) GetDaemonHost() string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.Hosts.Daemon
}

// TODO: Secrets should be stored in HSM when possible, or at least encrypted with the encryption key being in HSM or similar secure storage.

// Secrets values should never be accessed or modified directly as that could lead to race conditions.
type Secrets struct {
	lock         sync.RWMutex
	NodeJTI      string `json:"node_jti"`
	BackendToken string `json:"backend_token"`
}

func newSecrets() (*Secrets, error) {
	return &Secrets{
		lock:         sync.RWMutex{},
		NodeJTI:      "", // generated after node is connected to backend
		BackendToken: "",
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

func (s *Secrets) GetNodeJTI() string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.NodeJTI
}

func (s *Secrets) SetNodeJTI(nodeJTI string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.NodeJTI = nodeJTI
}

func (s *Secrets) GetBackendToken() string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.BackendToken
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
