package port

import (
	"context"
	"github.com/semho/hotel-booking/booking-service/internal/domain/model"
)

type BookingAPI interface {
	GetAvailableRooms(ctx context.Context, params model.SearchParams) ([]model.Room, error)
}

// RoomClient определяет интерфейс для взаимодействия с Room Service
type RoomClient interface {
	GetAvailableRooms(ctx context.Context, params model.SearchParams) ([]model.Room, error)
}
