package repository

import (
	"context"
	"strings"
	"time"

	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userRepositoryMail "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/external/mail"
	userRepositoryModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/mongo/model"

	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	domainUseCaseValidator "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/usecase"

	repositoryUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/utility"
	mongoMapper "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/data/repository/mongo"

	"github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"

	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	location                = "User.Data.Repository.MongoDB."
	updateIsNotSuccessful   = "update was not successful"
	delitionIsNotSuccessful = "delition was not successful"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) user.Repository {
	return &UserRepository{collection: collection}
}

func (userRepository *UserRepository) GetAllUsers(ctx context.Context, page int, limit int) (*userModel.Users, error) {
	if page == 0 || page < 0 || page > 100 {
		page = 1
	}
	if limit == 0 || limit < 0 || limit > 100 {
		limit = 10
	}
	skip := (page - 1) * limit
	option := options.FindOptions{}
	option.SetLimit(int64(limit))
	option.SetSkip(int64(skip))

	query := bson.M{}
	cursor, cursorFindError := userRepository.collection.Find(ctx, query, &option)
	if validator.IsErrorNotNil(cursorFindError) {
		getAllUsersEntityNotFoundError := domainError.NewEntityNotFoundError(location+"GetAllUsers.Find", cursorFindError.Error())
		logging.Logger(getAllUsersEntityNotFoundError)
		return nil, getAllUsersEntityNotFoundError
	}
	defer cursor.Close(ctx)

	var fetchedUsers = make([]*userRepositoryModel.UserRepository, 0, limit)
	for cursor.Next(ctx) {
		user := &userRepositoryModel.UserRepository{}
		cursorDecodeError := cursor.Decode(user)
		if validator.IsErrorNotNil(cursorDecodeError) {
			fetchedUserInternalError := domainError.NewInternalError(location+"GetAllUsers.cursor.decode", cursorDecodeError.Error())
			logging.Logger(fetchedUserInternalError)
			return nil, fetchedUserInternalError
		}
		fetchedUsers = append(fetchedUsers, user)
	}
	cursorError := cursor.Err()
	if validator.IsErrorNotNil(cursorError) {
		cursorInternalError := domainError.NewInternalError(location+"GetAllUsers.cursor.Err", cursorError.Error())
		logging.Logger(cursorInternalError)
		return nil, cursorInternalError
	}
	if len(fetchedUsers) == 0 {
		return &userModel.Users{
			Users: make([]*userModel.User, 0),
		}, nil
	}

	users := userRepositoryModel.UsersRepositoryToUsersMapper(fetchedUsers)
	users.Limit = limit
	return users, nil
}

func (userRepository *UserRepository) GetUserById(ctx context.Context, userID string) (*userModel.User, error) {
	userObjectID, _ := primitive.ObjectIDFromHex(userID)

	var fetchedUser *userRepositoryModel.UserRepository
	query := bson.M{"_id": userObjectID}
	userFindOneError := userRepository.collection.FindOne(ctx, query).Decode(&fetchedUser)
	if validator.IsErrorNotNil(userFindOneError) {
		userFindOneEntityNotFoundError := domainError.NewEntityNotFoundError(location+"GetUserById.FindOne.Decode", userFindOneError.Error())
		logging.Logger(userFindOneEntityNotFoundError)
		return nil, userFindOneEntityNotFoundError
	}

	user := userRepositoryModel.UserRepositoryToUserMapper(fetchedUser)
	return user, nil
}

func (userRepository *UserRepository) GetUserByEmail(ctx context.Context, email string) (*userModel.User, error) {
	var fetchedUser *userRepositoryModel.UserRepository
	query := bson.M{"email": strings.ToLower(email)}
	userFindOneError := userRepository.collection.FindOne(ctx, query).Decode(&fetchedUser)
	if validator.IsErrorNotNil(userFindOneError) {
		userFindOneEntityNotFoundError := domainError.NewEntityNotFoundError(location+"GetUserByEmail.FindOne.Decode", userFindOneError.Error())
		logging.Logger(userFindOneEntityNotFoundError)
		return nil, userFindOneEntityNotFoundError
	}

	user := userRepositoryModel.UserRepositoryToUserMapper(fetchedUser)
	return user, nil
}

