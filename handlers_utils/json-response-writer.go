package handlers_utils

import (
	"encoding/json"
	"net/http"

	"stories/helpers"
)

// SendJSONResponse sends a json response from a generic interface also setting the
// status code and the header to `application/json'
func SendJSONResponse(writer http.ResponseWriter, object interface{}, status int) {
	var response []byte
	var err error
	response, err = json.Marshal(object)
	if err != nil {
		RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
	} else {
		var header = writer.Header()
		// Set the correct content type
		header.Set("Content-Type", "application/json; charset=utf-8")
		// Finally start writing the response
		writer.WriteHeader(status)
		_, err = writer.Write(response)
		if err != nil {
			println(err.Error())
		}
	}
}
