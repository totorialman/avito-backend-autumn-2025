package user

import "github.com/totorialman/avito-backend-autumn-2025/internal/model"

type userService interface {
	SetActive(userID string, isActive bool) (*model.User, error)
	GetPRsByReviewer(userID string) ([]model.PullRequest, error)
}
