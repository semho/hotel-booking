package model

import (
	"time"

	pb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
)

// Используем типы из proto напрямую
type RoomType = pb.RoomType
type RoomStatus = pb.RoomStatus

type Room struct {
	ID       string
	Number   string
	Type     RoomType
	Price    string
	Capacity int
	Status   RoomStatus
}

type SearchParams struct {
	CheckIn  time.Time
	CheckOut time.Time
	Capacity *int32
	Type     *RoomType
}
