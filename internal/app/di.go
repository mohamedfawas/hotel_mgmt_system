package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/hotel_mgmt_system/internal/cache"
	"github.com/mohamedfawas/hotel_mgmt_system/internal/config"
	"github.com/mohamedfawas/hotel_mgmt_system/internal/db"
)

type App struct {
	Config     *config.Config
	DB         *db.Client
	Cache      *cache.Client
	Router     *gin.Engine
	httpServer *http.Server
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

	app.registerRoutes()

	return app, nil
}
