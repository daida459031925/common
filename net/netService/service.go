package netService

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	http.Server
}

func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		Server: http.Server{
			Addr:    addr,
			Handler: handler,

			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	log.Printf("Listening on: %s", s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.SetKeepAlivesEnabled(false)
	if err := s.Shutdown(ctx); err != nil {
		return err
	}
	log.Println("Server stopped")
	return nil
}
