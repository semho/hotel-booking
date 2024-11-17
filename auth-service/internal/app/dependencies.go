package app

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	grpcHandler "github.com/semho/hotel-booking/auth-service/internal/api/grpc"
	"github.com/semho/hotel-booking/auth-service/internal/config"
	"github.com/semho/hotel-booking/auth-service/internal/domain/service"
	"github.com/semho/hotel-booking/auth-service/internal/infrastructure/repository/postgres"
	"github.com/semho/hotel-booking/pkg/auth/jwt"
	"time"
)

type Deps struct {
	db          *sqlx.DB
	AuthHandler *grpcHandler.AuthHandler
}

func initDeps(cfg *config.Config) (*Deps, error) {
	// Инициализируем БД
	db, err := initDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	// Инициализируем JWT manager
	tokenManager := jwt.NewTokenManager(
		cfg.JWT.AccessTokenSecret,
		cfg.JWT.RefreshTokenSecret,
		time.Duration(cfg.JWT.AccessTokenTTL)*time.Minute,
		time.Duration(cfg.JWT.RefreshTokenTTL)*time.Hour*24,
	)

	// Инициализируем слои
	userRepo := postgres.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, tokenManager)
	authHandler := grpcHandler.NewAuthHandler(authService)

	return &Deps{
		db:          db,
		AuthHandler: authHandler,
	}, nil
}

func initDB(cfg *config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
