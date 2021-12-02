package errs

import "net/http"

type AppError struct {
	Code      int
	ThMessage string
	EnMessage string
}

// Create app error in english message
func (e AppError) Error() string {
	return e.EnMessage
}

// Create not found error in thai and english message
func NewNotFoundError(thMessage string, enMessage string) error {
	return AppError{
		Code:      http.StatusNotFound,
		ThMessage: thMessage,
		EnMessage: enMessage,
	}
}

// Create forbidden error in thai and english message
func NewForbiddenError(thMessage string, enMessage string) error {
	return AppError{
		Code:      http.StatusForbidden,
		ThMessage: thMessage,
		EnMessage: enMessage,
	}
}

// Create internal server error in thai and english message
func NewInternalServerError(thMessage string, enMessage string) error {
	return AppError{
		Code:      http.StatusInternalServerError,
		ThMessage: thMessage,
		EnMessage: enMessage,
	}
}

// Create bad request error in thai and english message
func NewBadRequestError(thMessage string, enMessage string) error {
	return AppError{
		Code:      http.StatusBadRequest,
		ThMessage: thMessage,
		EnMessage: enMessage,
	}
}

// Create service unavilable error in thai and english message
func NewServiceUnavailableError(thMessage string, enMessage string) error {
	return AppError{
		Code:      http.StatusServiceUnavailable,
		ThMessage: thMessage,
		EnMessage: enMessage,
	}
}
