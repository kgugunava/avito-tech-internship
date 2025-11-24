package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kgugunava/avito-tech-internship/internal/api/models"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) GetReviews(ctx context.Context, userId string) ([] models.PullRequestShort, error) {
	var prId int
	
	query := `
        SELECT pr.pr_id, pr.name, pr.author_id, pr.status
        FROM pull_request pr
        JOIN reviewers rev ON pr.pr_id = rev.pr_id
        JOIN users u ON rev.user_id = u.user_id
        WHERE u.user_id = $1
    `

    rows, err := r.pool.Query(ctx, query, userId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var reviews []models.PullRequestShort
    for rows.Next() {
        var pr models.PullRequestShort
        err := rows.Scan(&prId, &pr.PullRequestName, &pr.AuthorId, &pr.Status)
        if err != nil {
            return nil, err
        }
        reviews = append(reviews, pr)
    }

    if rows.Err() != nil {
        return nil, rows.Err()
    }

    return reviews, nil
}

func (r *UserRepository) SetIsActivePost(ctx context.Context, usersSetIsActiveRequest models.UsersSetIsActivePostRequest) (models.User, error) {
	var user models.User

	query := `
        UPDATE users
        SET is_active = $1
        WHERE user_id = $2
        RETURNING user_id, username, team_name, is_active;
    `

    err := r.pool.QueryRow(
        ctx,
        query,
        usersSetIsActiveRequest.IsActive,
        usersSetIsActiveRequest.UserId,
    ).Scan(
        &user.UserId,
        &user.Username,
        &user.TeamName,
        &user.IsActive,
    )

    if err != nil {
        return models.User{}, err
    }

    return user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userId string) (models.User, error) {
    var user models.User

    query := `
        SELECT user_id, username, team_id, is_active
        FROM users
        WHERE user_id = $1
    `
    row := r.pool.QueryRow(ctx, query, userId)

    var teamID *int

    err := row.Scan(
        &user.UserId,
        &user.Username,
        &teamID,
        &user.IsActive,
    )

    if err != nil {
        return models.User{}, err
    }

    return user, nil
}
