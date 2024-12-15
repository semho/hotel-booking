package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	pb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
	"github.com/shopspring/decimal"
)

// Используем типы из proto напрямую для упрощения поддержки контрактов и предотвращения рассинхронизации между сервисами
type RoomType = pb.RoomType
type RoomStatus = pb.RoomStatus

type Room struct {
	ID         uuid.UUID       `db:"id" json:"id"`
	RoomNumber string          `db:"room_number" json:"room_number"`
	Type       RoomType        `db:"type" json:"type"`
	Price      decimal.Decimal `db:"price" json:"price"`
	Capacity   int             `db:"capacity" json:"capacity"`
	Status     RoomStatus      `db:"status" json:"status"`
	Amenities  pq.StringArray  `db:"amenities" json:"amenities"`
	CreatedAt  time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at" json:"updated_at"`
}

type SearchParams struct {
	Capacity *int
	Type     *RoomType
	Status   *RoomStatus
}
