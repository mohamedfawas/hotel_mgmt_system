package rooms

type CreateRoomRequest struct {
	RoomNumber int    `json:"room_number"`
	RoomType   string `json:"room_type"`
}

type CreateRoomResponse struct {
	RoomNumber int    `json:"room_number"`
	RoomType   string `json:"room_type"`
}

type HotelListItem struct {
	HotelCode int64        `json:"hotel_code"`
	HotelName string       `json:"hotel_name"`
	Address   string       `json:"address"`
	Rooms     []RoomOutput `json:"rooms"`
}

type RoomOutput struct {
	RoomNumber int    `json:"room_number"`
	RoomType   string `json:"room_type"`
}
