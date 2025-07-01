package config

import (
	"encoding/json"
	"os"
	"sync"
)

const BasePath = "/etc/panelium"

// File names
const configFileName = "daemon-config.json"
const backendJWTPublicKeyFileName = "backend_jwt_public_key.pem"
const DatabaseFileName = "daemon.db"

// File Locations
const configLocation = BasePath + "/" + configFileName
const backendJWTPublicKeyLocation = BasePath + "/" + backendJWTPublicKeyFileName
const DatabaseLocation = BasePath + "/" + DatabaseFileName

var ConfigInstance *Config

func Init() error {
	if err := os.MkdirAll(BasePath, 0755); err != nil {
		return err
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

	config, err := loadConfig()
	if err != nil {
		return err
	}
	if err := config.Migrate(); err != nil {
		return err
	}
	ConfigInstance = config

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
