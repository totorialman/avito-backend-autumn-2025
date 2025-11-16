package team

import (
	"encoding/json"
	"net/http"

	"github.com/totorialman/avito-backend-autumn-2025/internal/handler"
	"github.com/totorialman/avito-backend-autumn-2025/internal/model"
)

type TeamHandler struct {
	teams teamService
}

func NewTeamHandler(teams teamService) *TeamHandler {
	return &TeamHandler{
		teams: teams,
	}
}

func (h *TeamHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
	var req AddTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	var members []model.User
	for _, m := range req.Members {
		members = append(members, model.User{
			ID:       m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}

	err := h.teams.CreateTeam(req.TeamName, members)
	if err != nil {
		handler.HandleError(w, err)
		return
	}

	resp := AddTeamResponse{
		Team: GetTeamResponse(req),
	}

	handler.WriteJSON(w, resp, http.StatusCreated)

}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		http.Error(w, "team_name is required", http.StatusBadRequest)
		return
	}

	team, err := h.teams.GetTeam(teamName)
	if err != nil {
		handler.HandleError(w, err)
		return
	}

	members := make([]TeamMember, 0)
	for _, u := range team.Members {
		members = append(members, TeamMember{
			UserID:   u.ID,
			Username: u.Username,
			IsActive: u.IsActive,
		})
	}

	resp := GetTeamResponse{
		TeamName: team.Name,
		Members:  members,
	}

	handler.WriteJSON(w, resp, http.StatusOK)
}
