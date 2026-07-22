package server

import (
	"context"
	"log"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func New(addr string, handler http.Handler) *Server {

	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s *Server) Start() error {

	log.Printf("Edge listening on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {

	log.Println("Shutting down server...")

	return s.httpServer.Shutdown(ctx)
}
