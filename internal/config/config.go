package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv              string `mapstructure:"APP_ENV"`
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

	KafkaBrokers       string `mapstructure:"KAFKA_BROKERS"`
	KafkaTopicPushJobs string `mapstructure:"KAFKA_TOPIC_PUSH_JOBS"`
	KafkaConsumerGroup string `mapstructure:"KAFKA_CONSUMER_GROUP"`
}

func (c *Config) GetPostgresDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.PostgresUser, c.PostgresPassword, c.PostgresHost, c.PostgresPort, c.PostgresDB, c.PostgresSSLMode,
	)
}

func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
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

	return &cfg, nil
}
