package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/semho/hotel-booking/pkg/errors"
	"github.com/semho/hotel-booking/pkg/logger"
	pb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
	"github.com/semho/hotel-booking/room-service/internal/domain/model"
	"github.com/semho/hotel-booking/room-service/internal/domain/port"
	"github.com/shopspring/decimal"
)

type RoomService struct {
	repo port.RoomRepository
}

func NewRoomService(repo port.RoomRepository) *RoomService {
	return &RoomService{
		repo: repo,
	}
}

func (s *RoomService) GetAvailableRooms(ctx context.Context, params model.SearchParams) ([]model.Room, error) {
	logger.Log.Info(
		"getting available rooms with params",
		"capacity", params.Capacity,
		"type", params.Type,
		"status", params.Status,
	)

	// Если статус не указан, устанавливаем AVAILABLE по умолчанию
	if params.Status == nil {
		availableStatus := pb.RoomStatus_ROOM_STATUS_AVAILABLE
		params.Status = &availableStatus
	}

	rooms, err := s.repo.GetAvailableRooms(ctx, params)
	if err != nil {
		logger.Log.Error(
			"failed to get rooms from repository",
			"error", err,
		)
		return nil, err
	}
	logger.Log.Info(
		"successfully retrieved rooms",
		"count", len(rooms),
	)

	return rooms, nil
}

func (s *RoomService) GetByID(ctx context.Context, id uuid.UUID) (*model.Room, error) {
	if id == uuid.Nil {
		return nil, errors.WithMessage(errors.ErrInvalidInput, "invalid room id")
	}

	room, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (s *RoomService) Create(ctx context.Context, room *model.Room) error {
	if room.RoomNumber == "" {
		return errors.WithMessage(errors.ErrInvalidInput, "room number is required")
	}

	if room.Price.Cmp(decimal.Zero) <= 0 {
		return errors.WithMessage(errors.ErrInvalidInput, "price must be greater than 0")
	}

	if room.Capacity <= 0 {
		return errors.WithMessage(errors.ErrInvalidInput, "capacity must be greater than 0")
	}

	if room.Type == 0 {
		room.Type = pb.RoomType_ROOM_TYPE_STANDARD
	}

	if room.Status == 0 {
		room.Status = pb.RoomStatus_ROOM_STATUS_AVAILABLE
	}

	return s.repo.Create(ctx, room)
}

func (s *RoomService) Update(ctx context.Context, room *model.Room) error {
	if room.ID == uuid.Nil {
		return errors.WithMessage(errors.ErrInvalidInput, "invalid room id")
	}

	if room.RoomNumber == "" {
		return errors.WithMessage(errors.ErrInvalidInput, "room number is required")
	}

	if room.Price.Cmp(decimal.Zero) <= 0 {
		return errors.WithMessage(errors.ErrInvalidInput, "price must be greater than 0")
	}

	if room.Capacity <= 0 {
		return errors.WithMessage(errors.ErrInvalidInput, "capacity must be greater than 0")
	}

	return s.repo.Update(ctx, room)
}

func (s *RoomService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.WithMessage(errors.ErrInvalidInput, "invalid room id")
	}

	return s.repo.Delete(ctx, id)
}

func (s *RoomService) GetRoomsCount(ctx context.Context, params model.SearchParams) (int32, error) {
	return s.repo.GetRoomsCount(ctx, params)
}

func (s *RoomService) GetFirstAvailableRoom(ctx context.Context, params model.SearchParams) (*model.Room, error) {
	return s.repo.GetFirstAvailableRoom(ctx, params)
}
