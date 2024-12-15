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

func ProtoToAvailableRooms(protoRooms []*roompb.Room) []response.AvailableRoom {
	rooms := make([]response.AvailableRoom, len(protoRooms))
	for i, pr := range protoRooms {
		rooms[i] = response.AvailableRoom{
			ID:         pr.Id,
			RoomNumber: pr.RoomNumber,
			Type:       pr.Type,
			Price:      pr.Price,
			Capacity:   int(pr.Capacity),
			Amenities:  pr.Amenities,
		}
	}
	return rooms
}
