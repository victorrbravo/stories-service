package handlers

import (
	"fmt"
	"log"
	"net/http"
	"github.com/victorrbravo/stories-service/generated/models"
	"github.com/victorrbravo/stories-service/users"

	"github.com/victorrbravo/stories-service/database"
	"github.com/victorrbravo/stories-service/handlers_utils"
	"github.com/victorrbravo/stories-service/helpers"
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
func (handler GetAllHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request, user users.User) {
	var err error
	//var users []models.TaUser
	var flows []models.TaFlow
	log.Println("GETQUERY.....handler.queries", handler.GetQuery())
	_, err = models.GetAllTaUsers(handler.GetQuery())
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
		return
	}
	flows, err = models.GetAllTaFlows(handler.GetQuery())
	fmt.Println("Get all executing...")
	if err != nil {
		fmt.Println("error flows")
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
		return
	}
	handlers_utils.SendJSONResponse(writer, flows, http.StatusOK)
}
