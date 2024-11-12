package app

import (
	"context"
	"fmt"
	"github.com/semho/hotel-booking/pkg/logger"
	pb "github.com/semho/hotel-booking/pkg/proto/room_v1"
	"github.com/semho/hotel-booking/room-service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	// Регистрируем сервисы
	pb.RegisterRoomServiceServer(grpcServer, deps.RoomHandler)

	reflection.Register(grpcServer)

	return &App{
		grpcServer: grpcServer,
		deps:       deps,
		cfg:        cfg,
	}, nil
}

func (a *App) Run() error {
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

func (a *App) Stop(_ context.Context) error {
	logger.Log.Info("shutting down gRPC server")
	a.grpcServer.GracefulStop()
	if err := a.deps.DB.Close(); err != nil {
		return fmt.Errorf("failed to close db connection: %w", err)
	}
	logger.Log.Info("gRPC server stopped")
	return nil
}
