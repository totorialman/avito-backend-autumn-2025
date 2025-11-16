package pr

import (
	"math/rand"
	"time"

	"github.com/totorialman/avito-backend-autumn-2025/internal/model"
)

type PRService struct {
	pullRequests pullRequestRepo
	users        userRepo
}

func NewPrService(pullRequests pullRequestRepo, users userRepo) *PRService {
	return &PRService{pullRequests: pullRequests, users: users}
}

func (s *PRService) CreatePR(prID, prName, authorID string) (*model.PullRequest, error) {
	existing, err := s.pullRequests.GetPR(prID)
	if err == nil && existing != nil {
		return nil, model.ErrPRExists
	}

	author, err := s.users.GetUser(authorID)
	if err != nil {
		return nil, model.ErrNotFound
	}

	candidates, err := s.users.GetActiveTeamUsersExcept(author.TeamName, authorID)
	if err != nil {
		return nil, err
	}

	reviewers := make([]string, 0, 2)
	maxReviewers := 2
	if len(candidates) < maxReviewers {
		maxReviewers = len(candidates)
	}
	for i := 0; i < maxReviewers; i++ {
		reviewers = append(reviewers, candidates[i].ID)
	}

	pr := model.PullRequest{
		ID:                prID,
		Name:              prName,
		AuthorID:          authorID,
		Status:            "OPEN",
		AssignedReviewers: reviewers,
		CreatedAt:         time.Now(),
	}

	if err := s.pullRequests.CreatePR(pr, reviewers); err != nil {
		return nil, err
	}

	return &pr, nil
}

func (s *PRService) MergePR(prID string) (*model.PullRequest, error) {
	pr, err := s.pullRequests.GetPR(prID)
	if err != nil {
		return nil, model.ErrNotFound
	}

	if pr.Status == "MERGED" {
		return pr, nil
	}

	now := time.Now()
	pr.Status = "MERGED"
	pr.MergedAt = &now

	if err := s.pullRequests.MergePR(prID); err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *PRService) ReassignReviewer(prID, oldUserID string) (*model.PullRequest, string, error) {
	pr, err := s.pullRequests.GetPR(prID)
	if err != nil {
		return nil, "", model.ErrNotFound
	}

	if pr.IsMerged() {
		return nil, "", model.ErrPRMerged
	}

	oldUser, err := s.users.GetUser(oldUserID)
	if err != nil {
		return nil, "", model.ErrNotFound
	}

	if !pr.IsReviewerAssigned(oldUserID) {
		return nil, "", model.ErrNotAssigned
	}

	candidates, err := s.users.GetActiveTeamUsersExcept(oldUser.TeamName, oldUserID)
	if err != nil {
		return nil, "", err
	}

	authorID := pr.AuthorID
	filteredCandidates := make([]model.User, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate.ID != authorID {
			filteredCandidates = append(filteredCandidates, candidate)
		}
	}

	existingReviewers := make(map[string]bool)
	for _, r := range pr.AssignedReviewers {
		existingReviewers[r] = true
	}

	finalCandidates := make([]model.User, 0)
	for _, c := range filteredCandidates {
		if c.ID == oldUserID {
			continue 
		}
		if !existingReviewers[c.ID] {
			finalCandidates = append(finalCandidates, c)
		}
	}

	if len(finalCandidates) == 0 {
		return nil, "", model.ErrNoCandidate
	}

	newReviewer := finalCandidates[rand.Intn(len(finalCandidates))]

	var newReviewers []string
	for _, r := range pr.AssignedReviewers {
		if r == oldUserID {
			newReviewers = append(newReviewers, newReviewer.ID)
		} else {
			newReviewers = append(newReviewers, r)
		}
	}

	if err := s.pullRequests.UpdateReviewers(prID, newReviewers); err != nil {
		return nil, "", err
	}

	pr.AssignedReviewers = newReviewers
	return pr, newReviewer.ID, nil
}

func (s *PRService) GetPRsByReviewer(userID string) ([]model.PullRequest, error) {
	return s.pullRequests.GetPRsByReviewer(userID)
}
