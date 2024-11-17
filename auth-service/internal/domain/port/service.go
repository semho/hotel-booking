package port

import (
	"context"
	"github.com/semho/hotel-booking/auth-service/internal/domain/model"
)

type AuthService interface {
	Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error)
	Login(ctx context.Context, email, password string) (*model.AuthResponse, error)
	ValidateAccessToken(ctx context.Context, token string) (*model.User, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*model.AuthResponse, error)
}
