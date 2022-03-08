package handlers_utils

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"stories/constants"
	"stories/crypto"
	"stories/helpers"
)

// RequestParameters an interface to validate parameters of a given request in a more controlled way
type RequestParameters interface {
	// IsValid will return an error from the implementation if the parameters fail to validate,
	// the error should explain what went wrong
	IsValid() error
}

func extractAccessTokenFromRequest(request *http.Request) ([]byte, error) {
	var err error
	var uri = request.URL
	var values = uri.Query()
	var keys []string
	var ok bool
	var decoded []byte
	keys, ok = values["token"]
	if !ok || len(keys[0]) < 1 {
		return nil, helpers.NewHttpError(http.StatusNotAcceptable, constants.OnlyOneTokenShouldBeSpecified)
	}
	decoded, err = hex.DecodeString(keys[0])
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

// GetToken extracts the access token from a http request
func GetToken(request *http.Request) (crypto.Token, error) {
	var decoded []byte
	var err error
	decoded, err = extractAccessTokenFromRequest(request)
	if err != nil {
		return crypto.Token{}, err
	}
	return crypto.NewToken(decoded), nil
}

// GetParameters extracts a json object from a http request's body and evaluates the `IsValid()' method
// afterwards returning an error if validation fails
func GetParameters(target RequestParameters, request *http.Request) error {
	var err error
	var content []byte
	var body io.ReadCloser
	// Make a pointer to keep the syntax clean
	body = request.Body
	// Try to read the body
	content, err = ioutil.ReadAll(body)
	if err != nil {
		return helpers.NewHttpError(http.StatusBadRequest, err.Error())
	}
	err = json.Unmarshal(content, &target)
	if err != nil {
		var message string
		// Build the error message.
		//
		// We want to let the caller see the json that they
		// sent, if it's json at all.
		message = fmt.Sprintf("%s: %s\n%s", constants.MalformedJSON, string(content), err.Error())
		// Return the ready made error to send back to the client
		return helpers.NewHttpError(http.StatusBadRequest, message)
	}
	// Validate the parameters according to the validation defined
	// by the mobile
	err = target.IsValid()
	if err != nil {
		return helpers.NewHttpError(http.StatusBadRequest, err.Error())
	}
	return nil
}