func (userRepository *UserRepository) CheckEmailDublicate(ctx context.Context, email string) error {
	var fetchedUser *userRepositoryModel.UserRepository
	query := bson.M{"email": strings.ToLower(email)}
	userFindOneError := userRepository.collection.FindOne(ctx, query).Decode(&fetchedUser)
	if validator.IsValueNil(fetchedUser) {
		return nil
	}
	if validator.IsErrorNotNil(userFindOneError) {
		userFindOneInternalError := domainError.NewInternalError(location+"CheckEmailDublicate.FindOne.Decode", userFindOneError.Error())
		logging.Logger(userFindOneInternalError)
		return userFindOneInternalError
	}
	userFindOneValidationError := domainError.NewValidationError("email", "required", domainUseCaseValidator.EmailAlreadyExists)
	logging.Logger(userFindOneValidationError)
	return userFindOneValidationError
}

func (userRepository *UserRepository) Register(ctx context.Context, userCreate *userModel.UserCreate) *common.Result[*userModel.User] {
	userCreateRepository := userRepositoryModel.UserCreateToUserCreateRepositoryMapper(userCreate)
	userCreateRepository.Password, _ = repositoryUtility.HashPassword(userCreate.Password)

	userCreateRepository.CreatedAt = time.Now()
	userCreateRepository.UpdatedAt = userCreate.CreatedAt

	insertOneResult, insertOneResultError := userRepository.collection.InsertOne(ctx, &userCreateRepository)
	if validator.IsErrorNotNil(insertOneResultError) {
		userCreateInternalError := domainError.NewInternalError(location+"Register.InsertOne", insertOneResultError.Error())
		logging.Logger(userCreateInternalError)
		return common.NewResultWithError[*userModel.User](userCreateInternalError)
	}

	// Create a unique index for the email field.
	option := options.Index()
	option.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: option}

	_, userIndexesCreateOneError := userRepository.collection.Indexes().CreateOne(ctx, index)
	if validator.IsErrorNotNil(userIndexesCreateOneError) {
		userCreateInternalError := domainError.NewInternalError(location+"Register.Indexes.CreateOne", userIndexesCreateOneError.Error())
		logging.Logger(userCreateInternalError)
		return common.NewResultWithError[*userModel.User](userCreateInternalError)
	}

	var createdUserRepository *userRepositoryModel.UserRepository
	query := bson.M{"_id": insertOneResult.InsertedID}
	userFindOneError := userRepository.collection.FindOne(ctx, query).Decode(&createdUserRepository)
	if validator.IsErrorNotNil(userFindOneError) {
		userCreateEntityNotFoundError := domainError.NewEntityNotFoundError(location+"Register.FindOne.Decode", userFindOneError.Error())
		logging.Logger(userCreateEntityNotFoundError)
		return common.NewResultWithError[*userModel.User](userCreateEntityNotFoundError)
	}

	createdUser := userRepositoryModel.UserRepositoryToUserMapper(createdUserRepository)
	return common.NewResultWithData[*userModel.User](createdUser)
}

func (userRepository *UserRepository) UpdateUserById(ctx context.Context, userID string, user *userModel.UserUpdate) (*userModel.User, error) {
	userUpdateRepository := userRepositoryModel.UserUpdateToUserUpdateRepositoryMapper(user)
	userUpdateRepository.UpdatedAt = time.Now()

	userUpdateRepositoryMappedToMongoDB, mongoMapperError := mongoMapper.MongoMappper(userUpdateRepository)
	if validator.IsErrorNotNil(mongoMapperError) {
		userUpdateError := domainError.NewInternalError(location+"UpdateUserById.MongoMapper", mongoMapperError.Error())
		logging.Logger(userUpdateError)
		return nil, userUpdateError
	}

	userIDMappedToMongoDB, _ := primitive.ObjectIDFromHex(userID)
	query := bson.D{{Key: "_id", Value: userIDMappedToMongoDB}}
	update := bson.D{{Key: "$set", Value: userUpdateRepositoryMappedToMongoDB}}
	result := userRepository.collection.FindOneAndUpdate(ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))
	var updatedUserRepository *userRepositoryModel.UserRepository
	updateUserRepositoryDecodeError := result.Decode(&updatedUserRepository)
	if validator.IsErrorNotNil(updateUserRepositoryDecodeError) {
		userUpdateError := domainError.NewInternalError(location+"UpdateUserById.Decode", updateUserRepositoryDecodeError.Error())
		logging.Logger(userUpdateError)
		return nil, userUpdateError
	}

	updatedUser := userRepositoryModel.UserRepositoryToUserMapper(updatedUserRepository)
	return updatedUser, nil
}

func (userRepository *UserRepository) DeleteUserById(ctx context.Context, userID string) error {
	userIDMappedToMongoDB, _ := primitive.ObjectIDFromHex(userID)
	query := bson.M{"_id": userIDMappedToMongoDB}
	result, userDeleteOneError := userRepository.collection.DeleteOne(ctx, query)
	if validator.IsErrorNotNil(userDeleteOneError) {
		deletedUserError := domainError.NewInternalError(location+"Delete.DeleteOne", userDeleteOneError.Error())
		logging.Logger(deletedUserError)
		return deletedUserError
	}
	if result.DeletedCount == 0 {
		deletedUserError := domainError.NewInternalError(location+"Delete.DeleteOne.DeletedCount", delitionIsNotSuccessful)
		logging.Logger(deletedUserError)
		return deletedUserError
	}
	return nil
}

