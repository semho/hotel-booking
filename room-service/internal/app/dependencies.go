package app

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/semho/hotel-booking/room-service/internal/api/grpc"
	"github.com/semho/hotel-booking/room-service/internal/config"
	"github.com/semho/hotel-booking/room-service/internal/domain/service"
	"github.com/semho/hotel-booking/room-service/internal/infrastructure/repository/postgres"
)

type Deps struct {
	DB          *sqlx.DB
	RoomHandler *grpc.RoomHandler
}

func initDeps(cfg *config.Config) (*Deps, error) {
	// Инициализируем БД
	db, err := initDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	// Инициализируем слои
	roomRepo := postgres.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo)
	roomHandler := grpc.NewRoomHandler(roomService)

	return &Deps{
		DB:          db,
		RoomHandler: roomHandler,
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
