package web

import "strings"

func GetParametersFromURI(uri string) (string, string) {
	var value string
	var action string
	var values []string
	values = strings.Split(uri, "/")
	if len(values) > 5 {
		value = values[4]
		action = values[5]
	}
	return value, action
}
