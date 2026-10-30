package main

import (
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/victorrbravo/stories-service/generators"
)

const Template string = `
package queries
// Generated file
// ALL EDITS WILL BE LOST

// Automatically generated queries
const (
{{- range .Queries }}
	{{ printf "%s = ` + "`%s`" + `" .Name .Sql }}
{{- end }} 
)
`

type query struct {
	Name string
	Sql  string
}

type queryHandler struct {
	template            *template.Template
	multipleSpacesRegex *regexp.Regexp
	operatorSpacesRegex *regexp.Regexp
	queries             []query
	sourcePath          string
	targetPath          string
}

func (handler *queryHandler) SetSourcePath(sourcePath string) {
	handler.sourcePath = sourcePath
}

func (handler *queryHandler) GetSourcePath() string {
	return handler.sourcePath
}

func (handler *queryHandler) GetTargetPath() string {
	return handler.targetPath
}

func (handler *queryHandler) Prepare() {
}

func (handler *queryHandler) Write() error {
	var file *os.File
	tpl := handler.template
	var err error
	file, err = os.Create(path.Join(handler.GetTargetPath(), "cmd.go"))
	if err != nil {
		return err
	}
	return tpl.Execute(file, struct {
		Queries []query
	}{
		Queries: handler.queries,
	})
}

func (handler *queryHandler) HandleFile(file os.FileInfo) error {
	// Just to avoid ugly syntax
	multipleSpacesRegex := handler.multipleSpacesRegex
	operatorSpacesRegex := handler.operatorSpacesRegex
	// Declare local stuff
	var err error
	var sql string
	var content string
	var name string

	name, err = generators.GetTypeName(file.Name())
	if err != nil {
		return err
	}
	content, err = generators.GetContent(handler.sourcePath, file.Name())
	if err != nil {
		return err
	}

	sql = strings.Replace(content, "\t\n\r", " ", -1)
	sql = multipleSpacesRegex.ReplaceAllString(sql, " ")
	sql = operatorSpacesRegex.ReplaceAllString(sql, "$1")
	sql = strings.Trim(sql, " \t\n\r")

	handler.queries = append(handler.queries, query{Name: name, Sql: sql})
	return nil
}

func main() {
	sourcePath := "../../data-model/queries"
	targetPath := "../../generated/queries"
	var err error
	parser := template.New("")
	handler := &queryHandler{
		sourcePath:          sourcePath,
		targetPath:          targetPath,
		template:            template.Must(parser.Parse(Template)),
		multipleSpacesRegex: regexp.MustCompile(`\s+`),
		operatorSpacesRegex: regexp.MustCompile(`(?:([()])\s+)`),
	}
	// Let the handler prepare itself
	handler.Prepare()
	// Generate
	err = generators.Generate(handler)
	if err != nil {
		log.Println(err)
	}
	err = handler.Write()
	if err != nil {
		log.Fatalln(err)
	}
}
