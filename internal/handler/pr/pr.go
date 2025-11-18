package pr

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/totorialman/avito-backend-autumn-2025/internal/handler"
)

type PrHandler struct {
	pullRequests pullRequestService
}

func NewPrHandler(pullRequests pullRequestService) *PrHandler {
	return &PrHandler{
		pullRequests: pullRequests,
	}
}

func (h *PrHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	var req CreatePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	log.Printf("CreatePR request: PRID=%s PRName=%s AuthorID=%s", req.PRID, req.PRName, req.AuthorID)

	pr, err := h.pullRequests.CreatePR(req.PRID, req.PRName, req.AuthorID)
	if err != nil {
		handler.HandleError(w, err)
		return
	}

	resp := CreatePRResponse{
		PR: PullRequestCreate{
			PRID:              pr.ID,
			PRName:            pr.Name,
			AuthorID:          pr.AuthorID,
			Status:            pr.Status,
			AssignedReviewers: pr.AssignedReviewers,
		},
	}

	handler.WriteJSON(w, resp, http.StatusCreated)
}

func (h *PrHandler) MergePR(w http.ResponseWriter, r *http.Request) {
	var req MergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	log.Printf("MergePR request: PRID=%s", req.PRID)

	pr, err := h.pullRequests.MergePR(req.PRID)
	if err != nil {
		handler.HandleError(w, err)
		return
	}

	var mergedAt *string
	if pr.MergedAt != nil {
		s := pr.MergedAt.Format(time.RFC3339)
		mergedAt = &s
	}

	resp := MergePRResponse{
		PR: PullRequestMerge{
			PRID:              pr.ID,
			PRName:            pr.Name,
			AuthorID:          pr.AuthorID,
			Status:            pr.Status,
			AssignedReviewers: pr.AssignedReviewers,
			MergedAt:          mergedAt,
		},
	}

	handler.WriteJSON(w, resp, http.StatusOK)
}

func (h *PrHandler) Reassign(w http.ResponseWriter, r *http.Request) {
	var req ReassignPRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	log.Printf("Reassign request: PRID=%s OldUserID=%s", req.PRID, req.OldUserID)

	pr, newReviewer, err := h.pullRequests.ReassignReviewer(req.PRID, req.OldUserID)
	if err != nil {
		handler.HandleError(w, err)
		return
	}

	resp := ReassignPRResponse{
		PR: PullRequestCreate{
			PRID:              pr.ID,
			PRName:            pr.Name,
			AuthorID:          pr.AuthorID,
			Status:            pr.Status,
			AssignedReviewers: pr.AssignedReviewers,
		},
		ReplacedBy: newReviewer,
	}

	handler.WriteJSON(w, resp, http.StatusOK)
}
