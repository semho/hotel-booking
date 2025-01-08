package mapper

import (
	"github.com/google/uuid"
	"github.com/semho/hotel-booking/booking-service/internal/domain/model"
	bookingpb "github.com/semho/hotel-booking/pkg/proto/booking_v1/booking"
	roompb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ProtoToSearchParamsInternal(req *bookingpb.GetAvailableRoomsRequest) model.SearchParams {
	var roomType *roompb.RoomType
	if req.Type != nil {
		rt := *req.Type
		roomType = &rt
	}

	return model.SearchParams{
		CheckIn:  req.CheckIn.AsTime(),
		CheckOut: req.CheckOut.AsTime(),
		Capacity: req.Capacity,
		Type:     roomType,
	}
}

func ProtoToSearchParams(req *bookingpb.GetAvailableRoomsRequest) model.SearchRoomsParams {
	var roomType *roompb.RoomType
	if req.Type != nil {
		rt := *req.Type
		roomType = &rt
	}

	return model.SearchRoomsParams{
		Capacity: req.Capacity,
		Type:     roomType,
	}
}

func ProtoToBooking(req *bookingpb.CreateBookingRequest) *model.Booking {
	var userID *uuid.UUID
	if req.UserId != nil {
		if id, err := uuid.Parse(*req.UserId); err == nil {
			userID = &id
		}
	}

	return &model.Booking{
		UserID:     userID,
		GuestName:  req.GuestName,
		GuestEmail: req.GuestEmail,
		GuestPhone: req.GuestPhone,
		CheckIn:    req.CheckIn.AsTime(),
		CheckOut:   req.CheckOut.AsTime(),
		// RoomID будет установлен в сервисе после выбора свободной комнаты
		// TotalPrice будет рассчитана в сервисе
		// CurrentStatus будет установлен в сервисе
	}
}

//TODO: подумать о композиции:
//func ProtoToSearchParams(input SearchParamsInput) model.SearchParams {
//	var roomType *roompb.RoomType
//	if input.Type != nil {
//		rt := *input.Type
//		roomType = &rt
//	}
//
//	return model.SearchParams{
//		CheckIn:  input.CheckIn.AsTime(),
//		CheckOut: input.CheckOut.AsTime(),
//		Capacity: input.Capacity,
//		Type:     roomType,
//	}
//}
//
//// Частные конверторы для каждого типа запроса
//func CreateBookingToSearchParamsInput(req *bookingpb.CreateBookingRequest) SearchParamsInput {
//	return SearchParamsInput{
//		CheckIn:  req.CheckIn,
//		CheckOut: req.CheckOut,
//		Capacity: req.Capacity,
//		Type:     req.Type,
//	}
//}
//
//func CheckAvailabilityToSearchParamsInput(req *bookingpb.CheckRoomsAvailabilityRequest) SearchParamsInput {
//	return SearchParamsInput{
//		CheckIn:  req.CheckIn,
//		CheckOut: req.CheckOut,
//		Capacity: req.Capacity,
//		Type:     req.Type,
//	}
//}

func BookingToProto(booking *model.Booking) *bookingpb.Booking {
	var userID *string
	if booking.UserID != nil {
		id := booking.UserID.String()
		userID = &id
	}

	var currentStatus bookingpb.BookingStatus
	if booking.CurrentStatus != nil {
		currentStatus = booking.CurrentStatus.Status
	}

	return &bookingpb.Booking{
		Id:            booking.ID.String(),
		RoomId:        booking.RoomID.String(),
		UserId:        userID,
		GuestName:     booking.GuestName,
		GuestEmail:    booking.GuestEmail,
		GuestPhone:    booking.GuestPhone,
		CheckIn:       timestamppb.New(booking.CheckIn),
		CheckOut:      timestamppb.New(booking.CheckOut),
		TotalPrice:    booking.TotalPrice,
		CreatedAt:     timestamppb.New(booking.CreatedAt),
		CurrentStatus: currentStatus,
	}
}
