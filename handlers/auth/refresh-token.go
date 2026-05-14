package auth

import (
	"errors"
	"net/http"
	"stories/database"
	"stories/handlers_utils"
	"stories/helpers"
	"stories/types"
)

type RefreshTokenHandler struct {
	RefreshToken types.String `json:"refresh_token"`
	Query        database.IQuery
}

func (handler RefreshTokenHandler) GetQuery() database.IQuery {
	return handler.Query
}

func (handler RefreshTokenHandler) IsValid() error {
	if handler.RefreshToken.IsEmpty() {
		return errors.New("refresh token is required")
	}
	return nil
}

func (handler RefreshTokenHandler) IsEnabled() bool {
	return true
}

func (handler *RefreshTokenHandler) SetQuery(query database.IQuery) {
	handler.Query = query
}

func (handler RefreshTokenHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	err := handlers_utils.GetParameters(&handler, request)
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, err.Error()))
		return
	}

	claims, err := helpers.VerifyRefreshToken(types.AsString(handler.RefreshToken))
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusUnauthorized, "invalid refresh token"))
		return
	}

	email, ok := claims["email"].(string)
	if !ok {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusUnauthorized, "invalid token claims"))
		return
	}

	// Issue new token
	tokenDetails, err := helpers.CreateJWTToken(email)
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
		return
	}

	handlers_utils.SendJSONResponse(writer, struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}{
		Token:        tokenDetails.AccessToken,
		RefreshToken: tokenDetails.RefreshToken,
	}, http.StatusOK)
}
