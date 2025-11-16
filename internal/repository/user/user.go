package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/totorialman/avito-backend-autumn-2025/internal/model"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) SetActive(userID string, isActive bool) (*model.User, error) {
	const setUserActiveQuery = `UPDATE users SET is_active = $1 WHERE user_id = $2`

	_, err := r.db.Exec(context.Background(),
		setUserActiveQuery,
		isActive, userID)
	if err != nil {
		return nil, err
	}
	return r.GetUser(userID)
}

func (r *UserRepository) GetUser(userID string) (*model.User, error) {
	const getUserQuery = `SELECT user_id, username, team_name, is_active FROM users WHERE user_id = $1`
	var u model.User

	err := r.db.QueryRow(context.Background(),
		getUserQuery,
		userID).Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetTeamUsers(teamName string) ([]model.User, error) {
	return r.getUsersByTeam(teamName, "")
}

func (r *UserRepository) GetActiveTeamUsersExcept(teamName, excludeUserID string) ([]model.User, error) {
	return r.getUsersByTeam(teamName, excludeUserID)
}

func (r *UserRepository) getUsersByTeam(teamName, excludeUserID string) ([]model.User, error) {
	const getActiveTeamUsersQuery = `
		SELECT user_id, username, is_active
		FROM users
		WHERE team_name = $1
		  AND is_active = true
		  AND ($2 = '' OR user_id != $2)`

	rows, err := r.db.Query(context.Background(), getActiveTeamUsersQuery, teamName, excludeUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		u.TeamName = teamName
		if err := rows.Scan(&u.ID, &u.Username, &u.IsActive); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
