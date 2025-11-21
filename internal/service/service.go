package service

import (
	"context"
	"pr-reviewer-service/internal/domain"
)

type UserServiceProvider interface {
	CreateTeam(ctx context.Context, teamName string) error
	CreateUser(ctx context.Context, user domain.User) error
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
}
