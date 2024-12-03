package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"time"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Server struct {
		Port     int    `toml:"port"`
		Username string `toml:"username"`
		Password string `toml:"password"`
	} `toml:"server"`
	Storage struct {
		Path string `toml:"path"`
	} `toml:"storage"`
	APIKeys map[string]APIKey `toml:"api_keys"`
}

type APIKey struct {
	Key         string `toml:"key"`
	Name        string `toml:"name"`
	Description string `toml:"description"`
	CreatedAt   string `toml:"created_at"`
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

func createDefaultConfig(configDir string) (*Config, error) {
	apiKey, err := generateRandomString(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	password, err := generateRandomString(16)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password: %w", err)
	}

	cfg := &Config{}
	cfg.Server.Port = 2010
	cfg.Server.Username = "admin"
	cfg.Server.Password = password
	cfg.Storage.Path = filepath.Join(configDir, "stats.json")

	cfg.APIKeys = make(map[string]APIKey)
	cfg.APIKeys["example"] = APIKey{
		Key:         apiKey,
		Name:        "Example Backup",
		Description: "Default API key generated on first run",
		CreatedAt:   time.Now().Format("2006-01-02"),
	}

	return cfg, nil
}

func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "gorestic")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	configFile := filepath.Join(configDir, "config.toml")

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		cfg, err := createDefaultConfig(configDir)
		if err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}

		f, err := os.Create(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create config file: %w", err)
		}
		defer f.Close()

		enc := toml.NewEncoder(f)
		enc.SetIndentTables(true)
		if err := enc.Encode(cfg); err != nil {
			return nil, fmt.Errorf("failed to write config: %w", err)
		}

		fmt.Printf("Created new config file: %s\n", configFile)
		fmt.Printf("Default admin password: %s\n", cfg.Server.Password)
		fmt.Printf("Default API key: %s\n", cfg.APIKeys["example"].Key)

		return cfg, nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	fmt.Println(cfg)

	dataDir := filepath.Dir(cfg.Storage.Path)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	return &cfg, nil
}
