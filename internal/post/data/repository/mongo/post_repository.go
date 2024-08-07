package repository

import (
	"context"
	"errors"
	"time"

	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	repository "github.com/yachnytskyi/golang-mongo-grpc/internal/post/data/repository/mongo/model"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/model"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/data/repository/mongo"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	location = "post.data.repository.mongo."
)

type PostRepository struct {
	Logger     interfaces.Logger
	collection *mongo.Collection
}

func NewPostRepository(logger interfaces.Logger, db *mongo.Database) interfaces.PostRepository {
	return &PostRepository{
		Logger:     logger,
		collection: db.Collection("posts"),
	}
}

func (postRepository *PostRepository) GetAllPosts(ctx context.Context, page int, limit int) (*post.Posts, error) {
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 100
	}

	skip := (page - 1) * limit

	option := options.FindOptions{}
	option.SetLimit(int64(limit))
	option.SetSkip(int64(skip))

	query := bson.M{}
	cursor, err := postRepository.collection.Find(ctx, query, &option)

	if validator.IsError(err) {
		return nil, err
	}

	defer cursor.Close(ctx)

	var fetchedPosts []*repository.PostRepository

	for cursor.Next(ctx) {
		post := &repository.PostRepository{}
		err := cursor.Decode(post)

		if validator.IsError(err) {
			return nil, err
		}

		fetchedPosts = append(fetchedPosts, post)
	}

	err = cursor.Err()
	if validator.IsError(err) {
		return nil, err
	}

	if len(fetchedPosts) == 0 {
		return &post.Posts{
			Posts: make([]*post.Post, 0),
		}, nil
	}

	return &post.Posts{
		Posts: repository.PostsRepositoryToPostsMapper(fetchedPosts),
	}, nil
}

func (postRepository *PostRepository) GetPostById(ctx context.Context, postID string) (*post.Post, error) {
	postObjectID := model.HexToObjectIDMapper(postRepository.Logger, location+"GetPostById", postID)
	if validator.IsError(postObjectID.Error) {
		return nil, postObjectID.Error
	}

	query := bson.M{"_id": postObjectID.Data}
	var fetchedPost *repository.PostRepository
	err := postRepository.collection.FindOne(ctx, query).Decode(&fetchedPost)
	if validator.IsError(err) {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no document with that Id exists")
		}

		return nil, err
	}

	return repository.PostRepositoryToPostMapper(fetchedPost), nil
}

func (postRepository *PostRepository) CreatePost(ctx context.Context, post *post.PostCreate) (*post.Post, error) {
	postMappedToRepository, postCreateToPostCreateRepositoryMapperError := repository.PostCreateToPostCreateRepositoryMapper(postRepository.Logger, post)
	if validator.IsError(postCreateToPostCreateRepositoryMapperError) {
		return nil, postCreateToPostCreateRepositoryMapperError
	}
	postMappedToRepository.CreatedAt = time.Now()
	postMappedToRepository.UpdatedAt = post.CreatedAt

	result, err := postRepository.collection.InsertOne(ctx, postMappedToRepository)
	if validator.IsError(err) {
		er, ok := err.(mongo.WriteException)
		if ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("post with that title already exists")
		}
		return nil, err
	}

	option := options.Index()
	option.SetUnique(true)

	index := mongo.IndexModel{Keys: bson.M{"title": 1}, Options: option}
	_, err = postRepository.collection.Indexes().CreateOne(ctx, index)
	if validator.IsError(err) {
		return nil, errors.New("could not create an index for a title")
	}

	var createdPost *repository.PostRepository
	query := bson.M{"_id": result.InsertedID}
	err = postRepository.collection.FindOne(ctx, query).Decode(&createdPost)
	if validator.IsError(err) {
		return nil, err
	}

	return repository.PostRepositoryToPostMapper(createdPost), nil
}

func (postRepository *PostRepository) UpdatePostById(ctx context.Context, postID string, postUpdate *post.PostUpdate) (*post.Post, error) {
	postUpdateRepository, postUpdateToPostUpdateRepositoryMapper := repository.PostUpdateToPostUpdateRepositoryMapper(postRepository.Logger, postUpdate)
	if validator.IsError(postUpdateToPostUpdateRepositoryMapper) {
		return nil, postUpdateToPostUpdateRepositoryMapper
	}

	postUpdateRepository.UpdatedAt = time.Now()

	// Map the user update repository to a BSON document for MongoDB update.
	postUpdateBson := model.DataToMongoDocumentMapper(postRepository.Logger, location+"UpdatePostById", postUpdateRepository)
	if validator.IsError(postUpdateBson.Error) {
		return nil, postUpdateBson.Error
	}

	postObjectID := model.HexToObjectIDMapper(postRepository.Logger, location+"UpdatePostById", postID)
	if validator.IsError(postObjectID.Error) {
		return nil, postObjectID.Error
	}

	query := bson.D{{Key: "_id", Value: postObjectID.Data}}
	update := bson.D{{Key: "$set", Value: postUpdateBson.Data}}
	result := postRepository.collection.FindOneAndUpdate(ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var updatedPost *post.Post
	err := result.Decode(&updatedPost)
	if validator.IsError(err) {
		return nil, errors.New("sorry, but this title already exists. Please choose another one")
	}

	return updatedPost, nil
}

func (postRepository *PostRepository) DeletePostByID(ctx context.Context, postID string) error {
	postObjectID := model.HexToObjectIDMapper(postRepository.Logger, location+"GetPostById", postID)
	if validator.IsError(postObjectID.Error) {
		return postObjectID.Error
	}

	query := bson.M{"_id": postObjectID.Data}
	result, err := postRepository.collection.DeleteOne(ctx, query)
	if validator.IsError(err) {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("no document with that Id exists")
	}

	return nil
}
