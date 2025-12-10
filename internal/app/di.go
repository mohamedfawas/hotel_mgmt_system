package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mohamedfawas/hotel_mgmt_system/internal/cache"
	"github.com/mohamedfawas/hotel_mgmt_system/internal/config"
	"github.com/mohamedfawas/hotel_mgmt_system/internal/db"
	"github.com/mohamedfawas/hotel_mgmt_system/internal/hotels"
	"github.com/mohamedfawas/hotel_mgmt_system/internal/rooms"
	"github.com/mohamedfawas/hotel_mgmt_system/internal/tenant"
)

type App struct {
	Config     *config.Config
	DB         *db.Client
	Cache      *cache.Client
	Router     *gin.Engine
	httpServer *http.Server

	TenantMw     *tenant.Middleware
	HotelHandler *hotels.Handler
	RoomHandler  *rooms.Handler
}

type roomListerAdapter struct {
	repo rooms.RoomRepository
}

func (r *roomListerAdapter) ListRoomsByHotelIDs(ctx context.Context, hotelIDs []uuid.UUID) (map[uuid.UUID][]hotels.RoomSummary, error) {
	roomRows, err := r.repo.ListRoomsByHotelIDs(ctx, hotelIDs)
	if err != nil {
		return nil, err
	}

	out := make(map[uuid.UUID][]hotels.RoomSummary)
	for _, rm := range roomRows {
		// defense-in-depth: ensure rm is not nil and hotel id exists
		if rm == nil {
			continue
		}
		out[rm.HotelID] = append(out[rm.HotelID], hotels.RoomSummary{
			RoomNumber: rm.RoomNumber,
			RoomType:   rm.RoomType,
		})
	}
	return out, nil
}

func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {
	dbClient, err := db.NewClient(ctx, db.Config{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		DBName:   cfg.Postgres.DBName,
		SSLMode:  cfg.Postgres.SSLMode,
	})
	if err != nil {
		return nil, fmt.Errorf("db initialization failed: %w", err)
	}

	cacheClient, err := cache.NewClient(ctx, cache.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		dbClient.Close()
		return nil, fmt.Errorf("cache initialization failed: %w", err)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	app := &App{
		Config: cfg,
		DB:     dbClient,
		Cache:  cacheClient,
		Router: router,
	}

	// initialize dependencies (repositories, services, handlers, middleware)
	if err := app.initDependencies(); err != nil {
		// if wiring fails, close clients and return error
		dbClient.Close()
		// you might also close cacheClient if it has Close()
		return nil, fmt.Errorf("failed to init dependencies: %w", err)
	}

	app.registerRoutes()

	return app, nil
}

func (a *App) initDependencies() error {
	// tenant middleware (uses cache + tenant repo)
	tenantRepo := tenant.NewTenantRepository(a.DB)
	a.TenantMw = tenant.New(a.Cache, tenantRepo)

	hotelRepo := hotels.NewRepository(a.DB)
	roomRepo := rooms.NewRepository(a.DB)
	roomsAdapter := &roomListerAdapter{repo: roomRepo}


	hotelSvc := hotels.NewService(hotelRepo, roomsAdapter)
	roomSvc := rooms.NewService(roomRepo, hotelRepo)


	a.HotelHandler = hotels.NewHandler(hotelSvc)
	a.RoomHandler = rooms.NewHandler(roomSvc)

	return nil
}
