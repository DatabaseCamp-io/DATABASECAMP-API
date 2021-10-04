package errs

import "net/http"

type AppError struct {
	Code      int
	ThMessage string
	EnMessage string
}

func (e AppError) Error() string {
	return e.EnMessage
}

func NewNotFoundError(thMessage string, enMessage string) error {
	return AppError{
		Code:      http.StatusNotFound,
		ThMessage: thMessage,
		EnMessage: enMessage,
	}
}
