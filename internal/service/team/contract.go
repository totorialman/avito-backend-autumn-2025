package team

import "github.com/totorialman/avito-backend-autumn-2025/internal/model"

type teamRepo interface {
	CreateTeam(team model.Team) error
	GetTeam(name string) (*model.Team, error)
}

type userRepo interface {
	SetActive(userID string, isActive bool) (*model.User, error)
	GetUser(userID string) (*model.User, error)
	GetTeamUsers(teamName string) ([]model.User, error)
	GetActiveTeamUsersExcept(teamName, excludeUserID string) ([]model.User, error)
}
