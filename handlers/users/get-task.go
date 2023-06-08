package users

import (
	"net/http"

	"stories/handlers/auth"
	"stories/handlers_utils"
)

type GetTask struct{}

func (handler GetTask) IsValid() error {
	return nil
}

// Check if endpoint is enabled
func (handler GetTask) IsEnabled() bool {
	return true
}

// Get currently authenticated user's information
func (handler GetTask) ServeHTTP(writer http.ResponseWriter, request *http.Request, user auth.ModelUser) {
	handlers_utils.SendJSONResponse(writer, user, http.StatusOK)
}
