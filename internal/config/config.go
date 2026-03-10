package config

import (
	"base-service/internal/errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

/*type Config struct {
	App   AppConfig   `mapstructure:"app"`
	DB    DBConfig    `mapstructure:"db"`
	Redis RedisConfig `mapstructure:"redis"`
	Kafka KafkaConfig `mapstructure:"kafka"`
}

type AppConfig struct {
	Name     string `mapstructure:"APP_NAME"`
	Host     string `mapstructure:"HTTP_HOST"`
	HttpPort int    `mapstructure:"HTTP_PORT"`
	GrpcPort int    `mapstructure:"GRPC_PORT"`
	LogLevel string `mapstructure:"LOG_LEVEL"`
	Debug    bool   `mapstructure:"DEBUG_MODE"`
}

type DBConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     int    `mapstructure:"DB_PORT"`
	Name     string `mapstructure:"DB_NAME"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	TLS      bool   `mapstructure:"DB_USE_TLS"`
	DSN      string // ← генерируется, не читается из окружения
}

type RedisConfig struct {
	Host   string `mapstructure:"REDIS_HOST"`
	Prefix string `mapstructure:"REDIS_PREFIX"`
}

type KafkaConfig struct {
	Host  string `mapstructure:"KAFKA_HOST"`
	Group string `mapstructure:"KAFKA_GROUP_ID"`
}*/

type Config struct {
	PostgresDSN  string `mapstructure:"POSTGRES_DSN"` // можно не задавать — мы строим сами
	AppName      string `mapstructure:"APP_NAME"`
	HttpHost     string `mapstructure:"HTTP_HOST"`
	HttpPort     int    `mapstructure:"HTTP_PORT"`
	AppGrpcPort  int    `mapstructure:"GRPC_PORT"`
	LogLevel     string `mapstructure:"LOG_LEVEL"`
	DebugMode    bool   `mapstructure:"DEBUG_MODE"`
	DbHost       string `mapstructure:"DB_HOST"`
	DbPort       int    `mapstructure:"DB_PORT"`
	DbName       string `mapstructure:"DB_NAME"`
	DbUser       string `mapstructure:"DB_USER"`
	DbPassword   string `mapstructure:"DB_PASSWORD"`
	DbTLS        bool   `mapstructure:"DB_USE_TLS"`
	RedisHost    string `mapstructure:"REDIS_HOST"`
	RedisPref    string `mapstructure:"REDIS_PREFIX"`
	KafkaHost    string `mapstructure:"KAFKA_HOST"`
	KafkaGroupID string `mapstructure:"KAFKA_GROUP_ID"`
}

func Load() (*Config, error) {
	log := slog.Default()
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			log.Error("failed to load .env file", "error", err)
		} else {
			log.Debug(".env file not found, using system environment")
		}
	}
	setDefaults()
	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Error("failed to unmarshal config", "config", err)
		return nil, errors.ErrUnmarshal
	}
	if err := validateRequired(&cfg); err != nil {
		log.Error("configuration validation failed", "error", err)
		return nil, errors.ErrValidRequired
	}
	sslmode := "disable"
	if cfg.DbTLS {
		sslmode = "require"
	}

	user := url.QueryEscape(cfg.DbUser)
	password := url.QueryEscape(cfg.DbPassword)

	cfg.PostgresDSN = fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		user, password, cfg.DbHost, cfg.DbPort, cfg.DbName, sslmode,
		"postgres://%s:%s@%s:%d/%s?sslmode=%s", user, password, cfg.DbHost, cfg.DbPort, cfg.DbName, sslmode,
	)

	log.Info("✓ config loaded",
		"app", cfg.AppName,
		"http", fmt.Sprintf("%s:%d", cfg.HttpHost, cfg.HttpPort),
		"db", fmt.Sprintf("%s@%s:%d/%s", cfg.DbUser, cfg.DbHost, cfg.DbPort, cfg.DbName),
	)
	return &cfg, nil
}

func setDefaults() {
	viper.SetDefault("APP_NAME", "base-service")
	viper.SetDefault("HTTP_HOST", "0.0.0.0") // Слушать все интерфейсы
	viper.SetDefault("HTTP_PORT", 8080)      // Стандартный HTTP порт
	viper.SetDefault("GRPC_PORT", 9090)      // Стандартный gRPC порт
	viper.SetDefault("LOG_LEVEL", "info")    // Уровень логирования
	viper.SetDefault("DEBUG_MODE", false)    // Режим отладки

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_NAME", "base_service")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("REDIS_HOST", "localhost:0000")
	viper.SetDefault("KAFKA_HOST", "localhost:0000")

	viper.SetDefault("DB_PORT", 5432)     // Стандартный порт PostgreSQL
	viper.SetDefault("DB_USE_TLS", false) // Без TLS по умолчанию

	viper.SetDefault("REDIS_PREFIX", "base-service")

	viper.SetDefault("KAFKA_GROUP_ID", "base-service-group")
}

func validateRequired(cfg *Config) error {
	var errs []string

	if cfg.DbHost == "" {
		errs = append(errs, "DB_HOST is required (example: localhost:5432)")
	}
	if cfg.DbName == "" {
		errs = append(errs, "DB_NAME is required")
	}
	if cfg.DbUser == "" {
		errs = append(errs, "DB_USER is required")
	}
	if cfg.DbPassword == "" {
		errs = append(errs, "DB_PASSWORD is required")
	}
	if cfg.RedisHost == "" {
		errs = append(errs, "REDIS_HOST is required (example: localhost:6379)")
	} else if !strings.Contains(cfg.RedisHost, ":") {
		errs = append(errs, fmt.Sprintf("REDIS_HOST must include port, got: %q", cfg.RedisHost))
	}

	if cfg.KafkaHost == "" {
		errs = append(errs, "KAFKA_HOST is required (example: localhost:9092)")
	} else if !strings.Contains(cfg.KafkaHost, ":") {
		errs = append(errs, fmt.Sprintf("KAFKA_HOST must include port, got: %q", cfg.KafkaHost))
	}

	if len(errs) > 0 {
		return fmt.Errorf(
			"configuration validation failed — these values MUST be set in environment variables or .env file:\n%s\n\n💡 Hint: Check your .env file or deployment configuration",
			strings.Join(errs, "\n"),
		)
	}
	return nil
}
