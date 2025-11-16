package user

type SetActiveResponse struct {
	User UserResponse `json:"user"`
}

type SetActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type UserResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type GetUserReviewsResponse struct {
	UserID       string    `json:"user_id"`
	PullRequests []PRShort `json:"pull_requests"`
}

type PRShort struct {
	PRID     string `json:"pull_request_id"`
	PRName   string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
	Status   string `json:"status"`
}
