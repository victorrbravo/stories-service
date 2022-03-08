package handlers

import (
	"log"
	"net/http"
	"stories/generated/models"

	"stories/database"
	"stories/handlers_utils"
	"stories/helpers"
)

type GetAllHandler struct {
	Query database.IQuery
}

func (GetAllHandler) IsValid() error {
	return nil
}

// Check if endpoint is enabled
func (handler GetAllHandler) IsEnabled() bool {
	return true
}

func (handler GetAllHandler) GetQuery() database.IQuery {
	log.Println("GetQuery was executed")
	return handler.Query
}

// Get all suggestions for you
func (handler GetAllHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var err error
	var users []models.User
	log.Println("GETQUERY.....handler.query", handler.GetQuery())
	users, err = models.GetAllUsers(handler.GetQuery())
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
		return
	}
	handlers_utils.SendJSONResponse(writer, users, http.StatusOK)
}
