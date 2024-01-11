package domain

import (
	"fmt"
	"regexp"
	"strings"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	location = "pkg.utility.validator.domain."
)

func SanitizeString(data string) string {
	return strings.TrimSpace(data)
}

func ToLowerString(data string) string {
	return strings.ToLower(data)
}

func SanitizeAndToLowerString(data string) string {
	data = strings.TrimSpace(data)
	return strings.ToLower(data)
}

func ValidateField(field, fieldName, fieldRegex, notification string) error {
	if IsStringLengthInvalid(field, constants.MinStringLength, constants.MaxStringLength) {
		notification = fmt.Sprintf(constants.StringAllowedLength, constants.MinStringLength, constants.MaxStringLength)
		validationError := domainError.NewValidationError(location+"ValidateField.IsStringLengthInvalid", fieldName, constants.FieldRequired, notification)
		logging.Logger(validationError)
		return validationError
	}
	if AreStringCharactersInvalid(field, fieldRegex) {
		validationError := domainError.NewValidationError(location+"ValidateField.AreStringCharactersInvalid", fieldName, constants.FieldRequired, notification)
		logging.Logger(validationError)
		return validationError
	}
	return nil
}

func ValidateOptionalField(field, fieldName, fieldType, fieldRegex, notification string) error {
	if IsStringLengthInvalid(field, constants.MinOptionalStringLength, constants.MaxStringLength) {
		notification = fmt.Sprintf(constants.StringOptionalAllowedLength, constants.MaxStringLength)
		validationError := domainError.NewValidationError(location+"ValidateOptionalField.IsStringLengthInvalid", fieldName, fieldType, notification)
		logging.Logger(validationError)
		return validationError
	}
	if AreStringCharactersInvalid(field, fieldRegex) {
		validationError := domainError.NewValidationError(location+"ValidateField.AreStringCharactersInvalid", fieldName, fieldType, notification)
		logging.Logger(validationError)
		return validationError
	}
	return nil
}

func IsStringLengthInvalid(checkedString string, minLength int, maxLength int) bool {
	if len(checkedString) < minLength || len(checkedString) > maxLength {
		return true
	}
	return false
}

func AreStringCharactersInvalid(checkedString string, regexString string) bool {
	return !regexp.MustCompile(regexString).MatchString(checkedString)
}

func IsOptionalStringLengthInvalid(checkedString string, minLength int, maxLength int) bool {
	if len(checkedString) < minLength || len(checkedString) > maxLength {
		return true
	}
	return false
}
