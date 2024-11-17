package main

import (
	"context"
	"github.com/semho/hotel-booking/api-gateway/internal/app"
	"github.com/semho/hotel-booking/api-gateway/internal/config"
	"github.com/semho/hotel-booking/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Инициализируем логгер
	logger.Init()

	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Создаем приложение
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	// Создаем канал для получения сигналов ОС
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем приложение в отдельной горутине
	errChan := make(chan error)
	go func() {
		if err = application.Run(); err != nil {
			errChan <- err
		}
	}()

	// Ожидаем сигнала завершения или ошибки
	select {
	case err = <-errChan:
		logger.Log.Error("Failed to run app", "error", err)
	case sig := <-sigChan:
		logger.Log.Info("Received signal", "signal", sig)
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = application.Stop(ctx); err != nil {
		logger.Log.Error("Failed to stop app", "error", err)
		os.Exit(1)
	}
}
