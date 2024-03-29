package domain

import (
	constant "github.com/yachnytskyi/golang-mongo-grpc/config/constant"
)

func HandleError(err error) error {
	switch errorType := err.(type) {
	case ValidationError:
		return errorType
	case ValidationErrors:
		return errorType
	case ErrorMessage:
		return errorType
	case EntityNotFoundError:
		return NewErrorMessage(constant.EntityNotFoundErrorNotification)
	case PaginationError:
		return errorType
	default:
		return NewErrorMessage(constant.InternalErrorNotification)
	}
}
