package request

import (
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
	roompb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
)

type CreateRoomRequest struct {
	RoomNumber string            `json:"room_number"`
	Type       roompb.RoomType   `json:"type"`
	Price      string            `json:"price"`
	Capacity   int               `json:"capacity"`
	Status     roompb.RoomStatus `json:"status"`
	Amenities  pq.StringArray    `json:"amenities"`
}

// UnmarshalJSON implements custom JSON unmarshaling for RoomType
func (r *CreateRoomRequest) UnmarshalJSON(data []byte) error {
	// Определяем временную структуру с интерфейсом для type и status
	type Alias struct {
		RoomNumber string         `json:"room_number"`
		Type       interface{}    `json:"type"`
		Price      string         `json:"price"`
		Capacity   int            `json:"capacity"`
		Status     interface{}    `json:"status"`
		Amenities  pq.StringArray `json:"amenities"`
	}

	var alias Alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	// Заполняем простые поля
	r.RoomNumber = alias.RoomNumber
	r.Price = alias.Price
	r.Capacity = alias.Capacity
	r.Amenities = alias.Amenities

	// Обрабатываем Type
	switch v := alias.Type.(type) {
	case string:
		if typeValue, ok := roompb.RoomType_value[v]; ok {
			r.Type = roompb.RoomType(typeValue)
		} else {
			return fmt.Errorf("invalid room type string: %s", v)
		}
	case float64: // JSON числа декодируются как float64
		if _, ok := roompb.RoomType_name[int32(v)]; ok {
			r.Type = roompb.RoomType(v)
		} else {
			return fmt.Errorf("invalid room type number: %v", v)
		}
	default:
		return fmt.Errorf("room type must be string or number, got %T", v)
	}

	// Обрабатываем Status
	switch v := alias.Status.(type) {
	case string:
		if statusValue, ok := roompb.RoomStatus_value[v]; ok {
			r.Status = roompb.RoomStatus(statusValue)
		} else {
			return fmt.Errorf("invalid room status string: %s", v)
		}
	case float64: // JSON числа декодируются как float64
		if _, ok := roompb.RoomStatus_name[int32(v)]; ok {
			r.Status = roompb.RoomStatus(v)
		} else {
			return fmt.Errorf("invalid room status number: %v", v)
		}
	default:
		return fmt.Errorf("room status must be string or number, got %T", v)
	}

	return nil
}
