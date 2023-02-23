package helpers

import "fmt"

type CommonError struct {
	Message string
}

func NewError(message string) CommonError {
	return CommonError{Message: message}
}

func (err CommonError) Error() string {
	return fmt.Sprintf("%s", err.Message)
}