func (userRepository *UserRepository) UpdateNewRegisteredUserById(ctx context.Context, userID string, key string, value string) (*userModel.User, error) {
	userIDMappedToMongoDB, _ := primitive.ObjectIDFromHex(userID)
	query := bson.D{{Key: "_id", Value: userIDMappedToMongoDB}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: key, Value: value}}}}
	result, userUpdateUpdateOneError := userRepository.collection.UpdateOne(ctx, query, update)
	if validator.IsErrorNotNil(userUpdateUpdateOneError) {
		updatedUserError := domainError.NewInternalError(location+"UpdateNewRegisteredUserById.UpdateOne", userUpdateUpdateOneError.Error())
		logging.Logger(updatedUserError)
		return nil, updatedUserError
	}
	if result.ModifiedCount == 0 {
		updatedUserError := domainError.NewInternalError(location+"UpdateNewRegisteredUserById.UpdateOne.ModifiedCount", updateIsNotSuccessful)
		logging.Logger(updatedUserError)
		return nil, updatedUserError
	}
	return nil, nil
}

func (userRepository *UserRepository) ResetUserPassword(ctx context.Context, firstKey string, firstValue string, secondKey string, passwordKey, password string) error {
	hashedPassword, _ := repositoryUtility.HashPassword(password)
	query := bson.D{{Key: firstKey, Value: firstValue}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: passwordKey, Value: hashedPassword}}}, {Key: "$unset", Value: bson.D{{Key: firstKey, Value: ""}, {Key: secondKey, Value: ""}}}}
	result, updateUserPasswordUpdateOneError := userRepository.collection.UpdateOne(ctx, query, update)
	if validator.IsErrorNotNil(updateUserPasswordUpdateOneError) {
		updatedUserPasswordError := domainError.NewInternalError(location+"ResetUserPassword.UpdateOne", updateUserPasswordUpdateOneError.Error())
		logging.Logger(updatedUserPasswordError)
		return updatedUserPasswordError
	}
	if result.ModifiedCount == 0 {
		updatedUserPasswordError := domainError.NewInternalError(location+"ResetUserPassword.UpdateOne.ModifiedCount", updateUserPasswordUpdateOneError.Error())
		logging.Logger(updatedUserPasswordError)
		return updatedUserPasswordError
	}
	return nil
}

func (userRepository *UserRepository) UpdatePasswordResetTokenUserByEmail(ctx context.Context, email string, firstKey string, firstValue string,
	secondKey string, SecondValue time.Time) error {

	query := bson.D{{Key: "email", Value: strings.ToLower(email)}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: firstKey, Value: firstValue}, {Key: secondKey, Value: secondKey}}}}
	result, updateUserResetTokenUpdateOneError := userRepository.collection.UpdateOne(ctx, query, update)
	if validator.IsErrorNotNil(updateUserResetTokenUpdateOneError) {
		updatedUserResetTokenError := domainError.NewInternalError(location+"UpdatePasswordResetTokenUserByEmail.UpdateOne", updateUserResetTokenUpdateOneError.Error())
		logging.Logger(updatedUserResetTokenError)
		return updatedUserResetTokenError
	}
	if result.ModifiedCount == 0 {
		updatedUserResetTokenError := domainError.NewInternalError(location+"UpdatePasswordResetTokenUserByEmail.UpdateOne.ModifiedCount", updateIsNotSuccessful)
		logging.Logger(updatedUserResetTokenError)
		return updatedUserResetTokenError
	}
	return nil
}

func (userRepository *UserRepository) SendEmailVerificationMessage(ctx context.Context, user *userModel.User, data *userModel.EmailData) error {
	sendEmailError := userRepositoryMail.SendEmail(ctx, user, data)
	if validator.IsErrorNotNil(sendEmailError) {
		logging.Logger(sendEmailError)
		return sendEmailError
	}
	return nil
}

func (userRepository *UserRepository) SendEmailForgottenPasswordMessage(ctx context.Context, user *userModel.User, data *userModel.EmailData) error {
	sendEmailError := userRepositoryMail.SendEmail(ctx, user, data)
	if validator.IsErrorNotNil(sendEmailError) {
		logging.Logger(sendEmailError)
		return sendEmailError
	}
	return nil
}
