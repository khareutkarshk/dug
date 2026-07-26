package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/khareutkarshk/dug/internal/app"
	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/internal/logger"
	"github.com/khareutkarshk/dug/internal/metrics"
)

func main() {

	metrics.Register()

	cfg, err := config.Load("configs/edge.yaml")
	if err != nil {
		log.Fatal(err)
	}

	logger.Log.Info(
		"config loaded",
		"routes", len(cfg.Routes),
	)

	edge, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = edge.EnableConfigReload("configs/edge.yaml")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := edge.Server.Start(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-quit

	signal.Stop(quit)

	logger.Log.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := edge.Server.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v", err)
		return
	}

	logger.Log.Info("Server gracefully stopped")
}
