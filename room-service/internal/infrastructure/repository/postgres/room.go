package postgres

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/semho/hotel-booking/pkg/errors"
	"github.com/semho/hotel-booking/pkg/logger"
	"github.com/semho/hotel-booking/room-service/internal/domain/model"
	"github.com/semho/hotel-booking/room-service/internal/domain/port"
)

const (
	tableRooms = "rooms"

	idColumn         = "id"
	roomNumberColumn = "room_number"
	typeColumn       = "type"
	priceColumn      = "price"
	capacityColumn   = "capacity"
	statusColumn     = "status"
	amenitiesColumn  = "amenities"
	createdAtColumn  = "created_at"
	updatedAtColumn  = "updated_at"
)

type roomRepository struct {
	db      *sqlx.DB
	builder squirrel.StatementBuilderType
}

func NewRoomRepository(db *sqlx.DB) port.RoomRepository {
	return &roomRepository{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *roomRepository) GetAvailableRooms(ctx context.Context, params model.SearchParams) ([]model.Room, error) {
	query := r.builder.Select(
		idColumn,
		roomNumberColumn,
		typeColumn,
		priceColumn,
		capacityColumn,
		statusColumn,
		amenitiesColumn,
		createdAtColumn,
		updatedAtColumn,
	).
		From(tableRooms).
		Where(squirrel.Eq{statusColumn: model.RoomStatusAvailable})

	// Применяем фильтры
	if params.Type != nil {
		query = query.Where(squirrel.Eq{typeColumn: *params.Type})
		logger.Log.Info("applying type filter", typeColumn, *params.Type)
	}
	if params.Capacity != nil {
		query = query.Where(squirrel.GtOrEq{capacityColumn: *params.Capacity})
		logger.Log.Info("applying capacity filter", capacityColumn, *params.Capacity)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	logger.Log.Info(
		"executing query",
		"sql", sql,
		"args", args,
	)

	var rooms []model.Room
	if err := r.db.SelectContext(ctx, &rooms, sql, args...); err != nil {
		logger.Log.Error(
			"database error",
			"error", err,
			"sql", sql,
			"args", args,
		)
		return nil, fmt.Errorf("failed to fetch rooms: %w", err)
	}

	logger.Log.Info(
		"found rooms",
		"count", len(rooms),
		"first_room_number", firstRoomNumber(rooms),
	)

	return rooms, nil
}

func firstRoomNumber(rooms []model.Room) string {
	if len(rooms) > 0 {
		return rooms[0].RoomNumber
	}
	return "<none>"
}

func (r *roomRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Room, error) {
	sql, args, err := r.builder.
		Select(
			idColumn,
			roomNumberColumn,
			typeColumn,
			priceColumn,
			capacityColumn,
			statusColumn,
			amenitiesColumn,
			createdAtColumn,
			updatedAtColumn,
		).
		From(tableRooms).
		Where(squirrel.Eq{idColumn: id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var room model.Room
	if err := r.db.GetContext(ctx, &room, sql, args...); err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *roomRepository) Create(ctx context.Context, room *model.Room) error {
	sql, args, err := r.builder.
		Insert(tableRooms).
		Columns(
			roomNumberColumn,
			typeColumn,
			priceColumn,
			capacityColumn,
			statusColumn,
			amenitiesColumn,
		).
		Values(
			room.RoomNumber,
			room.Type,
			room.Price,
			room.Capacity,
			room.Status,
			pq.Array(room.Amenities),
		).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()
	if err != nil {
		return err
	}

	return r.db.QueryRowContext(ctx, sql, args...).Scan(
		&room.ID,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
}

func (r *roomRepository) Update(ctx context.Context, room *model.Room) error {
	sql, args, err := r.builder.
		Update(tableRooms).
		Set(roomNumberColumn, room.RoomNumber).
		Set(typeColumn, room.Type).
		Set(priceColumn, room.Price).
		Set(capacityColumn, room.Capacity).
		Set(statusColumn, room.Status).
		Set(amenitiesColumn, pq.Array(room.Amenities)).
		Set(updatedAtColumn, squirrel.Expr("CURRENT_TIMESTAMP")).
		Where(squirrel.Eq{"id": room.ID}).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *roomRepository) Delete(ctx context.Context, id uuid.UUID) error {
	sql, args, err := r.builder.
		Delete(tableRooms).
		Where(squirrel.Eq{idColumn: id}).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}
