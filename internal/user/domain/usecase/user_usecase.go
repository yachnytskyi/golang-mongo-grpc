package usecase

import (
	"context"
	"time"

	"github.com/thanhpk/randstr"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	commonUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	verificationCodeLength int = 20
	resetTokenLength       int = 20
	userRole                   = "user"
)

type UserUseCase struct {
	userRepository user.UserRepository
}

func NewUserUseCase(userRepository user.UserRepository) user.UserUseCase {
	return UserUseCase{userRepository: userRepository}
}

// GetAllUsers retrieves a list of users based on the provided pagination parameters.
// The result is wrapped in a commonModel.Result containing either the user or an error.
func (userUseCase UserUseCase) GetAllUsers(ctx context.Context, paginationQuery commonModel.PaginationQuery) commonModel.Result[userModel.Users] {
	// Fetch the users.
	fetchedUsers := userUseCase.userRepository.GetAllUsers(ctx, paginationQuery)
	if validator.IsErrorNotNil(fetchedUsers.Error) {
		fetchedUsers.Error = domainError.HandleError(fetchedUsers.Error)
		return commonModel.NewResultOnFailure[userModel.Users](fetchedUsers.Error)
	}
	return fetchedUsers
}

// GetUserById retrieves a user by their ID using the user ID.
// The result is wrapped in a commonModel.Result containing either the user or an error.
func (userUseCase UserUseCase) GetUserById(ctx context.Context, userID string) commonModel.Result[userModel.User] {
	// Fetch the user.
	fetchedUser := userUseCase.userRepository.GetUserById(ctx, userID)
	if validator.IsErrorNotNil(fetchedUser.Error) {
		fetchedUser.Error = domainError.HandleError(fetchedUser.Error)
		return commonModel.NewResultOnFailure[userModel.User](fetchedUser.Error)
	}
	return fetchedUser
}

// GetUserByEmail retrieves a user by their ID using the provided user email.
// It performs email format validation and fetches the user from the repository.
// The result is wrapped in a commonModel.Result containing either the user or an error.
func (userUseCase UserUseCase) GetUserByEmail(ctx context.Context, email string) commonModel.Result[userModel.User] {
	// Validate the email.
	validateEmailError := validateEmail(email, emailRegex)
	if validator.IsValueNotNil(validateEmailError) {
		processedError := domainError.HandleError(validateEmailError)
		return commonModel.NewResultOnFailure[userModel.User](processedError)
	}

	// Fetch the user.
	fetchedUser := userUseCase.userRepository.GetUserByEmail(ctx, email)
	if validator.IsErrorNotNil(fetchedUser.Error) {
		fetchedUser.Error = domainError.HandleError(fetchedUser.Error)
		return commonModel.NewResultOnFailure[userModel.User](fetchedUser.Error)
	}
	return fetchedUser
}

// Register registers a new user based on the provided data, generates a verification token,
// and sends an email verification message. The result is wrapped in a commonModel.Result
// containing either the user or an error.
func (userUseCase UserUseCase) Register(ctx context.Context, userCreateData userModel.UserCreate) commonModel.Result[userModel.User] {
	// Validate the user creation data.
	userCreate := validateUserCreate(userCreateData)
	if validator.IsErrorNotNil(userCreate.Error) {
		userCreate.Error = domainError.HandleError(userCreate.Error)
		return commonModel.NewResultOnFailure[userModel.User](userCreate.Error)
	}

	// Check for duplicate email.
	checkEmailDuplicateError := userUseCase.userRepository.CheckEmailDuplicate(ctx, userCreate.Data.Email)
	if validator.IsErrorNotNil(checkEmailDuplicateError) {
		checkEmailDuplicateError = domainError.HandleError(checkEmailDuplicateError)
		return commonModel.NewResultOnFailure[userModel.User](checkEmailDuplicateError)
	}

	// Generate a verification token and set user properties.
	tokenValue := randstr.String(verificationCodeLength)
	encodedTokenValue := commonUtility.Encode(tokenValue)
	userCreate.Data.Role = userRole
	userCreate.Data.Verified = true
	userCreate.Data.VerificationCode = encodedTokenValue

	// Register the user.
	createdUser := userUseCase.userRepository.Register(ctx, userCreate.Data)
	if validator.IsErrorNotNil(createdUser.Error) {
		createdUser.Error = domainError.HandleError(createdUser.Error)
		return createdUser
	}

	// Prepare email data for user registration.
	// Send the email verification message and return the created user.
	emailData := prepareEmailDataForRegistration(ctx, createdUser.Data.Name, tokenValue)
	sendEmailVerificationMessageError := userUseCase.userRepository.SendEmailVerificationMessage(ctx, createdUser.Data, emailData)
	if validator.IsErrorNotNil(sendEmailVerificationMessageError) {
		sendEmailVerificationMessageError = domainError.HandleError(sendEmailVerificationMessageError)
		return commonModel.NewResultOnFailure[userModel.User](sendEmailVerificationMessageError)
	}
	return createdUser
}

// UpdateUserById updates a user's information based on the provided data.
// The result is wrapped in a commonModel.Result containing either the user or an error.
func (userUseCase UserUseCase) UpdateUserById(ctx context.Context, userUpdateData userModel.UserUpdate) commonModel.Result[userModel.User] {
	// Validate the user update data.
	userUpdate := validateUserUpdate(userUpdateData)
	if validator.IsErrorNotNil(userUpdate.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(userUpdate.Error))
	}

	// Update the user.
	updatedUser := userUseCase.userRepository.UpdateUserById(ctx, userUpdate.Data)
	if validator.IsErrorNotNil(updatedUser.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(updatedUser.Error))
	}
	return updatedUser
}

