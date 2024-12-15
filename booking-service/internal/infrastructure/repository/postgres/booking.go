package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/semho/hotel-booking/booking-service/internal/domain/model"
	"github.com/semho/hotel-booking/booking-service/internal/domain/port"
	"github.com/semho/hotel-booking/pkg/errors"
	pb "github.com/semho/hotel-booking/pkg/proto/booking_v1/booking"
)

const (
	tableBookings = "bookings"

	idColumn        = "id"
	roomIdColumn    = "room_id"
	userIdColumn    = "user_id"
	guestNameColumn = "guest_name"
	emailColumn     = "guest_email"
	phoneColumn     = "guest_phone"
	checkInColumn   = "check_in"
	checkOutColumn  = "check_out"
	statusColumn    = "status"
	priceColumn     = "total_price"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type bookingRepository struct {
	db      *sqlx.DB
	builder squirrel.StatementBuilderType
}

func NewBookingRepository(db *sqlx.DB) port.BookingRepository {
	return &bookingRepository{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *bookingRepository) GetBookingsForPeriod(ctx context.Context, checkIn, checkOut time.Time) (
	[]model.Booking,
	error,
) {
	// Ищем все бронирования, которые пересекаются с заданным периодом
	query := r.builder.
		Select(
			idColumn,
			roomIdColumn,
			userIdColumn,
			guestNameColumn,
			emailColumn,
			phoneColumn,
			checkInColumn,
			checkOutColumn,
			statusColumn,
			priceColumn,
			createdAtColumn,
			updatedAtColumn,
		).
		From(tableBookings).
		// Условие пересечения периодов
		Where(
			squirrel.And{
				squirrel.Lt{checkInColumn: checkOut},
				squirrel.Gt{checkOutColumn: checkIn},
				// Исключаем отмененные и завершенные бронирования
				squirrel.NotEq{statusColumn: pb.BookingStatus_BOOKING_STATUS_CANCELLED},
				squirrel.NotEq{statusColumn: pb.BookingStatus_BOOKING_STATUS_COMPLETED},
				squirrel.NotEq{statusColumn: pb.BookingStatus_BOOKING_STATUS_NO_SHOW},
			},
		)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var bookings []model.Booking
	if err := r.db.SelectContext(ctx, &bookings, sql, args...); err != nil {
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}

	return bookings, nil
}

func (r *bookingRepository) IsRoomAvailable(ctx context.Context, roomID uuid.UUID, checkIn, checkOut time.Time) (
	bool,
	error,
) {
	query := r.builder.
		Select("COUNT(*)").
		From(tableBookings).
		Where(
			squirrel.And{
				squirrel.Eq{roomIdColumn: roomID},
				squirrel.Lt{checkInColumn: checkOut},
				squirrel.Gt{checkOutColumn: checkIn},
				// Проверяем только активные бронирования
				squirrel.Or{
					squirrel.Eq{statusColumn: pb.BookingStatus_BOOKING_STATUS_PENDING},
					squirrel.Eq{statusColumn: pb.BookingStatus_BOOKING_STATUS_CONFIRMED},
				},
			},
		)

	sql, args, err := query.ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	var count int
	if err := r.db.GetContext(ctx, &count, sql, args...); err != nil {
		return false, fmt.Errorf("failed to check room availability: %w", err)
	}

	return count == 0, nil
}

func (r *bookingRepository) Create(ctx context.Context, booking *model.Booking) error {
	sql, args, err := r.builder.
		Insert(tableBookings).
		Columns(
			roomIdColumn,
			userIdColumn,
			guestNameColumn,
			emailColumn,
			phoneColumn,
			checkInColumn,
			checkOutColumn,
			statusColumn,
			priceColumn,
		).
		Values(
			booking.RoomID,
			booking.UserID,
			booking.GuestName,
			booking.GuestEmail,
			booking.GuestPhone,
			booking.CheckIn,
			booking.CheckOut,
			booking.Status,
			booking.TotalPrice,
		).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	err = r.db.QueryRowContext(ctx, sql, args...).Scan(
		&booking.ID,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}

	return nil
}

func (r *bookingRepository) Update(ctx context.Context, booking *model.Booking) error {
	booking.UpdatedAt = time.Now()

	sql, args, err := r.builder.
		Update(tableBookings).
		Set(roomIdColumn, booking.RoomID).
		Set(userIdColumn, booking.UserID).
		Set(guestNameColumn, booking.GuestName).
		Set(emailColumn, booking.GuestEmail).
		Set(phoneColumn, booking.GuestPhone).
		Set(checkInColumn, booking.CheckIn).
		Set(checkOutColumn, booking.CheckOut).
		Set(statusColumn, booking.Status).
		Set(priceColumn, booking.TotalPrice).
		Set(updatedAtColumn, booking.UpdatedAt).
		Where(squirrel.Eq{idColumn: booking.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return errors.ErrNotFound
	}

	return nil
}
