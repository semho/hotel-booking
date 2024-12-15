package mapper

import (
	"github.com/semho/hotel-booking/booking-service/internal/domain/model"
	roompb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
)

func SearchParamsToProto(params model.SearchParams) *roompb.GetAvailableRoomsRequest {
	var (
		roomType *roompb.RoomType
		status   *roompb.RoomStatus
	)

	if params.Type != nil {
		rt := *params.Type
		roomType = &rt
	}

	// По умолчанию ищем только доступные комнаты, т.к. запрос именно на свободные комнаты и фильтр по статусу лишний
	availableStatus := roompb.RoomStatus_ROOM_STATUS_AVAILABLE
	status = &availableStatus

	return &roompb.GetAvailableRoomsRequest{
		Capacity: params.Capacity,
		Type:     roomType,
		Status:   status,
	}
}

func ProtoToRooms(protoRooms []*roompb.Room) []model.Room {
	rooms := make([]model.Room, len(protoRooms))
	for i, pr := range protoRooms {
		rooms[i] = model.Room{
			ID:       pr.Id,
			Number:   pr.RoomNumber,
			Type:     pr.Type,
			Price:    pr.Price,
			Capacity: int(pr.Capacity),
			Status:   pr.Status,
		}
	}
	return rooms
}

func RoomsToProto(rooms []model.Room) []*roompb.Room {
	protoRooms := make([]*roompb.Room, len(rooms))
	for i, r := range rooms {
		protoRooms[i] = &roompb.Room{
			Id:         r.ID,
			RoomNumber: r.Number,
			Type:       r.Type,
			Price:      r.Price,
			Capacity:   int32(r.Capacity),
			Status:     r.Status,
		}
	}
	return protoRooms
}
