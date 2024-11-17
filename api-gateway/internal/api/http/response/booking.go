package response

type AvailableRoom struct {
	ID         string   `json:"id"`
	RoomNumber string   `json:"roomNumber"`
	Type       string   `json:"type"`
	Price      float64  `json:"price"`
	Capacity   int      `json:"capacity"`
	Amenities  []string `json:"amenities"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
