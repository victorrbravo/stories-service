package utils

import (
	"bytes"
	"reflect"
	"stories/types"
	"strconv"
	"strings"
	"text/template"
)

func ParseTemplate[T any](tmpls string, data *T) (string, error) {
	var outputBuffer bytes.Buffer
	var value string
	tmpl, err := template.New("Template").Parse(tmpls)
	if err != nil {
		return value, err
	}
	err = tmpl.Execute(&outputBuffer, data)
	if err != nil {
		return value, err
	}
	value = outputBuffer.String()
	return value, nil
}

func SetDefaults[T any](p *T) {
	for i := 0; i < reflect.TypeOf(*p).NumField(); i++ {
		field := reflect.TypeOf(*p).Field(i)
		fieldName := field.Name
		valueOf := reflect.ValueOf(p).Elem()
		valueField := valueOf.FieldByName(fieldName).Interface()
		if value, ok := field.Tag.Lookup("default"); ok {
			currentType := reflect.TypeOf(valueField).String()
			switch currentType {
			case "types.String":
				currentValue := valueField.(types.String)
				currentString := currentValue.String()
				if strings.Compare(currentString, "nil") == 0 {
					valueOf.FieldByName(field.Name).Set(reflect.ValueOf(types.NewString(value)))
				}
			case "int":
				currentValue := valueField.(int)
				if currentValue == 0 {
					valueInt, _ := strconv.Atoi(value)
					valueOf.FieldByName(field.Name).Set(reflect.ValueOf(valueInt))
				}
			}
		}
	}
}
