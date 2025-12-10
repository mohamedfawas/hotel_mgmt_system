package tenant

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/hotel_mgmt_system/pkg/apiresponse"
	"github.com/mohamedfawas/hotel_mgmt_system/pkg/apperror"
	"github.com/mohamedfawas/hotel_mgmt_system/pkg/constants"
)

type Cache interface {
	Exists(ctx context.Context, key string) (bool, error)
	Set(ctx context.Context, key string, value interface{}) error
}

type TenantRepository interface {
	Exists(ctx context.Context, id string) (bool, error)
}

type Middleware struct {
	cache   Cache
	tenants TenantRepository
}

func New(cache Cache, tenants TenantRepository) *Middleware {
	return &Middleware{
		cache:   cache,
		tenants: tenants,
	}
}

func (m *Middleware) ResolveTenant() gin.HandlerFunc {
	return func(c *gin.Context) {

		businessID := c.GetHeader(constants.BusinessIDHeader)
		if businessID == "" {
			apiresponse.Error(c, apperror.ErrBusinessIDMissing, nil)
			c.Abort()
			return
		}

		cacheKey := "business:" + businessID

		exists, err := m.cache.Exists(c.Request.Context(), cacheKey)
		if err != nil {
			apiresponse.Error(c, err, nil)
			c.Abort()
			return
		}

		if exists {
			c.Set(constants.ContextBusinessID, businessID)
			c.Next()
			return
		}

		dbExists, err := m.tenants.Exists(c.Request.Context(), businessID)
		if err != nil {
			apiresponse.Error(c, err, nil)
			c.Abort()
			return
		}

		if !dbExists {
			apiresponse.Error(c, apperror.ErrBusinessNotFound, nil)
			c.Abort()
			return
		}

		_ = m.cache.Set(c.Request.Context(), cacheKey, "1")

		c.Set(constants.ContextBusinessID, businessID)
		c.Next()
	}
}
