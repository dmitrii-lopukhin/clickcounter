package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/rs/zerolog/log"
)

type ClickRepository interface {
	InsertStats(stats []ClickStat) error
	GetStats(bannerID int, from, to time.Time) ([]ClickStat, error)
}

type clickRepository struct {
	db *sql.DB
}

func NewClickRepository(db *sql.DB) ClickRepository {
	return &clickRepository{db: db}
}

func (r *clickRepository) InsertStats(stats []ClickStat) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO banner_clicks (timestamp, banner_id, count)
        VALUES ($1, $2, $3)
        ON CONFLICT (timestamp, banner_id)
        DO UPDATE SET count = banner_clicks.count + EXCLUDED.count
    `)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, stat := range stats {
		_, err = stmt.ExecContext(ctx, stat.Timestamp, stat.BannerID, stat.Count)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *clickRepository) GetStats(bannerID int, from, to time.Time) ([]ClickStat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `
        SELECT timestamp, banner_id, count
        FROM banner_clicks
        WHERE banner_id = $1 AND timestamp >= $2 AND timestamp <= $3
        ORDER BY timestamp
    `, bannerID, from, to)
	if err != nil {
		log.Error().Err(err).Msg("GetStats query error")
		return nil, err
	}
	defer rows.Close()

	var stats []ClickStat
	for rows.Next() {
		var stat ClickStat
		if err := rows.Scan(&stat.Timestamp, &stat.BannerID, &stat.Count); err != nil {
			log.Error().Err(err).Msg("GetStats scan error")
			return nil, err
		}
		stats = append(stats, stat)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("GetStats rows iteration error")
		return nil, err
	}

	return stats, nil
}