// DeleteUserById deletes a user based on the provided user ID.
func (userUseCase UserUseCase) DeleteUserById(ctx context.Context, userID string) error {
	deletedUser := userUseCase.userRepository.DeleteUserById(ctx, userID)
	if validator.IsErrorNotNil(deletedUser) {
		return domainError.HandleError(deletedUser)
	}

	// User deletion was successful. Return nil to indicate no error.
	return nil
}

func (userUseCase UserUseCase) Login(ctx context.Context, userLoginData userModel.UserLogin) (string, error) {
	userLogin := validateUserLogin(userLoginData)
	if validator.IsErrorNotNil(userLogin.Error) {
		handledError := domainError.HandleError(userLogin.Error)
		return "", handledError
	}

	fetchedUser := userUseCase.userRepository.GetUserByEmail(ctx, userLogin.Data.Email)
	if validator.IsErrorNotNil(fetchedUser.Error) {
		return "", fetchedUser.Error
	}
	arePasswordsNotEqualError := arePasswordsNotEqual(fetchedUser.Data.Password, userLoginData.Password)
	if validator.IsValueNotNil(arePasswordsNotEqualError) {
		arePasswordsNotEqualError.Notification = invalidEmailOrPassword
		return "", arePasswordsNotEqualError
	}
	return fetchedUser.Data.UserID, nil
}

func (userUseCase UserUseCase) UpdatePasswordResetTokenUserByEmail(ctx context.Context, email string, firstKey string, firstValue string,
	secondKey string, secondValue time.Time) error {

	validateEmailError := validateEmail(email, emailRegex)
	if validator.IsValueNotNil(validateEmailError) {
		handledError := domainError.HandleError(validateEmailError)
		return handledError
	}
	updatedUserError := userUseCase.userRepository.UpdatePasswordResetTokenUserByEmail(ctx, email, firstKey, firstValue, secondKey, secondValue)
	if validator.IsErrorNotNil(updatedUserError) {
		updatedUserError = domainError.HandleError(updatedUserError)
		return updatedUserError
	}

	// Generate verification code.
	tokenValue := randstr.String(resetTokenLength)
	encodedTokenValue := commonUtility.Encode(tokenValue)
	tokenExpirationTime := time.Now().Add(time.Minute * 15)

	// Update the user.
	fetchedUser := userUseCase.GetUserByEmail(ctx, email)
	if validator.IsErrorNotNil(fetchedUser.Error) {
		fetchedUserError := domainError.HandleError(fetchedUser.Error)
		return fetchedUserError
	}
	updatedUserPasswordError := userUseCase.userRepository.UpdatePasswordResetTokenUserByEmail(ctx, fetchedUser.Data.Email, "passwordResetToken", encodedTokenValue, "passwordResetAt", tokenExpirationTime)
	if validator.IsErrorNotNil(updatedUserPasswordError) {
		updatedUserPasswordError = domainError.HandleError(updatedUserPasswordError)
		return updatedUserPasswordError
	}

	emailData := prepareEmailDataForUpdatePasswordResetToken(ctx, fetchedUser.Data.Name, tokenValue)
	sendEmailForgottenPasswordMessageError := userUseCase.userRepository.SendEmailForgottenPasswordMessage(ctx, fetchedUser.Data, emailData)
	if validator.IsErrorNotNil(sendEmailForgottenPasswordMessageError) {
		sendEmailForgottenPasswordMessageError = domainError.HandleError(sendEmailForgottenPasswordMessageError)
		return sendEmailForgottenPasswordMessageError
	}
	return nil
}

func (userUseCase UserUseCase) ResetUserPassword(ctx context.Context, firstKey string, firstValue string, secondKey string, passwordKey, password string) error {
	updatedUser := userUseCase.userRepository.ResetUserPassword(ctx, firstKey, firstValue, secondKey, passwordKey, password)
	return updatedUser
}

// prepareEmailData is a helper function to create an EmailData model for sending an email.
// It takes the context, user name, token value, email subject, URL, template name, and template path as input.
// It constructs an EmailData model and returns it in a Result.
func prepareEmailData(ctx context.Context, userName, tokenValue, subject, url, templateName, templatePath string) userModel.EmailData {
	applicationConfig := config.AppConfig
	userFirstName := domainUtility.UserFirstName(userName)
	emailData := userModel.NewEmailData(applicationConfig.Email.ClientOriginUrl+url+tokenValue, templateName, templatePath, userFirstName, subject)
	return emailData
}

// prepareEmailDataForRegister is a helper function to prepare an EmailData model specifically for user registration.
// It takes the context, user name, and token value as input and uses the constants for email subject, URL, template name, and template path.
// It internally calls prepareEmailData with the appropriate parameters and returns the result.
func prepareEmailDataForRegistration(ctx context.Context, userName, tokenValue string) userModel.EmailData {
	applicationConfig := config.AppConfig
	return prepareEmailData(ctx, userName, tokenValue, constants.EmailConfirmationSubject, constants.EmailConfirmationUrl,
		applicationConfig.Email.UserConfirmationTemplateName, applicationConfig.Email.UserConfirmationTemplatePath)
}

// prepareEmailDataForUserUpdate is a helper function to prepare an EmailData model specifically for updating user information.
// It takes the context, user name, and token value as input and uses the constants for email subject, URL, template name, and template path.
// It internally calls prepareEmailData with the appropriate parameters and returns the result.
func prepareEmailDataForUpdatePasswordResetToken(ctx context.Context, userName, tokenValue string) userModel.EmailData {
	applicationConfig := config.AppConfig
	return prepareEmailData(ctx, userName, tokenValue, constants.ForgottenPasswordSubject, constants.ForgottenPasswordUrl,
		applicationConfig.Email.ForgottenPasswordTemplateName, applicationConfig.Email.ForgottenPasswordTemplatePath)
}
