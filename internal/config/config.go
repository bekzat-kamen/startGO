package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	LogLevel string
	Database *DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "INFO"),

		Database: &DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Username: getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_NAME", "neondb"),
			SSLMode:  getEnv("SSL_MODE", "require"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
