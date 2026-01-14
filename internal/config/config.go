package config

import (
	"os"
	"strconv"
)

type Config struct {
	PostgresDSN string
	RedisAddr   string
	KafkaBroker string
	HTTPPort    int
}

func Load() *Config {
	return &Config{
		DBName:      getEnv("POSTGRES_DSN", "postgres://user:pass@localhost:5432/db?sslmode=disable"),
		DBUser:      getEnv("POSTGRES_DSN", "postgres://user:pass@localhost:5432/db?sslmode=disable"),
		DBPassword:  getEnv("POSTGRES_DSN", "postgres://user:pass@localhost:5432/db?sslmode=disable"),
		DBTLS:       getEnv("POSTGRES_DSN", "postgres://user:pass@localhost:5432/db?sslmode=disable"),
		PostgresDSN: getEnv("POSTGRES_DSN", "postgres://user:pass@localhost:5432/db?sslmode=disable"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		KafkaBroker: getEnv("KAFKA_BROKER", "localhost:9092"),
		HTTPPort:    getEnvAsInt("HTTP_PORT", 8080),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
