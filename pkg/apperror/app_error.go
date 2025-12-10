package apperror

import (
	"errors"
)

type AppError struct {
	Err            error  // Underlying error
	Code           string // Machine readable error code
	HTTPStatusCode int
	PublicMsg      string // User Friendly msg
}

func (e *AppError) Error() string { return e.Err.Error() }
func (e *AppError) Unwrap() error { return e.Err }

func IsAppError(err error) bool {
	var ae *AppError

	return errors.As(err, &ae)
}

func ShouldLogError(err error) bool {
	if err == nil {
		return false
	}

	if IsAppError(err) {
		return false
	}

	return true
}
