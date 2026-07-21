package server

import (
	"log"
	"net/http"
)

func StartServer(addr string, handler http.Handler) error {
	log.Printf("Server is starting on %s\n", addr)

	return http.ListenAndServe(addr, handler)
}
