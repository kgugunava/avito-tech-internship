package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	api_models "github.com/kgugunava/avito-tech-internship/internal/api/models"
)

type TeamRepository struct {
	pool *pgxpool.Pool
}

var ErrTeamNotFound = errors.New("team not found")

func NewTeamRepository(pool *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{pool: pool}
}

func (r *TeamRepository) IsTeamExists(ctx context.Context, teamName string) (bool, error) {
	var exists bool

	err := r.pool.QueryRow(ctx, `
        SELECT EXISTS(
            SELECT 1 FROM teams WHERE team_name = $1
        )
    `, teamName).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *TeamRepository) CreateTeam(ctx context.Context, team api_models.Team) (api_models.Team, error) {
	var teamId int32
	
	tx, err := r.pool.Begin(ctx)
    if err != nil {
        return api_models.Team{}, err
    }
    defer tx.Rollback(ctx)

    err = tx.QueryRow(ctx,
        `INSERT INTO teams (team_name) VALUES ($1) RETURNING team_id`,
        team.TeamName,
    ).Scan(&teamId)
    if err != nil {
        return api_models.Team{}, err
    }

    for _, m := range team.Members {
        _, err := tx.Exec(ctx,
            `INSERT INTO users (user_id, username, is_active, team_id)
             VALUES ($1, $2, $3, $4)
             ON CONFLICT (user_id)
             DO UPDATE SET username = EXCLUDED.username, is_active = EXCLUDED.is_active, team_id = EXCLUDED.team_id`,
            m.UserId, m.Username, m.IsActive, teamId,
        )
        if err != nil {
            return api_models.Team{}, err
        }
    }

    if err := tx.Commit(ctx); err != nil {
        return api_models.Team{}, err
    }

    return team, nil
}

func (r *TeamRepository) GetTeamByName(ctx context.Context, teamName string) (api_models.Team, error) {
	var team api_models.Team
	var teamId int

    err := r.pool.QueryRow(ctx,
        `SELECT team_id, team_name 
         FROM teams 
         WHERE team_name = $1`,
        teamName,
    ).Scan(&teamId, &team.TeamName)

    if err != nil {
        return api_models.Team{}, err
    }

    rows, err := r.pool.Query(ctx,
        `SELECT user_id, username, is_active
         FROM users
         WHERE team_id = $1`,
        teamId,
    )
    if err != nil {
        return api_models.Team{}, err
    }
    defer rows.Close()

    members := []api_models.TeamMember{}

    for rows.Next() {
        var m api_models.TeamMember
        if err := rows.Scan(&m.UserId, &m.Username, &m.IsActive); err != nil {
            return api_models.Team{}, err
        }
        members = append(members, m)
    }

    team.Members = members

    return team, nil
}