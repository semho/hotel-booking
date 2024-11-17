package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	DB          DBConfig          `mapstructure:"db"`
	GRPC        GRPCConfig        `mapstructure:"grpc"`
	RoomService RoomServiceConfig `mapstructure:"room_service"`
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

type RoomServiceConfig struct {
	Address string `mapstructure:"address"`
}

func Load() (*Config, error) {
	v := viper.New()

	// 1. Настройка переменных окружения (без префикса)
	v.AutomaticEnv()

	// 2. Настройка конфига
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")

	// 3. Явное сопоставление переменных окружения с полями конфига
	v.BindEnv("db.host", "DB_HOST")
	v.BindEnv("db.port", "DB_PORT")
	v.BindEnv("db.user", "DB_USER")
	v.BindEnv("db.password", "DB_PASSWORD")
	v.BindEnv("db.name", "DB_NAME")
	v.BindEnv("grpc.port", "GRPC_PORT")
	v.BindEnv("room_service.address", "ROOM_SERVICE_ADDR")

	// 4. Загрузка конфига
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
