package main

import (
	"log"

	"github.com/khareutkarshk/dug/edge/internal/router"
	"github.com/khareutkarshk/dug/edge/internal/server"
)

func main() {
	r, err := router.NewRouter()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("🐶 Dug Edge is starting...")
	log.Fatal(server.StartServer(":8080", r))
}
