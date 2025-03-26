package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	HTTPPort        string
	DBConnString    string
	FlushInterval   time.Duration
	ShutdownTimeout time.Duration
}

func Load() *Config {
	return &Config{
		HTTPPort:        getEnv("HTTP_PORT", "8081"),
		DBConnString:    getEnv("DB_CONN", "postgres://user:password@localhost:5432/dbname?sslmode=disable"),
		FlushInterval:   getEnvDuration("FLUSH_INTERVAL", time.Minute),
		ShutdownTimeout: getEnvDuration("SHUTDOWN_TIMEOUT", 5*time.Second),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	dur, err := time.ParseDuration(val)
	if err != nil {
		log.Printf("invalid duration for %s: %v", key, err)
		return fallback
	}
	return dur
}
