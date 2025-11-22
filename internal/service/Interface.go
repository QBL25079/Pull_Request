package service

import (
	"context"
	"pr-reviewer-service/internal/domain"
)

type UserServiceProvider interface {
	CreateTeam(ctx context.Context, teamName string) error
	UpdateTeamName(ctx context.Context, oldName, newName string) error

	CreateUser(ctx context.Context, user domain.User) error
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)

	CreatePullRequest(ctx context.Context, pr domain.PullRequest) error
	GetPullRequestByID(ctx context.Context, prID string) (*domain.PullRequest, error)
}
