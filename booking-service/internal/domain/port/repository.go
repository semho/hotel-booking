package port

import (
	"context"
	"github.com/google/uuid"
	"github.com/semho/hotel-booking/booking-service/internal/domain/model"
	"time"
)

type BookingRepository interface {
	GetBookingsForPeriod(ctx context.Context, checkIn, checkOut time.Time) ([]model.Booking, error)
	IsRoomAvailable(ctx context.Context, roomID uuid.UUID, checkIn, checkOut time.Time) (bool, error)
	Create(ctx context.Context, booking *model.Booking) error
	// другие методы...
}
