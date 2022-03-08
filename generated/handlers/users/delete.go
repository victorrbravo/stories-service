package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"stories/database"
	"stories/generated/models"
	"stories/handlers_utils"
	"stories/helpers"
	"stories/types"
	"strconv"
)

type DeleteHandler struct {
	Query database.IQuery
}

// Check if endpoint is enabled
func (handler DeleteHandler) IsEnabled() bool {
	return true
}

func (handler DeleteHandler) GetQuery() database.IQuery {
	log.Println("GetQuery was executed")
	return handler.Query
}

// Delete follow with given id
func (handler DeleteHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var err error
	var element *models.User
	var query database.IQuery
	var id string
	var intId int

	query = handler.GetQuery()
	var values url.Values
	values = request.URL.Query()
	id = values.Get("id")
	intId, err = strconv.Atoi(id)
	log.Println("ID: ", intId)
	element, err = models.GetUserByID(query, types.NewInt(intId))
	if err != nil {
		if err == sql.ErrNoRows {
			handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusNotFound, err.Error()))
			return
		}
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
		return
	}
	err = element.Delete(query)
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
		return
	}
	handlers_utils.SendJSONResponse(writer, element, http.StatusOK)
}
