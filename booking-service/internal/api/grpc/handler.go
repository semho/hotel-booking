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
	params := mapper.ProtoToSearchParams(req)

	// Получаем комнаты через клиент room service
	rooms, err := h.roomClient.GetAvailableRooms(ctx, params)
	if err != nil {
		return nil, err
	}

	// Проверяем доступность через сервис бронирований
	availableRooms, err := h.bookingService.GetAvailableRooms(ctx, params, rooms)
	if err != nil {
		return nil, err
	}

	// Конвертируем обратно в proto
	return &bookingpb.GetAvailableRoomsResponse{
		Rooms: mapper.RoomsToProto(availableRooms),
	}, nil
}
