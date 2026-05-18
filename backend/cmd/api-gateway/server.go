// Package main (server.go) provides a thin HTTP server wrapper used by the api-gateway
// entrypoint to listen on the configured port and perform graceful shutdown.
package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// Server wraps net/http.Server with Pebble api-gateway listen and shutdown helpers.
type Server struct {
	httpServer *http.Server
}

// NewServer binds the Chi router to ":port" and returns a Server ready for Start.
func NewServer(port string, router *chi.Mux) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: router,
		},
	}
}

// Start blocks on ListenAndServe until the server stops or returns a non-ErrServerClosed error.
func (s *Server) Start() error {
	log.Info().Str("port", s.httpServer.Addr).Msg("api-gateway listening")
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop gracefully drains in-flight requests using http.Server.Shutdown with the given context.
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
