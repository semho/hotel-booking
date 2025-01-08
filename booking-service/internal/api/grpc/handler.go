package grpc

import (
	"context"
	"github.com/semho/hotel-booking/booking-service/internal/api/grpc/mapper"
	"github.com/semho/hotel-booking/booking-service/internal/domain/port"
	bookingpb "github.com/semho/hotel-booking/pkg/proto/booking_v1/booking"
)

type BookingHandler struct {
	bookingpb.UnimplementedBookingServiceServer
	bookingService port.BookingService
	roomClient     port.RoomClient
}

func NewBookingHandler(bookingService port.BookingService, roomClient port.RoomClient) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
		roomClient:     roomClient,
	}
}

func (h *BookingHandler) GetAvailableRooms(
	ctx context.Context,
	req *bookingpb.GetAvailableRoomsRequest,
) (*bookingpb.GetAvailableRoomsResponse, error) {
	// Конвертируем параметры запроса в доменную модель
	paramsForRoom := mapper.ProtoToSearchParams(req)
	// Получаем комнаты через клиент room service
	rooms, err := h.roomClient.GetAvailableRooms(ctx, paramsForRoom)
	if err != nil {
		return nil, err
	}
	// Проверяем доступность через сервис бронирований
	params := mapper.ProtoToSearchParamsInternal(req)
	availableRooms, err := h.bookingService.GetAvailableRooms(ctx, params, rooms)
	if err != nil {
		return nil, err
	}
	// Конвертируем обратно в proto
	return &bookingpb.GetAvailableRoomsResponse{
		Rooms: mapper.RoomsToProto(availableRooms),
	}, nil
}

func (h *BookingHandler) CreateBooking(
	ctx context.Context,
	req *bookingpb.CreateBookingRequest,
) (*bookingpb.CreateBookingResponse, error) {
	booking := mapper.ProtoToBooking(req)

	err := h.bookingService.CreateBooking(ctx, booking, req.GetType(), req.GetCapacity())
	if err != nil {
		return nil, err
	}

	return &bookingpb.CreateBookingResponse{
		Booking: mapper.BookingToProto(booking),
	}, nil
}
