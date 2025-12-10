package hotels

import (
	"context"

	"github.com/google/uuid"
)

type HotelService interface {
	CreateHotel(ctx context.Context, hotel *Hotel) (*Hotel, error)
	ListHotelsWithRooms(ctx context.Context, businessID uuid.UUID, limit, offset int) ([]*HotelWithRooms, error)
}

type RoomsLister interface {
	ListRoomsByHotelIDs(ctx context.Context, hotelIDs []uuid.UUID) (map[uuid.UUID][]RoomSummary, error)
}

type service struct {
	repository  HotelRepository
	roomsLister RoomsLister
}

func NewService(repository HotelRepository, roomsLister RoomsLister) HotelService {
	return &service{
		repository:  repository,
		roomsLister: roomsLister,
	}
}

func (s *service) CreateHotel(ctx context.Context, hotel *Hotel) (*Hotel, error) {
	newHotel := &Hotel{
		ID:         uuid.New(),
		BusinessID: hotel.BusinessID,
		Name:       hotel.Name,
		Address:    hotel.Address,
	}

	createdHotel, err := s.repository.CreateHotel(ctx, newHotel)
	if err != nil {
		return nil, err
	}

	return createdHotel, nil
}

func (s *service) ListHotelsWithRooms(ctx context.Context, businessID uuid.UUID, limit, offset int) ([]*HotelWithRooms, error) {
	hotelsList, err := s.repository.ListHotelsByBusiness(ctx, businessID, limit, offset)
	if err != nil {
		return nil, err
	}
	if len(hotelsList) == 0 {
		return []*HotelWithRooms{}, nil
	}

	var hotelIDs []uuid.UUID
	for _, h := range hotelsList {
		hotelIDs = append(hotelIDs, h.ID)
	}

	// fetch rooms via the adapter
	roomsMap, err := s.roomsLister.ListRoomsByHotelIDs(ctx, hotelIDs)
	if err != nil {
		return nil, err
	}

	var out []*HotelWithRooms
	for _, h := range hotelsList {
		out = append(out, &HotelWithRooms{
			HotelCode: h.HotelCode,
			Name:      h.Name,
			Address:   h.Address,
			Rooms:     roomsMap[h.ID], // if key missing, zero-value slice ([]RoomSummary{}) is returned
		})
	}
	return out, nil
}
