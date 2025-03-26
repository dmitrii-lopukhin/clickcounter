package handler

import (
	"net/http"
	"strconv"

	"github.com/dmitrii-lopukhin/clicks-counter/internal/app"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type CounterHandler struct {
	service app.ClickService
	log     zerolog.Logger
}

func NewCounterHandler(service app.ClickService, log zerolog.Logger) *CounterHandler {
	return &CounterHandler{
		service: service,
		log:     log,
	}
}

func (h *CounterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bannerIDStr := chi.URLParam(r, "bannerID")
	bannerID, err := strconv.Atoi(bannerIDStr)
	if err != nil {
		http.Error(w, "Invalid bannerID", http.StatusBadRequest)
		return
	}

	h.service.HandleClick(bannerID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
