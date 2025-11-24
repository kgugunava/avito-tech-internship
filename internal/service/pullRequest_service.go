package service

import (
	"context"
	"time"

	"github.com/kgugunava/avito-tech-internship/internal/adapters/postgres"
	"github.com/kgugunava/avito-tech-internship/internal/api/models"
)

type PullRequestService struct {
	pullRequestRepo *postgres.PullRequestRepository
}

func NewPullRequestService(pullRequestRepo *postgres.PullRequestRepository) *PullRequestService {
	return &PullRequestService{pullRequestRepo: pullRequestRepo}
}

func (s *PullRequestService) Create(ctx context.Context, req models.PullRequestCreatePostRequest) (models.PullRequest, models.ErrorResponse) {
    exists, err := s.pullRequestRepo.PRExists(ctx, req.PullRequestId)
    if err != nil {
        return models.PullRequest{}, models.ErrorResponse{
			Error: models.ErrorResponseError{
				Code: "INERNAL_ERROR",
				Message: err.Error(),
			},
		}
    }
    if exists {
        return models.PullRequest{}, models.ErrorResponse{
			Error: models.ErrorResponseError{
				Code: "PR_EXISTS",
				Message: "PR id already exists",
			},
		}
    }

    teamId, err := s.pullRequestRepo.GetUserTeam(ctx, req.AuthorId)
    if err != nil {
        return models.PullRequest{}, models.ErrorResponse{
			Error: models.ErrorResponseError{
				Code: "AUTHOR_NOT_FOUND",
				Message: "author not found",
			},
		}
    }
    if teamId == 0 {
        return models.PullRequest{}, models.ErrorResponse{
			Error: models.ErrorResponseError{
				Code: "TEAM_NOT_FOUND",
				Message: "team not found",
			},
		}
    }

    reviewers, err := s.pullRequestRepo.GetTeamReviewers(ctx, teamId, req.AuthorId)
    if err != nil {
        return models.PullRequest{}, models.ErrorResponse{
			Error: models.ErrorResponseError{
				Code: "INERNAL_ERROR",
				Message: err.Error(),
			},
		}
    }

    if len(reviewers) > 2 {
        reviewers = reviewers[:2]
    }

    err = s.pullRequestRepo.CreatePR(ctx, req)
    if err != nil {
        return models.PullRequest{}, models.ErrorResponse{
			Error: models.ErrorResponseError{
				Code: "INERNAL_ERROR",
				Message: err.Error(),
			},
		}
    }

    err = s.pullRequestRepo.AssignReviewers(ctx, req.PullRequestId, reviewers)
    if err != nil {
        return models.PullRequest{}, models.ErrorResponse{
			Error: models.ErrorResponseError{
				Code: "INERNAL_ERROR",
				Message: err.Error(),
			},
		}
    }

    return models.PullRequest{
        PullRequestId:     req.PullRequestId,
        PullRequestName:   req.PullRequestName,
        AuthorId:          req.AuthorId,
        Status:            "OPEN",
        AssignedReviewers: reviewers,
    }, models.ErrorResponse{}
}

func (s *PullRequestService) Merge(ctx context.Context, req models.PullRequestMergePostRequest) (models.PullRequest, models.ErrorResponse) {
	mergedAt := time.Now().UTC()

    pr, err := s.pullRequestRepo.SetMerged(ctx, req.PullRequestId, mergedAt)
    if err != nil {
        return pr, models.ErrorResponse{
			Error: models.ErrorResponseError{
				Code: "NOT FOUND",
				Message: err.Error(),
			},
		}
    }

    return pr, models.ErrorResponse{}
}

func (s *PullRequestService) Reassign(ctx context.Context, req models.PullRequestReassignPostRequest) (models.PullRequest, models.ErrorResponse, string) {
	pr, err := s.pullRequestRepo.GetByID(ctx, req.PullRequestId)
    if err != nil {
        return pr, models.ErrorResponse{
			Error: models.ErrorResponseError{
				Code: "NOT_FOUND",
				Message: "pull request not found",
			},
		}, ""
    }

    if pr.Status == "MERGED" {
        return pr, models.ErrorResponse{
			Error: models.ErrorResponseError{
				Code: "PR_MERGED",
				Message: "cannot reassign on merged PR",
			},
		}, ""
    }

    isAssigned := false
    for _, r := range pr.AssignedReviewers {
        if r == req.OldUserId {
            isAssigned = true
            break
        }
    }
    if !isAssigned {
        return pr, models.ErrorResponse{
			Error: models.ErrorResponseError{
				Code: "NOT_ASSIGNED",
				Message: "reviewer is not assigned to this PR",
			},
		}, ""
    }

    teamName, err := s.pullRequestRepo.GetUserTeam(ctx, req.OldUserId)
    if err != nil {
        return pr, models.ErrorResponse{
			Error: models.ErrorResponseError{
				Code: "NOT_FOUND",
				Message: "team not found",
			},
		}, ""
    }

    newReviewer, err := s.pullRequestRepo.FindReplacement(ctx, req.OldUserId, teamName)
    if err != nil {
        return pr, models.ErrorResponse{
			Error: models.ErrorResponseError{
				Code: "NO_CANDIDATE",
				Message: "no active replacement candidate in team",
			},
		}, ""
    }

    err = s.pullRequestRepo.ReplaceReviewer(ctx, req.PullRequestId, req.OldUserId, newReviewer)
    if err != nil {
        return pr, models.ErrorResponse{
			Error: models.ErrorResponseError{
				Code: "INTERNAL_ERROR", 
				Message: err.Error(),
			},
		}, ""
    }

    for i := range pr.AssignedReviewers {
        if pr.AssignedReviewers[i] == req.OldUserId {
            pr.AssignedReviewers[i] = newReviewer
            break
        }
    }

    return pr, models.ErrorResponse{}, newReviewer
}