package grpc

import (
	"context"
	"github.com/semho/hotel-booking/pkg/logger"
	pb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
	"github.com/semho/hotel-booking/room-service/internal/api/grpc/mapper"
	"github.com/semho/hotel-booking/room-service/internal/domain/port"
)

type RoomHandler struct {
	pb.UnimplementedRoomServiceServer
	roomService port.RoomService
}

func NewRoomHandler(roomService port.RoomService) *RoomHandler {
	return &RoomHandler{
		roomService: roomService,
	}
}

func (h *RoomHandler) GetAvailableRooms(
	ctx context.Context,
	req *pb.GetAvailableRoomsRequest,
) (*pb.GetAvailableRoomsResponse, error) {
	logger.Log.Info(
		"received GetAvailableRooms request",
		"capacity", req.Capacity,
		"type", req.Type,
		"status", req.Status,
	)

	params := mapper.ToSearchParams(req)
	logger.Log.Info(
		"mapped request to search params",
		"params", params,
	)

	rooms, err := h.roomService.GetAvailableRooms(ctx, params)
	if err != nil {
		logger.Log.Error(
			"failed to get available rooms",
			"error", err,
		)
		return nil, mapper.ToDomainError(err)
	}

	logger.Log.Info(
		"found rooms",
		"count", len(rooms),
	)

	response := &pb.GetAvailableRoomsResponse{
		Rooms: make([]*pb.Room, len(rooms)),
	}

	for i, room := range rooms {
		response.Rooms[i] = mapper.ToProtoRoom(room)
	}

	return response, nil
}
