package domain

import (
	"time"
)

const (
	PRStatusOpen     PRStatus = "OPEN"
	PRStatusInReview PRStatus = "IN_REVIEW"
	PRStatusApproved PRStatus = "APPROVED"
	PRStatusRejected PRStatus = "REJECTED"
	PRStatusMerged   PRStatus = "MERGED"
	PRStatusClosed   PRStatus = "CLOSED"
)

type TeamMember struct {
	TeamName  string       `json:"team_name" db:"team_name"`
	Members   []TeamMember `json:"members"`
	CreatedAt time.Time    `json:"createdAt" db:"created_at"`
}
type PRStatus string
type User struct {
	UserID    string    `json:"user_id" db:"user_id"`
	Username  string    `json:"username" db:"username"`
	TeamName  string    `json:"team_name" db:"team_name"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type Team struct {
	TeamName  string       `json:"team_name" db:"team_name"`
	Members   []TeamMember `json:"members"`
	CreatedAt time.Time    `json:"createdAt" db:"created_at"`
}

type PullRequest struct {
	PullRequestID   string `json:"pull_request_id" db:"pull_request_id"`
	PullRequestName string `json:"pull_request_name" db:"pull_request_name"`
	AuthorID        string `json:"author_id" db:"author_id"`

	Status            PRStatus `json:"status" db:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`

	Reviewer1ID *string `json:"-" db:"reviewer1_id"`
	Reviewer2ID *string `json:"-" db:"reviewer2_id"`

	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
	MergedAt  *time.Time `json:"mergedAt" db:"merged_at"`
}

type UserTeam struct {
	UserID int `db:"user_id"`
	TeamID int `db:"team_id"`
}

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id" db:"pull_request_id"`
	PullRequestName string `json:"pull_request_name" db:"pull_request_name"`
	AuthorID        string `json:"author_id" db:"author_id"`
	Status          string `json:"status" db:"status"`
}
