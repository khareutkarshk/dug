package main

import (
	"fmt"
	"log"

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
	log.Fatal(server.StartServer(addr, r))
}
