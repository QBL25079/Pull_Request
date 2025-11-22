// internal/service/service.go
package service

import (
	"context"
	"pr-reviewer-service/internal/domain"
	"pr-reviewer-service/internal/repository"
)

type userService struct {
	repo repository.Repository
}

func New(repo repository.Repository) UserServiceProvider {
	return &userService{repo: repo}
}

func (s *userService) CreateTeam(ctx context.Context, teamName string) error {
	return s.repo.CreateTeam(ctx, teamName)
}

func (s *userService) UpdateTeamName(ctx context.Context, oldName, newName string) error {
	return s.repo.UpdateTeamName(ctx, oldName, newName)
}

func (s *userService) CreateUser(ctx context.Context, user domain.User) error {
	return s.repo.CreateUser(ctx, user)
}

func (s *userService) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *userService) CreatePullRequest(ctx context.Context, pr domain.PullRequest) error {
	return s.repo.CreatePullRequest(ctx, pr)
}

func (s *userService) GetPullRequestByID(ctx context.Context, prID string) (*domain.PullRequest, error) {
	return s.repo.GetPullRequestByID(ctx, prID)
}
