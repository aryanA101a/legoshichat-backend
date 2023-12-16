package error

import (
	"gofr.dev/pkg/errors"
)

type Error struct {
	Err struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (ae Error) Error() string {
	return ae.Err.Message
}
func NewError(message string) Error {
	return Error{Err: struct {
		Message string `json:"message"`
	}{Message: message}}
}

func HttpStatusError(statusCode int, message string) errors.Raw {
	var statusMessage string

	switch statusCode {
	case 400:
		statusMessage = "Bad Request"
	case 401:
		statusMessage = "Unauthorized"
	case 403:
		statusMessage = "Forbidden"
	case 404:
		statusMessage = "Not Found"
	case 500:
		statusMessage = "Internal Server Error"
	default:
		statusMessage = "Unknown Status Code"
	}

	return errors.Raw{StatusCode: statusCode, Err: NewError(statusMessage + " - " + message)}
}
