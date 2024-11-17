package model

import (
	"github.com/google/uuid"
	"time"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "PENDING"
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusCancelled BookingStatus = "CANCELLED"
	BookingStatusCompleted BookingStatus = "COMPLETED"
	BookingStatusNoShow    BookingStatus = "NO_SHOW"
)

type Booking struct {
	ID         uuid.UUID     `db:"id" json:"id"`
	RoomID     uuid.UUID     `db:"room_id" json:"room_id"`
	UserID     *uuid.UUID    `db:"user_id" json:"user_id,omitempty"` // может быть nil для анонимных бронирований
	GuestName  string        `db:"guest_name" json:"guest_name"`
	GuestEmail string        `db:"guest_email" json:"guest_email"`
	GuestPhone string        `db:"guest_phone" json:"guest_phone"`
	CheckIn    time.Time     `db:"check_in" json:"check_in"`
	CheckOut   time.Time     `db:"check_out" json:"check_out"`
	Status     BookingStatus `db:"status" json:"status"`
	TotalPrice float64       `db:"total_price" json:"total_price"`
	CreatedAt  time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time     `db:"updated_at" json:"updated_at"`
}
