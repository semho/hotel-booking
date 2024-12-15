package mapper

import (
	"github.com/semho/hotel-booking/api-gateway/internal/api/http/request"
	"github.com/semho/hotel-booking/api-gateway/internal/api/http/response"
	roompb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
)

func ToHTTPRoom(protoRoom *roompb.Room) response.CreateRoomResponse {
	return response.CreateRoomResponse{
		ID:        protoRoom.Id,
		Number:    protoRoom.RoomNumber,
		Type:      protoRoom.Type,
		Price:     protoRoom.Price,
		Capacity:  int(protoRoom.Capacity),
		Status:    protoRoom.Status,
		Amenities: protoRoom.Amenities,
	}
}

func HttpToProto(req request.CreateRoomRequest) *roompb.CreateRoomRequest {
	return &roompb.CreateRoomRequest{
		RoomNumber: req.RoomNumber,
		Type:       req.Type,
		Price:      req.Price,
		Capacity:   int32(req.Capacity),
		Status:     req.Status,
		Amenities:  req.Amenities,
	}
}
