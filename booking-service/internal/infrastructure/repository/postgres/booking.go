package postgres

import (
	"context"
	"fmt"
	pb "github.com/semho/hotel-booking/pkg/proto/booking_v1/booking"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/semho/hotel-booking/booking-service/internal/domain/model"
	"github.com/semho/hotel-booking/booking-service/internal/domain/port"
)

const (
	bookingsTable      = "bookings"
	statusHistoryTable = "booking_status_history"

	// Columns for bookings
	idColumn        = "id"
	roomIdColumn    = "room_id"
	userIdColumn    = "user_id"
	guestNameColumn = "guest_name"
	emailColumn     = "guest_email"
	phoneColumn     = "guest_phone"
	checkInColumn   = "check_in"
	checkOutColumn  = "check_out"
	priceColumn     = "total_price"
	createdAtColumn = "created_at"

	// Columns for status history
	statusIdColumn  = "id"
	bookingIdColumn = "booking_id"
	statusColumn    = "status"
	reasonColumn    = "reason"
	changedByColumn = "changed_by"
	changedAtColumn = "changed_at"
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

// обертка транзакций
func (r *bookingRepository) getExecutor(ctx context.Context) port.SQLExecutor {
	tx, ok := ctx.Value("tx").(*sqlx.Tx)
	if ok {
		return tx
	}
	return r.db
}

func (r *bookingRepository) GetBookingsForPeriod(ctx context.Context, checkIn, checkOut time.Time) (
	[]model.Booking,
	error,
) {
	query := r.builder.
		Select(
			"b.*",
			"bsh.id as status_id",
			"bsh.status as status_status",
			"bsh.reason as status_reason",
			"bsh.changed_by as status_changed_by",
			"bsh.changed_at as status_changed_at",
		).
		From(fmt.Sprintf("%s AS b", bookingsTable)).
		LeftJoin(
			fmt.Sprintf(
				"(SELECT DISTINCT ON (booking_id) * "+
					"FROM %s ORDER BY booking_id, changed_at DESC) AS bsh ON bsh.booking_id = b.id",
				statusHistoryTable,
			),
		).
		Where(
			squirrel.And{
				squirrel.Lt{fmt.Sprintf("%s.%s", "b", checkInColumn): checkOut},
				squirrel.Gt{fmt.Sprintf("%s.%s", "b", checkOutColumn): checkIn},
				// фильтр по активным статусам
				squirrel.Or{
					squirrel.Eq{"bsh.status": pb.BookingStatus_BOOKING_STATUS_PENDING},
					squirrel.Eq{"bsh.status": pb.BookingStatus_BOOKING_STATUS_CONFIRMED},
				},
			},
		).
		OrderBy(fmt.Sprintf("%s.%s", "b", createdAtColumn))

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var rows []model.BookingRow
	if err := r.getExecutor(ctx).SelectContext(ctx, &rows, sql, args...); err != nil {
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}

	result := make([]model.Booking, len(rows))
	for i, row := range rows {
		result[i] = row.Booking
		result[i].CurrentStatus = &model.BookingStatusHistory{
			ID:        row.StatusID,
			BookingID: row.ID,
			Status:    row.StatusStatus,
			Reason:    row.StatusReason,
			ChangedBy: row.StatusChangedBy,
			ChangedAt: row.StatusChangedAt,
		}
	}

	return result, nil
}

func (r *bookingRepository) Create(ctx context.Context, booking *model.Booking) error {
	query := r.builder.
		Insert(bookingsTable).
		Columns(
			roomIdColumn,
			userIdColumn,
			guestNameColumn,
			emailColumn,
			phoneColumn,
			checkInColumn,
			checkOutColumn,
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
			booking.TotalPrice,
		).
		Suffix("RETURNING id, created_at")

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	return r.getExecutor(ctx).QueryRowContext(ctx, sql, args...).Scan(&booking.ID, &booking.CreatedAt)
}

func (r *bookingRepository) AddBookingStatus(
	ctx context.Context,
	bookingID uuid.UUID,
	status *model.BookingStatusHistory,
) error {
	query := r.builder.
		Insert(statusHistoryTable).
		Columns(
			bookingIdColumn,
			statusColumn,
			reasonColumn,
			changedByColumn,
		).
		Values(
			bookingID,
			status.Status,
			status.Reason,
			status.ChangedBy,
		).
		Suffix("RETURNING id, changed_at")

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	return r.getExecutor(ctx).QueryRowContext(ctx, sql, args...).Scan(&status.ID, &status.ChangedAt)
}

func (r *bookingRepository) GetBookingWithStatus(ctx context.Context, bookingID uuid.UUID) (
	*model.BookingWithStatus,
	error,
) {
	query := r.builder.
		Select(
			fmt.Sprintf("%s.*", bookingsTable),
			fmt.Sprintf("%s.status as current_status", statusHistoryTable),
			fmt.Sprintf("%s.changed_at as status_changed_at", statusHistoryTable),
		).
		From(bookingsTable).
		Join(
			fmt.Sprintf(
				"(SELECT DISTINCT ON (booking_id) booking_id, status, changed_at "+
					"FROM %s ORDER BY booking_id, changed_at DESC) AS %s ON %s.booking_id = %s.id",
				statusHistoryTable, statusHistoryTable, statusHistoryTable, bookingsTable,
			),
		).
		Where(squirrel.Eq{fmt.Sprintf("%s.id", bookingsTable): bookingID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var booking model.BookingWithStatus
	if err := r.getExecutor(ctx).GetContext(ctx, &booking, sql, args...); err != nil {
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	return &booking, nil
}
func (r *bookingRepository) GetBookingStatusHistory(
	ctx context.Context,
	bookingID uuid.UUID,
) ([]model.BookingStatusHistory, error) {
	query := r.builder.
		Select("*").
		From(statusHistoryTable).
		Where(squirrel.Eq{bookingIdColumn: bookingID}).
		OrderBy(fmt.Sprintf("%s DESC", changedAtColumn))

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var history []model.BookingStatusHistory
	if err := r.getExecutor(ctx).SelectContext(ctx, &history, sql, args...); err != nil {
		return nil, fmt.Errorf("failed to get status history: %w", err)
	}

	return history, nil
}

func (r *bookingRepository) GetBookedRoomIDs(
	ctx context.Context,
	roomIDs []uuid.UUID,
	checkIn, checkOut time.Time,
	forUpdate bool,
) ([]uuid.UUID, error) {
	// Подзапрос для получения room_id
	subQuery := r.builder.
		Select("b.room_id").
		From(fmt.Sprintf("%s AS b", bookingsTable)).
		Where(
			squirrel.And{
				squirrel.Lt{"b.check_in": checkOut},
				squirrel.Gt{"b.check_out": checkIn},
				squirrel.Eq{"b.room_id": roomIDs},
				squirrel.Expr(
					"EXISTS (SELECT 1 FROM booking_status_history WHERE "+
						"booking_id = b.id AND status IN (?, ?) "+
						"AND changed_at = (SELECT MAX(changed_at) FROM booking_status_history WHERE booking_id = b.id))",
					pb.BookingStatus_BOOKING_STATUS_PENDING,
					pb.BookingStatus_BOOKING_STATUS_CONFIRMED,
				),
			},
		)

	// Основной запрос с блокировкой
	query := r.builder.
		Select("t.room_id").
		FromSelect(subQuery, "t")

	if forUpdate {
		query = query.Suffix("FOR UPDATE")
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var bookedRoomIDs []uuid.UUID
	if err = r.getExecutor(ctx).SelectContext(ctx, &bookedRoomIDs, sql, args...); err != nil {
		return nil, fmt.Errorf("failed to get booked room IDs: %w", err)
	}

	return bookedRoomIDs, nil
}
