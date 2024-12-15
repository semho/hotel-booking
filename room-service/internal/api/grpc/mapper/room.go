package mapper

import (
	"time"

	"github.com/semho/hotel-booking/pkg/errors"
	pb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
	"github.com/semho/hotel-booking/room-service/internal/domain/model"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToSearchParams(req *pb.GetAvailableRoomsRequest) model.SearchParams {
	params := model.SearchParams{}

	if req.Capacity != nil {
		capacity := int(*req.Capacity)
		params.Capacity = &capacity
	}

	if req.Type != nil {
		roomType := *req.Type
		params.Type = &roomType
	}

	if req.Status != nil {
		roomStatus := *req.Status
		params.Status = &roomStatus
	}

	return params
}

func ToProtoRoom(r model.Room) *pb.Room {
	return &pb.Room{
		Id:         r.ID.String(),
		RoomNumber: r.RoomNumber,
		Type:       r.Type,
		Price:      r.Price.String(),
		Capacity:   int32(r.Capacity),
		Status:     r.Status,
		Amenities:  r.Amenities,
	}
}

func ToDomainError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.IsNotFound(err):
		return status.Error(codes.NotFound, err.Error())
	case errors.IsInvalidInput(err):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.IsConflict(err):
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}

func ProtoToRoom(protoRoom *pb.CreateRoomRequest, price decimal.Decimal) *model.Room {
	return &model.Room{
		RoomNumber: protoRoom.RoomNumber,
		Type:       protoRoom.Type,
		Price:      price,
		Capacity:   int(protoRoom.Capacity),
		Status:     protoRoom.Status,
		Amenities:  protoRoom.Amenities,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
