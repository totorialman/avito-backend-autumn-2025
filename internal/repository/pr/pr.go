package pr

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/totorialman/avito-backend-autumn-2025/internal/model"
)

type PrRepository struct {
	db *pgxpool.Pool
}

func NewPrRepository(db *pgxpool.Pool) *PrRepository {
	return &PrRepository{db: db}
}

func (r *PrRepository) CreatePR(pr model.PullRequest, reviewers []string) error {
	const (
		createPRQuery = `INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at)
		 VALUES ($1, $2, $3, $4, $5)`
		insertReviewerQuery = `INSERT INTO reviewers (pull_request_id, user_id) VALUES ($1, $2)`
	)

	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(),
		createPRQuery,
		pr.ID, pr.Name, pr.AuthorID, pr.Status, pr.CreatedAt)
	if err != nil {
		return err
	}

	for _, reviewer := range reviewers {
		_, err = tx.Exec(context.Background(),
			insertReviewerQuery,
			pr.ID, reviewer)
		if err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
}

func (r *PrRepository) GetPR(prID string) (*model.PullRequest, error) {
	const (
		getPRQuery = `SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at
		 FROM pull_requests WHERE pull_request_id = $1`
		getReviewersByPRQuery = `SELECT user_id FROM reviewers WHERE pull_request_id = $1`
	)
	var pr model.PullRequest
	var mergedAt sql.NullTime

	err := r.db.QueryRow(context.Background(),
		getPRQuery,
		prID).Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &mergedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if mergedAt.Valid {
		pr.MergedAt = &mergedAt.Time
	}

	rows, err := r.db.Query(context.Background(),
		getReviewersByPRQuery,
		prID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, uid)
	}
	pr.AssignedReviewers = reviewers

	return &pr, nil
}

func (r *PrRepository) UpdateReviewers(prID string, reviewers []string) error {
	const (
		deleteReviewersByPRQuery = `DELETE FROM reviewers WHERE pull_request_id = $1`
		insertReviewerQuery      = `INSERT INTO reviewers (pull_request_id, user_id) VALUES ($1, $2)`
	)

	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), deleteReviewersByPRQuery, prID)
	if err != nil {
		return err
	}

	for _, uid := range reviewers {
		_, err = tx.Exec(context.Background(), insertReviewerQuery, prID, uid)
		if err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
}

func (r *PrRepository) MergePR(prID string) error {
	const mergePRQuery = `UPDATE pull_requests SET status = 'MERGED', merged_at = $1 WHERE pull_request_id = $2`

	now := time.Now()
	_, err := r.db.Exec(context.Background(),
		mergePRQuery,
		now, prID)
	return err
}

func (r *PrRepository) GetPRsByReviewer(userID string) ([]model.PullRequest, error) {
	const (
		getPRsByReviewerQuery = `SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status, pr.created_at, pr.merged_at
		 FROM pull_requests pr
		 JOIN reviewers r ON pr.pull_request_id = r.pull_request_id
		 WHERE r.user_id = $1`
		getReviewersByPRQuery = `SELECT user_id FROM reviewers WHERE pull_request_id = $1`
	)

	rows, err := r.db.Query(context.Background(),
		getPRsByReviewerQuery,
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var prs []model.PullRequest
	for rows.Next() {
		var pr model.PullRequest
		var mergedAt sql.NullTime
		if err := rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &mergedAt); err != nil {
			return nil, err
		}
		if mergedAt.Valid {
			pr.MergedAt = &mergedAt.Time
		}

		reviewersRows, err := r.db.Query(context.Background(),
			getReviewersByPRQuery, pr.ID)
		if err != nil {
			return nil, err
		}
		var reviewers []string
		for reviewersRows.Next() {
			var uid string
			_ = reviewersRows.Scan(&uid)
			reviewers = append(reviewers, uid)
		}
		reviewersRows.Close()
		pr.AssignedReviewers = reviewers

		prs = append(prs, pr)
	}
	return prs, nil
}
