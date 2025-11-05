package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	Port         string
	BaseURL      string
	DockerSocket string
	DataDir      string
	LogLevel     string
}

// Load loads the configuration from environment variables
func Load() *Config {
	return &Config{
		Port:         getEnv("PORT", "8080"),
		BaseURL:      getEnv("BASE_URL", "http://localhost:8080"),
		DockerSocket: getEnv("DOCKER_SOCKET", "/var/run/docker.sock"),
		DataDir:      getEnv("DATA_DIR", "./data"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
