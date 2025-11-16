package team

type AddTeamRequest struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

type AddTeamResponse struct {
	Team GetTeamResponse `json:"team"`
}

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type GetTeamResponse struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}
