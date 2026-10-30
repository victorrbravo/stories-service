package auth

import (
	"fmt"
	"log"
	"net/http"
	"github.com/victorrbravo/stories-service/constants"
	"github.com/victorrbravo/stories-service/crypto"
	"github.com/victorrbravo/stories-service/database"
	"github.com/victorrbravo/stories-service/generated/models"
	"github.com/victorrbravo/stories-service/handlers_utils"
	"github.com/victorrbravo/stories-service/helpers"
	"github.com/victorrbravo/stories-service/types"

	"golang.org/x/crypto/bcrypt"
)

type ChangePasswordHandler struct {
	Code        types.String `json:"code"`
	Password    types.String `json:"password" validate:"required"`
	NewPassword types.String `json:"new_password" validate:"required"`
	Query       database.IQuery
}

func (handler ChangePasswordHandler) IsValid() error {
	return nil
}

// Check if endpoint is enabled
func (handler ChangePasswordHandler) IsEnabled() bool {
	return true
}

// Change user's password
func (handler ChangePasswordHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Println("change password handler")
	var err error
	err = handlers_utils.GetParameters(&handler, request)
	if err != nil {
		log.Println("Input has errors err", err.Error())
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, err.Error()))
		return
	}
	searchCode := handler.Code
	log.Println("searchCode", searchCode)
	taSessions, err := models.GetTaSessionByToken(handler.Query, searchCode)
	log.Println("taSessions err", err)
	log.Println("len taSessions", taSessions)
	if err != nil {
		log.Printf("Error getting token %s", searchCode)
		handlers_utils.ResponseWithInternalError(writer)
		return
	}
	if len(taSessions) == 0 {
		log.Printf("Token %s not found", searchCode)
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, fmt.Sprintf(constants.InvalidToken, searchCode)))
		return
	}
	log.Println("checking passwords")
	if handler.Password.IsEmpty() && handler.NewPassword.IsEmpty() {
		handlers_utils.SendJSONResponse(writer, struct {
			Status  int    `json:"Status"`
			Message string `json:"Message"`
		}{
			Status:  http.StatusOK,
			Message: "Code verified successfully",
		}, http.StatusOK)
	}

	taSession := taSessions[0]
	pkUser := taSession.FkUser
	taUser, err := models.GetTaUserByPkUser(handler.Query, pkUser)
	if err != nil {
		log.Printf("Error getting user %v", pkUser)
		handlers_utils.ResponseWithInternalError(writer)
		return
	}
	pass := taUser.GetPassword()
	err = crypto.ComparePasswords([]byte(types.AsString(pass)), types.AsString(handler.Password))
	log.Println("crypto crypto.ComparePasswords err", err)
	if err != nil {
		handlers_utils.RespondWithError(
			writer,
			helpers.NewHttpError(http.StatusUnauthorized, constants.UserNotFoundOrPasswordMismatch),
		)
		return
	}
	passcrypt, _ := bcrypt.GenerateFromPassword([]byte(types.AsString(handler.NewPassword)), 0)
	password := string(passcrypt)
	err = taUser.Update(handler.Query,
		types.NilString(),
		types.NilString(),
		types.NilString(),
		types.NewString(password),
		types.NilString(),
		types.NilBool,
		types.NilString(),
	)
	if err != nil {
		log.Println(fmt.Sprintf("Error updating user %+v", pkUser))
		handlers_utils.ResponseWithInternalError(writer)
		return
	}
	handlers_utils.SendJSONResponse(writer, struct {
		Status  int    `json:"Status"`
		Message string `json:"Message"`
	}{
		Status:  http.StatusOK,
		Message: "Password changed successfully",
	}, http.StatusOK)

}
