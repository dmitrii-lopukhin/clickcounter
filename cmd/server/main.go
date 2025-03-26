package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dmitrii-lopukhin/clicks-counter/internal/app"
	"github.com/dmitrii-lopukhin/clicks-counter/internal/config"
	"github.com/dmitrii-lopukhin/clicks-counter/internal/handler"
	"github.com/dmitrii-lopukhin/clicks-counter/internal/repository"
	"github.com/dmitrii-lopukhin/clicks-counter/internal/server"
	"github.com/dmitrii-lopukhin/clicks-counter/internal/storage"
	"github.com/dmitrii-lopukhin/clicks-counter/pkg/logger"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()
	logg := logger.New()

	db, err := sql.Open("postgres", cfg.DBConnString)
	if err != nil {
		logg.Fatal().Err(err).Msg("Failed to connect to DB")
	}
	defer db.Close()

	clickRepo := repository.NewClickRepository(db)

	buffer := storage.NewClickBuffer()

	clickService := app.NewClickService(clickRepo, buffer)

	counterHandler := handler.NewCounterHandler(clickService, logg)
	statsHandler := handler.NewStatsHandler(clickService, logg)

	r := server.NewRouter(counterHandler, statsHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: r,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go storage.StartFlusher(ctx, buffer, clickService.FlushStats, cfg.FlushInterval)

	go func() {
		logg.Info().Str("addr", ":"+cfg.HTTPPort).Msg("Starting HTTP server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.Fatal().Err(err).Msg("HTTP server error")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logg.Info().Msg("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, cfg.ShutdownTimeout)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logg.Error().Err(err).Msg("Server Shutdown error")
	}
}
