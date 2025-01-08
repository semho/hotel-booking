package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/semho/hotel-booking/pkg/proto/room_v1/room"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
	"strings"
	"time"

	"github.com/semho/hotel-booking/booking-service/internal/domain/model"
	"github.com/semho/hotel-booking/booking-service/internal/domain/port"
	"github.com/semho/hotel-booking/pkg/errors"
	pb "github.com/semho/hotel-booking/pkg/proto/booking_v1/booking"
)

type bookingService struct {
	uow         port.BookingUnitOfWork
	bookingRepo port.BookingRepository
	roomClient  port.RoomClient
}

func NewBookingService(
	bookingRepo port.BookingRepository,
	uow port.BookingUnitOfWork,
	roomClient port.RoomClient,
) port.BookingService {
	return &bookingService{
		uow:         uow,
		bookingRepo: bookingRepo,
		roomClient:  roomClient,
	}
}

// Валидация статуса
func (s *bookingService) validateStatus(status model.BookingStatus) error {
	switch status {
	case pb.BookingStatus_BOOKING_STATUS_PENDING,
		pb.BookingStatus_BOOKING_STATUS_CONFIRMED,
		pb.BookingStatus_BOOKING_STATUS_CANCELLED,
		pb.BookingStatus_BOOKING_STATUS_COMPLETED,
		pb.BookingStatus_BOOKING_STATUS_NO_SHOW:
		return nil
	default:
		return errors.WithMessage(
			errors.ErrInvalidInput,
			fmt.Sprintf("invalid booking status: %d", status),
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

func (s *bookingService) calculateTotalPrice(room *model.Room, checkIn, checkOut time.Time) (float64, error) {
	price, err := decimal.NewFromString(room.Price)
	if err != nil {
		return 0, fmt.Errorf("failed to parse room price: %w", err)
	}

	// Вычисляем количество ночей проживания (округляем вверх для неполных суток)
	nights := decimal.NewFromInt(int64(math.Ceil(checkOut.Sub(checkIn).Hours() / 24)))
	// Рассчитываем полную стоимость проживания (базовая цена * количество ночей)
	totalPrice := price.Mul(nights)

	res, ok := totalPrice.Float64()
	if !ok {
		return 0, fmt.Errorf("failed to convert total price to float64")
	}
	return res, nil
}

func (s *bookingService) CreateBooking(
	ctx context.Context,
	booking *model.Booking,
	roomType room.RoomType,
	roomCapacity int32,
) error {
	return s.uow.WithinTransaction(
		ctx, func(txCtx context.Context) error {
			// Валидация входных данных
			if err := s.validateBooking(
				booking.GuestName,
				booking.GuestEmail,
				timestamppb.New(booking.CheckIn),
				timestamppb.New(booking.CheckOut),
			); err != nil {
				return err
			}

			// 1. Получаем список подходящих комнат из room-service
			params := model.SearchRoomsParams{
				Capacity: &roomCapacity,
				Type:     &roomType,
			}
			rooms, err := s.roomClient.GetAvailableRooms(ctx, params)
			if err != nil {
				return err
			}

			// 2. Формируем список room_id
			roomIDs := make([]uuid.UUID, len(rooms), len(rooms))
			var roomUUID uuid.UUID
			for i, room := range rooms {
				roomUUID, err = uuid.Parse(room.ID)
				if err != nil {
					return err
				}
				roomIDs[i] = roomUUID
			}

			// 3. Получаем список броней
			bookedRoomIDs, err := s.bookingRepo.GetBookedRoomIDs(
				txCtx,
				roomIDs,
				booking.CheckIn,
				booking.CheckOut,
				true, // с блокировкой
			)
			if err != nil {
				return err
			}

			// 4. Создаем множество доступных room_id
			availableRoomIds := make(map[uuid.UUID]struct{}, len(roomIDs))
			for _, id := range roomIDs {
				availableRoomIds[id] = struct{}{}
			}

			// 5. Удаляем забронированные комнаты из доступных
			for _, bookedRoomID := range bookedRoomIDs {
				delete(availableRoomIds, bookedRoomID)
			}

			if len(availableRoomIds) == 0 {
				return fmt.Errorf("no rooms available")
			}

			// 6. Выбираем первую свободную комнату
			var selectedRoomId uuid.UUID
			for roomId := range availableRoomIds {
				selectedRoomId = roomId
				break
			}

			// 7. Получаем детали комнаты для расчета цены
			selectedRoom, err := s.roomClient.GetRoomInfo(ctx, selectedRoomId)
			if err != nil {
				return err
			}

			// 8. Рассчитываем полную стоимость
			totalPrice, err := s.calculateTotalPrice(selectedRoom, booking.CheckIn, booking.CheckOut)
			if err != nil {
				return err
			}

			// 9. Заполняем оставшиеся поля бронирования
			booking.RoomID = selectedRoomId
			booking.TotalPrice = totalPrice

			// 10. Создаем бронь
			if err = s.bookingRepo.Create(txCtx, booking); err != nil {
				return err
			}

			// 11. Создаем начальный статус PENDING
			statusHistory := &model.BookingStatusHistory{
				BookingID: booking.ID,
				Status:    pb.BookingStatus_BOOKING_STATUS_PENDING,
				ChangedBy: "system",
				Reason:    "Initial booking creation",
			}

			booking.CurrentStatus = statusHistory

			return s.bookingRepo.AddBookingStatus(txCtx, booking.ID, statusHistory)

		},
	)
}

func (s *bookingService) UpdateBookingStatus(
	ctx context.Context,
	bookingID uuid.UUID,
	status pb.BookingStatus,
	reason string,
	changedBy string,
) error {
	// Проверяем существование брони
	_, err := s.bookingRepo.GetBookingWithStatus(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("failed to get booking: %w", err)
	}

	// Создаем новую запись в истории статусов
	statusHistory := &model.BookingStatusHistory{
		BookingID: bookingID,
		Status:    status,
		Reason:    reason,
		ChangedBy: changedBy,
	}

	if err = s.bookingRepo.AddBookingStatus(ctx, bookingID, statusHistory); err != nil {
		return fmt.Errorf("failed to update booking status: %w", err)
	}

	return nil
}

func (s *bookingService) validateBooking(username, email string, checkIn, checkOut *timestamppb.Timestamp) error {
	// Проверка обязательных полей
	if username == "" {
		return errors.WithMessage(
			errors.ErrInvalidInput,
			"guest name is required",
		)
	}

	if email == "" {
		return errors.WithMessage(
			errors.ErrInvalidInput,
			"guest email is required",
		)
	}

	// Базовая валидация email
	if !strings.Contains(email, "@") {
		return errors.WithMessage(
			errors.ErrInvalidInput,
			"invalid email format",
		)
	}

	// Проверка дат
	if checkIn == nil || checkOut == nil {
		return errors.WithMessage(
			errors.ErrInvalidInput,
			"check-in and check-out dates are required",
		)
	}

	checkInTime := checkIn.AsTime()
	checkOutTime := checkOut.AsTime()

	if checkOutTime.Before(checkInTime) {
		return errors.WithMessage(
			errors.ErrInvalidInput,
			"check-out date cannot be before check-in date",
		)
	}

	return nil
}
