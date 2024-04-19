package auth

import (
	"log"
	"net/http"
	"stories/types"
)

type ChangePasswordHandler struct {
	Password    types.String `json:"password" validate:"required"`
	NewPassword types.String `json:"new_password" validate:"required"`
}

func (handler ChangePasswordHandler) IsValid() error {
	return nil
}

// Check if endpoint is enabled
func (handler ChangePasswordHandler) IsEnabled() bool {
	return true
}

// Change user's password
func (handler ChangePasswordHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Println("change password handler")
}
