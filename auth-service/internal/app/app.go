package app

import (
	"context"
	"fmt"
	"github.com/semho/hotel-booking/auth-service/internal/config"
	"github.com/semho/hotel-booking/pkg/logger"
	pb "github.com/semho/hotel-booking/pkg/proto/auth_v1/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type App struct {
	grpcServer *grpc.Server
	deps       *Deps
	cfg        *config.Config
}

func New(cfg *config.Config) (*App, error) {
	// Инициализируем зависимости
	deps, err := initDeps(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init dependencies: %w", err)
	}

	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()

	// Регистрируем сервисы
	pb.RegisterAuthServiceServer(grpcServer, deps.AuthHandler)

	// Включаем reflection для удобства разработки
	reflection.Register(grpcServer)

	return &App{
		grpcServer: grpcServer,
		deps:       deps,
		cfg:        cfg,
	}, nil
}

func (a *App) Run() error {
	// Запускаем gRPC сервер
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.cfg.GRPC.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	logger.Log.Info("starting gRPC server", "port", a.cfg.GRPC.Port)

	if err := a.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	// Graceful shutdown
	logger.Log.Info("stopping gRPC server")

	// Останавливаем gRPC сервер
	a.grpcServer.GracefulStop()

	// Закрываем соединение с БД
	if err := a.deps.db.Close(); err != nil {
		return fmt.Errorf("failed to close db connection: %w", err)
	}

	logger.Log.Info("gRPC server stopped")
	return nil
}
