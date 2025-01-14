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
	Components   DivStruct
}

var ZeusColorsToHTMLColors = map[uint8]string{
	0:     "#000000",
	Red:   "#ff0000",
	Blue:  "#0000ff",
	Brown: "#804000",
	Black: "#000000",
}

var HTMLCode string

func GenerateHTMLComponent(components DivStruct) string {
	var texts []string
	for i := range components.Components {
		switch component := components.Components[i].(type) {
		case *H1Struct:
			texts = append(texts, fmt.Sprintf("<h1 style=\"color: %s\">%s</h1>", ZeusColorsToHTMLColors[component.Properties.Color], component.Content))
		case *H2Struct:
			texts = append(texts, fmt.Sprintf("<h2 style=\"color: %s\">%s</h2>", ZeusColorsToHTMLColors[component.Properties.Color], component.Content))
		case *H3Struct:
			texts = append(texts, fmt.Sprintf("<h3 style=\"color: %s\">%s</h3>", ZeusColorsToHTMLColors[component.Properties.Color], component.Content))
		case *H4Struct:
			texts = append(texts, fmt.Sprintf("<h4 style=\"color: %s\">%s</h4>", ZeusColorsToHTMLColors[component.Properties.Color], component.Content))
		case *H5Struct:
			texts = append(texts, fmt.Sprintf("<h5 style=\"color: %s\">%s</h5>", ZeusColorsToHTMLColors[component.Properties.Color], component.Content))
		case *H6Struct:
			texts = append(texts, fmt.Sprintf("<h6 style=\"color: %s\">%s</h6>", ZeusColorsToHTMLColors[component.Properties.Color], component.Content))
		case *PStruct:
			texts = append(texts, fmt.Sprintf("<p style=\"color: %s\">%s</p>", ZeusColorsToHTMLColors[component.Properties.Color], component.Content))
		case *ButtonStruct:
			texts = append(texts, fmt.Sprintf("<button style=\"background-color:%s\">%s</button>", ZeusColorsToHTMLColors[component.Properties.Color], component.Text))
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
