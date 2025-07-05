package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration settings
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

// LoadConfig loads configuration from environment variables and .env file
func LoadConfig() (*Config, error) {
	// Load .env file if it exists (for local development)
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
		}
	}

	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "todo_user"),
		DBPassword: getEnv("DB_PASSWORD", "todo_password"),
		DBName:     getEnv("DB_NAME", "todo_db"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}

	// Basic validation: ensure critical DB configs are not empty
	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" || cfg.DBPort == "" {
		return nil, fmt.Errorf("missing critical database environment variables: DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT must be set")
	}

	return cfg, nil
}

// getEnv retrieves an environment variable or returns a fallback value
func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
