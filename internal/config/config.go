package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/vanvanni/gorestic/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	DB     *gorm.DB
	Server ServerConfig
}

type ServerConfig struct {
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

func initializeDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := models.AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

func Load() (*Config, error) {
	configDir := getConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	configFile := filepath.Join(configDir, "server.toml")
	cfg, err := loadConfig(configFile)
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(configDir, "gorestic.sqlite")
	db, err := initializeDB(dbPath)
	if err != nil {
		return nil, err
	}
	cfg.DB = db

	return cfg, nil
}

func loadConfig(configFile string) (*Config, error) {
	log.Printf("loading config: %s", configFile)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return createDefaultConfig(configFile)
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

func createDefaultConfig(configFile string) (*Config, error) {
	password, err := generateRandomString(16)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password: %w", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:     2010,
			Username: "admin",
			Password: password,
		},
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

	log.Printf("created new config file: %s\n", configFile)
	log.Printf("default admin password: %s\n", cfg.Server.Password)

	return cfg, nil
}

func getConfigDir() string {
	if os.Getenv("GORESTIC_DOCKER") == "true" {
		return "/home/app/.config/gorestic"
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		if home := os.Getenv("HOME"); home != "" {
			return filepath.Join(home, ".config", "gorestic")
		}
		return filepath.Join(".", ".config", "gorestic")
	}

	return filepath.Join(homeDir, ".config", "gorestic")
}

func (c *Config) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
