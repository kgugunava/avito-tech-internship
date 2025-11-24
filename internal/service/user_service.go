package service

import (
	"context"

	"github.com/kgugunava/avito-tech-internship/internal/adapters/postgres"
	"github.com/kgugunava/avito-tech-internship/internal/api/models"
)

type UserService struct {
	userRepo *postgres.UserRepository
}

func NewUserService(userRepo *postgres.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetReviews(ctx context.Context, userId string) (models.UsersGetReviewGet200Response, models.ErrorResponse) {
	_, err := s.userRepo.GetUserByID(ctx, userId)
	if err != nil {
		return models.UsersGetReviewGet200Response{}, models.ErrorResponse{
			Error: models.ErrorResponseError{
				Code: "USER_NOT_FOUND",
				Message: "user not found",
			},
		}
	}

	pullRequests, err := s.userRepo.GetReviews(ctx, userId)
	if err != nil {
		return models.UsersGetReviewGet200Response{}, models.ErrorResponse{
			Error: models.ErrorResponseError{
				Code: "INTERNAL_ERROR",
				Message: err.Error(),
			},
		}
	}
	
	return models.UsersGetReviewGet200Response{
		UserId: userId,
		PullRequests: pullRequests,
	}, models.ErrorResponse{}
}

func (s *UserService) SetIsActivePost(ctx context.Context, usersSetIsActiveRequest models.UsersSetIsActivePostRequest) (models.User, error) {
	return s.userRepo.SetIsActivePost(ctx, usersSetIsActiveRequest)
}