package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(counterHandler http.Handler, statsHandler http.Handler) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/counter", func(r chi.Router) {
		r.Get("/{bannerID}", counterHandler.ServeHTTP)
	})

	r.Route("/stats", func(r chi.Router) {
		r.Post("/{bannerID}", statsHandler.ServeHTTP)
	})

	return r
}
