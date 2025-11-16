package user

import (
	"encoding/json"
	"net/http"

	"github.com/totorialman/avito-backend-autumn-2025/internal/handler"
)

type UserHandler struct {
	users userService
}

func NewUserHandler(users userService) *UserHandler {
	return &UserHandler{
		users: users,
	}
}

func (h *UserHandler) SetActive(w http.ResponseWriter, r *http.Request) {
	var req SetActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.users.SetActive(req.UserID, req.IsActive)
	if err != nil {
		handler.HandleError(w, err)
		return
	}

	resp := SetActiveResponse{
		User: UserResponse{
			UserID:   user.ID,
			Username: user.Username,
			TeamName: user.TeamName,
			IsActive: user.IsActive,
		},
	}

	handler.WriteJSON(w, resp, http.StatusOK)
}

func (h *UserHandler) GetUserReviews(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	prs, err := h.users.GetPRsByReviewer(userID)
	if err != nil {
		handler.HandleError(w, err)
		return
	}

	shortPRs := make([]PRShort, 0, len(prs))
	for _, pr := range prs {
		shortPRs = append(shortPRs, PRShort{
			PRID:     pr.ID,
			PRName:   pr.Name,
			AuthorID: pr.AuthorID,
			Status:   pr.Status,
		})
	}

	resp := GetUserReviewsResponse{
		UserID:       userID,
		PullRequests: shortPRs,
	}

	handler.WriteJSON(w, resp, http.StatusOK)
}
