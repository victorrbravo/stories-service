package frontend

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
)

type PageObject struct {
	PageName     string
	PageLanguage string
	Components   DivStruct
}

var ZeusColorsToHTMLColors = map[uint8]string{
	Red:   "#ff0000",
	Blue:  "#0000ff",
	Brown: "#804000",
	Black: "#000000",
}

var HTMLCode string

func GenerateHTMLStyle(properties reflect.Value) (styleCode string) {
	styleCode = `style="`
	colorField := properties.FieldByName("Color")
	if colorField.IsValid() && colorField.Kind() == reflect.Uint8 {
		if color, ok := colorField.Interface().(uint8); ok && color != 0 {
			if htmlColor, exists := ZeusColorsToHTMLColors[color]; exists {
				styleCode += fmt.Sprintf(`color:%s; `, htmlColor)
			}
		}
	}

	backgroundColorField := properties.FieldByName("BackgroundColor")
	if backgroundColorField.IsValid() && backgroundColorField.Kind() == reflect.Uint8 {
		if backgroundColor, ok := backgroundColorField.Interface().(uint8); ok && backgroundColor != 0 {
			if htmlBackgroundColor, exists := ZeusColorsToHTMLColors[backgroundColor]; exists {
				styleCode += fmt.Sprintf(`background-color:%s; `, htmlBackgroundColor)
			}
		}
	}
	styleCode = strings.TrimSpace(styleCode)
	if styleCode == `style="` {
		styleCode = ""
		return
	}
	styleCode += `"`
	return
}

func GenerateHTMLComponentProperties[T any](component T) (propertiesString string) {
	propertiesString = ""
	htmlComponentValueOfStruct := reflect.ValueOf(component)
	if htmlComponentValueOfStruct.Kind() != reflect.Struct {
		return
	}

	propertiesField := htmlComponentValueOfStruct.FieldByName("Properties")
	if !propertiesField.IsValid() || propertiesField.Kind() != reflect.Struct {
		return
	}

	propertiesString += GenerateHTMLStyle(propertiesField)
	classField := propertiesField.FieldByName("Class")
	if classField.IsValid() && classField.Kind() == reflect.String {
		if className, ok := classField.Interface().(string); ok && strings.TrimSpace(className) != "" {
			propertiesString += fmt.Sprintf(`class="%s" `, className)
		}
	}

	idField := propertiesField.FieldByName("ID")
	if idField.IsValid() && idField.Kind() == reflect.String {
		if idName, ok := idField.Interface().(string); ok && strings.TrimSpace(idName) != "" {
			propertiesString += fmt.Sprintf(`id="%s" `, idName)
		}
	}

	propertiesString = strings.TrimSpace(propertiesString)
	return
}

func GenerateHTMLComponent(components DivStruct) string {
	var texts []string
	for i := range components.Components {
		switch component := components.Components[i].(type) {
		case *H1Struct:
			texts = append(texts, fmt.Sprintf("<h1 %s>%s</h1>", GenerateHTMLComponentProperties(*component), component.Content))
		case *H2Struct:
			texts = append(texts, fmt.Sprintf("<h2 %s>%s</h2>", GenerateHTMLComponentProperties(*component), component.Content))
		case *H3Struct:
			texts = append(texts, fmt.Sprintf("<h3 %s>%s</h3>", GenerateHTMLComponentProperties(*component), component.Content))
		case *H4Struct:
			texts = append(texts, fmt.Sprintf("<h4 %s>%s</h4>", GenerateHTMLComponentProperties(*component), component.Content))
		case *H5Struct:
			texts = append(texts, fmt.Sprintf("<h5 %s>%s</h5>", GenerateHTMLComponentProperties(*component), component.Content))
		case *H6Struct:
			texts = append(texts, fmt.Sprintf("<h6 %s>%s</h6>", GenerateHTMLComponentProperties(*component), component.Content))
		case *PStruct:
			texts = append(texts, fmt.Sprintf("<p %s>%s</p>", GenerateHTMLComponentProperties(*component), component.Content))
		case *ButtonStruct:
			texts = append(texts, fmt.Sprintf("<button %s>%s</button>", GenerateHTMLComponentProperties(*component), component.Text))
		case *DivStruct:
			texts = append(texts, fmt.Sprintf("<div>\n%s\n</div>", GenerateHTMLComponent(*component)))
		default:
		}
	}
	return strings.Join(texts, "\n")
}

func GenerateHTMLCode(pageObject PageObject) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	file, err := os.ReadFile(filepath.Join(dir, "code.tmpl"))
	if err != nil {
		panic(err)
	}
	tmpl, err := template.New("code").Funcs(template.FuncMap{
		"generateHTMLComponent": GenerateHTMLComponent,
	}).Parse(string(file))
	if err != nil {
		panic(err)
	}
	fileIndex, err := os.Create("./index.html")
	if err != nil {
		panic(err)
	}
	defer fileIndex.Close()
	err = tmpl.Execute(fileIndex, pageObject)
	if err != nil {
		panic(err)
	}
}
