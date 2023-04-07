package user

import (
	"context"

	"github.com/yachnytskyi/golang-mongo-grpc/models"
)

type Service interface {
	Register(ctx context.Context, user *models.UserCreate) (*models.UserFullResponse, error)
	Login(ctx context.Context, user *models.UserSignIn) (*models.UserFullResponse, error)
	UpdateNewRegisteredUserById(ctx context.Context, userID string, key string, value string) (*models.UserFullResponse, error)
	UpdateUserById(ctx context.Context, userID string, user *models.UserUpdate) (*models.UserFullResponse, error)
	GetUserById(ctx context.Context, userID string) (*models.UserFullResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserFullResponse, error)
	DeleteUserById(ctx context.Context, userID string) error
}
