package post

import (
	"context"

	postModel "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/model"
)

type Repository interface {
	GetPostById(ctx context.Context, postID string) (*postModel.Post, error)
	GetAllPosts(ctx context.Context, page int, limit int) ([]*postModel.Post, error)
	CreatePost(ctx context.Context, user *postModel.PostCreate) (*postModel.Post, error)
	UpdatePostById(ctx context.Context, postID string, post *postModel.PostUpdate) (*postModel.Post, error)
	DeletePostByID(ctx context.Context, postID string) error
}
