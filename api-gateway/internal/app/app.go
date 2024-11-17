package app

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/semho/hotel-booking/api-gateway/internal/config"
	"github.com/semho/hotel-booking/pkg/logger"
	"net/http"
	"time"
)

type App struct {
	httpServer *http.Server
	deps       *Deps
}

func New(cfg *config.Config) (*App, error) {
	// Инициализируем зависимости
	deps, err := initDeps(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init dependencies: %w", err)
	}

	// Создаем роутер
	router := chi.NewRouter()

	// Добавляем CORS middleware перед другими middleware
	router.Use(
		cors.Handler(
			cors.Options{
				AllowedOrigins:   []string{"*"}, // TODO: Для прода заменить на конкретные домены
				AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
				ExposedHeaders:   []string{"Link"},
				AllowCredentials: true,
				MaxAge:           300, // Maximum value not ignored by any of major browsers
			},
		),
	)

	// Добавляем middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// Регистрируем обработчики
	deps.BookingHandler.RegisterRoutes(router)

	// Создаем HTTP сервер
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &App{
		httpServer: httpServer,
		deps:       deps,
	}, nil
}

func (a *App) Run() error {
	logger.Log.Info("starting HTTP server", "port", a.httpServer.Addr)
	if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	logger.Log.Info("shutting down HTTP server")

	if err := a.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}
	logger.Log.Info("HTTP server stopped")

	// Закрываем клиентские соединения
	if err := a.deps.Close(); err != nil {
		return fmt.Errorf("failed to close dependencies: %w", err)
	}

	return nil
}
