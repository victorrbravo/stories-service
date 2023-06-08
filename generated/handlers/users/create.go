package handlers

import (
	"fmt"
	"log"
	"net/http"

	"stories/constants"
	"stories/database"
	"stories/generated/models"
	"stories/handlers_utils"
	"stories/helpers"
	"stories/types"
)

type CreateHandler struct {
	FirstName   types.String `json:"first_name"`
	LastName    types.String `json:"last_name"`
	Email       types.String `json:"email"`
	AccountType types.String `json:"account_type"`
	Query       database.IQuery
}

func (handler CreateHandler) IsValid() error {
	if types.IsNil(handler.FirstName) {
		return fmt.Errorf(constants.RequiredFieldIsEmpty, "name")
	}
	return nil
}

func (handler CreateHandler) GetQuery() database.IQuery {
	log.Println("SetQuery was executed")
	return handler.Query
}

// Check if endpoint is enabled
func (handler CreateHandler) IsEnabled() bool {
	return true
}

// Create a new currency
func (handler CreateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var err error
	err = handlers_utils.GetParameters(&handler, request)
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, err.Error()))
		return
	}
	user, err := models.CreateTaUser(handler.GetQuery(), handler.FirstName, handler.LastName, handler.Email,
		types.NewString("pass"),
		handler.AccountType)
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
		return
	}
	handlers_utils.SendJSONResponse(writer, user, http.StatusOK)
}
