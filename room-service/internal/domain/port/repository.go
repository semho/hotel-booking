package port

import (
	"context"
	"github.com/google/uuid"
	"github.com/semho/hotel-booking/room-service/internal/domain/model"
)

type RoomRepository interface {
	GetAvailableRooms(ctx context.Context, params model.SearchParams) ([]model.Room, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Room, error)
	Create(ctx context.Context, room *model.Room) error
	Update(ctx context.Context, room *model.Room) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetRoomsCount(ctx context.Context, params model.SearchParams) (int32, error)
	GetFirstAvailableRoom(ctx context.Context, params model.SearchParams) (*model.Room, error)
}
