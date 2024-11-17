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

type CreateBookingResponse struct {
	ID       string    `json:"id"`
	RoomID   string    `json:"roomId"`
	UserInfo *UserInfo `json:"userInfo,omitempty"`
	Status   string    `json:"status"`
	Message  string    `json:"message"`
}
