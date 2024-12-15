package response

import (
	"github.com/lib/pq"
	roompb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
)

type CreateRoomResponse struct {
	ID        string            `json:"id"`
	Number    string            `json:"number"`
	Type      roompb.RoomType   `json:"type"`
	Price     string            `json:"price"`
	Capacity  int               `json:"capacity"`
	Status    roompb.RoomStatus `json:"status"`
	Amenities pq.StringArray    `json:"amenities"`
}
