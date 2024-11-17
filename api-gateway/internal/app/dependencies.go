package app

import (
	"fmt"
	"github.com/semho/hotel-booking/api-gateway/internal/api/http/handler"
	"github.com/semho/hotel-booking/api-gateway/internal/api/http/middleware"
	"github.com/semho/hotel-booking/api-gateway/internal/config"
	"github.com/semho/hotel-booking/pkg/logger"
	authpb "github.com/semho/hotel-booking/pkg/proto/auth_v1/auth"
	bookingpb "github.com/semho/hotel-booking/pkg/proto/booking_v1/booking"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Deps struct {
	BookingHandler *handler.BookingHandler
	AuthHandler    *handler.AuthHandler
	bookingConn    *grpc.ClientConn
	authConn       *grpc.ClientConn
}

func initDeps(cfg *config.Config) (*Deps, error) {
	logger.Log.Info("connecting to auth service", "address", cfg.AuthService.Address)
	// Устанавливаем соединение с booking service
	bookingConn, err := grpc.NewClient(
		cfg.BookingService.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to booking service: %w", err)
	}

	// Устанавливаем соединение с auth service
	authConn, err := grpc.NewClient(
		cfg.AuthService.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	// Создаем gRPC клиент
	bookingClient := bookingpb.NewBookingServiceClient(bookingConn)
	authClient := authpb.NewAuthServiceClient(authConn)

	// Создаем middleware
	authMiddleware := middleware.NewAuthMiddleware(authClient)

	// Создаем HTTP хендлер
	bookingHandler := handler.NewBookingHandler(bookingClient, authMiddleware)
	authHandler := handler.NewAuthHandler(authClient)

	return &Deps{
		BookingHandler: bookingHandler,
		AuthHandler:    authHandler,
		bookingConn:    bookingConn,
		authConn:       authConn,
	}, nil
}

func (d *Deps) Close() error {
	if err := d.bookingConn.Close(); err != nil {
		logger.Log.Error("failed to close booking service connection", "error", err)
		return err
	}
	if err := d.authConn.Close(); err != nil {
		logger.Log.Error("failed to close auth service connection", "error", err)
	}
	return nil
}
