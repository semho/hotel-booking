package mapper

import (
	"github.com/semho/hotel-booking/booking-service/internal/domain/model"
	bookingpb "github.com/semho/hotel-booking/pkg/proto/booking_v1/booking"
	roompb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
)

func ProtoToSearchParams(req *bookingpb.GetAvailableRoomsRequest) model.SearchParams {
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
