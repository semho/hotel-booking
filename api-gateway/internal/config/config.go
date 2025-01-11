package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Environment    string               `mapstructure:"environment"`
	HTTP           HTTPConfig           `mapstructure:"http"`
	BookingService BookingServiceConfig `mapstructure:"booking_service"`
	AuthService    AuthServiceConfig    `mapstructure:"auth_service"`
	RoomService    RoomServiceConfig    `mapstructure:"room_service"`
	CORS           CORSConfig           `mapstructure:"cors"`
}

type HTTPConfig struct {
	Port int `mapstructure:"port"`
}

type BookingServiceConfig struct {
	Address string `mapstructure:"address"`
}

type RoomServiceConfig struct {
	Address string `mapstructure:"address"`
}

type AuthServiceConfig struct {
	Address string `mapstructure:"address"`
}

type CORSConfig struct {
	Origins        []string `mapstructure:"origins"`
	Methods        []string `mapstructure:"methods"`
	Headers        []string `mapstructure:"headers"`
	ExposedHeaders []string `mapstructure:"exposed_headers"`
	Credentials    bool     `mapstructure:"credentials"`
	MaxAge         int      `mapstructure:"max_age"`
	Debug          bool     `mapstructure:"debug"`
}

func Load() (*Config, error) {
	v := viper.New()

	// Настройка поиска файла конфигурации
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")

	environment := os.Getenv("APP_ENV")
	if environment == "" {
		environment = "development"
	}

	v.SetDefault("environment", environment)
	// Настройка переменных окружения
	v.AutomaticEnv()
	// Убираем префикс APP
	// v.SetEnvPrefix("APP")

	// Явное сопоставление переменных окружения с полями конфига
	if environment == "production" {
		v.BindEnv("http.port", "HTTP_PORT")
		v.BindEnv("booking_service.address", "BOOKING_SERVICE_ADDR")
		v.BindEnv("auth_service.address", "AUTH_SERVICE_ADDR")
		v.BindEnv("room_service.address", "ROOM_SERVICE_ADDR")
	}

	// Загрузка конфигурации
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
