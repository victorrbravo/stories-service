package handlers_utils

import (
	"context"
	"errors"
	"log"
	"net/http"
	"stories/pkg/utils"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"stories/database"
	"stories/generated/models"
	"stories/types"

	"stories/constants"
	"stories/crypto"
	"stories/helpers"
	"stories/users"
)

// AuthorizedHTTPHandler an interface that adds one parameter to the built-in `ServeHTTP', namely the
// authorized user that previously authenticated through the authentication endpoints
type AuthorizedHTTPHandler interface {
	GetQuery() database.IQuery
	IsEnabled() bool
	ServeHTTP(writer http.ResponseWriter, request *http.Request, user users.User)
}

// RequestAuthorizer is kind of a middleware to help authorizing requests and subsequently calling
// the authorized handler
type RequestAuthorizer struct {
	next  AuthorizedHTTPHandler
	Query database.IQuery
	Ctx   context.Context
}

// NewRequestAuthorizer creates new authorizer with the given the next handler that should handle
// the request
func NewRequestAuthorizer(next AuthorizedHTTPHandler, query database.IQuery) *RequestAuthorizer {
	return &RequestAuthorizer{next: next, Query: query}
}

func getToken(header http.Header) (crypto.Token, error) {
	var matches []string
	var value string
	value = header.Get("Authorization")
	matches = strings.Split(value, " ")
	if len(matches) != 2 {
		return crypto.Token{}, errors.New(constants.InvalidAuthorizationHeader)
	}
	// And now use the token string
	return crypto.NewTokenFromString(strings.Trim(matches[1], " "))
}

func (handler *RequestAuthorizer) GetQuery() database.IQuery {
	return handler.Query
}

func (handler *RequestAuthorizer) SetQuery(query database.IQuery) {
	handler.Query = query
}

// Check if endpoint is enabled
func (handler *RequestAuthorizer) IsEnabled() bool {
	return constants.WorkingModeValueNormal == 0
}

func getAuthorizationField(request *http.Request) (string, error) {
	value := strings.TrimSpace(request.Header.Get("Authorization"))
	log.Printf("value was gotten")
	if !strings.HasPrefix(value, utils.BearerPrefix) {
		return value, errors.New(utils.HeaderMessageInvalidAuthorizationHeader)
	}
	result := strings.TrimPrefix(value, utils.BearerPrefix)
	// And now use the token string
	return result, nil
}

func getJWTToken(request *http.Request) (string, error) {
	var err error
	var email string
	var code string
	var jwtToken jwt.MapClaims
	code, err = getAuthorizationField(request)
	if err != nil {
		return code, err
	}
	jwtToken, err = helpers.TokenValid(strings.TrimSuffix(code, "\\r\\n"))
	log.Println("token valid err", err)
	if err != nil {
		return "", err
	}
	if _, ok := jwtToken["email"]; !ok {
		return "", err
	}
	email = jwtToken["email"].(string)
	log.Println("getJWTToken email", email)

	return email, nil
}

func (handler *RequestAuthorizer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var email string
	var err error
	query := handler.GetQuery()
	email, err = getJWTToken(request)
	log.Println("after email", email)
	ctx := context.Background()
	emailCtx := context.WithValue(ctx, "email", email)
	log.Println("handler.Ctx", handler.Ctx, " email", email, " err", err)
	log.Println("\"handler.Ctx ServeHTTP getJWTToken email", email, " err", err)
	if err != nil {
		RespondWithError(writer, helpers.NewHttpError(http.StatusForbidden, constants.YouHaveNoAccessToResource))
		return
	}
	users, err := models.GetTaUserByEmail(query, types.NewString(email))
	log.Println("models.TaUser{}.GetTaUserByEmail checking err", err)
	if err != nil {
		RespondWithError(writer, helpers.NewHttpError(http.StatusForbidden, constants.YouHaveNoAccessToResource))
	} else if len(users) > 0 {
		taUser := users[0]
		handler.Ctx = context.WithValue(emailCtx, "user_id", types.AsInt(taUser.PkUser))
		log.Println("models.TaUser{}.GetTaUserByEmail checking taUser", taUser)
	} else {
		RespondWithError(writer, helpers.NewHttpError(http.StatusForbidden, constants.YouHaveNoAccessToResource))
	}
}
