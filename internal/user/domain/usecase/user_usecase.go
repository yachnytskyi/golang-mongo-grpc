package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/thanhpk/randstr"
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

type UserUseCase struct {
	userRepository user.Repository
}

func NewUserUseCase(userRepository user.Repository) user.UseCase {
	return &UserUseCase{userRepository: userRepository}
}

func (userUseCase *UserUseCase) GetAllUsers(ctx context.Context, page int, limit int) (*userModel.Users, error) {
	fetchedUsers, err := userUseCase.userRepository.GetAllUsers(ctx, page, limit)

	return fetchedUsers, err
}

func (userUseCase *UserUseCase) GetUserById(ctx context.Context, userID string) (*userModel.User, error) {
	fetchedUser, err := userUseCase.userRepository.GetUserById(ctx, userID)

	return fetchedUser, err
}

func (userUseCase *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*userModel.User, error) {
	fetchedUser, err := userUseCase.userRepository.GetUserByEmail(ctx, email)

	return fetchedUser, err
}

func (userUseCase *UserUseCase) Register(ctx context.Context, userCreate *userModel.UserCreate) error {
	if userUseCase.userRepository.CheckEmailDublicate(ctx, userCreate.Email) {
		userCreateValidationError := &domainError.ValidationError{
			Field:        "email",
			FieldType:    "required",
			Notification: EmailAlreadyExists,
		}

		domainError.HandleError(userCreateValidationError)
		return userCreateValidationError
	}
	userCreateValidationErrors := UserCreateValidator(userCreate)
	if validator.IsErrorNotNil(userCreateValidationErrors) {
		userCreateValidationErrors = domainError.HandleError(userCreateValidationErrors)
		return userCreateValidationErrors
	}
	userCreate.Verified = true
	userCreate.Role = "user"
	userCreate.Password, _ = domainUtility.HashPassword(userCreate.Password)
	createdUser, userCreateError := userUseCase.userRepository.Register(ctx, userCreate)
	if validator.IsErrorNotNil(userCreateError) {
		userCreateError = domainError.HandleError(userCreateError)
		return userCreateError
	}

	// Generate verification code.
	code := randstr.String(20)
	verificationCode := utility.Encode(code)

	_, userUpdateError := userUseCase.userRepository.UpdateNewRegisteredUserById(ctx, createdUser.UserID, "verificationCode", verificationCode)
	if validator.IsErrorNotNil(userUpdateError) {
		userUpdateError = domainError.HandleError(userUpdateError)
		return userUpdateError
	}

	// Send an email.
	loadConfig, loadConfigError := config.LoadConfig(".")
	if validator.IsErrorNotNil(loadConfigError) {
		var loadConfigInternalError *domainError.InternalError = new(domainError.InternalError)
		loadConfigInternalError.Location = "User.Domain.UserUseCase.Registration.LoadConfig"
		loadConfigInternalError.Reason = loadConfigError.Error()
		domainError.HandleError(loadConfigInternalError)
		return loadConfigInternalError
	}

	userFirstName := domainUtility.UserFirstName(createdUser.Name)
	emailData := userModel.EmailData{
		URL:       loadConfig.Origin + "/verifyemail/" + code,
		FirstName: userFirstName,
		Subject:   "Your account verification code",
	}

	if validator.IsErrorNotNil(userUseCase.userRepository.SendEmailVerificationMessage(createdUser, &emailData, config.TemplateName)) {
		sendEmailVerificationMessageError := &domainError.ValidationError{
			Notification: "domainError.InternalErrorNotification",
		}
		domainError.HandleError(sendEmailVerificationMessageError)
		return sendEmailVerificationMessageError
	}

	return nil
}

func (userUseCase *UserUseCase) UpdateUserById(ctx context.Context, userID string, user *userModel.UserUpdate) (*userModel.User, error) {
	userUpdateValidationErrors := UserUpdateValidator(user)
	if validator.IsErrorNotNil(userUpdateValidationErrors) {
		return nil, userUpdateValidationErrors
	}

	updatedUser, userUpdateError := userUseCase.userRepository.UpdateUserById(ctx, userID, user)
	if validator.IsErrorNotNil(userUpdateError) {
		return nil, userUpdateError
	}
	return updatedUser, nil
}

