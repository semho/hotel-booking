package service

import (
	"context"
	"fmt"
	"github.com/semho/hotel-booking/booking-service/internal/domain/model"
	"github.com/semho/hotel-booking/booking-service/internal/domain/port"
	"github.com/semho/hotel-booking/pkg/errors"
)

type bookingService struct {
	bookingRepo port.BookingRepository
}

func NewBookingService(bookingRepo port.BookingRepository) port.BookingService {
	return &bookingService{
		bookingRepo: bookingRepo,
	}
}

// Валидация статуса
func (s *bookingService) validateStatus(status model.BookingStatus) error {
	switch status {
	case model.BookingStatusPending,
		model.BookingStatusConfirmed,
		model.BookingStatusCancelled,
		model.BookingStatusCompleted,
		model.BookingStatusNoShow:
		return nil
	default:
		return errors.WithMessage(
			errors.ErrInvalidInput,
			fmt.Sprintf("invalid booking status: %s", status),
		)
	}
}

func (s *bookingService) GetAvailableRooms(
	ctx context.Context,
	params model.SearchParams,
	rooms []model.Room,
) ([]model.Room, error) {
	// Получаем бронирования на период
	bookings, err := s.bookingRepo.GetBookingsForPeriod(
		ctx,
		params.CheckIn,
		params.CheckOut,
	)
	if err != nil {
		return nil, err
	}

	return s.filterAvailableRooms(rooms, bookings), nil
}

func (s *bookingService) filterAvailableRooms(rooms []model.Room, bookings []model.Booking) []model.Room {
	bookedRooms := make(map[string]struct{})
	for _, booking := range bookings {
		bookedRooms[booking.RoomID.String()] = struct{}{}
	}

	var availableRooms []model.Room
	for _, room := range rooms {
		if _, booked := bookedRooms[room.ID]; !booked {
			availableRooms = append(availableRooms, room)
		}
	}

	return availableRooms
}

func (s *bookingService) Create(ctx context.Context, booking *model.Booking) error {
	if err := s.validateStatus(booking.Status); err != nil {
		return err
	}

	// Другие проверки бизнес-логики
	if booking.CheckOut.Before(booking.CheckIn) {
		return errors.WithMessage(
			errors.ErrInvalidInput,
			"check-out date cannot be before check-in date",
		)
	}

	// Проверка доступности номера
	available, err := s.bookingRepo.IsRoomAvailable(ctx, booking.RoomID, booking.CheckIn, booking.CheckOut)
	if err != nil {
		return fmt.Errorf("failed to check room availability: %w", err)
	}
	if !available {
		return errors.WithMessage(
			errors.ErrConflict,
			"room is not available for selected dates",
		)
	}

	return s.bookingRepo.Create(ctx, booking)
}
