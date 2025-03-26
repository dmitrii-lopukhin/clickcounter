package server

import (
	"context"
	"net/http"

	"github.com/dmitrii-lopukhin/clicks-counter/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type Server struct {
	cfg    *config.Config
	log    zerolog.Logger
	router http.Handler
	http   *http.Server
}

func New(cfg *config.Config, log zerolog.Logger) *Server {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: r,
	}

	return &Server{
		cfg:    cfg,
		log:    log,
		router: r,
		http:   srv,
	}
}

func (s *Server) Start() error {
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) {
	s.http.Shutdown(ctx)
}
