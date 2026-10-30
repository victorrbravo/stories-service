package social

import (
	"log"
	"net/http"

	"github.com/victorrbravo/stories-service/database"
	"github.com/victorrbravo/stories-service/handlers_utils"
	"github.com/victorrbravo/stories-service/helpers"
	"github.com/victorrbravo/stories-service/social"
	"github.com/victorrbravo/stories-service/types"
	"github.com/victorrbravo/stories-service/users"
)

type AddStepHandler struct {
	Name      types.String    `json:"name"`
	TokenLink types.String    `json:"token_link"`
	Pattern   types.String    `json:"pattern"`
	Query     database.IQuery `json:"queries,omitempty"`
}

func (handler AddStepHandler) GetQuery() database.IQuery {
	return handler.Query
}

func (handler AddStepHandler) IsValid() error {
	return nil
}

// Check if endpoint is enabled
func (handler AddStepHandler) IsEnabled() bool {
	return true
}

// Create a new suggestion for you. It's only an administrative action.
func (handler AddStepHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request, user users.User) {
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
