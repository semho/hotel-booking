package app

import (
	"fmt"

	"github.com/semho/hotel-booking/api-gateway/internal/api/http/handler"
	"github.com/semho/hotel-booking/api-gateway/internal/api/http/middleware"
	"github.com/semho/hotel-booking/api-gateway/internal/config"
	"github.com/semho/hotel-booking/pkg/logger"
	authpb "github.com/semho/hotel-booking/pkg/proto/auth_v1/auth"
	bookingpb "github.com/semho/hotel-booking/pkg/proto/booking_v1/booking"
	roompb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Deps struct {
	BookingHandler *handler.BookingHandler
	RoomHandler    *handler.RoomHandler
	AuthHandler    *handler.AuthHandler
	bookingConn    *grpc.ClientConn
	authConn       *grpc.ClientConn
	roomConn       *grpc.ClientConn
}

func initDeps(cfg *config.Config) (*Deps, error) {
	// Устанавливаем соединение с booking service
	logger.Log.Info("connecting to booking service", "address", cfg.BookingService.Address)
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

	// Устанавливаем соединение с room service
	roomConn, err := grpc.NewClient(
		cfg.RoomService.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to room service: %w", err)
	}

	// Создаем gRPC клиент
	bookingClient := bookingpb.NewBookingServiceClient(bookingConn)
	authClient := authpb.NewAuthServiceClient(authConn)
	roomClient := roompb.NewRoomServiceClient(roomConn)

	// Создаем middleware
	authMiddleware := middleware.NewAuthMiddleware(authClient)

	// Создаем HTTP хендлер
	bookingHandler := handler.NewBookingHandler(bookingClient, authMiddleware)
	authHandler := handler.NewAuthHandler(authClient)
	roomHandler := handler.NewRoomHandler(roomClient, authMiddleware)

	return &Deps{
		BookingHandler: bookingHandler,
		AuthHandler:    authHandler,
		RoomHandler:    roomHandler,
		bookingConn:    bookingConn,
		authConn:       authConn,
		roomConn:       roomConn,
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
	if err := d.roomConn.Close(); err != nil {
		logger.Log.Error("failed to close room service connection", "error", err)
	}
	return nil
}
