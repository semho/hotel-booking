package request

import (
	"github.com/semho/hotel-booking/pkg/errors"
	roompb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
	"time"
)

type SearchParams struct {
	CheckIn  time.Time
	CheckOut time.Time
	Capacity *int32
	Type     *roompb.RoomType
}

func (p *SearchParams) Validate() error {
	if p.CheckIn.IsZero() {
		return errors.WithMessage(errors.ErrInvalidInput, "check-in date is required")
	}
	if p.CheckOut.IsZero() {
		return errors.WithMessage(errors.ErrInvalidInput, "check-out date is required")
	}
	if p.CheckOut.Before(p.CheckIn) {
		return errors.WithMessage(errors.ErrInvalidInput, "check-out date must be after check-in date")
	}
	if p.Capacity != nil && *p.Capacity <= 0 {
		return errors.WithMessage(errors.ErrInvalidInput, "capacity must be positive")
	}

	return nil
}

type CreateBooking struct {
	RoomID     string    `json:"roomId"`
	GuestName  string    `json:"guestName"`
	GuestEmail string    `json:"guestEmail"`
	GuestPhone string    `json:"guestPhone"`
	CheckIn    time.Time `json:"checkIn"`
	CheckOut   time.Time `json:"checkOut"`
}

func (req *CreateBooking) Validate() error {
	if req.RoomID == "" {
		return errors.WithMessage(errors.ErrInvalidInput, "room ID is required")
	}
	if req.GuestName == "" {
		return errors.WithMessage(errors.ErrInvalidInput, "guest name is required")
	}
	if req.GuestEmail == "" {
		return errors.WithMessage(errors.ErrInvalidInput, "guest email is required")
	}
	if req.CheckIn.IsZero() || req.CheckOut.IsZero() {
		return errors.WithMessage(errors.ErrInvalidInput, "check-in and check-out dates are required")
	}
	if req.CheckOut.Before(req.CheckIn) {
		return errors.WithMessage(errors.ErrInvalidInput, "check-out date must be after check-in date")
	}
	return nil
}

type CreateBookingRequest struct {
	RoomID     string `json:"roomId"`
	CheckIn    string `json:"checkIn"`
	CheckOut   string `json:"checkOut"`
	GuestName  string `json:"guestName"`
	GuestEmail string `json:"guestEmail"`
	GuestPhone string `json:"guestPhone"`
}
