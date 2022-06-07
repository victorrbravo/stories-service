package social

import (
	"fmt"
	"log"
	"net/http"
	"stories/database"
	"stories/handlers_utils"
	"stories/helpers"
	"stories/social"
	"stories/types"
	"stories/users"
)

type AddStepHandler struct {
	Comment types.String `json:"comment"`
	Query   database.IQuery
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
	var comment social.Comment
	obj.EntityName = "ta_flow"
	obj.Keyfield = "name"
	log.Println("GetPathParameters name request.RequestURI", request.RequestURI)
	strId, err = handlers_utils.GetPathParameter("nameId", request)
	log.Println("*strId |", strId, "|")
	err = handlers_utils.GetParameters(&handler, request)
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, err.Error()))
		return
	}
	comment, err = obj.AddComment(types.NewString(strId), fmt.Sprintf("%s", handler.Comment), user, handler.GetQuery())
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, err.Error()))
		return
	}
	handlers_utils.SendJSONResponse(writer, comment, http.StatusOK)
}
