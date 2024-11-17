package port

import (
	"context"
	pb "github.com/semho/hotel-booking/pkg/proto/auth_v1/auth"
)

type AuthAPI interface {
	Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error)
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error)
	Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error)
	Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.AuthResponse, error)
}
