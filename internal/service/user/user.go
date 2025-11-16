package user

import (
	"github.com/totorialman/avito-backend-autumn-2025/internal/model"
)

type UserService struct {
	users        userRepo
	pullRequests pullRequestRepo
}

func NewUserService(users userRepo, pullRequests pullRequestRepo) *UserService {
	return &UserService{users: users, pullRequests: pullRequests}
}

func (s *UserService) SetActive(userID string, isActive bool) (*model.User, error) {
	user, err := s.users.SetActive(userID, isActive)
	if err != nil {
		return nil, model.ErrNotFound
	}
	return user, nil
}

func (s *UserService) GetPRsByReviewer(userID string) ([]model.PullRequest, error) {
	return s.pullRequests.GetPRsByReviewer(userID)
}
