package team

import "github.com/totorialman/avito-backend-autumn-2025/internal/model"

type teamService interface {
	CreateTeam(teamName string, members []model.User) error
	GetTeam(teamName string) (*model.Team, error)
}
