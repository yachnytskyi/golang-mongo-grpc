package model

import (
	"time"

	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
)

// [GET].
type TokenView struct {
	Token string `json:"token"`
}

// [GET].
type UsersView struct {
	UsersView          []UserView                   `json:"users"`
	PaginationResponse httpModel.PaginationResponse `json:"pagination_response"`
}

// [GET].
type UserView struct {
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// [POST].
type UserCreateView struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

// [PUT].
type UserUpdateView struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

// [POST].
type UserLoginView struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// [GET].
type UserTokenView struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// [POST].
type UserForgottenPasswordView struct {
	Email string `json:"email"`
}

// [POST].
type UserResetPasswordView struct {
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

// [GET].
type UserWelcomeMessageView struct {
	Notification string `json:"notification"`
}

func NewWelcomeMessageView(notification string) UserWelcomeMessageView {
	return UserWelcomeMessageView{
		Notification: notification,
	}
}
