package model

import (
	"time"

	"github.com/google/uuid"
	pb "github.com/semho/hotel-booking/pkg/proto/booking_v1/booking"
)

// Используем тип из proto напрямую
type BookingStatus = pb.BookingStatus

// Вспомогательная структура для объединения брони с текущим статусом
type BookingWithStatus struct {
	Booking
	CurrentStatus BookingStatusHistory `db:"status" json:"current_status"`
}

type Booking struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	RoomID     uuid.UUID  `db:"room_id" json:"room_id"`
	UserID     *uuid.UUID `db:"user_id" json:"user_id,omitempty"` // может быть nil для анонимных бронирований
	GuestName  string     `db:"guest_name" json:"guest_name"`
	GuestEmail string     `db:"guest_email" json:"guest_email"`
	GuestPhone string     `db:"guest_phone" json:"guest_phone"`
	CheckIn    time.Time  `db:"check_in" json:"check_in"`
	CheckOut   time.Time  `db:"check_out" json:"check_out"`
	TotalPrice float64    `db:"total_price" json:"total_price"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`

	// Добавляем поле для текущего статуса, которое не хранится в БД
	CurrentStatus *BookingStatusHistory `db:"-" json:"current_status,omitempty"`
}

type BookingStatusHistory struct {
	ID        uuid.UUID     `db:"id" json:"id"`
	BookingID uuid.UUID     `db:"booking_id" json:"booking_id"`
	Status    BookingStatus `db:"status" json:"status"`
	Reason    string        `db:"reason" json:"reason"`
	ChangedBy string        `db:"changed_by" json:"changed_by"`
	ChangedAt time.Time     `db:"changed_at" json:"changed_at"`
}

type SearchParams struct {
	CheckIn  time.Time
	CheckOut time.Time
	Capacity *int32
	Type     *RoomType
}

type BookingRow struct {
	Booking
	StatusID        uuid.UUID     `db:"status_id"`
	StatusStatus    BookingStatus `db:"status_status"`
	StatusReason    string        `db:"status_reason"`
	StatusChangedBy string        `db:"status_changed_by"`
	StatusChangedAt time.Time     `db:"status_changed_at"`
}
