package pr

type CreatePRRequest struct {
	PRID     string `json:"pull_request_id"`
	PRName   string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
}

type MergePRRequest struct {
	PRID string `json:"pull_request_id"`
}

type ReassignPRRequest struct {
	PRID      string `json:"pull_request_id"`
	OldUserID string `json:"old_reviewer_id"`
}

type CreatePRResponse struct {
	PR PullRequestCreate `json:"pr"`
}

type MergePRResponse struct {
	PR PullRequestMerge `json:"pr"`
}

type ReassignPRResponse struct {
	PR         PullRequestCreate `json:"pr"`
	ReplacedBy string            `json:"replaced_by"`
}

type PullRequestMerge struct {
	PRID              string   `json:"pull_request_id"`
	PRName            string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	MergedAt          *string  `json:"mergedAt"`
}

type PullRequestCreate struct {
	PRID              string   `json:"pull_request_id"`
	PRName            string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
}
