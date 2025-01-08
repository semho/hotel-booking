package grpc

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/semho/hotel-booking/pkg/logger"
	pb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
	"github.com/semho/hotel-booking/room-service/internal/api/grpc/mapper"
	"github.com/semho/hotel-booking/room-service/internal/domain/port"
	"github.com/shopspring/decimal"
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

func (h *RoomHandler) CreateRoom(ctx context.Context, req *pb.CreateRoomRequest) (*pb.CreateRoomResponse, error) {
	logger.Log.Info("create room request", "room number", req.RoomNumber)
	price, err := decimal.NewFromString(req.Price)
	if err != nil {
		logger.Log.Error(
			"failed to create rooms, invalid price format",
			"error", err,
		)
		return nil, fmt.Errorf("invalid price format: %w", err)
	}

	room := mapper.ProtoToRoom(req, price)

	err = h.roomService.Create(ctx, room)
	if err != nil {
		logger.Log.Error(
			"failed to create room in service",
			"error", err,
		)
		return nil, fmt.Errorf("failed to create room: %w", err)
	}

	// Маппинг модели комнаты обратно в формат протобаф
	protoRoom := mapper.ToProtoRoom(*room)

	logger.Log.Info("room successfully created", "room id", protoRoom.Id)
	return &pb.CreateRoomResponse{
		Room: protoRoom,
	}, nil
}

func (h *RoomHandler) GetAvailableRooms(
	ctx context.Context,
	req *pb.GetAvailableRoomsRequest,
) (*pb.GetAvailableRoomsResponse, error) {

	params := mapper.ToSearchParams(req)

	rooms, err := h.roomService.GetAvailableRooms(ctx, params)
	if err != nil {
		logger.Log.Error(
			"failed to get available rooms",
			"error", err,
		)
		return nil, mapper.ToDomainError(err)
	}

	response := &pb.GetAvailableRoomsResponse{
		Rooms: make([]*pb.Room, len(rooms)),
	}

	for i, room := range rooms {
		response.Rooms[i] = mapper.ToProtoRoom(room)
	}

	return response, nil
}

func (h *RoomHandler) GetRoomsCount(
	ctx context.Context,
	req *pb.GetAvailableRoomsRequest,
) (*pb.GetRoomsCountResponse, error) {
	params := mapper.ToSearchParams(req)

	count, err := h.roomService.GetRoomsCount(ctx, params)
	if err != nil {
		return nil, err
	}

	return &pb.GetRoomsCountResponse{
		Count: count,
	}, nil
}

func (h *RoomHandler) GetRoom(ctx context.Context, req *pb.GetRoomRequest) (*pb.GetRoomResponse, error) {
	roomID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid room_id: %w", err)
	}

	room, err := h.roomService.GetByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	return &pb.GetRoomResponse{
		Room: mapper.ToProtoRoom(*room),
	}, nil
}

func (h *RoomHandler) GetFirstAvailableRoom(ctx context.Context, req *pb.GetAvailableRoomsRequest) (
	*pb.GetRoomResponse,
	error,
) {
	room, err := h.roomService.GetFirstAvailableRoom(ctx, mapper.ToSearchParams(req))
	if err != nil {
		return nil, err
	}

	return &pb.GetRoomResponse{
		Room: mapper.ToProtoRoom(*room),
	}, nil
}
