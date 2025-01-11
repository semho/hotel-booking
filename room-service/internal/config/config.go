package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Environment string     `mapstructure:"environment"`
	DB          DBConfig   `mapstructure:"db"`
	GRPC        GRPCConfig `mapstructure:"grpc"`
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

func Load() (*Config, error) {
	v := viper.New()

	environment := os.Getenv("APP_ENV")
	if environment == "" {
		environment = "development"
	}

	v.SetDefault("environment", environment)
	// 1. Настройка переменных окружения
	v.AutomaticEnv()

	// 2. Настройка конфига
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")

	// 3. Явное сопоставление переменных окружения с полями конфига
	if environment == "production" {
		v.BindEnv("db.host", "DB_HOST")
		v.BindEnv("db.port", "DB_PORT")
		v.BindEnv("db.user", "DB_USER")
		v.BindEnv("db.name", "DB_NAME")
		v.BindEnv("db.password", "DB_PASSWORD")
		v.BindEnv("grpc.port", "GRPC_PORT")
	}

	// 4. Загрузка конфига
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := v.UnmarshalKey("environments."+environment, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
