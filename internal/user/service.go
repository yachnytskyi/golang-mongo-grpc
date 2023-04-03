package user

import (
	"context"

	"github.com/yachnytskyi/golang-mongo-grpc/models"
)

type Service interface {
	Register(ctx context.Context, user *models.UserCreate) (*models.UserFullResponse, error)
	UserGetById(ctx context.Context, userID string) (*models.UserFullResponse, error)
	UserGetByEmail(ctx context.Context, email string) (*models.UserFullResponse, error)
}
