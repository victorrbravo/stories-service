package handlers_utils

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"stories/types"
)

func GetPathParameter(key string, request *http.Request) (string, error) {
	var form url.Values
	form = request.Form
	log.Println("url values form", form)
	if values, ok := form[key]; ok {
		if len(values) < 1 {
			return "", fmt.Errorf("`%s' was found in the parameter list, but it's empty", key)
		}
		return values[0], nil
	}
	return "", fmt.Errorf("trying to read path parameter `%s' but it's not there", key)
}

func GetPathParameterAsUUID(key string, request *http.Request) (types.UUID, error) {
	var value string
	var err error
	value, err = GetPathParameter(key, request)
	if err != nil {
		return types.UUID{}, err
	}
	return types.NewUUID(value)
}
