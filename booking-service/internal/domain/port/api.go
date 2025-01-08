package port

import (
	"context"
	"github.com/google/uuid"
	"github.com/semho/hotel-booking/booking-service/internal/domain/model"
)

type BookingAPI interface {
	GetAvailableRooms(ctx context.Context, params model.SearchParams) ([]model.Room, error)
	// бронь (статус PENDING)
	CreateBooking(ctx context.Context, booking *model.Booking) error
}

// RoomClient определяет интерфейс для взаимодействия с Room Service
type RoomClient interface {
	GetAvailableRooms(ctx context.Context, params model.SearchRoomsParams) ([]model.Room, error)
	GetRoomsCount(ctx context.Context, params model.SearchRoomsParams) (int32, error)
	GetRoomInfo(ctx context.Context, roomID uuid.UUID) (*model.Room, error)
	GetFirstAvailableRoom(ctx context.Context, params model.SearchRoomsParams) (*model.Room, error)
}
