package handlers_utils

import (
	"errors"
	"net/http"
	"stories/database"
	"strings"

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
	query database.IQuery
}

// NewRequestAuthorizer creates new authorizer with the given the next handler that should handle
// the request
func NewRequestAuthorizer(next AuthorizedHTTPHandler) RequestAuthorizer {
	return RequestAuthorizer{next: next}
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

func (handler RequestAuthorizer) GetQuery() database.IQuery {
	return nil
}

// Check if endpoint is enabled
func (handler RequestAuthorizer) IsEnabled() bool {
	return constants.WorkingModeValueNormal == 0
}

func (handler RequestAuthorizer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var next AuthorizedHTTPHandler
	next = handler.next
	var user users.User
	var err error

	user, err = users.DummyUser{}, nil
	if err != nil {
		RespondWithError(writer, helpers.NewHttpError(http.StatusForbidden, constants.YouHaveNoAccessToResource))
	} else if user != nil {
		if next.IsEnabled() {
			next.ServeHTTP(writer, request, user)
		} else {
			RespondWithError(writer, helpers.NewHttpError(http.StatusMethodNotAllowed, constants.ThisEndpointIsDisabled))
		}
	} else {
		RespondWithError(writer, helpers.NewHttpError(http.StatusForbidden, constants.YouHaveNoAccessToResource))
	}
}
