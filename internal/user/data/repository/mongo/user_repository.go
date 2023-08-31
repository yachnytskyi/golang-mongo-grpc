package repository

import (
	"context"
	"strings"
	"time"

	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userRepositoryMail "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/external/mail"
	userRepositoryModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/mongo/model"

	repositoryUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/utility"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	utility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"

	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/error/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	sendEmailVerificationMessageSendEmail            = "User.Data.Repository.MongoDB.SendEmailVerificationMessage.SendEmail"
	forgottenPasswordMessageSendEmail                = "User.Data.Repository.MongoDB.ForgottenPasswordMessage.SendEmail"
	getAllUsersFind                                  = "User.Data.Repository.MongoDB.GetAllUsers.Find"
	getAllUsersCursorDecode                          = "User.Data.Repository.MongoDB.GetAllUsers.cursor.decode"
	getAllUsersCursor                                = "User.Data.Repository.MongoDB.GetAllUsers.cursor.Err"
	getUserByIdFindOne                               = "User.Data.Repository.MongoDB.GetUserById.FindOne"
	getUserByEmailFindOne                            = "User.Data.Repository.MongoDB.GetUserByEmail.FindOne"
	registerFindOne                                  = "User.Data.Repository.MongoDB.Register.FindOne"
	registerInsertOne                                = "User.Data.Repository.MongoDB.Register.InsertOne"
	registerIndexesCreateOne                         = "User.Data.Repository.MongoDB.Register.Indexes.CreateOne"
	updateUserByIdMongoMapper                        = "User.Data.Repository.MongoDB.UpdateUserById.MongoMapper"
	updateUserByIdDecode                             = "User.Data.Repository.MongoDB.UpdateUserById.Decode"
	deleteUserByIdDeleteOne                          = "User.Data.Repository.MongoDB.Delete.DeleteOne"
	deleteUserByIdDeletedCount                       = "User.Data.Repository.MongoDB.Delete.DeletedCount"
	updateNewRegisteredUserByIdUpdateOne             = "User.Data.Repository.MongoDB.UpdateNewRegisteredUserById.UpdateOne"
	updateNewRegisteredUserByIdModifiedCount         = "User.Data.Repository.MongoDB.UpdateNewRegisteredUserById.ModifiedCount"
	resetUserPasswordUpdateOne                       = "User.Data.Repository.MongoDB.ResetUserPassword.UpdateOne"
	resetUserPasswordModifiedCount                   = "User.Data.Repository.MongoDB.ResetUserPassword.ModifiedCount"
	updatePasswordResetTokenUserByEmailUpdateOne     = "User.Data.Repository.MongoDB.UpdatePasswordResetTokenUserByEmail.UpdateOne"
	updatePasswordResetTokenUserByEmailModifiedCount = "User.Data.Repository.MongoDB.UpdatePasswordResetTokenUserByEmail.ModifiedCount"
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

	if cursorFindError != nil {
		getAllUsersEntityNotFoundError := domainError.NewEntityNotFoundError(getAllUsersFind, cursorFindError.Error())
		logging.Logger(getAllUsersEntityNotFoundError)
		return nil, getAllUsersEntityNotFoundError
	}
	defer cursor.Close(ctx)

	var fetchedUsers = make([]*userRepositoryModel.UserRepository, 0, limit)
	for cursor.Next(ctx) {
		user := &userRepositoryModel.UserRepository{}

		cursorDecodeError := cursor.Decode(user)
		if cursorDecodeError != nil {
			fetchedUserInternalError := domainError.NewInternalError(getAllUsersCursorDecode, cursorDecodeError.Error())
			logging.Logger(fetchedUserInternalError)
			return nil, fetchedUserInternalError
		}
		fetchedUsers = append(fetchedUsers, user)
	}
	cursorError := cursor.Err()
	if cursorError != nil {
		cursorInternalError := domainError.NewInternalError(getAllUsersCursor, cursorError.Error())
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
	if userFindOneError != nil {
		userFindOneEntityNotFoundError := domainError.NewEntityNotFoundError(getUserByIdFindOne, userFindOneError.Error())
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
	if userFindOneError != nil {
		userFindOneEntityNotFoundError := domainError.NewEntityNotFoundError(getUserByEmailFindOne, userFindOneError.Error())
		logging.Logger(userFindOneEntityNotFoundError)
		return nil, userFindOneEntityNotFoundError
	}

	user := userRepositoryModel.UserRepositoryToUserMapper(fetchedUser)
	return user, nil
}

func (userRepository *UserRepository) CheckEmailDublicate(ctx context.Context, email string) bool {
	var fetchedUser *userRepositoryModel.UserRepository
	query := bson.M{"email": strings.ToLower(email)}
	userFindOneError := userRepository.collection.FindOne(ctx, query).Decode(&fetchedUser)
	if fetchedUser == nil {
		return false
	}
	if userFindOneError != nil {
		return false
	}
	return true
}

func (userRepository *UserRepository) SendEmailVerificationMessage(user *userModel.User, data *userModel.EmailData, templateName string) error {
	sendEmailError := userRepositoryMail.SendEmail(user, data, templateName)
	if sendEmailError != nil {
		logging.Logger(domainError.NewInternalError(sendEmailVerificationMessageSendEmail, sendEmailError.Error()))
		return sendEmailError
	}
	return nil
}

func (userRepository *UserRepository) SendEmailForgottenPasswordMessage(user *userModel.User, data *userModel.EmailData, templateName string) error {
	sendEmailError := userRepositoryMail.SendEmail(user, data, templateName)
	if sendEmailError != nil {
		logging.Logger(domainError.NewInternalError(forgottenPasswordMessageSendEmail, sendEmailError.Error()))
		return sendEmailError
	}
	return nil
}

func (userRepository *UserRepository) Register(ctx context.Context, userCreate *userModel.UserCreate) (*userModel.User, error) {
	userCreateRepository := userRepositoryModel.UserCreateToUserCreateRepositoryMapper(userCreate)
	userCreateRepository.CreatedAt = time.Now()
	userCreateRepository.UpdatedAt = userCreate.CreatedAt
	userCreateRepository.Password, _ = repositoryUtility.HashPassword(userCreate.Password)

	insertOneResult, insertOneResultError := userRepository.collection.InsertOne(ctx, &userCreateRepository)
	if insertOneResultError != nil {
		userCreateInternalError := domainError.NewInternalError(registerInsertOne, insertOneResultError.Error())
		logging.Logger(userCreateInternalError)
		return nil, userCreateInternalError
	}

	// Create a unique index for the email field.
	option := options.Index()
	option.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: option}

	_, userIndexesCreateOneError := userRepository.collection.Indexes().CreateOne(ctx, index)
	if userIndexesCreateOneError != nil {
		userCreateInternalError := domainError.NewInternalError(registerIndexesCreateOne, userIndexesCreateOneError.Error())
		logging.Logger(userCreateInternalError)
		return nil, userCreateInternalError
	}

	var createdUser *userModel.User
	query := bson.M{"_id": insertOneResult.InsertedID}
	userFindOneError := userRepository.collection.FindOne(ctx, query).Decode(&createdUser)
	if userFindOneError != nil {
		userCreateEntityNotFoundError := domainError.NewEntityNotFoundError(registerFindOne, userFindOneError.Error())
		logging.Logger(userCreateEntityNotFoundError)
		return nil, userCreateEntityNotFoundError
	}
	return createdUser, nil
}

func (userRepository *UserRepository) UpdateUserById(ctx context.Context, userID string, user *userModel.UserUpdate) (*userModel.User, error) {
	userUpdateRepository := userRepositoryModel.UserUpdateToUserUpdateRepositoryMapper(user)
	userUpdateRepository.UpdatedAt = time.Now()

	userUpdateRepositoryMappedToMongoDB, mongoMapperError := utility.MongoMappper(userUpdateRepository)
	if mongoMapperError != nil {
		userUpdateError := domainError.NewInternalError(updateUserByIdMongoMapper, mongoMapperError.Error())
		logging.Logger(userUpdateError)
		return nil, userUpdateError
	}

	userIDMappedToMongoDB, _ := primitive.ObjectIDFromHex(userID)
	query := bson.D{{Key: "_id", Value: userIDMappedToMongoDB}}
	update := bson.D{{Key: "$set", Value: userUpdateRepositoryMappedToMongoDB}}
	result := userRepository.collection.FindOneAndUpdate(ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))
	var updatedUserRepository *userRepositoryModel.UserRepository

	updateUserRepositoryDecodeError := result.Decode(&updatedUserRepository)
	if updateUserRepositoryDecodeError != nil {
		userUpdateError := domainError.NewInternalError(updateUserByIdDecode, updateUserRepositoryDecodeError.Error())
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

	if userDeleteOneError != nil {
		deletedUserError := domainError.NewInternalError(deleteUserByIdDeleteOne, userDeleteOneError.Error())
		logging.Logger(deletedUserError)
		return deletedUserError
	}
	if result.DeletedCount == 0 {
		deletedUserError := domainError.NewInternalError(deleteUserByIdDeletedCount, userDeleteOneError.Error())
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

	if userUpdateUpdateOneError != nil {
		updatedUserError := domainError.NewInternalError(updateNewRegisteredUserByIdUpdateOne, userUpdateUpdateOneError.Error())
		logging.Logger(updatedUserError)
		return nil, updatedUserError
	}
	if result.ModifiedCount == 0 {
		updatedUserError := domainError.NewInternalError(updateNewRegisteredUserByIdUpdateOne, userUpdateUpdateOneError.Error())
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

	if updateUserPasswordUpdateOneError != nil {
		updatedUserPasswordError := domainError.NewInternalError(resetUserPasswordUpdateOne, updateUserPasswordUpdateOneError.Error())
		logging.Logger(updatedUserPasswordError)
		return updatedUserPasswordError
	}
	if result.ModifiedCount == 0 {
		updatedUserPasswordError := domainError.NewInternalError(resetUserPasswordUpdateOne, updateUserPasswordUpdateOneError.Error())
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

	if updateUserResetTokenUpdateOneError != nil {
		updatedUserResetTokenError := domainError.NewInternalError(resetUserPasswordUpdateOne, updateUserResetTokenUpdateOneError.Error())
		logging.Logger(updatedUserResetTokenError)
		return updatedUserResetTokenError
	}
	if result.ModifiedCount == 0 {
		updatedUserResetTokenError := domainError.NewInternalError(resetUserPasswordUpdateOne, updateUserResetTokenUpdateOneError.Error())
		logging.Logger(updatedUserResetTokenError)
		return updatedUserResetTokenError
	}
	return nil
}
