package rooms

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamedfawas/hotel_mgmt_system/internal/hotels"
	"github.com/mohamedfawas/hotel_mgmt_system/pkg/apperror"
)

type RoomService interface {
	CreateRoom(ctx context.Context, hotelCode int64, room *Room) (*Room, error)
}

type service struct {
	repository RoomRepository
	hotelRepo  hotels.HotelRepository
}

func NewService(repository RoomRepository, hotelRepo hotels.HotelRepository) RoomService {
	return &service{
		repository: repository,
		hotelRepo:  hotelRepo,
	}
}

func (s *service) CreateRoom(ctx context.Context, hotelCode int64, room *Room) (*Room, error) {

	hotel, err := s.hotelRepo.GetHotelByHotelCode(ctx, hotelCode)
	if err != nil {
		return nil, err
	}
	if hotel == nil {
		return nil, apperror.ErrRequestedHotelNotFound
	}

	if hotel.BusinessID != room.BusinessID {
		return nil, apperror.ErrTenantNotAuthorized
	}
	room.HotelID = hotel.ID
	room.ID = uuid.New()
	createdRoom, err := s.repository.CreateRoom(ctx, room)
	if err != nil {
		return nil, err
	}
	return createdRoom, nil

}
