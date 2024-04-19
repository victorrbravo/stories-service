package auth

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"stories/database"
	"stories/generated/models"
	"stories/handlers_utils"
	"stories/helpers"
	"stories/types"
)

type SignupUserHandler struct {
	Email                string       `json:"email" validate:"required,email"`
	Password             types.String `json:"password" validate:"required"`
	PasswordConfirmation types.String `json:"password_confirmation"  validate:"required"`
	Name                 types.String `json:"name"  validate:"required"`
	LastName             types.String `json:"last_name"  validate:"required"`
	Phone                string       `json:"phone" validate:"required,e164"`
	Gender               types.Gender `json:"gender"`
	Query                database.IQuery
}

func (handler SignupUserHandler) IsValid() error {
	return nil
}

// Check if endpoint is enabled
func (handler SignupUserHandler) IsEnabled() bool {
	return true
}

func (handler SignupUserHandler) GetQuery() database.IQuery {
	log.Println("SetQuery was executed")
	return handler.Query
}

// Get currently authenticated user's information
func (handler SignupUserHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if handler.Query.GetDatabase() == nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, "Database not initialized"))
		return
	}
	err := handlers_utils.GetParameters(&handler, request)
	if err != nil {
		log.Println("SignupUserHandler ServeHTTP GetParameters error", err)
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, err.Error()))
		return
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(handler)
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, err.Error()))
		return
	}
	if handler.Password != handler.PasswordConfirmation {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, "Passwords do not match"))
		return
	}
	passcrypt, _ := bcrypt.GenerateFromPassword([]byte(types.AsString(handler.Password)), 0)
	password := string(passcrypt)
	taUser, err := models.CreateTaUser(handler.Query,
		handler.Name,
		handler.LastName,
		types.NewString(handler.Email),
		types.NewString(password),
		types.NewString("user"),
		types.NewBool(false),
	)
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
		return
	}
	log.Println(" created taUser ID:", taUser.PkUser)
	handlers_utils.SendJSONResponse(writer, struct {
		Message string `json:"message"`
		Email   string `json:"email"`
	}{Message: fmt.Sprintf("User %s %s created successfully", handler.Name, handler.LastName),
		Email: handler.Email}, http.StatusOK)
}
