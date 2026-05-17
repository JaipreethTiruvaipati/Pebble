package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// Server wraps the HTTP server.
type Server struct {
	httpServer *http.Server
}

// NewServer creates a new configured server.
func NewServer(port string, router *chi.Mux) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: router,
		},
	}
}

// Start runs the HTTP server.
func (s *Server) Start() error {
	log.Info().Str("port", s.httpServer.Addr).Msg("api-gateway listening")
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop gracefully shuts down the server.
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
