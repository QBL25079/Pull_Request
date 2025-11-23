package service

import (
	"context"
	"pr-reviewer-service/internal/domain"
	"pr-reviewer-service/internal/repository"
)

type TeamService struct {
	repo repository.Repository
}

func NewTeamService(repo repository.Repository) *TeamService {
	return &TeamService{repo: repo}
}

func (s *TeamService) CreateTeam(t *domain.Team) error {
	return s.repo.CreateTeam(context.Background(), *t)
}

func (s *TeamService) ListTeams() ([]domain.Team, error) {
	return s.repo.ListTeams(context.Background())
}
