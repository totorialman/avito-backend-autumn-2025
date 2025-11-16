package team

import (
	"github.com/totorialman/avito-backend-autumn-2025/internal/model"
)

type TeamService struct {
	teams teamRepo
	users userRepo
}

func NewTeamService(teams teamRepo, users userRepo) *TeamService {
	return &TeamService{teams: teams, users: users}
}

func (s *TeamService) CreateTeam(teamName string, members []model.User) error {
	_, err := s.teams.GetTeam(teamName)
	if err == nil {
		return model.ErrTeamExists
	}

	team := model.Team{Name: teamName, Members: members}
	return s.teams.CreateTeam(team)
}

func (s *TeamService) GetTeam(teamName string) (*model.Team, error) {
	team, err := s.teams.GetTeam(teamName)
	if err != nil {
		return nil, model.ErrNotFound
	}
	return team, nil
}
