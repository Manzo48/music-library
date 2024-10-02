// pkg/error_message/errors.go
package error_message

import (
   // "github.com/sirupsen/logrus"
)

type ErrorResponse struct {
    Message string `json:"message"`
}

var (
    ErrNotFound    = NewError("Record not found")
    ErrInternal    = NewError("Internal server error")
    ErrBadRequest  = NewError("Bad request")
    ErrExternalAPI = NewError("Failed to fetch data from external API")
)

func NewError(message string) error {
    return &CustomError{Message: message}
}

type CustomError struct {
    Message string
}

func (e *CustomError) Error() string {
    return e.Message
}
