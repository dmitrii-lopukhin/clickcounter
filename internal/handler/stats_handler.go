package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/dmitrii-lopukhin/clicks-counter/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type StatsHandler struct {
	service app.ClickService
	log     zerolog.Logger
}

func NewStatsHandler(service app.ClickService, log zerolog.Logger) *StatsHandler {
	return &StatsHandler{
		service: service,
		log:     log,
	}
}

type StatsRequest struct {
	TsFrom time.Time `json:"tsFrom"`
	TsTo   time.Time `json:"tsTo"`
}

func (h *StatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bannerIDStr := chi.URLParam(r, "bannerID")
	bannerID, err := strconv.Atoi(bannerIDStr)
	if err != nil {
		http.Error(w, "Invalid bannerID", http.StatusBadRequest)
		return
	}

	var req StatsRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	stats, err := h.service.GetStats(bannerID, req.TsFrom, req.TsTo)
	if err != nil {
		h.log.Error().Err(err).Msg("Error retrieving stats in StatsHandler")
		http.Error(w, "Error retrieving stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
