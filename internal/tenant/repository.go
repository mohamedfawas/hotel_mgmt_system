package tenant

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/mohamedfawas/hotel_mgmt_system/internal/db"
)

type tenantRepository struct {
	db *db.Client
}

func NewTenantRepository(db *db.Client) TenantRepository {
	return &tenantRepository{
		db: db,
	}
}

func (r *tenantRepository) Exists(ctx context.Context, id string) (bool, error) {
	const query = `
        SELECT 1 
        FROM businesses 
        WHERE id = $1
    `

	var dummy int
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&dummy)

	if err == pgx.ErrNoRows {
		// Tenant does NOT exist
		return false, nil
	}

	if err != nil {
		// Any DB error (connection, syntax, etc.)
		return false, fmt.Errorf("tenant lookup failed: %w", err)
	}

	return true, nil
}
