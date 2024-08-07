package model

import (
	"time"

	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/domain"
)

type Users struct {
	Users              []User
	PaginationResponse common.PaginationResponse
}

type User struct {
	model.BaseEntity
	Name     string
	Email    string
	Password string
	Role     string
	Verified bool
}

type UserCreate struct {
	Name             string
	Email            string
	Password         string
	PasswordConfirm  string
	Role             string
	Verified         bool
	VerificationCode string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type UserUpdate struct {
	ID        string
	Name      string
	UpdatedAt time.Time
}

type UserLogin struct {
	Email    string
	Password string
}

type UserToken struct {
	AccessToken  string
	RefreshToken string
}

type UserForgottenPassword struct {
	Email       string
	ResetToken  string
	ResetExpiry time.Time
}

type UserResetPassword struct {
	ResetToken      string
	Password        string
	PasswordConfirm string
}

type UserResetExpiry struct {
	ResetExpiry time.Time
}

func NewUsers(users []User, paginationResponse common.PaginationResponse) Users {
	return Users{
		Users:              users,
		PaginationResponse: paginationResponse,
	}
}

func NewUser(id string, createdAt, updatedAt time.Time, name, email, password, role string, verified bool) User {
	return User{
		BaseEntity: model.NewBaseEntity(id, createdAt, updatedAt),
		Name:       name,
		Email:      email,
		Password:   password,
		Role:       role,
		Verified:   verified,
	}
}

func NewUserCreate(name, email, password, passwordConfirm string) UserCreate {
	return UserCreate{
		Name:            name,
		Email:           email,
		Password:        password,
		PasswordConfirm: passwordConfirm,
	}
}

func NewUserUpdate(id, name string) UserUpdate {
	return UserUpdate{
		ID:   id,
		Name: name,
	}
}

func NewUserLogin(email, password string) UserLogin {
	return UserLogin{
		Email:    email,
		Password: password,
	}
}

func NewUserToken(accessToken, refreshToken string) UserToken {
	return UserToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func NewUserForgottenPassword(email string) UserForgottenPassword {
	return UserForgottenPassword{
		Email: email,
	}
}

func NewUserResetPassword(resetToken, password, passwordConfirm string) UserResetPassword {
	return UserResetPassword{
		ResetToken:      resetToken,
		Password:        password,
		PasswordConfirm: passwordConfirm,
	}
}

func NewUserResetExpiry(resetExpiry time.Time) UserResetExpiry {
	return UserResetExpiry{
		ResetExpiry: resetExpiry,
	}
}
