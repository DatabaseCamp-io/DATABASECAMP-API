package errs

// errs.go
/**
 * 	This file used to manage error of the application
 */

import "net/http"

/**
 * This class represent code and message of the error
 */
type AppError struct {
	Code      int    // Code of the error
	ThMessage string // Message of the error in Thai
	EnMessage string // Message of the error in English
}

// Implement build in error class
func (e AppError) Error() string {
	return e.EnMessage
}

/**
 * Create not found error by Thai and English message
 *
 * @param thMessage error message in Thai
 * @param enMessage error message in English
 *
 * @return error
 */
func NewNotFoundError(thMessage string, enMessage string) error {
	return AppError{
		Code:      http.StatusNotFound,
		ThMessage: thMessage,
		EnMessage: enMessage,
	}
}

/**
 * Create forbidden error error by Thai and English message
 *
 * @param thMessage error message in Thai
 * @param enMessage error message in English
 *
 * @return error
 */
func NewForbiddenError(thMessage string, enMessage string) error {
	return AppError{
		Code:      http.StatusForbidden,
		ThMessage: thMessage,
		EnMessage: enMessage,
	}
}

/**
 * Create internal server error error by Thai and English message
 *
 * @param thMessage error message in Thai
 * @param enMessage error message in English
 *
 * @return error
 */
func NewInternalServerError(thMessage string, enMessage string) error {
	return AppError{
		Code:      http.StatusInternalServerError,
		ThMessage: thMessage,
		EnMessage: enMessage,
	}
}

/**
 * Create bad request error error by Thai and English message
 *
 * @param thMessage error message in Thai
 * @param enMessage error message in English
 *
 * @return error
 */
func NewBadRequestError(thMessage string, enMessage string) error {
	return AppError{
		Code:      http.StatusBadRequest,
		ThMessage: thMessage,
		EnMessage: enMessage,
	}
}

/**
 * Create service unavilable error error by Thai and English message
 *
 * @param thMessage error message in Thai
 * @param enMessage error message in English
 *
 * @return error
 */
func NewServiceUnavailableError(thMessage string, enMessage string) error {
	return AppError{
		Code:      http.StatusServiceUnavailable,
		ThMessage: thMessage,
		EnMessage: enMessage,
	}
}
