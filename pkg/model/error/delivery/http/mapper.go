package http

import (
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

func ValidationErrorToHttpValidationErrorViewMapper(validationError domainError.ValidationError) HttpValidationErrorView {
	return HttpValidationErrorView{
		Field:        validationError.Field,
		FieldType:    validationError.FieldType,
		Notification: validationError.Notification,
	}
}

func ValidationErrorsToHttpValidationErrorsViewMapper(validationErrors domainError.ValidationErrors) HttpValidationErrorsView {
	httpValidationErrorsView := make([]HttpValidationErrorView, 0, len(validationErrors.ValidationErrors))
	for _, validationError := range validationErrors.ValidationErrors {
		httpValidationErrorView := HttpValidationErrorView{}
		httpValidationErrorView.Field = validationError.Field
		httpValidationErrorView.FieldType = validationError.FieldType
		httpValidationErrorView.Notification = validationError.Notification
		httpValidationErrorsView = append(httpValidationErrorsView, httpValidationErrorView)
	}
	return HttpValidationErrorsView{HttpValidationErrorsView: httpValidationErrorsView}
}

func ErrorMessageToHttpErrorMessageViewMapper(errorMessage domainError.ErrorMessage) HttpErrorMessageView {
	return HttpErrorMessageView{
		Notification: errorMessage.Notification,
	}
}

func PaginationErrorToHttpPaginationErrorView(errorMessage domainError.PaginationError) HttpPaginationErrorView {
	return HttpPaginationErrorView{
		CurrentPage:  errorMessage.CurrentPage,
		TotalPages:   errorMessage.TotalPages,
		Notification: errorMessage.Notification,
	}
}
