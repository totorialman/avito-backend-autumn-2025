package stats

import "github.com/totorialman/avito-backend-autumn-2025/internal/model"

type statsService interface {
	GetReviewerAssignmentStats() ([]model.ReviewerAssignmentStat, error)
}
