package storage

import (
	"sync"
	"time"

	"github.com/dmitrii-lopukhin/clicks-counter/internal/repository"
)

type ClickBuffer struct {
	mu     sync.Mutex
	clicks map[int]map[time.Time]int
}

func NewClickBuffer() *ClickBuffer {
	return &ClickBuffer{
		clicks: make(map[int]map[time.Time]int),
	}
}

func (cb *ClickBuffer) AddClick(bannerID int) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now().Truncate(time.Minute)
	if _, exists := cb.clicks[bannerID]; !exists {
		cb.clicks[bannerID] = make(map[time.Time]int)
	}
	cb.clicks[bannerID][now]++
}

func (cb *ClickBuffer) Flush() []repository.ClickStat {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	var stats []repository.ClickStat
	for bannerID, tsMap := range cb.clicks {
		for ts, count := range tsMap {
			stats = append(stats, repository.ClickStat{
				Timestamp: ts,
				BannerID:  bannerID,
				Count:     count,
			})
		}
	}
	cb.clicks = make(map[int]map[time.Time]int)
	return stats
}
