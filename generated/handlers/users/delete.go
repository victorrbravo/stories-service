package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/victorrbravo/stories-service/database"
	"github.com/victorrbravo/stories-service/generated/models"
	"github.com/victorrbravo/stories-service/handlers_utils"
	"github.com/victorrbravo/stories-service/helpers"
	"github.com/victorrbravo/stories-service/types"
	"github.com/victorrbravo/stories-service/users"
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
func (handler DeleteHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request, user users.User) {
	var err error
	var element *models.TaUser
	var query database.IQuery
	var id string
	var intId int

	query = handler.GetQuery()
	var values url.Values
	values = request.URL.Query()
	id = values.Get("id")
	intId, err = strconv.Atoi(id)
	log.Println("ID: ", intId)
	element, err = models.GetTaUserByPkUser(query, types.NewInt(intId))
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
