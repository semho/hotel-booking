package mapper

import (
	"github.com/semho/hotel-booking/booking-service/internal/domain/model"
	roompb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
)

// Маппинг типов комнат
func roomTypeToProto(t model.RoomType) roompb.RoomType {
	switch t {
	case model.RoomTypeStandard:
		return roompb.RoomType_ROOM_TYPE_STANDARD
	case model.RoomTypeDeluxe:
		return roompb.RoomType_ROOM_TYPE_DELUXE
	case model.RoomTypeSuite:
		return roompb.RoomType_ROOM_TYPE_SUITE
	default:
		return roompb.RoomType_ROOM_TYPE_UNSPECIFIED
	}
}

func protoToRoomType(t roompb.RoomType) model.RoomType {
	switch t {
	case roompb.RoomType_ROOM_TYPE_STANDARD:
		return model.RoomTypeStandard
	case roompb.RoomType_ROOM_TYPE_DELUXE:
		return model.RoomTypeDeluxe
	case roompb.RoomType_ROOM_TYPE_SUITE:
		return model.RoomTypeSuite
	default:
		return model.RoomTypeStandard // или обработать ошибку
	}
}

// Маппинг статусов
func roomStatusToProto(s model.RoomStatus) roompb.RoomStatus {
	switch s {
	case model.RoomStatusAvailable:
		return roompb.RoomStatus_ROOM_STATUS_AVAILABLE
	case model.RoomStatusRepair:
		return roompb.RoomStatus_ROOM_STATUS_REPAIR
	case model.RoomStatusMaintenance:
		return roompb.RoomStatus_ROOM_STATUS_MAINTENANCE
	case model.RoomStatusOutOfService:
		return roompb.RoomStatus_ROOM_STATUS_OUT_OF_SERVICE
	default:
		return roompb.RoomStatus_ROOM_STATUS_UNSPECIFIED
	}
}

func protoToRoomStatus(s roompb.RoomStatus) model.RoomStatus {
	switch s {
	case roompb.RoomStatus_ROOM_STATUS_AVAILABLE:
		return model.RoomStatusAvailable
	case roompb.RoomStatus_ROOM_STATUS_REPAIR:
		return model.RoomStatusRepair
	case roompb.RoomStatus_ROOM_STATUS_MAINTENANCE:
		return model.RoomStatusMaintenance
	case roompb.RoomStatus_ROOM_STATUS_OUT_OF_SERVICE:
		return model.RoomStatusOutOfService
	default:
		return model.RoomStatusAvailable // или обработать ошибку
	}
}

func SearchParamsToProto(params model.SearchParams) *roompb.GetAvailableRoomsRequest {
	var (
		roomType *roompb.RoomType
		status   *roompb.RoomStatus
	)

	if params.Type != nil {
		rt := roomTypeToProto(*params.Type)
		roomType = &rt
	}

	// По умолчанию ищем только доступные комнаты
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
			Type:     protoToRoomType(pr.Type),
			Price:    pr.Price,
			Capacity: int(pr.Capacity),
			Status:   protoToRoomStatus(pr.Status),
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
			Type:       roomTypeToProto(r.Type),
			Price:      r.Price,
			Capacity:   int32(r.Capacity),
			Status:     roomStatusToProto(r.Status),
		}
	}
	return protoRooms
}
