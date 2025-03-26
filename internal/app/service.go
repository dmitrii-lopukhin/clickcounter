package app

import (
	"time"

	"github.com/dmitrii-lopukhin/clicks-counter/internal/repository"
	"github.com/dmitrii-lopukhin/clicks-counter/internal/storage"
)

type clickService struct {
	repo   repository.ClickRepository
	buffer *storage.ClickBuffer
}

func NewClickService(repo repository.ClickRepository, buffer *storage.ClickBuffer) ClickService {
	return &clickService{
		repo:   repo,
		buffer: buffer,
	}
}

func (s *clickService) HandleClick(bannerID int) {
	s.buffer.AddClick(bannerID)
}

func (s *clickService) FlushStats(stats []repository.ClickStat) {
	if len(stats) == 0 {
		return
	}
	err := s.repo.InsertStats(stats)
	if err != nil {
	}
}

func (s *clickService) GetStats(bannerID int, from, to time.Time) ([]repository.ClickStat, error) {
	return s.repo.GetStats(bannerID, from, to)
}
