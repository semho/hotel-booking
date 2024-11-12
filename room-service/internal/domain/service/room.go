package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/semho/hotel-booking/pkg/errors"
	"github.com/semho/hotel-booking/pkg/logger"
	"github.com/semho/hotel-booking/room-service/internal/domain/model"
	"github.com/semho/hotel-booking/room-service/internal/domain/port"
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
		"check_in", params.CheckIn,
		"check_out", params.CheckOut,
		"capacity", params.Capacity,
		"type", params.Type,
	)

	if params.CheckIn.IsZero() || params.CheckOut.IsZero() {
		return nil, errors.WithMessage(errors.ErrInvalidInput, "check-in and check-out dates are required")
	}

	if params.CheckIn.After(params.CheckOut) {
		return nil, errors.WithMessage(errors.ErrInvalidInput, "check-in date must be before check-out date")
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

	if room.Price <= 0 {
		return errors.WithMessage(errors.ErrInvalidInput, "price must be greater than 0")
	}

	if room.Capacity <= 0 {
		return errors.WithMessage(errors.ErrInvalidInput, "capacity must be greater than 0")
	}

	if room.Status == "" {
		room.Status = model.RoomStatusAvailable
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

	if room.Price <= 0 {
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
