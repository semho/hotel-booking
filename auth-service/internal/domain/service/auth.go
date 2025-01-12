package service

import (
	"context"

	"github.com/semho/hotel-booking/auth-service/internal/domain/model"
	"github.com/semho/hotel-booking/auth-service/internal/domain/port"
	"github.com/semho/hotel-booking/pkg/auth/jwt"
	"github.com/semho/hotel-booking/pkg/errors"
	"github.com/semho/hotel-booking/pkg/logger"
	pb "github.com/semho/hotel-booking/pkg/proto/auth_v1/auth"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo     port.UserRepository
	tokenManager *jwt.TokenManager
}

func NewAuthService(userRepo port.UserRepository, tokenManager *jwt.TokenManager) port.AuthService {
	return &authService{
		userRepo:     userRepo,
		tokenManager: tokenManager,
	}
}

func (s *authService) Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error) {
	// Проверяем, не существует ли уже пользователь
	if _, err := s.userRepo.GetByEmail(ctx, req.Email); err == nil {
		return nil, errors.WithMessage(errors.ErrConflict, "user already exists")
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создаем пользователя
	user := &model.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      pb.UserRole_USER_ROLE_USER,
	}

	if err = s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Создаем токены
	accessToken, accessExp, err := s.tokenManager.CreateAccessToken(
		user.ID,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExp, err := s.tokenManager.CreateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessExp,
		RefreshTokenExpiresAt: refreshExp,
		User:                  user,
	}, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*model.AuthResponse, error) {
	// Получаем пользователя
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.WithMessage(errors.ErrInvalidInput, "invalid email or password")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	); err != nil {
		return nil, errors.WithMessage(errors.ErrInvalidInput, "invalid email or password")
	}

	// Создаем токены
	accessToken, accessExp, err := s.tokenManager.CreateAccessToken(
		user.ID,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExp, err := s.tokenManager.CreateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	logger.Log.Info(
		"auth service: tokens created",
		"access_token", accessToken,
		"refresh_token", refreshToken,
	)

	return &model.AuthResponse{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessExp,
		RefreshTokenExpiresAt: refreshExp,
		User:                  user,
	}, nil
}

func (s *authService) ValidateAccessToken(ctx context.Context, token string) (*model.User, error) {
	// Валидируем токен
	claims, err := s.tokenManager.ValidateAccessToken(token)
	if err != nil {
		return nil, errors.WithMessage(errors.ErrUnauthorized, "invalid token")
	}

	// Получаем пользователя
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.WithMessage(errors.ErrUnauthorized, "user not found")
	}

	return user, nil
}

func (s *authService) RefreshTokens(ctx context.Context, refreshToken string) (*model.AuthResponse, error) {
	// Валидируем refresh token
	userID, err := s.tokenManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.WithMessage(errors.ErrUnauthorized, "invalid refresh token")
	}

	// Получаем пользователя
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.WithMessage(errors.ErrUnauthorized, "user not found")
	}

	// Создаем новые токены
	accessToken, accessExp, err := s.tokenManager.CreateAccessToken(
		user.ID,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	newRefreshToken, refreshExp, err := s.tokenManager.CreateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		AccessToken:           accessToken,
		RefreshToken:          newRefreshToken,
		AccessTokenExpiresAt:  accessExp,
		RefreshTokenExpiresAt: refreshExp,
		User:                  user,
	}, nil
}
