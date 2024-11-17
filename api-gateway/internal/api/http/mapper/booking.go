package mapper

import (
	"github.com/semho/hotel-booking/api-gateway/internal/api/http/response"
	roompb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func TimeToProtoTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func ProtoRoomTypeToString(t roompb.RoomType) string {
	switch t {
	case roompb.RoomType_ROOM_TYPE_STANDARD:
		return "STANDARD"
	case roompb.RoomType_ROOM_TYPE_DELUXE:
		return "DELUXE"
	case roompb.RoomType_ROOM_TYPE_SUITE:
		return "SUITE"
	default:
		return "UNKNOWN"
	}
}

func StringToProtoRoomType(s string) roompb.RoomType {
	switch s {
	case "STANDARD":
		return roompb.RoomType_ROOM_TYPE_STANDARD
	case "DELUXE":
		return roompb.RoomType_ROOM_TYPE_DELUXE
	case "SUITE":
		return roompb.RoomType_ROOM_TYPE_SUITE
	default:
		return roompb.RoomType_ROOM_TYPE_UNSPECIFIED
	}
}

func ProtoToAvailableRooms(protoRooms []*roompb.Room) []response.AvailableRoom {
	rooms := make([]response.AvailableRoom, len(protoRooms))
	for i, pr := range protoRooms {
		rooms[i] = response.AvailableRoom{
			ID:         pr.Id,
			RoomNumber: pr.RoomNumber,
			Type:       ProtoRoomTypeToString(pr.Type),
			Price:      pr.Price,
			Capacity:   int(pr.Capacity),
			Amenities:  pr.Amenities,
		}
	}
	return rooms
}
