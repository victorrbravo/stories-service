package utils

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"reflect"
	"regexp"
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

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
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

func FormatYaml(input string) string {
	// Define the regular expression pattern
	pattern := `(\r?\n)?\w+:`

	// Compile the regular expression
	regexpPattern := regexp.MustCompile(pattern)

	// Find the matching expression in the input string
	results := regexpPattern.FindAll([]byte(input), -1)
	fReplace := func(x []byte) string {
		pattern := `^\r?\n`
		regexpPattern := regexp.MustCompile(pattern)
		value := regexpPattern.ReplaceAllString(string(x), "")
		value = "\r\n" + value
		return value
	}
	values := Map(results, fReplace)
	for i, v := range results {
		input = strings.ReplaceAll(input, string(v), values[i])
	}

	return input
}

func EncodeImageBase64(path string) (string, error) {
	imageData, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	// Convert the image data to base64
	base64String := base64.StdEncoding.EncodeToString(imageData)
	return base64String, nil
}
