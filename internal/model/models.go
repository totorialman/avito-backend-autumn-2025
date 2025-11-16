package model

import "time"

type Team struct {
	Name    string
	Members []User
}

type User struct {
	ID       string
	Username string
	TeamName string
	IsActive bool
}

type PullRequest struct {
	ID                string
	Name              string
	AuthorID          string
	Status            string
	AssignedReviewers []string
	CreatedAt         time.Time
	MergedAt          *time.Time
}

type ReviewerAssignmentStat struct {
	UserID      string
	Username    string
	TeamName    string
	Assignments int
}

func (pr *PullRequest) IsMerged() bool {
	return pr.Status == "MERGED"
}

func (pr *PullRequest) IsReviewerAssigned(userID string) bool {
	for _, r := range pr.AssignedReviewers {
		if r == userID {
			return true
		}
	}
	return false
}
