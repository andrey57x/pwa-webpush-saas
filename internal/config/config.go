package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv              string `mapstructure:"APP_ENV"`
	APIHost             string `mapstructure:"API_HOST"`
	APIPort             string `mapstructure:"API_PORT"`
	WorkerPoolSize      int    `mapstructure:"WORKER_POOL_SIZE"`
	JWTSecret           string `mapstructure:"JWT_SECRET"`
	MasterEncryptionKey string `mapstructure:"MASTER_ENCRYPTION_KEY"`
	ContactEmail        string `mapstructure:"CONTACT_EMAIL"`

	PostgresHost     string `mapstructure:"POSTGRES_HOST"`
	PostgresPort     string `mapstructure:"POSTGRES_PORT"`
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDB       string `mapstructure:"POSTGRES_DB"`
	PostgresSSLMode  string `mapstructure:"POSTGRES_SSLMODE"`

	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`

	KafkaHost          string `mapstructure:"KAFKA_HOST"`
	KafkaPort          string `mapstructure:"KAFKA_PORT"`
	KafkaBrokers       string `mapstructure:"KAFKA_BROKERS"`
	KafkaTopicPushJobs string `mapstructure:"KAFKA_TOPIC_PUSH_JOBS"`
	KafkaConsumerGroup string `mapstructure:"KAFKA_CONSUMER_GROUP"`
}

func (c *Config) GetPostgresDSN() string {
	escapedPassword := url.QueryEscape(c.PostgresPassword)
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.PostgresUser, escapedPassword, c.PostgresHost, c.PostgresPort, c.PostgresDB, c.PostgresSSLMode,
	)
}

func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

func (c *Config) GetKafkaBrokers() []string {
	if c.KafkaBrokers != "" {
		return strings.Split(c.KafkaBrokers, ",")
	}
	if c.KafkaHost != "" && c.KafkaPort != "" {
		return []string{fmt.Sprintf("%s:%s", c.KafkaHost, c.KafkaPort)}
	}
	return []string{"localhost:9092"}
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if cfg.APIHost == "" {
		cfg.APIHost = "0.0.0.0"
	}
	if cfg.APIPort == "" {
		cfg.APIPort = "8080"
	}

	return &cfg, nil
}
