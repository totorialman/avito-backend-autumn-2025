package stats

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/totorialman/avito-backend-autumn-2025/internal/model"
)

type StatsRepository struct {
	db *pgxpool.Pool
}

func NewStatsRepository(db *pgxpool.Pool) *StatsRepository {
	return &StatsRepository{db: db}
}

func (r *StatsRepository) GetReviewerAssignmentStats(ctx context.Context) ([]model.ReviewerAssignmentStat, error) {
	const query = `
		SELECT u.user_id, u.username, u.team_name, COUNT(r.user_id)
		FROM reviewers r
		JOIN users u ON r.user_id = u.user_id
		GROUP BY u.user_id, u.username, u.team_name
		ORDER BY COUNT(r.user_id) DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.ReviewerAssignmentStat
	for rows.Next() {
		var stat model.ReviewerAssignmentStat
		if err := rows.Scan(&stat.UserID, &stat.Username, &stat.TeamName, &stat.Assignments); err != nil {
			return nil, err
		}
		result = append(result, stat)
	}
	return result, nil
}
