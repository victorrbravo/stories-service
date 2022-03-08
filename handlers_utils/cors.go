package handlers_utils

import "net/http"

type CORS struct {
	handler http.Handler
}

func NewCORS(handler http.Handler) CORS {
	return CORS{handler}
}

func (cors CORS) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var header http.Header = writer.Header()
	/// FIXME: must use a value from the settings or config file
	//         instead if *
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization")
	header.Add("Access-Control-Allow-Credentials", "true")
	// Call the handler now
	cors.handler.ServeHTTP(writer, request)
}
