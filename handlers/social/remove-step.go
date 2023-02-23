package social

import (
	"log"
	"net/http"

	"stories/database"
	"stories/handlers_utils"
	"stories/helpers"
	"stories/social"
	"stories/types"
	"stories/users"
)

type ReplaceStepHandler struct {
	Name      types.String    `json:"name"`
	TokenLink types.String    `json:"token_link"`
	Pattern   types.String    `json:"pattern"`
	Query     database.IQuery `json:"query,omitempty"`
}

func (handler ReplaceStepHandler) GetQuery() database.IQuery {
	return handler.Query
}

func (handler ReplaceStepHandler) IsValid() error {
	return nil
}

// Check if endpoint is enabled
func (handler ReplaceStepHandler) IsEnabled() bool {
	return true
}

// Create a new suggestion for you. It's only an administrative action.
func (handler ReplaceStepHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request, user users.User) {
	var err error
	var strId string
	var obj social.Social
	var task social.Task
	obj.EntityName = "ta_flow"
	obj.Keyfield = "name"
	log.Println("GetPathParameters name request.RequestURI", request.RequestURI)
	strId, err = handlers_utils.GetPathParameter("nameId", request)
	log.Println("*strId |", strId, "|")
	err = handlers_utils.GetParameters(&handler, request)
	log.Println("*GetParameters err |", err, "|")
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, err.Error()))
		return
	}

	task, err = obj.AddTask(strId, types.AsString(handler.Name), types.AsString(handler.TokenLink),
		types.AsString(handler.Pattern), user, handler.GetQuery())

	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, err.Error()))
		return
	}

	handlers_utils.SendJSONResponse(writer, task, http.StatusOK)
}
