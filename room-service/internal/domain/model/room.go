package model

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
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
	RoomStatusAvailable   RoomStatus = "AVAILABLE"
	RoomStatusOccupied    RoomStatus = "OCCUPIED"
	RoomStatusMaintenance RoomStatus = "MAINTENANCE"
)

type Room struct {
	ID         uuid.UUID      `db:"id" json:"id"`
	RoomNumber string         `db:"room_number" json:"room_number"`
	Type       RoomType       `db:"type" json:"type"`
	Price      float64        `db:"price" json:"price"`
	Capacity   int            `db:"capacity" json:"capacity"`
	Status     RoomStatus     `db:"status" json:"status"`
	Amenities  pq.StringArray `db:"amenities" json:"amenities"`
	CreatedAt  time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time      `db:"updated_at" json:"updated_at"`
}

type SearchParams struct {
	CheckIn  time.Time
	CheckOut time.Time
	Capacity *int
	Type     *RoomType
}
