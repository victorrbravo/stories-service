package handlers_utils

import (
	"fmt"
	"net/http"

	"github.com/victorrbravo/stories-service/helpers"
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

func ResponseWithInternalError(writer http.ResponseWriter) {
	err := helpers.NewHttpError(http.StatusInternalServerError, fmt.Errorf("internal error").Error())
	SendJSONResponse(writer, err, err.Status)
}
