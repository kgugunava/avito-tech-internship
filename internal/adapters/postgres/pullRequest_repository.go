package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"

	"github.com/kgugunava/avito-tech-internship/internal/api/models"
)

type PullRequestRepository struct {
	pool *pgxpool.Pool
}

func NewPullRequestRepository(pool *pgxpool.Pool) *PullRequestRepository {
	return &PullRequestRepository{pool: pool}
}

func (r *PullRequestRepository) PRExists(ctx context.Context, prId string) (bool, error) {
	var exists bool
    err := r.pool.QueryRow(ctx,
        `SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1)`,
        prId,
    ).Scan(&exists)
    return exists, err
}

func (r *PullRequestRepository) GetUserTeam(ctx context.Context, userId string) (int, error) {
	var teamId *int

    err := r.pool.QueryRow(ctx,
        `SELECT team_id FROM users WHERE user_id = $1`,
        userId,
    ).Scan(&teamId)

    if errors.Is(err, pgx.ErrNoRows) {
        return 0, errors.New("author not found")
    }

    if err != nil {
        return 0, err
    }

    return *teamId, nil
}

func (r *PullRequestRepository) GetTeamReviewers(ctx context.Context, teamId int, userId string) ([]string, error) {
    rows, err := r.pool.Query(ctx,
        `SELECT user_id FROM users
         WHERE team_id = $1 AND user_id != $2 AND is_active = TRUE
         LIMIT 2`,
        teamId, userId,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var reviewers []string
    for rows.Next() {
        var u string
        rows.Scan(&u)
        reviewers = append(reviewers, u)
    }

    return reviewers, nil
}

func (r *PullRequestRepository) CreatePR(ctx context.Context, req models.PullRequestCreatePostRequest) error {
    _, err := r.pool.Exec(ctx,
        `INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status)
         VALUES ($1, $2, $3, 'OPEN')`,
        req.PullRequestId,
        req.PullRequestName,
        req.AuthorId,
    )
    return err
}

func (r *PullRequestRepository) AssignReviewers(ctx context.Context, prId string, reviewers []string) error {
    for _, rID := range reviewers {
        _, err := r.pool.Exec(ctx,
            `INSERT INTO pull_request_reviewers (pull_request_id, reviewer_id) VALUES ($1, $2)`,
            prId, rID,
        )
        if err != nil {
            return err
        }
    }
    return nil
}

func (r *PullRequestRepository) GetByID(ctx context.Context, prID string) (models.PullRequest, error) {
    var pr models.PullRequest

    query := `
        SELECT pull_request_id,
               pull_request_name,
               author_id,
               status,
               merged_at
        FROM pull_requests
        WHERE pull_request_id = $1
    `

    err := r.pool.QueryRow(ctx, query, prID).Scan(
        &pr.PullRequestId,
        &pr.PullRequestName,
        &pr.AuthorId,
        &pr.Status,
        &pr.MergedAt,
    )
    if err != nil {
        return pr, errors.New("pull request not found")
    }

    reviewersRows, err := r.pool.Query(ctx,
        `SELECT reviewer_id FROM pull_request_reviewers WHERE pull_request_id = $1`,
        prID,
    )
    if err != nil {
        return pr, err
    }
    defer reviewersRows.Close()

    var reviewers []string
    for reviewersRows.Next() {
        var id string
        reviewersRows.Scan(&id)
        reviewers = append(reviewers, id)
    }

    pr.AssignedReviewers = reviewers

    return pr, nil
}

func (r *PullRequestRepository) SetMerged(ctx context.Context, prID string, mergedAt time.Time) (models.PullRequest, error) {
    pr, err := r.GetByID(ctx, prID)
    if err != nil {
        return pr, err
    }

    if pr.Status == "MERGED" {
        return pr, nil
    }

     query := `
        UPDATE pull_requests
        SET status = 'MERGED',
            merged_at = $1
        WHERE pull_request_id = $2
    `

    _, err = r.pool.Exec(ctx, query, mergedAt, prID)
    if err != nil {
        return pr, err
    }

    pr.Status = "MERGED"
    pr.MergedAt = &mergedAt

    return pr, nil
}

func (r *PullRequestRepository) FindReplacement(ctx context.Context, oldReviewerID string, teamName int) (string, error) {
    var newUserId string

    query := `
        SELECT user_id
        FROM users
        WHERE team_id = $1
          AND user_id <> $2
          AND is_active = TRUE
        LIMIT 1
    `

    err := r.pool.QueryRow(ctx, query, teamName, oldReviewerID).Scan(&newUserId)
    if err != nil {
        return "", errors.New("no active replacement candidate in team")
    }

    return newUserId, nil
}

func (r *PullRequestRepository) ReplaceReviewer(ctx context.Context, prID, oldReviewer, newReviewer string) error {
    query := `
        UPDATE pull_request_reviewers
        SET reviewer_id = $1
        WHERE pull_request_id = $2 AND reviewer_id = $3
    `
    result, err := r.pool.Exec(ctx, query, newReviewer, prID, oldReviewer)
    if err != nil {
        return err
    }
    if result.RowsAffected() == 0 {
        return errors.New("reviewer is not assigned to this PR")
    }

    return nil
}