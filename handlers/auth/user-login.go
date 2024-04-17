package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"stories/constants"
	"stories/crypto"
	"stories/database"
	"stories/generated/models"
	"stories/handlers_utils"
	"stories/helpers"
	"stories/types"
)

type UserLoginHandler struct {
	Email    types.String `json:"email"`
	Password types.String `json:"password"`
	Query    database.IQuery
	User     ModelUser
}

func (handler UserLoginHandler) GetQuery() database.IQuery {
	return handler.Query
}

func (handler UserLoginHandler) IsValid() error {
	if types.IsNil(handler.Email) {
		return errors.New(constants.EmailCanNotBeEmpty)
	}
	if types.IsNil(handler.Password) {
		return errors.New(constants.PasswordCanNotBeEmpty)
	}
	return nil
}

func (handler UserLoginHandler) IsEnabled() bool {
	return true
}

func (handler *UserLoginHandler) SetQuery(query database.IQuery) {
	handler.Query = query
}

func (handler *UserLoginHandler) SetUser(user ModelUser) {
	handler.User = user
}

// Issue a bearer token to authenticate HTTP requests
// Credentials are evaluated and a new _JWT_ is created that allows the caller
// to authenticate their requests
func (handler UserLoginHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var err error
	fmt.Println("serverhttp login")
	err = handlers_utils.GetParameters(&handler, request)
	fmt.Println("serverhttp login err", err)
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, err.Error()))
		return
	}
	fmt.Println("serverhttp login Email searching by |", handler.Email, "|")
	fmt.Println("serverhttp login Password searching by |", handler.Password, "|")

	objects, err := models.GetTaUserByEmail(handler.GetQuery(), handler.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			handlers_utils.RespondWithError(
				writer,
				helpers.NewHttpError(http.StatusUnauthorized, constants.UserNotFoundOrPasswordMismatch),
			)
			return
		}
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
		return
	}
	taUser := objects[0]
	fmt.Println("CONVERT taUser", taUser, "err", err)
	pass := taUser.GetPassword()
	fmt.Printf("Compare pass |%s|\n", pass)
	fmt.Printf("Compare handler.Password |%s|\n", handler.Password)
	passcrypt, _ := bcrypt.GenerateFromPassword([]byte(types.AsString(handler.Password)), 0)
	fmt.Printf(
		"pass bcrypt |%s|\n", string(passcrypt))
	err = crypto.ComparePasswords([]byte(types.AsString(pass)), types.AsString(handler.Password))
	fmt.Println("crypto crypto.ComparePasswords err", err)
	if err != nil {
		handlers_utils.RespondWithError(
			writer,
			helpers.NewHttpError(http.StatusUnauthorized, constants.UserNotFoundOrPasswordMismatch),
		)

		return
	}
	// token, err = user.CreateSession(handler.queries)
	var tokenDetails *helpers.TokenDetails
	tokenDetails, err = helpers.CreateJWTToken(types.AsString(handler.Email))
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
		return
	}
	log.Println("tokenDetails.AccessToken", tokenDetails.AccessToken)
	// token.SessionCount = user.GetSessionCount()
	// Update auth token
	// user.SetAuthToken(token.Token)
	// Send the response to the client
	handlers_utils.SendJSONResponse(writer, tokenDetails.AccessToken, http.StatusOK)
}
