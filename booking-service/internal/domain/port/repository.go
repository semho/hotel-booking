package port

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/semho/hotel-booking/booking-service/internal/domain/model"
	"time"
)

type SQLExecutor interface {
	sqlx.ExtContext
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type BookingRepository interface {
	GetBookingsForPeriod(ctx context.Context, checkIn, checkOut time.Time) ([]model.Booking, error)
	Create(ctx context.Context, booking *model.Booking) error
	// Добавление статуса в историю
	AddBookingStatus(ctx context.Context, bookingID uuid.UUID, status *model.BookingStatusHistory) error
	// Получение брони с текущим статусом
	GetBookingWithStatus(ctx context.Context, bookingID uuid.UUID) (*model.BookingWithStatus, error)
	// Получение истории статусов брони
	GetBookingStatusHistory(ctx context.Context, bookingID uuid.UUID) ([]model.BookingStatusHistory, error)
	GetBookedRoomIDs(ctx context.Context, roomIDs []uuid.UUID, checkIn, checkOut time.Time, forUpdate bool) (
		[]uuid.UUID,
		error,
	)
}
