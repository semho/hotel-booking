package room

import (
	"context"
	"github.com/semho/hotel-booking/booking-service/internal/api/grpc/mapper"
	"github.com/semho/hotel-booking/booking-service/internal/domain/model"
	"github.com/semho/hotel-booking/booking-service/internal/domain/port"
	roompb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
)

type roomClient struct {
	client roompb.RoomServiceClient
}

func NewRoomClient(client roompb.RoomServiceClient) port.RoomClient {
	return &roomClient{
		client: client,
	}
}

func (c *roomClient) GetAvailableRooms(ctx context.Context, params model.SearchParams) ([]model.Room, error) {
	// Конвертируем параметры в proto
	req := mapper.SearchParamsToProto(params)

	// Делаем запрос к room service
	resp, err := c.client.GetAvailableRooms(ctx, req)
	if err != nil {
		return nil, err
	}

	// Конвертируем ответ в доменную модель
	return mapper.ProtoToRooms(resp.Rooms), nil
}
