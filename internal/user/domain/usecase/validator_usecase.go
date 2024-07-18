package usecase

import (
	"fmt"
	"net"
	"strings"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/domain"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/domain"
	logger "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logger"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	domainValidator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator/domain"
	bcrypt "golang.org/x/crypto/bcrypt"
)

// Constants used for various validation messages and field names.
const (
	// Regex Patterns for validating email, username, and password.
	emailRegex    = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[\\\\\\\\/=\\\\{\\|}]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[\\\\+\\-\\/=\\\\_{\\|}]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.||[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.||[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	usernameRegex = `^[a-zA-z0-9-_ \t]*$`
	passwordRegex = `^[a-zA-z0-9-_*,.]*$`

	// Error Messages for invalid inputs.
	passwordAllowedCharacters = "Sorry, only letters (a-z), numbers(0-9), the asterics, hyphen and underscore characters are allowed."
	emailAllowedCharacters    = "Sorry, only letters (a-z), numbers(0-9) and periods (.) are allowed, you cannot use a period in the end and more than one in a row."
	invalidEmailDomain        = "Email domain does not exist."
	passwordsDoNotMatch       = "Passwords do not match."
	invalidEmailOrPassword    = "Invalid email or password."

	// Field Names used in validation.
	usernameField         = "name"
	EmailField            = "email"
	passwordField         = "password"
	emailOrPasswordFields = "email or password"
	resetTokenField       = "reset token"
)

var (
	emailValidator = domainModel.CommonValidator{
		FieldName:    EmailField,
		FieldRegex:   emailRegex,
		MinLength:    constants.MinStringLength,
		MaxLength:    constants.MaxStringLength,
		Notification: emailAllowedCharacters,
	}
	usernameValidator = domainModel.CommonValidator{
		FieldName:    usernameField,
		FieldRegex:   usernameRegex,
		MinLength:    constants.MinStringLength,
		MaxLength:    constants.MaxStringLength,
		Notification: constants.StringAllowedCharacters,
	}
	passwordValidator = domainModel.CommonValidator{
		FieldName:    passwordField,
		FieldRegex:   usernameRegex,
		MinLength:    constants.MinStringLength,
		MaxLength:    constants.MaxStringLength,
		Notification: passwordAllowedCharacters,
	}
	tokenValidator = domainModel.CommonValidator{
		FieldName:  resetTokenField,
		FieldRegex: usernameRegex,
		MinLength:  resetTokenLength,
		MaxLength:  resetTokenLength,
	}
	// Add more validators for other fields as needed.
)

func validateUserCreate(userCreate userModel.UserCreate) common.Result[userModel.UserCreate] {
	validationErrors := make([]error, 0, 4)
	userCreate.Email = domainUtility.SanitizeAndToLowerString(userCreate.Email)
	userCreate.Name = domainUtility.SanitizeString(userCreate.Name)
	userCreate.Password = domainUtility.SanitizeString(userCreate.Password)
	userCreate.PasswordConfirm = domainUtility.SanitizeString(userCreate.PasswordConfirm)

	validationErrors = validateEmail(location+"validateUserCreate", userCreate.Email, validationErrors)
	validationErrors = domainValidator.ValidateField(location+"validateUserCreate", userCreate.Name, usernameValidator, validationErrors)
	validationErrors = validatePassword(location+"validateUserCreate", userCreate.Password, userCreate.PasswordConfirm, validationErrors)
	if validator.IsSliceNotEmpty(validationErrors) {
		return common.NewResultOnFailure[userModel.UserCreate](domainError.NewValidationErrors(validationErrors))
	}

	return common.NewResultOnSuccess[userModel.UserCreate](userCreate)
}

func validateUserUpdate(userUpdate userModel.UserUpdate) common.Result[userModel.UserUpdate] {
	validationErrors := make([]error, 0, 1)
	userUpdate.Name = domainUtility.SanitizeString(userUpdate.Name)

	validationErrors = domainValidator.ValidateField(location+"validateUserUpdate", userUpdate.Name, usernameValidator, validationErrors)
	if validator.IsSliceNotEmpty(validationErrors) {
		return common.NewResultOnFailure[userModel.UserUpdate](domainError.NewValidationErrors(validationErrors))
	}

	return common.NewResultOnSuccess[userModel.UserUpdate](userUpdate)
}

func validateUserLogin(userLogin userModel.UserLogin) common.Result[userModel.UserLogin] {
	validationErrors := make([]error, 0, 2)
	userLogin.Email = domainUtility.SanitizeAndToLowerString(userLogin.Email)
	userLogin.Password = domainUtility.SanitizeString(userLogin.Password)

	validationErrors = validateEmail(location+"validateUserLogin", userLogin.Email, validationErrors)
	validationErrors = domainValidator.ValidateField(location+"validateUserLogin", userLogin.Password, usernameValidator, validationErrors)

	if validator.IsSliceNotEmpty(validationErrors) {
		return common.NewResultOnFailure[userModel.UserLogin](domainError.NewValidationErrors(validationErrors))
	}

	return common.NewResultOnSuccess[userModel.UserLogin](userLogin)
}

func validateUserForgottenPassword(userForgottenPassword userModel.UserForgottenPassword) common.Result[userModel.UserForgottenPassword] {
	validationErrors := make([]error, 0, 2)
	userForgottenPassword.Email = domainUtility.SanitizeAndToLowerString(userForgottenPassword.Email)

	validationErrors = validateEmail(location+"validateUserForgottenPassword", userForgottenPassword.Email, validationErrors)
	if validator.IsSliceNotEmpty(validationErrors) {
		return common.NewResultOnFailure[userModel.UserForgottenPassword](domainError.NewValidationErrors(validationErrors))
	}

	return common.NewResultOnSuccess[userModel.UserForgottenPassword](userForgottenPassword)
}

func validateUserResetPassword(userResetPassword userModel.UserResetPassword) common.Result[userModel.UserResetPassword] {
	validationErrors := make([]error, 0, 2)
	userResetPassword.ResetToken = domainUtility.SanitizeString(userResetPassword.ResetToken)
	userResetPassword.Password = domainUtility.SanitizeString(userResetPassword.Password)
	userResetPassword.PasswordConfirm = domainUtility.SanitizeString(userResetPassword.PasswordConfirm)

	validationErrors = validatePassword(location+"validateUserResetPassword", userResetPassword.Password, userResetPassword.PasswordConfirm, validationErrors)
	validationErrors = domainValidator.ValidateField(location+"validateUserResetPassword", userResetPassword.ResetToken, tokenValidator, validationErrors)

	if validator.IsSliceNotEmpty(validationErrors) {
		return common.NewResultOnFailure[userModel.UserResetPassword](domainError.NewValidationErrors(validationErrors))
	}

	return common.NewResultOnSuccess[userModel.UserResetPassword](userResetPassword)
}

func validateEmail(location, email string, validationErrors []error) []error {
	errors := validationErrors

	validateFieldError := validateField(location+".validateEmail", email, emailValidator)
	if validator.IsError(validateFieldError) {
		errors = append(errors, validateFieldError)
		return errors
	}

	checkEmailDomainError := checkEmailDomain(location+".validateEmail", email)
	if validator.IsError(checkEmailDomainError) {
		errors = append(errors, checkEmailDomainError)
	}

	return errors
}

func validatePassword(location, password, passwordConfirm string, validationErrors []error) []error {
	errors := validationErrors

	validateFieldError := validateField(location+".validatePassword", password, passwordValidator)
	if validator.IsError(validateFieldError) {
		errors = append(errors, validateFieldError)
	}

	if password != passwordConfirm {
		validationError := domainError.NewValidationError(
			location+".validatePassword",
			passwordValidator.FieldName,
			constants.FieldRequired,
			passwordsDoNotMatch,
		)

		logger.Logger(validationError)
		errors = append(errors, validationError)
	}

	return errors
}

// checkEmailDomain checks if the email domain exists by resolving DNS records.
func checkEmailDomain(location, emailString string) error {
	host := strings.Split(emailString, "@")[1]
	_, lookupMXError := net.LookupMX(host)
	if validator.IsError(lookupMXError) {
		validationError := domainError.NewValidationError(
			location+".checkEmailDomain",
			EmailField,
			constants.FieldRequired,
			invalidEmailDomain,
		)

		logger.Logger(validationError)
		return validationError
	}

	return nil
}

func checkPasswords(location, hashedPassword string, checkedPassword string) error {
	if validator.IsError(bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(checkedPassword))) {
		validationError := domainError.NewValidationError(
			location+".checkPasswords.CompareHashAndPassword",
			emailOrPasswordFields,
			constants.FieldRequired,
			passwordsDoNotMatch,
		)

		logger.Logger(validationError)
		validationError.Notification = invalidEmailOrPassword
		return validationError
	}

	return nil
}

func checkEmail(location, email string) error {
	validateFieldError := validateField(location+".checkEmail", email, emailValidator)
	if validator.IsError(validateFieldError) {
		return validateFieldError
	}

	return checkEmailDomain(location+".checkEmail", email)
}

func validateField(location, fieldValue string, commonValidator domainModel.CommonValidator) error {
	if domainValidator.IsStringLengthInvalid(fieldValue, commonValidator.MinLength, commonValidator.MaxLength) {
		notification := fmt.Sprintf(constants.StringAllowedLength, commonValidator.MinLength, commonValidator.MaxLength)
		validationError := domainError.NewValidationError(
			location+".validateField.IsStringLengthInvalid",
			commonValidator.FieldName,
			constants.FieldRequired,
			notification,
		)

		logger.Logger(validationError)
		return validationError
	}
	if domainValidator.AreStringCharactersInvalid(fieldValue, commonValidator.FieldRegex) {
		validationError := domainError.NewValidationError(
			location+".validateField.AreStringCharactersInvalid",
			commonValidator.FieldName,
			constants.FieldRequired,
			commonValidator.Notification,
		)

		logger.Logger(validationError)
		return validationError
	}

	return nil
}
