package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	SMTP     SMTPConfig
	Queue    QueueConfig
}

type ServerConfig struct {
	Port    int
	Timeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret         string
	ExpirationTime time.Duration
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	FromName string
	FromAddr string
}

type QueueConfig struct {
	WorkerCount  int
	MaxRetries   int
	RetryBackoff time.Duration
	RateLimit    int
}

func Load() (*Config, error) {
	configPath := "config"
	configName := "config"
	configType := "yaml"
	if os.Getenv("CONFIG_PATH") != "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found, using defaults and environment variables")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

func setDefaults() {
	viper.SetDefault("server.port", 3000)
	viper.SetDefault("server.timeout", "30s")

	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.name", "listmonk")

	viper.SetDefault("jwt.secret", "your-secret-key-change-me")
	viper.SetDefault("jwt.expirationTime", "24h")

	viper.SetDefault("smtp.host", "localhost")
	viper.SetDefault("smtp.port", 25)
	viper.SetDefault("smtp.username", "")
	viper.SetDefault("smtp.password", "")
	viper.SetDefault("smtp.fromName", "Listmonk Clone")
	viper.SetDefault("smtp.fromAddr", "noreply@example.com")

	viper.SetDefault("queue.workerCount", 5)
	viper.SetDefault("queue.maxRetries", 3)
	viper.SetDefault("queue.retryBackoff", "5m")
	viper.SetDefault("queue.rateLimit", 10)
}
