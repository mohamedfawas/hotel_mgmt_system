package hotels

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mohamedfawas/hotel_mgmt_system/internal/db"
)

type HotelRepository interface {
	CreateHotel(ctx context.Context, hotel *Hotel) (*Hotel, error)
	GetHotelByHotelCode(ctx context.Context, hotelCode int64) (*Hotel, error)
	ListHotelsByBusiness(ctx context.Context, businessID uuid.UUID, limit, offset int) ([]*Hotel, error)
}

type repository struct {
	db *db.Client
}

func NewRepository(dbClient *db.Client) HotelRepository {
	return &repository{db: dbClient}
}

func (r *repository) CreateHotel(ctx context.Context, hotel *Hotel) (*Hotel, error) {
	const query = `
		INSERT INTO hotels (id, business_id, name, address)
		VALUES ($1, $2, $3, $4)
		RETURNING hotel_code, name, address
	`

	var createdHotel Hotel

	err := r.db.Pool.QueryRow(
		ctx,
		query,
		hotel.ID,
		hotel.BusinessID,
		hotel.Name,
		hotel.Address,
	).Scan(
		&createdHotel.HotelCode,
		&createdHotel.Name,
		&createdHotel.Address,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create hotel: %w", err)
	}

	return &createdHotel, nil
}

func (r *repository) GetHotelByHotelCode(ctx context.Context, hotelCode int64) (*Hotel, error) {
	const query = `
		SELECT id, business_id, hotel_code, name, address
		FROM hotels
		WHERE hotel_code = $1
	`

	var fetchedHotel Hotel
	row := r.db.Pool.QueryRow(ctx, query, hotelCode)
	err := row.Scan(
		&fetchedHotel.ID,
		&fetchedHotel.BusinessID,
		&fetchedHotel.HotelCode,
		&fetchedHotel.Name,
		&fetchedHotel.Address)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query hotel by code: %w", err)
	}
	return &fetchedHotel, nil
}

func (r *repository) ListHotelsByBusiness(ctx context.Context, businessID uuid.UUID, limit, offset int) ([]*Hotel, error) {
	const query = `
		SELECT id, business_id, hotel_code, name, address
		FROM hotels
		WHERE business_id = $1
		ORDER BY hotel_code DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Pool.Query(ctx, query, businessID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query hotels: %w", err)
	}
	defer rows.Close()
	var res []*Hotel
	for rows.Next() {
		var h Hotel
		if err := rows.Scan(&h.ID, &h.BusinessID, &h.HotelCode, &h.Name, &h.Address); err != nil {
			return nil, fmt.Errorf("failed to scan hotel row: %w", err)
		}
		res = append(res, &h)
	}
	return res, nil
}
