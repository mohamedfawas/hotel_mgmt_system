package rooms

import "github.com/google/uuid"

type Room struct {
	ID         uuid.UUID
	BusinessID uuid.UUID
	HotelID    uuid.UUID
	RoomNumber int
	RoomType   string
}
