package rooms

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/mohamedfawas/hotel_mgmt_system/internal/db"
)

type RoomRepository interface {
	CreateRoom(ctx context.Context, room *Room) (*Room, error)
	ListRoomsByHotelIDs(ctx context.Context, hotelIDs []uuid.UUID) ([]*Room, error)
}

type repository struct {
	db *db.Client
}

func NewRepository(dbClient *db.Client) RoomRepository {
	return &repository{db: dbClient}
}

func (r *repository) CreateRoom(ctx context.Context, room *Room) (*Room, error) {
	const query = `
		INSERT INTO rooms (id, business_id, hotel_id, room_number, room_type)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING room_number, room_type
	`

	var createdRoom Room
	err := r.db.Pool.QueryRow(
		ctx,
		query,
		room.ID,
		room.BusinessID,
		room.HotelID,
		room.RoomNumber,
		room.RoomType,
	).Scan(
		&createdRoom.RoomNumber,
		&createdRoom.RoomType,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create room: %w", err)
	}
	createdRoom.ID = room.ID
	createdRoom.HotelID = room.HotelID
	createdRoom.BusinessID = room.BusinessID
	return &createdRoom, nil
}

func (r *repository) ListRoomsByHotelIDs(ctx context.Context, hotelIDs []uuid.UUID) ([]*Room, error) {
	if len(hotelIDs) == 0 {
		return []*Room{}, nil
	}
	const query = `
		SELECT id, business_id, hotel_id, room_number, room_type
		FROM rooms
		WHERE hotel_id = ANY($1)
	`
	rows, err := r.db.Pool.Query(ctx, query, hotelIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to query rooms: %w", err)
	}
	defer rows.Close()
	var res []*Room
	for rows.Next() {
		var rm Room
		if err := rows.Scan(&rm.ID, &rm.BusinessID, &rm.HotelID, &rm.RoomNumber, &rm.RoomType); err != nil {
			return nil, fmt.Errorf("failed to scan room row: %w", err)
		}
		res = append(res, &rm)
	}
	return res, nil
}
