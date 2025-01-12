// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.28.3
// source: booking/booking.proto

package booking

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

// BookingServiceClient is the client API for BookingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BookingServiceClient interface {
	GetAvailableRooms(ctx context.Context, in *GetAvailableRoomsRequest, opts ...grpc.CallOption) (*GetAvailableRoomsResponse, error)
	CreateBooking(ctx context.Context, in *CreateBookingRequest, opts ...grpc.CallOption) (*CreateBookingResponse, error)
	// UpdateBookingStatus updates booking status
	UpdateBookingStatus(ctx context.Context, in *UpdateBookingStatusRequest, opts ...grpc.CallOption) (*UpdateBookingStatusResponse, error)
}

type bookingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBookingServiceClient(cc grpc.ClientConnInterface) BookingServiceClient {
	return &bookingServiceClient{cc}
}

func (c *bookingServiceClient) GetAvailableRooms(ctx context.Context, in *GetAvailableRoomsRequest, opts ...grpc.CallOption) (*GetAvailableRoomsResponse, error) {
	out := new(GetAvailableRoomsResponse)
	err := c.cc.Invoke(ctx, "/hotel.booking.v1.BookingService/GetAvailableRooms", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookingServiceClient) CreateBooking(ctx context.Context, in *CreateBookingRequest, opts ...grpc.CallOption) (*CreateBookingResponse, error) {
	out := new(CreateBookingResponse)
	err := c.cc.Invoke(ctx, "/hotel.booking.v1.BookingService/CreateBooking", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookingServiceClient) UpdateBookingStatus(ctx context.Context, in *UpdateBookingStatusRequest, opts ...grpc.CallOption) (*UpdateBookingStatusResponse, error) {
	out := new(UpdateBookingStatusResponse)
	err := c.cc.Invoke(ctx, "/hotel.booking.v1.BookingService/UpdateBookingStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BookingServiceServer is the server API for BookingService service.
// All implementations must embed UnimplementedBookingServiceServer
// for forward compatibility
type BookingServiceServer interface {
	GetAvailableRooms(context.Context, *GetAvailableRoomsRequest) (*GetAvailableRoomsResponse, error)
	CreateBooking(context.Context, *CreateBookingRequest) (*CreateBookingResponse, error)
	// UpdateBookingStatus updates booking status
	UpdateBookingStatus(context.Context, *UpdateBookingStatusRequest) (*UpdateBookingStatusResponse, error)
	mustEmbedUnimplementedBookingServiceServer()
}

// UnimplementedBookingServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBookingServiceServer struct {
}

func (UnimplementedBookingServiceServer) GetAvailableRooms(context.Context, *GetAvailableRoomsRequest) (*GetAvailableRoomsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAvailableRooms not implemented")
}
func (UnimplementedBookingServiceServer) CreateBooking(context.Context, *CreateBookingRequest) (*CreateBookingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBooking not implemented")
}
func (UnimplementedBookingServiceServer) UpdateBookingStatus(context.Context, *UpdateBookingStatusRequest) (*UpdateBookingStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateBookingStatus not implemented")
}
func (UnimplementedBookingServiceServer) mustEmbedUnimplementedBookingServiceServer() {}

// UnsafeBookingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BookingServiceServer will
// result in compilation errors.
type UnsafeBookingServiceServer interface {
	mustEmbedUnimplementedBookingServiceServer()
}

func RegisterBookingServiceServer(s grpc.ServiceRegistrar, srv BookingServiceServer) {
	s.RegisterService(&BookingService_ServiceDesc, srv)
}

func _BookingService_GetAvailableRooms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAvailableRoomsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookingServiceServer).GetAvailableRooms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hotel.booking.v1.BookingService/GetAvailableRooms",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookingServiceServer).GetAvailableRooms(ctx, req.(*GetAvailableRoomsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookingService_CreateBooking_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBookingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookingServiceServer).CreateBooking(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hotel.booking.v1.BookingService/CreateBooking",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookingServiceServer).CreateBooking(ctx, req.(*CreateBookingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookingService_UpdateBookingStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateBookingStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookingServiceServer).UpdateBookingStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hotel.booking.v1.BookingService/UpdateBookingStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookingServiceServer).UpdateBookingStatus(ctx, req.(*UpdateBookingStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BookingService_ServiceDesc is the grpc.ServiceDesc for BookingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BookingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "hotel.booking.v1.BookingService",
	HandlerType: (*BookingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAvailableRooms",
			Handler:    _BookingService_GetAvailableRooms_Handler,
		},
		{
			MethodName: "CreateBooking",
			Handler:    _BookingService_CreateBooking_Handler,
		},
		{
			MethodName: "UpdateBookingStatus",
			Handler:    _BookingService_UpdateBookingStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "booking/booking.proto",
}
