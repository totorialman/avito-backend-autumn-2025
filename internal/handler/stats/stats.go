package stats

import (
	"net/http"

	"github.com/totorialman/avito-backend-autumn-2025/internal/handler"
)

type StatsHandler struct {
	stats statsService
}

func NewStatsHandler(stats statsService) *StatsHandler {
	return &StatsHandler{stats: stats}
}

func (h *StatsHandler) GetReviewerAssignmentStats(w http.ResponseWriter, r *http.Request) {
	domainStats, err := h.stats.GetReviewerAssignmentStats()
	if err != nil {
		handler.HandleError(w, err)
		return
	}

	resp := make([]GetReviewerAssignmentStatsResponse, len(domainStats))
	for i, s := range domainStats {
		resp[i] = GetReviewerAssignmentStatsResponse{
			UserID:      s.UserID,
			Username:    s.Username,
			TeamName:    s.TeamName,
			Assignments: s.Assignments,
		}
	}

	handler.WriteJSON(w, resp, http.StatusOK)
}