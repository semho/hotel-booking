package port

import (
	"context"
	pb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
)

type RoomAPI interface {
	// получаем список доступных комнат
	GetAvailableRooms(ctx context.Context, req *pb.GetAvailableRoomsRequest) (*pb.GetAvailableRoomsResponse, error)
	CreateRoom(ctx context.Context, req *pb.CreateRoomRequest) (*pb.Room, error)
	GetRoomsCount(ctx context.Context, req *pb.GetAvailableRoomsRequest) (*pb.GetRoomsCountResponse, error)
	GetRoom(ctx context.Context, req *pb.GetRoomRequest) (*pb.CreateRoomResponse, error)
	GetFirstAvailableRoom(ctx context.Context, req *pb.GetRoomRequest) (*pb.CreateRoomResponse, error) //TODO: удалить
	// UpdateRoom(ctx context.Context, req *pb.UpdateRoomRequest) (*pb.Room, error)
	// DeleteRoom(ctx context.Context, req *pb.DeleteRoomRequest) (*pb.Empty, error)
}
