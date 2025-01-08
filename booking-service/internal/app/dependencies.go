package app

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	grpcHandler "github.com/semho/hotel-booking/booking-service/internal/api/grpc"
	"github.com/semho/hotel-booking/booking-service/internal/config"
	"github.com/semho/hotel-booking/booking-service/internal/domain/port"
	"github.com/semho/hotel-booking/booking-service/internal/domain/service"
	"github.com/semho/hotel-booking/booking-service/internal/infrastructure/client/room"
	"github.com/semho/hotel-booking/booking-service/internal/infrastructure/repository/postgres"
	"github.com/semho/hotel-booking/booking-service/internal/infrastructure/unitofwork"
	roompb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Deps struct {
	DB             *sqlx.DB
	BookingHandler *grpcHandler.BookingHandler
	RoomClient     roompb.RoomServiceClient
	BookingUoW     port.BookingUnitOfWork
}

func initDeps(cfg *config.Config) (*Deps, error) {
	// Инициализируем БД
	db, err := initDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	// Инициализируем клиент room service
	roomConn, err := grpc.NewClient(
		cfg.RoomService.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to room service: %w", err)
	}
	roomClient := roompb.NewRoomServiceClient(roomConn)

	// Инициализируем слои
	bookingRepo := postgres.NewBookingRepository(db)
	roomClientWrapper := room.NewRoomClient(roomClient)
	bookingUoW := unitofwork.NewBookingUnitOfWork(db)

	bookingService := service.NewBookingService(bookingRepo, bookingUoW, roomClientWrapper)
	bookingHandler := grpcHandler.NewBookingHandler(bookingService, roomClientWrapper)

	return &Deps{
		DB:             db,
		BookingHandler: bookingHandler,
		RoomClient:     roomClient,
		BookingUoW:     bookingUoW,
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
