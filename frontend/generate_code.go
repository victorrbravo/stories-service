package frontend

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type PageObject struct {
	PageName     string
	PageLanguage string
	Components   ViewStruct
}

var ZeusColorsToHTMLColors = map[uint8]string{
	Red:   "#ff0000",
	Blue:  "#0000ff",
	Brown: "#804000",
	Dark:  "#000000",
}

var HTMLCode string

func GenerateHTMLComponent(components ViewStruct) string {
	var texts []string
	for i := range components.Components {
		switch component := components.Components[i].(type) {
		case *TextStruct:
			texts = append(texts, fmt.Sprintf("<p style=\"color: %s\">%s</p>", ZeusColorsToHTMLColors[component.Properties.Color], component.Content))
		case *ButtonStruct:
			texts = append(texts, fmt.Sprintf("<button style=\"background-color:%s\">%s</button>", ZeusColorsToHTMLColors[component.Properties.Color], component.Text))
		case *ViewStruct:
			texts = append(texts, fmt.Sprintf("<article>\n%s\n</article>", GenerateHTMLComponent(*component)))
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
