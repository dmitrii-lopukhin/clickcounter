package storage

import (
	"context"
	"time"

	"github.com/dmitrii-lopukhin/clicks-counter/internal/repository"
)

func StartFlusher(ctx context.Context, buffer *ClickBuffer, flushFunc func(stats []repository.ClickStat), interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stats := buffer.Flush()
			if len(stats) > 0 {
				flushFunc(stats)
			}
		case <-ctx.Done():
			stats := buffer.Flush()
			if len(stats) > 0 {
				flushFunc(stats)
			}
			return
		}
	}
}
