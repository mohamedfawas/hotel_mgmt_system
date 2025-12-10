package hotels

import "github.com/google/uuid"

type Hotel struct {
	ID         uuid.UUID
	BusinessID uuid.UUID
	HotelCode  int64
	Name       string
	Address    string
}

type HotelWithRooms struct {
	HotelCode int64         `json:"hotel_code"`
	Name      string        `json:"name"`
	Address   string        `json:"address"`
	Rooms     []RoomSummary `json:"rooms"`
}

type RoomSummary struct {
	RoomNumber int    `json:"room_number"`
	RoomType   string `json:"room_type"`
}
