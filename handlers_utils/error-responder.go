package handlers_utils

import (
	"net/http"

	"stories/helpers"
)

// RespondWithError conveniently sends an error to the client discriminating between our special
// HTTPError and the regular go errors
func RespondWithError(writer http.ResponseWriter, err error) {
	switch base := err.(type) {
	case helpers.HTTPError:
		SendJSONResponse(writer, base, base.Status)
	default:
		writer.WriteHeader(http.StatusInternalServerError)
		// Generic error instead
		_, _ = writer.Write([]byte("unknown error"))
	}
}
