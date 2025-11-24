package service

import (
	"context"
	"errors"

	"github.com/kgugunava/avito-tech-internship/internal/adapters/postgres"
	api_models "github.com/kgugunava/avito-tech-internship/internal/api/models"
)

type TeamService struct {
	teamRepo *postgres.TeamRepository
}

var ErrTeamExists = errors.New("team already exists")


func NewTeamService(teamRepo *postgres.TeamRepository) *TeamService {
	return &TeamService{teamRepo: teamRepo}
}

func (s *TeamService) IsTeamExists(ctx context.Context, teamName string) (bool, error) {
	return s.teamRepo.IsTeamExists(ctx, teamName)
}

func (s *TeamService) CreateNewTeam(ctx context.Context, team api_models.Team) (api_models.Team, error) {
	exists, err := s.teamRepo.IsTeamExists(ctx, team.TeamName)
    if err != nil {
        return api_models.Team{}, err
    }
    if exists {
        return api_models.Team{}, postgres.ErrTeamNotFound
    }

    createdTeam, err := s.teamRepo.CreateTeam(ctx, team)
    if err != nil {
        return api_models.Team{}, err
    }

    return createdTeam, nil
}

func (s* TeamService) GetTeamByName(ctx context.Context, teamName string) (api_models.Team, error) {
    exists, err := s.teamRepo.IsTeamExists(ctx, teamName)
    if err != nil {
        return api_models.Team{}, err
    }

    if !exists {
        return api_models.Team{}, postgres.ErrTeamNotFound
    }

    team, err := s.teamRepo.GetTeamByName(ctx, teamName)
    if err != nil {
        return api_models.Team{}, err
    }

    return team, nil
}