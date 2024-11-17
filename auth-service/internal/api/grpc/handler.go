package grpc

import (
	"context"
	"github.com/semho/hotel-booking/auth-service/internal/api/grpc/mapper"
	"github.com/semho/hotel-booking/auth-service/internal/domain/model"
	"github.com/semho/hotel-booking/auth-service/internal/domain/port"
	"github.com/semho/hotel-booking/pkg/logger"
	pb "github.com/semho/hotel-booking/pkg/proto/auth_v1/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService port.AuthService
}

func NewAuthHandler(authService port.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	logger.Log.Info(
		"auth Register",
		"user_email", req.Email,
	)
	domainReq := &model.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
	if req.Phone != nil {
		domainReq.Phone = req.Phone
	}

	response, err := h.authService.Register(ctx, domainReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return mapper.ToProtoAuthResponse(response), nil
}

func (h *AuthHandler) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	logger.Log.Info(
		"auth Validate",
		"AccessToken", req.AccessToken,
	)

	// Валидируем токен через сервис
	user, err := h.authService.ValidateAccessToken(ctx, req.AccessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	return &pb.ValidateResponse{
		Valid: true,
		User:  mapper.ToProtoUserInfo(user),
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	logger.Log.Info(
		"auth Login",
		"user_email", req.Email,
	)
	domainReq := &model.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	response, err := h.authService.Login(ctx, domainReq.Email, domainReq.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to login user: %v", err)
	}

	return mapper.ToProtoAuthResponse(response), nil
}
