package port

import (
	"context"
	"github.com/google/uuid"
	"github.com/semho/hotel-booking/auth-service/internal/domain/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
