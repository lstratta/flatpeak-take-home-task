package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Host           string
	Port           string
	AllowedOrigins []string
}

func New() Config {
	host := getEnvDefault("HOST", "localhost")
	port := getEnvDefault("PORT", "7777")
	ao := strings.Split(getEnvDefault("ALLOWED_ORIGINS", fmt.Sprintf("http://%s:%s,https://%s:%s", host, port, host, port)), ",")

	return Config{
		Host:           host,
		Port:           port,
		AllowedOrigins: ao,
	}
}

// this function allows easy switching between a local dev
// and a cloud provider environment variables
func getEnvDefault(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return def
}
