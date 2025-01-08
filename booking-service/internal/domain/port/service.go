package port

import (
	"context"
	"github.com/google/uuid"
	"github.com/semho/hotel-booking/booking-service/internal/domain/model"
	pb "github.com/semho/hotel-booking/pkg/proto/booking_v1/booking"
	"github.com/semho/hotel-booking/pkg/proto/room_v1/room"
)

type BookingUnitOfWork interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type BookingService interface {
	GetAvailableRooms(ctx context.Context, params model.SearchParams, rooms []model.Room) ([]model.Room, error)
	// Создание брони со статусом
	CreateBooking(ctx context.Context, booking *model.Booking, roomType room.RoomType, roomCapacity int32) error

	// Обновление статуса брони
	UpdateBookingStatus(
		ctx context.Context,
		bookingID uuid.UUID,
		status pb.BookingStatus,
		reason string,
		changedBy string,
	) error
}
