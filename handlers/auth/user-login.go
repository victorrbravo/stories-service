package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"stories/pkg/utils"

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
	log.Println("serverhttp login")
	err = handlers_utils.GetParameters(&handler, request)
	log.Println("serverhttp login err", err)
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, err.Error()))
		return
	}
	log.Println("serverhttp login Email searching by |", handler.Email, "|")
	objects, err := models.GetTaUserByEmail(handler.GetQuery(), handler.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
	log.Println("CONVERT taUser", taUser, "err", err)
	if !taUser.Activated.ToBool() {
		log.Println("User not activated yet email:", handler.Email)
		handlers_utils.RespondWithError(
			writer,
			helpers.NewHttpError(http.StatusUnauthorized, UserNotActivated),
		)
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
	// token, err = user.CreateSession(handler.queries)
	var tokenDetails *helpers.TokenDetails
	tokenDetails, err = helpers.CreateJWTToken(types.AsString(handler.Email))
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
		return
	}
	handlers_utils.SendJSONResponse(writer, struct {
		Message       string `json:"message"`
		Token         string `json:"token"`
		Refreshtoken  string `json:"refreshtoken"`
		Expires       int64  `json:"expires"`
		ExpiresString string `json:"expires_string"`
		Username      string `json:"username"`
	}{
		Message:       fmt.Sprintf("User %s logged in successfully", handler.Email),
		Token:         tokenDetails.AccessToken,
		Refreshtoken:  tokenDetails.RefreshToken,
		Expires:       tokenDetails.AtExpires,
		ExpiresString: utils.UnixToRFC3339(tokenDetails.AtExpires),
		Username:      taUser.Username.String(),
	}, http.StatusOK)
}
