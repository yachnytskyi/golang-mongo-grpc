package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userRepositoryModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/model"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"

	repositoryUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/utility"
	utility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) user.Repository {
	return &UserRepository{collection: collection}
}

func (userRepository *UserRepository) GetUserById(ctx context.Context, userID string) (*userModel.User, error) {
	userIDMappedToMongoDB, _ := primitive.ObjectIDFromHex(userID)

	var fetchedUser *userModel.User
	query := bson.M{"_id": userIDMappedToMongoDB}
	err := userRepository.collection.FindOne(ctx, query).Decode(&fetchedUser)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &userModel.User{}, err
		}
		return nil, err
	}

	return fetchedUser, nil
}

func (userRepository *UserRepository) GetUserByEmail(ctx context.Context, email string) (*userModel.User, error) {
	var fetchedUser *userModel.User

	query := bson.M{"email": strings.ToLower(email)}
	err := userRepository.collection.FindOne(ctx, query).Decode(&fetchedUser)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &userModel.User{}, err
		}
		return nil, err
	}

	return fetchedUser, nil
}

func (userRepository *UserRepository) Register(ctx context.Context, user *userModel.UserCreate) (*userModel.User, error) {
	userMappedToRepository := userRepositoryModel.UserCreateToUserCreateRepositoryMapper(user)
	userMappedToRepository.CreatedAt = time.Now()
	userMappedToRepository.UpdatedAt = user.CreatedAt
	userMappedToRepository.Email = strings.ToLower(user.Email)
	userMappedToRepository.Verified = true
	userMappedToRepository.Role = "user"
	userMappedToRepository.Password, _ = repositoryUtility.HashPassword(user.Password)

	result, err := userRepository.collection.InsertOne(ctx, &userMappedToRepository)

	if err != nil {
		if err, ok := err.(mongo.WriteException); ok && err.WriteErrors[0].Code == 11000 {
			return nil, errors.New("user with this email already exists")
		}
		return nil, err
	}

	// Create a unique index for the email field.
	option := options.Index()
	option.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: option}

	if _, err := userRepository.collection.Indexes().CreateOne(ctx, index); err != nil {
		return nil, errors.New("could not create an index for an email")
	}

	var createdUser *userModel.User
	query := bson.M{"_id": result.InsertedID}

	err = userRepository.collection.FindOne(ctx, query).Decode(&createdUser)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (userRepository *UserRepository) UpdateUserById(ctx context.Context, userID string, user *userModel.UserUpdate) (*userModel.User, error) {
	userMappedToRepository := userRepositoryModel.UserUpdateToUserUpdateRepositoryMapper(user)
	userMappedToRepository.UpdatedAt = time.Now()

	userMappedToMongoDB, err := utility.MongoMappper(userMappedToRepository)

	if err != nil {
		return &userModel.User{}, err
	}

	userIDMappedToMongoDB, _ := primitive.ObjectIDFromHex(userID)

	query := bson.D{{Key: "_id", Value: userIDMappedToMongoDB}}
	update := bson.D{{Key: "$set", Value: userMappedToMongoDB}}
	result := userRepository.collection.FindOneAndUpdate(ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var updatedUser *userModel.User

	if err := result.Decode(&updatedUser); err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (userRepository *UserRepository) UpdateNewRegisteredUserById(ctx context.Context, userID string, key string, value string) (*userModel.User, error) {
	userIDMappedToMongoDB, _ := primitive.ObjectIDFromHex(userID)

	query := bson.D{{Key: "_id", Value: userIDMappedToMongoDB}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: key, Value: value}}}}
	result, err := userRepository.collection.UpdateOne(ctx, query, update)

	if err != nil {
		return &userModel.User{}, err
	}

	if result.ModifiedCount == 0 {
		return &userModel.User{}, err
	}

	return &userModel.User{}, nil
}

func (userRepository *UserRepository) UpdatePasswordResetTokenUserByEmail(ctx context.Context, email string, firstKey string, firstValue string,
	secondKey string, SecondValue time.Time) error {

	query := bson.D{{Key: "email", Value: strings.ToLower(email)}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: firstKey, Value: firstValue}, {Key: secondKey, Value: secondKey}}}}
	result, err := userRepository.collection.UpdateOne(ctx, query, update)

	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return err
	}

	return nil
}

func (userRepository *UserRepository) ResetUserPassword(ctx context.Context, firstKey string, firstValue string, secondKey string, passwordKey, password string) error {
	hashedPassword, _ := repositoryUtility.HashPassword(password)

	query := bson.D{{Key: firstKey, Value: firstValue}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: passwordKey, Value: hashedPassword}}}, {Key: "$unset", Value: bson.D{{Key: firstKey, Value: ""}, {Key: secondKey, Value: ""}}}}

	result, err := userRepository.collection.UpdateOne(ctx, query, update)

	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return err
	}

	return nil
}

func (userRepository *UserRepository) DeleteUserById(ctx context.Context, userID string) error {
	userIDMappedToMongoDB, _ := primitive.ObjectIDFromHex(userID)

	query := bson.M{"_id": userIDMappedToMongoDB}
	result, err := userRepository.collection.DeleteOne(ctx, query)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return err
	}

	return nil
}
