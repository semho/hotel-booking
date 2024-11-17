package port

import (
	"context"
	"github.com/semho/hotel-booking/booking-service/internal/domain/model"
)

type BookingService interface {
	GetAvailableRooms(ctx context.Context, params model.SearchParams, rooms []model.Room) ([]model.Room, error)
}
