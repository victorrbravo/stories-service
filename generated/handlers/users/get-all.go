package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"stories/generated/models"
	"stories/pkg/workflow/common"
	"stories/users"

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
func (handler GetAllHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request, user users.User) {
	var err error
	//var users []models.TaUser
	var flows []models.TaFlow
	log.Println("GETQUERY.....handler.query", handler.GetQuery())
	_, err = models.GetAllTaUsers(handler.GetQuery())
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
		return
	}
	flows, err = models.GetAllTaFlows(handler.GetQuery())
	wf, err := common.CreateExample(context.Background(), "TouristPlace")
	fmt.Println("CreateExample err", err)
	data, err := json.Marshal(wf)
	fmt.Println("CreateExample wf", string(data))
	if err != nil {
		fmt.Println("error flows")
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
		return
	}
	handlers_utils.SendJSONResponse(writer, flows, http.StatusOK)
}
