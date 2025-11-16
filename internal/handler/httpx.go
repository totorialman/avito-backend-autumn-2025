package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/totorialman/avito-backend-autumn-2025/internal/model"
)

func WriteJSON(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func HandleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	var (
		code    string
		message string
		status  int = http.StatusInternalServerError
	)

	switch {
	case errors.Is(err, model.ErrTeamExists):
		code = "TEAM_EXISTS"
		message = "team_name already exists"
		status = http.StatusBadRequest

	case errors.Is(err, model.ErrPRExists):
		code = "PR_EXISTS"
		message = "PR id already exists"
		status = http.StatusConflict

	case errors.Is(err, model.ErrPRMerged):
		code = "PR_MERGED"
		message = "cannot reassign on merged PR"
		status = http.StatusConflict

	case errors.Is(err, model.ErrNotAssigned):
		code = "NOT_ASSIGNED"
		message = "reviewer is not assigned to this PR"
		status = http.StatusConflict

	case errors.Is(err, model.ErrNoCandidate):
		code = "NO_CANDIDATE"
		message = "no active replacement candidate in team"
		status = http.StatusConflict

	case errors.Is(err, model.ErrNotFound):
		code = "NOT_FOUND"
		message = "resource not found"
		status = http.StatusNotFound

	default:
		code = "INTERNAL_ERROR"
		message = "internal server error"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := ErrorResponse{}
	resp.Error.Code = code
	resp.Error.Message = message

	_ = json.NewEncoder(w).Encode(resp)
}
