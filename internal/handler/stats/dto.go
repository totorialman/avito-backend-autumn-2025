package stats

type GetReviewerAssignmentStatsResponse struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	TeamName    string `json:"team_name"`
	Assignments int    `json:"assignments"`
}
