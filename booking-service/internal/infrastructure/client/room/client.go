package room

import (
	"context"
	"fmt"
	"github.com/google/uuid"
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

func (c *roomClient) GetAvailableRooms(ctx context.Context, params model.SearchRoomsParams) ([]model.Room, error) {
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

func (c *roomClient) GetRoomsCount(ctx context.Context, params model.SearchRoomsParams) (int32, error) {
	req := mapper.SearchParamsToProto(params)
	resp, err := c.client.GetRoomsCount(ctx, req)
	if err != nil {
		return 0, err
	}
	return resp.Count, nil
}

func (c *roomClient) GetRoomInfo(ctx context.Context, roomID uuid.UUID) (*model.Room, error) {
	req := &roompb.GetRoomRequest{
		Id: roomID.String(),
	}
	resp, err := c.client.GetRoom(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get room info: %w", err)
	}

	return mapper.ProtoToRoom(resp.Room), nil

}

func (c *roomClient) GetFirstAvailableRoom(ctx context.Context, params model.SearchRoomsParams) (*model.Room, error) {
	status := roompb.RoomStatus_ROOM_STATUS_AVAILABLE
	req := &roompb.GetAvailableRoomsRequest{
		Capacity: params.Capacity,
		Type:     params.Type,
		Status:   &status,
	}

	resp, err := c.client.GetFirstAvailableRoom(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error in getting a free room: %w", err)
	}

	return mapper.ProtoToRoom(resp.Room), nil
}
