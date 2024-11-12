package main

import (
	"context"
	"github.com/semho/hotel-booking/pkg/logger"
	"github.com/semho/hotel-booking/room-service/internal/app"
	"github.com/semho/hotel-booking/room-service/internal/config"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
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
	logger.Log.Info("Starting servers")
	errChan := make(chan error)
	go func() {
		if err = application.Run(); err != nil {
			errChan <- err
		}
	}()

	// Ожидаем сигнала завершения или ошибки
	select {
	case err = <-errChan:
		log.Printf("Failed to run app: %v", err)
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = application.Stop(ctx); err != nil {
		log.Printf("Failed to stop app: %v", err)
	}
}
