package stats

import (
	"context"

	"github.com/totorialman/avito-backend-autumn-2025/internal/model"
)

type statsRepository interface {
	GetReviewerAssignmentStats(ctx context.Context) ([]model.ReviewerAssignmentStat, error)
}
