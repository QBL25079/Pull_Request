package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"pr-reviewer-service/internal/domain"
	"time"

	"github.com/lib/pq"
)

type Repository interface {
	CreateTeam(ctx context.Context, teamName string) error
	UpdateTeamName(ctx context.Context, oldName, newName string) error
	CreateUser(ctx context.Context, user domain.User) error
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
	CreatePullRequest(ctx context.Context, pr domain.PullRequest) error
	GetPullRequestByID(ctx context.Context, prID string) (*domain.PullRequest, error)
}

type postgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepository(db *sql.DB) *postgresRepository {
	return &postgresRepository{DB: db}
}

func (r *postgresRepository) CreateTeam(ctx context.Context, teamName string) error {
	query := `INSERT  INTO team  (team_name) VALUES ($1)`

	_, err := r.DB.ExecContext(ctx, query, teamName)

	if err != nil {
		return fmt.Errorf("repository: failed to create team %s: %w", teamName, err)
	}

	return nil
}

func (r *postgresRepository) UpdateTeamName(ctx context.Context, oldName, newName string) error {
	query := `UPDATE team SET team_name = $2 WHERE team_name = $1`

	result, err := r.DB.ExecContext(ctx, query, oldName, newName)

	if err != nil {
		return fmt.Errorf("repository: error to update team name from %s to %s: %w", oldName, newName, err)
	}
	rAffected, _ := result.RowsAffected()
	if rAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *postgresRepository) CreateUser(ctx context.Context, user domain.User) error {
	query := `
		INSERT INTO users (user_id, username, team_name, is_active, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`

	_, err := r.DB.ExecContext(ctx, query,
		user.UserID,
		user.Username,
		user.TeamName,
		user.IsActive,
	)
	if err != nil {
		return fmt.Errorf("repository: failed to create user %s: %w", user.UserID, err)
	}
	return nil
}

func (r *postgresRepository) CreatePullRequest(ctx context.Context, pr domain.PullRequest) error {

	query := `INSERT INTO pull_request (pull_request_id, pull_request_name, author_id, status, reviewer1_id, reviewer2_id) VALUES ($1, $2, $3, $4, $5, $6)`

	var r1, r2 sql.NullString

	if pr.Reviewer1ID != nil && *pr.Reviewer1ID != "" {
		r1.String = *pr.Reviewer1ID
		r1.Valid = true
	}
	if pr.Reviewer2ID != nil && *pr.Reviewer2ID != "" {
		r2.String = *pr.Reviewer2ID
		r2.Valid = true
	}

	_, err := r.DB.ExecContext(ctx, query, pr.Reviewer1ID, pr.Reviewer2ID, pr.AuthorID, string(pr.Status), r1, r2)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			return fmt.Errorf("repository: PostgreSQL error (%s) when creating PR %s: %w", pgErr.Code, pr.PullRequestID, err)
		}
		return fmt.Errorf("repository: failed to create pull request %s: %w", pr.PullRequestID, err)
	}
	return nil
}

func (r *postgresRepository) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	var u domain.User
	query := `SELECT user_id, username, team_name, is_active, created_at 
			  FROM users WHERE user_id = $1`
	err := r.DB.QueryRowContext(ctx, query, userID).Scan(
		&u.UserID, &u.Username, &u.TeamName, &u.IsActive, &u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("repository: failed to get user %s: %w", userID, err)
	}
	return &u, nil
}

func (r *postgresRepository) GetPullRequestByID(ctx context.Context, prID string) (*domain.PullRequest, error) {
	pr := domain.PullRequest{}

	var r1ID, r2ID sql.NullString
	var mergedAt sql.NullTime
	var createdAt time.Time
	var status string

	query := `SELECT pull_request_id, pull_request_name, author_id, status, reviewer1_id, reviewer2_id, created_at, merged_at FROM pull_request WHERE pull_request_id = $1`

	row := r.DB.QueryRowContext(ctx, query, prID)

	err := row.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &status, &r1ID, &r2ID, &createdAt, &mergedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("repository: failed to get pull request %s: %w", prID, err)
	}

	pr.Status = domain.PRStatus(status)
	pr.CreatedAt = &createdAt

	if r1ID.Valid {
		pr.Reviewer1ID = &r1ID.String
		pr.AssignedReviewers = append(pr.AssignedReviewers, r1ID.String)
	}
	if r2ID.Valid {
		pr.Reviewer2ID = &r2ID.String
		pr.AssignedReviewers = append(pr.AssignedReviewers, r2ID.String)
	}

	if mergedAt.Valid {
		pr.MergedAt = &mergedAt.Time
	} else {
		pr.MergedAt = nil
	}

	return &pr, nil
}
