package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/internal/logger"
	"github.com/khareutkarshk/dug/internal/metrics"
	"github.com/khareutkarshk/dug/internal/router"
	"github.com/khareutkarshk/dug/internal/server"
)

func main() {

	// Load configuration
	cfg, err := config.Load("configs/edge.yaml")

	metrics.Register()

	if err != nil {
		log.Fatal(err)
	}

	logger.Log.Info(
		"config loaded", "routes", len(cfg.Routes),
	)

	r, err := router.NewRouter(cfg)

	if err != nil {
		log.Fatal(err)
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	srv := server.New(addr, r)

	go func() {
		if err := srv.Start(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {

			log.Fatal(err)
		}
	}()

	// create a channel to listen for OS signals
	quit := make(chan os.Signal, 1)

	// Notify the channel on interrupt or terminate signals
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Wait until a signal is received
	<-quit

	signal.Stop(quit)

	logger.Log.Info("Shutting signal received, shutting down server...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Shutdown returned error: %T: %v", err, err)
		return
	}

	logger.Log.Info("Server gracefully stopped")

}
