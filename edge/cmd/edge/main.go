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

	"github.com/khareutkarshk/dug/edge/internal/config"
	"github.com/khareutkarshk/dug/edge/internal/router"
	"github.com/khareutkarshk/dug/edge/internal/server"
)

func main() {

	// Load configuration
	cfg, err := config.Load("configs/edge.yaml")

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Loaded %d routes", len(cfg.Routes))

	r, err := router.NewRouter(cfg.Routes)

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

	log.Println("Received shutdown signal, shutting down gracefully...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Shutdown returned error: %T: %v", err, err)
		return
	}

	log.Println("Server stopped gracefully")

}
