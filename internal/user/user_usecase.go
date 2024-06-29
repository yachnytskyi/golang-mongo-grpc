package user

import (
	"context"

	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
)

type UserUseCase interface {
	GetAllUsers(ctx context.Context, paginationQuery commonModel.PaginationQuery) commonModel.Result[userModel.Users]
	GetUserById(ctx context.Context, userID string) commonModel.Result[userModel.User]
	GetUserByEmail(ctx context.Context, email string) commonModel.Result[userModel.User]
	Register(ctx context.Context, user userModel.UserCreate) commonModel.Result[userModel.User]
	UpdateCurrentUser(ctx context.Context, user userModel.UserUpdate) commonModel.Result[userModel.User]
	DeleteUserById(ctx context.Context, userID string) error
	Login(ctx context.Context, user userModel.UserLogin) commonModel.Result[userModel.UserToken]
	RefreshAccessToken(ctx context.Context, user userModel.User) commonModel.Result[userModel.UserToken]
	ForgottenPassword(ctx context.Context, userForgottenPasswordView userModel.UserForgottenPassword) error
	ResetUserPassword(ctx context.Context, firstKey string, firstValue string, secondKey, passwordKey, password string) error
}
