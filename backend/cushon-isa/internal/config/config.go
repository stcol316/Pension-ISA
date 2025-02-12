// config/config.go
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Server
	Port        string
	Environment string

	// JWT
	JWTSecret string
}

func Load() (*Config, error) {
	// Load .env file only in development
	if os.Getenv("GO_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			fmt.Printf("Warning: .env file not found: %v\n", err)
		}
	}

	config := &Config{
		// Database
		DBHost:     getEnvWithDefault("DB_HOST", "localhost"),
		DBPort:     getEnvWithDefault("DB_PORT", "5433"),
		DBUser:     getEnvWithDefault("DB_USER", "dev_user"),
		DBPassword: requireEnv("DB_PASSWORD"),
		DBName:     getEnvWithDefault("DB_NAME", "dev_db"),

		// Server
		Port:        getEnvWithDefault("PORT", "8080"),
		Environment: getEnvWithDefault("GO_ENV", "development"),

		// JWT
		JWTSecret: requireEnv("JWT_SECRET"),
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Validate() error {
	required := []struct {
		name, value string
	}{
		{"DB_PASSWORD", c.DBPassword},
		{"JWT_SECRET", c.JWTSecret},
	}

	for _, r := range required {
		if r.value == "" {
			return fmt.Errorf("required environment variable %s is not set", r.name)
		}
	}

	return nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func requireEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		fmt.Printf("Warning: required environment variable %s is not set\n", key)
	}
	return value
}
