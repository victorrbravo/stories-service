package auth

import (
	"net/http"

	"stories/handlers_utils"
	"stories/types"
)

type ModelUser interface {
	GetEmail() types.String
	GetPassword() types.String
}

type GetUserInfoHandler struct{}

func (handler GetUserInfoHandler) IsValid() error {
	return nil
}

// Check if endpoint is enabled
func (handler GetUserInfoHandler) IsEnabled() bool {
	return true
}

// Get currently authenticated user's information
func (handler GetUserInfoHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request, user ModelUser) {
	handlers_utils.SendJSONResponse(writer, user, http.StatusOK)
}
