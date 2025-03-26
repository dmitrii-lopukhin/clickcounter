package app

import (
	"time"

	"github.com/dmitrii-lopukhin/clicks-counter/internal/repository"
)

type ClickService interface {
	HandleClick(bannerID int)
	FlushStats(stats []repository.ClickStat)
	GetStats(bannerID int, from, to time.Time) ([]repository.ClickStat, error)
}
