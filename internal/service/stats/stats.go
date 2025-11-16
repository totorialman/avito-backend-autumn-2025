package stats

import (
	"context"

	"github.com/totorialman/avito-backend-autumn-2025/internal/model"
)

type StatsService struct {
	stats statsRepository
}

func NewStatsService(statsRepo statsRepository) *StatsService {
	return &StatsService{stats: statsRepo}
}

func (u *StatsService) GetReviewerAssignmentStats() ([]model.ReviewerAssignmentStat, error) {
	return u.stats.GetReviewerAssignmentStats(context.Background())
}