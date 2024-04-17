package auth

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
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
	Gender               types.Gender `json:"gender"  validate:"required"`
}

func (handler SignupUserHandler) IsValid() error {
	return nil
}

// Check if endpoint is enabled
func (handler SignupUserHandler) IsEnabled() bool {
	return true
}

// Get currently authenticated user's information
func (handler SignupUserHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	err := handlers_utils.GetParameters(&handler, request)
	if err != nil {
		log.Println("SignupUserHandler ServeHTTP GetParameters error", err)
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, err.Error()))
		return
	}

	fmt.Println("SignupHandler Name", handler.Name)
	fmt.Println("SignupHandler Email", handler.Email)
	fmt.Println("SignupHandler Password", handler.Password)
	fmt.Println("SignupHandler RepeatPassword", handler.PasswordConfirmation)
	fmt.Println("SignupHandler LastName", handler.LastName)
	fmt.Println("SignupHandler Phone", handler.Phone)
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
	objResult := "{}"
	handlers_utils.SendJSONResponse(writer, objResult, http.StatusOK)
}
