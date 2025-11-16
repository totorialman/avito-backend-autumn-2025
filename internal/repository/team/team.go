package team

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/totorialman/avito-backend-autumn-2025/internal/model"
)

type TeamRepository struct {
	db *pgxpool.Pool
}

func NewTeamRepository(db *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) CreateTeam(team model.Team) error {
	const (
		insertTeamQuery = `INSERT INTO teams (team_name) VALUES ($1)`
		upsertUserQuery = `INSERT INTO users (user_id, username, team_name, is_active) VALUES ($1, $2, $3, $4)
			 ON CONFLICT (user_id) DO UPDATE SET username = $2, team_name = $3, is_active = $4`
	)
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), insertTeamQuery, team.Name)
	if err != nil {
		return err
	}

	for _, user := range team.Members {
		_, err = tx.Exec(context.Background(),
			upsertUserQuery,
			user.ID, user.Username, team.Name, user.IsActive)
		if err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
}

func (r *TeamRepository) GetTeam(name string) (*model.Team, error) {
	const (
		getTeamExistsQuery  = `SELECT team_name FROM teams WHERE team_name = $1`
		getUsersByTeamQuery = `SELECT user_id, username, is_active FROM users WHERE team_name = $1`
	)
	var exists string

	err := r.db.QueryRow(context.Background(), getTeamExistsQuery, name).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(context.Background(),
		getUsersByTeamQuery,
		name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []model.User
	for rows.Next() {
		var u model.User
		u.TeamName = name
		if err := rows.Scan(&u.ID, &u.Username, &u.IsActive); err != nil {
			return nil, err
		}
		members = append(members, u)
	}

	return &model.Team{Name: name, Members: members}, nil
}
