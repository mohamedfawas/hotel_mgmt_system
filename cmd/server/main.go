package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mohamedfawas/hotel_mgmt_system/internal/app"
	"github.com/mohamedfawas/hotel_mgmt_system/internal/config"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// initialize app (wires DB, cache, router)
	application, err := app.NewApp(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- application.Run()
	}()

	log.Printf("server started at the port : %s", cfg.HTTP.Port)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Printf("signal received: %v\n", sig)
	case err := <-errCh:
		log.Printf("server error: %v\n", err)
	}

	// begin graceful shutdown
	shutdownTimeout := 15 * time.Second
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		log.Printf("error during shutdown: %v", err)
	} else {
		log.Println("graceful shutdown complete")
	}
}
