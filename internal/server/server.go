package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/khareutkarshk/dug/internal/config"
)

type Server struct {
	httpServer *http.Server
	cfg        *config.Config
	listener   net.Listener
}

func New(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		cfg: cfg,
		httpServer: &http.Server{
			Addr:         ":" + strconv.Itoa(cfg.Server.Port),
			Handler:      handler,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
	}
}

func (s *Server) Start() error {

	var err error

	s.listener, err = net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return err
	}

	log.Printf("Edge listening on %s", s.listener.Addr())

	if s.cfg.Server.TLS.Enabled {
		log.Println("HTTPS enabled")

		return s.httpServer.ServeTLS(
			s.listener,
			s.cfg.Server.TLS.CertFile,
			s.cfg.Server.TLS.KeyFile,
		)
	}

	log.Println("HTTP enabled")

	return s.httpServer.Serve(s.listener)
}

func (s *Server) Shutdown(ctx context.Context) error {

	log.Println("Shutting down server...")

	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Addr() string {
	if s.listener == nil {
		return ""
	}
	return s.listener.Addr().String()
}
