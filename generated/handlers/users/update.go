package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"stories/generated/models"
	"stories/users"

	"stories/database"
	"stories/handlers_utils"
	"stories/helpers"
	"stories/types"
)

type UpdateHandler struct {
	FirstName   types.String `json:"first_name"`
	LastName    types.String `json:"last_name"`
	Email       types.String `json:"email"`
	AccountType types.String `json:"account_type"`
	Query       database.IQuery
}

func (handler UpdateHandler) IsValid() error {
	return nil
}

// Check if endpoint is enabled
func (handler UpdateHandler) IsEnabled() bool {
	return true
}

func (handler UpdateHandler) GetQuery() database.IQuery {
	return handler.Query
}

// Edit ProfileSettings
func (handler UpdateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request, user users.User) {
	var err error
	var element *models.TaUser
	var id string
	values := request.URL.Query()
	id = values.Get("id")
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, err.Error()))
		return
	}
	err = handlers_utils.GetParameters(&handler, request)
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, err.Error()))
		return
	}
	valueInt, _ := strconv.Atoi(id)
	element, err = models.GetTaUserByPkUser(handler.GetQuery(), types.NewInt(valueInt))
	if err != nil {
		if err == sql.ErrNoRows {
			handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusNotFound, err.Error()))
			return
		}
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
		return
	}
	var pass types.String
	err = element.Update(
		handler.GetQuery(),
		element.FirstName,
		element.LastName,
		handler.Email,
		pass,
		handler.AccountType,
	)
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
		return
	}
	handlers_utils.SendJSONResponse(writer, element, http.StatusOK)
}
