package helpers

import (
	"fmt"
	"net/http"
	"strings"
)

type HTTPError struct {
	Status  int
	Message string
}

type MethodNotAllowedError struct {
	Status  int
	Message string
	Methods []string
}

func (err HTTPError) Error() string {
	return fmt.Sprintf("%d: %s", err.Status, err.Message)
}

func NewHttpError(status int, message string) HTTPError {
	return HTTPError{Status: status, Message: message}
}

func (err HTTPError) IsNil() bool {
	return err.Status == 0
}

func NewMethodNotAllowedError(methods []string) MethodNotAllowedError {
	return MethodNotAllowedError{
		Status:  http.StatusMethodNotAllowed,
		Message: "method not allowed",
		Methods: methods,
	}
}

func (err MethodNotAllowedError) Error() string {
	return fmt.Sprintf("%d: %s - %s", err.Status, err.Message, strings.Join(err.Methods, ", "))
}
