package config

import (
	"os"
	"time"
)

type Config struct {
	Port            string
	StaticDir       string
	ShutdownTimeout time.Duration
}

func Load() Config {
	return Config{
		Port:            getEnv("PORT", "8080"),
		StaticDir:       getEnv("STATIC_DIR", ""),
		ShutdownTimeout: 10 * time.Second,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
