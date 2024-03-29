package model

import (
	"time"

	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UsersRepository struct {
	Users              []UserRepository
	PaginationResponse commonModel.PaginationResponse
}

type UserRepository struct {
	UserID    primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	Role      string             `bson:"role"`
	Verified  bool               `bson:"verified"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type UserCreateRepository struct {
	Name             string    `bson:"name"`
	Email            string    `bson:"email"`
	Password         string    `bson:"password"`
	Role             string    `bson:"role"`
	Verified         bool      `bson:"verified"`
	VerificationCode string    `bson:"verification_code"`
	CreatedAt        time.Time `bson:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at"`
}

type UserUpdateRepository struct {
	Name      string    `bson:"name"`
	UpdatedAt time.Time `bson:"updated_at"`
}
