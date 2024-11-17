package model

import (
	"time"
)

type RoomType string
type RoomStatus string

const (
	RoomTypeStandard RoomType = "STANDARD"
	RoomTypeDeluxe   RoomType = "DELUXE"
	RoomTypeSuite    RoomType = "SUITE"
)

const (
	RoomStatusAvailable    RoomStatus = "AVAILABLE"
	RoomStatusRepair       RoomStatus = "REPAIR"
	RoomStatusMaintenance  RoomStatus = "MAINTENANCE"
	RoomStatusOutOfService RoomStatus = "OUT_OF_SERVICE"
)

type Room struct {
	ID       string
	Number   string
	Type     RoomType
	Price    float64
	Capacity int
	Status   RoomStatus
}

type SearchParams struct {
	CheckIn  time.Time
	CheckOut time.Time
	Capacity *int32
	Type     *RoomType
}
