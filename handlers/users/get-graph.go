package users

import (
	"net/http"

	"github.com/victorrbravo/stories-service/handlers/auth"
	"github.com/victorrbravo/stories-service/handlers_utils"
)

type GetGraph struct{}

func (handler GetGraph) IsValid() error {
	return nil
}

// Check if endpoint is enabled
func (handler GetGraph) IsEnabled() bool {
	return true
}

// Get currently authenticated user's information
func (handler GetGraph) ServeHTTP(writer http.ResponseWriter, request *http.Request, user auth.ModelUser) {
	handlers_utils.SendJSONResponse(writer, user, http.StatusOK)
}
