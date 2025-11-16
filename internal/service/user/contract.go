package user

import "github.com/totorialman/avito-backend-autumn-2025/internal/model"

type userRepo interface {
	SetActive(userID string, isActive bool) (*model.User, error)
	GetUser(userID string) (*model.User, error)
	GetTeamUsers(teamName string) ([]model.User, error)
	GetActiveTeamUsersExcept(teamName, excludeUserID string) ([]model.User, error)
}

type pullRequestRepo interface {
	CreatePR(pr model.PullRequest, reviewers []string) error
	GetPR(prID string) (*model.PullRequest, error)
	UpdateReviewers(prID string, reviewers []string) error
	MergePR(prID string) error
	GetPRsByReviewer(userID string) ([]model.PullRequest, error)
}