func (userUseCase *UserUseCase) DeleteUserById(ctx context.Context, userID string) error {
	deletedUser := userUseCase.userRepository.DeleteUserById(ctx, userID)

	return deletedUser
}

func (userUseCase *UserUseCase) Login(ctx context.Context, user *userModel.UserLogin) (string, error) {
	err := UserLoginValidator(user)
	if validator.IsErrorNotNil(err) {
		return "", err
	}

	fetchedUser, err := userUseCase.userRepository.GetUserByEmail(ctx, user.Email)

	// Will return wrong email or password.
	if validator.IsErrorNotNil(err) {
		return "", fmt.Errorf("invalid email or password")
	}

	// Verify password - we previously created this method.
	matchPasswords := domainUtility.VerifyPassword(fetchedUser.Password, user.Password)
	if validator.IsErrorNotNil(matchPasswords) {
		return "", fmt.Errorf("invalid email or password")
	}
	return fetchedUser.UserID, err
}

func (userUseCase *UserUseCase) UpdateNewRegisteredUserById(ctx context.Context, userID string, key string, value string) (*userModel.User, error) {
	updatedUser, err := userUseCase.userRepository.UpdateNewRegisteredUserById(ctx, userID, key, value)

	return updatedUser, err
}

func (userUseCase *UserUseCase) UpdatePasswordResetTokenUserByEmail(ctx context.Context, email string, firstKey string, firstValue string,
	secondKey string, secondValue time.Time) error {
	updatedUserError := userUseCase.userRepository.UpdatePasswordResetTokenUserByEmail(ctx, email, firstKey, firstValue, secondKey, secondValue)
	if validator.IsErrorNotNil(updatedUserError) {
		updatedUserError = domainError.HandleError(updatedUserError)
		return updatedUserError
	}

	// Generate verification code.
	resetToken := randstr.String(20)
	passwordResetToken := utility.Encode(resetToken)
	passwordResetAt := time.Now().Add(time.Minute * 15)

	// Update the user.
	fetchedUser, fetchedUserError := userUseCase.GetUserByEmail(ctx, email)
	if validator.IsErrorNotNil(fetchedUserError) {
		fetchedUserError = domainError.HandleError(fetchedUserError)
		return fetchedUserError
	}

	updatedUserPasswordError := userUseCase.userRepository.UpdatePasswordResetTokenUserByEmail(ctx, fetchedUser.Email, "passwordResetToken", passwordResetToken, "passwordResetAt", passwordResetAt)
	if validator.IsErrorNotNil(updatedUserPasswordError) {
		updatedUserPasswordError = domainError.HandleError(updatedUserPasswordError)
		return updatedUserPasswordError
	}
	userFirstName := domainUtility.UserFirstName(fetchedUser.Name)

	// Send an email.
	loadConfig, loadConfigError := config.LoadConfig(".")
	if validator.IsErrorNotNil(loadConfigError) {
		var sendEmailInternalError *domainError.InternalError = new(domainError.InternalError)
		sendEmailInternalError.Location = "User.Domain.UserUseCase.UpdatePasswordResetTokenUserByEmail.LoadConfig"
		sendEmailInternalError.Reason = loadConfigError.Error()
		domainError.HandleError(sendEmailInternalError)
		return sendEmailInternalError
	}

	emailData := userModel.EmailData{
		URL:       loadConfig.Origin + "/reset-password/" + resetToken,
		FirstName: userFirstName,
		Subject:   "Your password reset token (it is valid for 15 minutes)",
	}

	if validator.IsErrorNotNil(userUseCase.userRepository.SendEmailForgottenPasswordMessage(fetchedUser, &emailData, config.TemplateName)) {
		sendEmailForgottenPasswordMessage := &domainError.ValidationError{
			Notification: "domainError.InternalErrorNotification",
		}
		domainError.HandleError(sendEmailForgottenPasswordMessage)
		return sendEmailForgottenPasswordMessage
	}
	return nil
}

func (userUseCase *UserUseCase) ResetUserPassword(ctx context.Context, firstKey string, firstValue string, secondKey string, passwordKey, password string) error {
	hashedPassword, _ := domainUtility.HashPassword(password)
	updatedUser := userUseCase.userRepository.ResetUserPassword(ctx, firstKey, firstValue, secondKey, passwordKey, hashedPassword)
	return updatedUser

}
