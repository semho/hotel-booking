package port

import (
	"context"
	pb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
)

type RoomAPI interface {
	// получаем список доступных комнат
	GetAvailableRooms(ctx context.Context, req *pb.GetAvailableRoomsRequest) (*pb.GetAvailableRoomsResponse, error)

	// другие методы:
	// CreateRoom(ctx context.Context, req *pb.CreateRoomRequest) (*pb.Room, error)
	// UpdateRoom(ctx context.Context, req *pb.UpdateRoomRequest) (*pb.Room, error)
	// DeleteRoom(ctx context.Context, req *pb.DeleteRoomRequest) (*pb.Empty, error)
}
