package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostRepository struct {
	PostID    primitive.ObjectID `bson:"_id"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Title     string             `bson:"title"`
	Content   string             `bson:"content"`
	Image     string             `bson:"image"`
	User      string             `bson:"user"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type PostCreateRepository struct {
	UserID    primitive.ObjectID `bson:"user_id"`
	Title     string             `bson:"title"`
	Content   string             `bson:"content"`
	Image     string             `bson:"image"`
	User      string             `bson:"user"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type PostUpdateRepository struct {
	PostID    primitive.ObjectID `bson:"_id"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Title     string             `bson:"title"`
	Content   string             `bson:"content"`
	Image     string             `bson:"image"`
	User      string             `bson:"user"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}
