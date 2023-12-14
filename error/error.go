package error

import (
	"gofr.dev/pkg/errors"
)

type ApiError struct {
	Err struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (ae ApiError) Error() string {
	return ae.Err.Message
}
func NewApiError(message string) ApiError {
	return ApiError{Err: struct {
		Message string `json:"message"`
	}{Message: message}}
}

func HttpStatusError(statusCode int, message string) errors.Raw {
	var statusMessage string

	switch statusCode {
	case 400:
		statusMessage = "Bad Request"
	case 500:
		statusMessage = "Internal Server Error"
	default:
		statusMessage = "Unknown Status Code"
	}

	return errors.Raw{StatusCode: statusCode, Err: NewApiError(statusMessage + " - " + message)}
}
