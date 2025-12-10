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

	// blocks until server stops or returns an error
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("listen and serve: %w", err)
	}

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	var firstErr error

	// 1) HTTP shutdown
	if a.httpServer != nil {
		if err := a.httpServer.Shutdown(ctx); err != nil {
			log.Printf("http shutdown error: %v", err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	// 2) close DB
	if a.DB != nil {
		a.DB.Close()
	}

	// 3) close cache
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
