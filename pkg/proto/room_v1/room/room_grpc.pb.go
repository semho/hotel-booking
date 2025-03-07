// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.28.3
// source: room/room.proto

package room

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// RoomServiceClient is the client API for RoomService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RoomServiceClient interface {
	// GetAvailableRooms returns list of available rooms based on search criteria
	GetAvailableRooms(ctx context.Context, in *GetAvailableRoomsRequest, opts ...grpc.CallOption) (*GetAvailableRoomsResponse, error)
	CreateRoom(ctx context.Context, in *CreateRoomRequest, opts ...grpc.CallOption) (*CreateRoomResponse, error)
	GetRoomsCount(ctx context.Context, in *GetAvailableRoomsRequest, opts ...grpc.CallOption) (*GetRoomsCountResponse, error)
	GetRoom(ctx context.Context, in *GetRoomRequest, opts ...grpc.CallOption) (*GetRoomResponse, error)
	GetFirstAvailableRoom(ctx context.Context, in *GetAvailableRoomsRequest, opts ...grpc.CallOption) (*GetRoomResponse, error)
}

type roomServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRoomServiceClient(cc grpc.ClientConnInterface) RoomServiceClient {
	return &roomServiceClient{cc}
}

func (c *roomServiceClient) GetAvailableRooms(ctx context.Context, in *GetAvailableRoomsRequest, opts ...grpc.CallOption) (*GetAvailableRoomsResponse, error) {
	out := new(GetAvailableRoomsResponse)
	err := c.cc.Invoke(ctx, "/hotel.room.v1.RoomService/GetAvailableRooms", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomServiceClient) CreateRoom(ctx context.Context, in *CreateRoomRequest, opts ...grpc.CallOption) (*CreateRoomResponse, error) {
	out := new(CreateRoomResponse)
	err := c.cc.Invoke(ctx, "/hotel.room.v1.RoomService/CreateRoom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomServiceClient) GetRoomsCount(ctx context.Context, in *GetAvailableRoomsRequest, opts ...grpc.CallOption) (*GetRoomsCountResponse, error) {
	out := new(GetRoomsCountResponse)
	err := c.cc.Invoke(ctx, "/hotel.room.v1.RoomService/GetRoomsCount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomServiceClient) GetRoom(ctx context.Context, in *GetRoomRequest, opts ...grpc.CallOption) (*GetRoomResponse, error) {
	out := new(GetRoomResponse)
	err := c.cc.Invoke(ctx, "/hotel.room.v1.RoomService/GetRoom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomServiceClient) GetFirstAvailableRoom(ctx context.Context, in *GetAvailableRoomsRequest, opts ...grpc.CallOption) (*GetRoomResponse, error) {
	out := new(GetRoomResponse)
	err := c.cc.Invoke(ctx, "/hotel.room.v1.RoomService/GetFirstAvailableRoom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RoomServiceServer is the server API for RoomService service.
// All implementations must embed UnimplementedRoomServiceServer
// for forward compatibility
type RoomServiceServer interface {
	// GetAvailableRooms returns list of available rooms based on search criteria
	GetAvailableRooms(context.Context, *GetAvailableRoomsRequest) (*GetAvailableRoomsResponse, error)
	CreateRoom(context.Context, *CreateRoomRequest) (*CreateRoomResponse, error)
	GetRoomsCount(context.Context, *GetAvailableRoomsRequest) (*GetRoomsCountResponse, error)
	GetRoom(context.Context, *GetRoomRequest) (*GetRoomResponse, error)
	GetFirstAvailableRoom(context.Context, *GetAvailableRoomsRequest) (*GetRoomResponse, error)
	mustEmbedUnimplementedRoomServiceServer()
}

// UnimplementedRoomServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRoomServiceServer struct {
}

func (UnimplementedRoomServiceServer) GetAvailableRooms(context.Context, *GetAvailableRoomsRequest) (*GetAvailableRoomsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAvailableRooms not implemented")
}
func (UnimplementedRoomServiceServer) CreateRoom(context.Context, *CreateRoomRequest) (*CreateRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRoom not implemented")
}
func (UnimplementedRoomServiceServer) GetRoomsCount(context.Context, *GetAvailableRoomsRequest) (*GetRoomsCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoomsCount not implemented")
}
func (UnimplementedRoomServiceServer) GetRoom(context.Context, *GetRoomRequest) (*GetRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoom not implemented")
}
func (UnimplementedRoomServiceServer) GetFirstAvailableRoom(context.Context, *GetAvailableRoomsRequest) (*GetRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFirstAvailableRoom not implemented")
}
func (UnimplementedRoomServiceServer) mustEmbedUnimplementedRoomServiceServer() {}

// UnsafeRoomServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RoomServiceServer will
// result in compilation errors.
type UnsafeRoomServiceServer interface {
	mustEmbedUnimplementedRoomServiceServer()
}

func RegisterRoomServiceServer(s grpc.ServiceRegistrar, srv RoomServiceServer) {
	s.RegisterService(&RoomService_ServiceDesc, srv)
}

func _RoomService_GetAvailableRooms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAvailableRoomsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServiceServer).GetAvailableRooms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hotel.room.v1.RoomService/GetAvailableRooms",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServiceServer).GetAvailableRooms(ctx, req.(*GetAvailableRoomsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomService_CreateRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServiceServer).CreateRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hotel.room.v1.RoomService/CreateRoom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServiceServer).CreateRoom(ctx, req.(*CreateRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomService_GetRoomsCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAvailableRoomsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServiceServer).GetRoomsCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hotel.room.v1.RoomService/GetRoomsCount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServiceServer).GetRoomsCount(ctx, req.(*GetAvailableRoomsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomService_GetRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServiceServer).GetRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hotel.room.v1.RoomService/GetRoom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServiceServer).GetRoom(ctx, req.(*GetRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomService_GetFirstAvailableRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAvailableRoomsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServiceServer).GetFirstAvailableRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hotel.room.v1.RoomService/GetFirstAvailableRoom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServiceServer).GetFirstAvailableRoom(ctx, req.(*GetAvailableRoomsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RoomService_ServiceDesc is the grpc.ServiceDesc for RoomService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RoomService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "hotel.room.v1.RoomService",
	HandlerType: (*RoomServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAvailableRooms",
			Handler:    _RoomService_GetAvailableRooms_Handler,
		},
		{
			MethodName: "CreateRoom",
			Handler:    _RoomService_CreateRoom_Handler,
		},
		{
			MethodName: "GetRoomsCount",
			Handler:    _RoomService_GetRoomsCount_Handler,
		},
		{
			MethodName: "GetRoom",
			Handler:    _RoomService_GetRoom_Handler,
		},
		{
			MethodName: "GetFirstAvailableRoom",
			Handler:    _RoomService_GetFirstAvailableRoom_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "room/room.proto",
}
