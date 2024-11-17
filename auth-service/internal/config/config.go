package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	DB   DBConfig   `mapstructure:"db"`
	GRPC GRPCConfig `mapstructure:"grpc"`
	JWT  JWTConfig  `mapstructure:"jwt"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type GRPCConfig struct {
	Port int `mapstructure:"port"`
}

type JWTConfig struct {
	AccessTokenSecret  string `mapstructure:"access_token_secret"`
	RefreshTokenSecret string `mapstructure:"refresh_token_secret"`
	AccessTokenTTL     int    `mapstructure:"access_token_ttl"`  // в минутах
	RefreshTokenTTL    int    `mapstructure:"refresh_token_ttl"` // в днях
}

func Load() (*Config, error) {
	v := viper.New()

	// Настройка конфига
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")

	// Настройка переменных окружения
	v.AutomaticEnv()
	v.SetEnvPrefix("AUTH")

	// Явное сопоставление переменных окружения с полями конфига
	v.BindEnv("db.host", "DB_HOST")
	v.BindEnv("db.port", "DB_PORT")
	v.BindEnv("db.user", "DB_USER")
	v.BindEnv("db.password", "DB_PASSWORD")
	v.BindEnv("db.name", "DB_NAME")
	v.BindEnv("grpc.port", "GRPC_PORT")
	v.BindEnv("jwt.access_token_secret", "JWT_ACCESS_SECRET")
	v.BindEnv("jwt.refresh_token_secret", "JWT_REFRESH_SECRET")
	v.BindEnv("jwt.access_token_ttl", "JWT_ACCESS_TTL")
	v.BindEnv("jwt.refresh_token_ttl", "JWT_REFRESH_TTL")

	// Загрузка конфигурации
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
