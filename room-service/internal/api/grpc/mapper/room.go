package mapper

import (
	"github.com/semho/hotel-booking/pkg/errors"
	pb "github.com/semho/hotel-booking/pkg/proto/room_v1"
	"github.com/semho/hotel-booking/room-service/internal/domain/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToSearchParams(req *pb.GetAvailableRoomsRequest) model.SearchParams {
	params := model.SearchParams{
		CheckIn:  req.CheckIn.AsTime(),
		CheckOut: req.CheckOut.AsTime(),
	}

	if req.Capacity != nil {
		capacity := int(*req.Capacity)
		params.Capacity = &capacity
	}

	if req.Type != nil && *req.Type != pb.RoomType_ROOM_TYPE_UNSPECIFIED {
		roomType := ToModelRoomType(*req.Type)
		if roomType != "" {
			params.Type = &roomType
		}
	}

	return params
}

func ToModelRoomType(t pb.RoomType) model.RoomType {
	switch t {
	case pb.RoomType_ROOM_TYPE_STANDARD:
		return model.RoomTypeStandard
	case pb.RoomType_ROOM_TYPE_DELUXE:
		return model.RoomTypeDeluxe
	case pb.RoomType_ROOM_TYPE_SUITE:
		return model.RoomTypeSuite
	default:
		return ""
	}
}

func ToProtoRoomType(t model.RoomType) pb.RoomType {
	switch t {
	case model.RoomTypeStandard:
		return pb.RoomType_ROOM_TYPE_STANDARD
	case model.RoomTypeDeluxe:
		return pb.RoomType_ROOM_TYPE_DELUXE
	case model.RoomTypeSuite:
		return pb.RoomType_ROOM_TYPE_SUITE
	default:
		return pb.RoomType_ROOM_TYPE_UNSPECIFIED
	}
}

func ToProtoRoomStatus(s model.RoomStatus) pb.RoomStatus {
	switch s {
	case model.RoomStatusAvailable:
		return pb.RoomStatus_ROOM_STATUS_AVAILABLE
	case model.RoomStatusOccupied:
		return pb.RoomStatus_ROOM_STATUS_OCCUPIED
	case model.RoomStatusMaintenance:
		return pb.RoomStatus_ROOM_STATUS_MAINTENANCE
	default:
		return pb.RoomStatus_ROOM_STATUS_UNSPECIFIED
	}
}

func ToProtoRoom(r model.Room) *pb.Room {
	return &pb.Room{
		Id:         r.ID.String(),
		RoomNumber: r.RoomNumber,
		Type:       ToProtoRoomType(r.Type),
		Price:      r.Price,
		Capacity:   int32(r.Capacity),
		Status:     ToProtoRoomStatus(r.Status),
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
