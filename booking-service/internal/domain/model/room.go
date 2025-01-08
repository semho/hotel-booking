package model

import (
	pb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
)

// Используем типы из proto напрямую
type RoomType = pb.RoomType
type RoomStatus = pb.RoomStatus

type Room struct {
	ID        string
	Number    string
	Type      RoomType
	Price     string
	Capacity  int
	Status    RoomStatus
	Amenities []string
}

type SearchRoomsParams struct {
	Capacity *int32
	Type     *RoomType
	Status   *RoomStatus
}
