package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	HTTP           HTTPConfig           `mapstructure:"http"`
	BookingService BookingServiceConfig `mapstructure:"booking_service"`
	AuthService    BookingServiceConfig `mapstructure:"auth_service"`
}

type HTTPConfig struct {
	Port int `mapstructure:"port"`
}

type BookingServiceConfig struct {
	Address string `mapstructure:"address"`
}

type AuthService struct {
	Address string `mapstructure:"address"`
}

func Load() (*Config, error) {
	v := viper.New()

	// Настройка поиска файла конфигурации
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")

	// Настройка переменных окружения
	v.AutomaticEnv()
	v.SetEnvPrefix("APP")

	// Явное сопоставление переменных окружения с полями конфига
	v.BindEnv("http.port", "HTTP_PORT")
	v.BindEnv("booking_service.address", "BOOKING_SERVICE_ADDR")
	v.BindEnv("auth_service.address", "AUTH_SERVICE_ADDR")

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
