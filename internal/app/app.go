package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (a *App) Run() error {
	addr := ":" + a.Config.HTTP.Port

	srv := &http.Server{
		Addr:         addr,
		Handler:      a.Router,
		ReadTimeout:  time.Duration(a.Config.HTTP.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(a.Config.HTTP.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(a.Config.HTTP.IdleTimeout) * time.Second,
	}

	// store the server so Shutdown can access it
	a.httpServer = srv

	log.Printf("Server running on %s", addr)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("listen and serve: %w", err)
	}

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	var firstErr error

	if a.httpServer != nil {
		if err := a.httpServer.Shutdown(ctx); err != nil {
			log.Printf("http shutdown error: %v", err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	if a.DB != nil {
		a.DB.Close()
	}

	if a.Cache != nil {
		if err := a.Cache.Close(); err != nil {
			log.Printf("cache close error: %v", err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	return firstErr
}
