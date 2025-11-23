package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"pr-reviewer-service/internal/domain"
)

type Repository interface {
	CreateUser(ctx context.Context, user domain.User) error
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) error
	DeleteUser(ctx context.Context, userID string) error
	ListUsers(ctx context.Context, teamName string) ([]domain.User, error)
	GetPullRequestByID(ctx context.Context, prID string) (*domain.PullRequest, error)
	CreatePullRequest(ctx context.Context, pr domain.PullRequest) error

	CreateTeam(ctx context.Context, team domain.Team) error
	ListTeams(ctx context.Context) ([]domain.Team, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user domain.User) error {
	query := `
		INSERT INTO users (user_id, username, team_name, is_active, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`
	_, err := r.db.ExecContext(ctx, query, user.UserID, user.Username, user.TeamName, user.IsActive)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *repository) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	var u domain.User
	query := `SELECT user_id, username, team_name, is_active, created_at FROM users WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&u.UserID, &u.Username, &u.TeamName, &u.IsActive, &u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &u, nil
}

func (r *repository) UpdateUser(ctx context.Context, user domain.User) error {
	query := `
		UPDATE users SET username = $1, team_name = $2, is_active = $3
		WHERE user_id = $4
	`
	result, err := r.db.ExecContext(ctx, query, user.Username, user.TeamName, user.IsActive, user.UserID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (r *repository) DeleteUser(ctx context.Context, userID string) error {
	query := `DELETE FROM users WHERE user_id = $1`
	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (r *repository) ListUsers(ctx context.Context, teamName string) ([]domain.User, error) {
	query := `SELECT user_id, username, team_name, is_active, created_at FROM users`
	var args []interface{}
	if teamName != "" {
		query += ` WHERE team_name = $1`
		args = append(args, teamName)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.UserID, &u.Username, &u.TeamName, &u.IsActive, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repository) GetPullRequestByID(ctx context.Context, prID string) (*domain.PullRequest, error) {
	var pr domain.PullRequest
	var reviewer1ID, reviewer2ID sql.NullString
	var mergedAt sql.NullTime

	query := `
		SELECT 
			pull_request_id, pull_request_name, author_id, status,
			reviewer1_id, reviewer2_id, created_at, merged_at
		FROM pull_request 
		WHERE pull_request_id = $1
	`

	err := r.db.QueryRowContext(ctx, query, prID).Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&pr.Status,
		&reviewer1ID,
		&reviewer2ID,
		&pr.CreatedAt,
		&mergedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get pull request: %w", err)
	}

	pr.AssignedReviewers = []string{}
	if reviewer1ID.Valid {
		pr.Reviewer1ID = &reviewer1ID.String
		pr.AssignedReviewers = append(pr.AssignedReviewers, reviewer1ID.String)
	}
	if reviewer2ID.Valid {
		pr.Reviewer2ID = &reviewer2ID.String
		pr.AssignedReviewers = append(pr.AssignedReviewers, reviewer2ID.String)
	}
	if mergedAt.Valid {
		pr.MergedAt = &mergedAt.Time
	}

	return &pr, nil
}
func (r *repository) CreatePullRequest(ctx context.Context, pr domain.PullRequest) error {
	var reviewer1ID, reviewer2ID sql.NullString
	if pr.Reviewer1ID != nil && *pr.Reviewer1ID != "" {
		reviewer1ID.Valid = true
		reviewer1ID.String = *pr.Reviewer1ID
	}
	if pr.Reviewer2ID != nil && *pr.Reviewer2ID != "" {
		reviewer2ID.Valid = true
		reviewer2ID.String = *pr.Reviewer2ID
	}

	query := `
		INSERT INTO pull_request (
			pull_request_id, pull_request_name, author_id, status,
			reviewer1_id, reviewer2_id, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`

	_, err := r.db.ExecContext(ctx, query,
		pr.PullRequestID, pr.PullRequestName, pr.AuthorID, string(pr.Status),
		reviewer1ID, reviewer2ID,
	)
	if err != nil {
		return fmt.Errorf("failed to create pull request: %w", err)
	}
	return nil
}
func (r *repository) CreateTeam(ctx context.Context, team domain.Team) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO teams (name) VALUES ($1)`, team.TeamName)
	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}
	return nil
}

func (r *repository) ListTeams(ctx context.Context) ([]domain.Team, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT name FROM teams`)
	if err != nil {
		return nil, fmt.Errorf("failed to list teams: %w", err)
	}
	defer rows.Close()

	var teams []domain.Team
	for rows.Next() {
		var t domain.Team
		if err := rows.Scan(&t.TeamName); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, nil
}
