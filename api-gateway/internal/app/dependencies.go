package app

import (
	"fmt"
	"github.com/semho/hotel-booking/api-gateway/internal/api/http/handler"
	"github.com/semho/hotel-booking/api-gateway/internal/config"
	"github.com/semho/hotel-booking/pkg/logger"
	bookingpb "github.com/semho/hotel-booking/pkg/proto/booking_v1/booking"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Deps struct {
	BookingHandler *handler.BookingHandler
	bookingConn    *grpc.ClientConn
}

func initDeps(cfg *config.Config) (*Deps, error) {
	// Устанавливаем соединение с booking service
	bookingConn, err := grpc.NewClient(
		cfg.BookingService.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to booking service: %w", err)
	}

	// Создаем gRPC клиент
	bookingClient := bookingpb.NewBookingServiceClient(bookingConn)

	// Создаем HTTP хендлер
	bookingHandler := handler.NewBookingHandler(bookingClient)

	return &Deps{
		BookingHandler: bookingHandler,
		bookingConn:    bookingConn,
	}, nil
}

func (d *Deps) Close() error {
	if err := d.bookingConn.Close(); err != nil {
		logger.Log.Error("failed to close booking service connection", "error", err)
		return err
	}
	return nil
}
