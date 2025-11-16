package pr

import "github.com/totorialman/avito-backend-autumn-2025/internal/model"

type pullRequestService interface {
	CreatePR(prID, prName, authorID string) (*model.PullRequest, error)
	MergePR(prID string) (*model.PullRequest, error)
	ReassignReviewer(prID, oldUserID string) (*model.PullRequest, string, error)
	GetPRsByReviewer(userID string) ([]model.PullRequest, error)
}
