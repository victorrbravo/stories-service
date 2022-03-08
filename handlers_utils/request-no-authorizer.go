package handlers_utils

import (
	"net/http"
	"stories/constants"
	"stories/database"
	"stories/helpers"
)

// NoAuthorizedHTTPHandler an interface that adds one parameter to the built-in `ServeHTTP', namely the
// authorized user that previously authenticated through the authentication endpoints
type NoAuthorizedHTTPHandler interface {
	GetQuery() database.IQuery
	IsEnabled() bool
	ServeHTTP(writer http.ResponseWriter, request *http.Request)
}

// RequestAuthorizer is kind of a middleware to help authorizing requests and subsequently calling
// the authorized handler
type RequestNoAuthorizer struct {
	next  NoAuthorizedHTTPHandler
	query database.IQuery
}

// NewRequestAuthorizer creates new authorizer with the given the next handler that should handle
// the request
func NewRequestNoAuthorizer(next NoAuthorizedHTTPHandler) RequestNoAuthorizer {
	return RequestNoAuthorizer{next: next}
}

func (handler RequestNoAuthorizer) GetQuery() database.IQuery {
	return nil
}

func (handler RequestNoAuthorizer) IsEnabled() bool {
	return true
}

func (handler RequestNoAuthorizer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var next NoAuthorizedHTTPHandler
	next = handler.next
	// We finally determined that this token was valid and also have
	// the mobile that owns it, so we can call the handler
	if next.IsEnabled() {
		next.ServeHTTP(writer, request)
	} else {
		RespondWithError(writer, helpers.NewHttpError(http.StatusMethodNotAllowed, constants.ThisEndpointIsDisabled))
	}

}
