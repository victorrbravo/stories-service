package models
import (
	"database/sql"

{{- if .NeedsPostgis }}
    "github.com/cridenour/go-postgis"
{{- end }}

{{- if .UsesEnums }}
	"stories/enums"
{{- end }}
	"stories/database"
	"stories/scanners"
{{- if .UsesTypes }}
	"stories/types"
{{- end }}
)
/*
 * Generated queries (these can be used perfectly instead of writing them by hand
 */
const (
    {{ .Model }}CreateQuery = "INSERT INTO {{ .Schema }}\"{{ .TableName.ToSnakeCases }}\" ({{- range $index, $attr := onlyNonKeys .Attrs }}{{- if $index }}, {{- end }}{{ checkUpdateOrDelete $attr.FieldName $attr.Name }}{{- end }}) VALUES ({{- range $index, $attr := onlyNonKeys .Attrs }}{{- if $index }}, {{ end }}${{ add $index 1 }}{{- end }}) RETURNING ({{- range $index, $attr := onlyKeys .Attrs }}{{- if $index }}, {{- end }}{{ checkUpdateOrDelete $attr.FieldName $attr.Name }}{{- end }})"
    {{- if checkIsNotEmpty .OrderBy }}
    All{{ .Model }}Query = "SELECT {{- range $index, $attr := .Attrs }} {{- if $index }},{{- end }} {{ checkForSelect .FieldName .Name }} {{- end }} FROM {{ .Schema }}\"{{ .TableName.ToSnakeCases }}\" ORDER BY \"{{ .OrderBy.ToSnakeCase }}\" {{ .Order }}"
    {{- else }}
    All{{ .Model }}Query = "SELECT {{- range $index, $attr := .Attrs }} {{- if $index }},{{- end }} {{ checkForSelect .FieldName .Name }} {{- end }} FROM {{ .Schema }}\"{{ .TableName.ToSnakeCases }}\""
    {{- end }}
    {{- range .Attrs }}
    {{- if .Key }}
    {{ $.Model }}By{{ .Name }}Query = "SELECT {{- range $index, $attr := $.Attrs }} {{- if $index }},{{- end }} {{ (checkForSelect .FieldName .Name) }} {{- end }} FROM {{ $.Schema }}\"{{ $.TableName.ToSnakeCases }}\" WHERE {{ checkUpdateOrDelete  .FieldName .Name }} = $1"
    {{- end }}
    {{- end }}
    {{- range .Attrs }}
    {{- if .KeyForMany }}
    {{ $.Model }}By{{ .Name }}Query = "SELECT {{- range $index, $attr := $.Attrs }} {{- if $index }},{{- end }} {{ (checkForSelect .FieldName .Name) }} {{- end }} FROM {{ $.Schema }}\"{{ $.TableName.ToSnakeCases }}\" WHERE {{ checkUpdateOrDelete  .FieldName .Name }} = $1"
    {{- end }}
    {{- end }}
    {{- range .Attrs }}
    {{- if .Searchable }}
    {{ $.Model }}By{{ .Name }}Query = "SELECT {{- range $index, $attr := $.Attrs }} {{- if $index }},{{- end }} {{ checkForSelect .FieldName .Name }} {{- end }} FROM {{ $.Schema }}\"{{ $.TableName.ToSnakeCases }}\" WHERE {{ checkUpdateOrDelete  .FieldName .Name }} = $1"
    {{- end }}
    {{- end }}
    {{- range .Attrs }}
    {{- if .Discrete }}
    Get{{ $.Model }}CategoriesFrom{{ .Name }}Query = "SELECT {{ checkForSelect .FieldName .Name }} FROM {{ $.Schema }}\"{{ $.TableName.ToSnakeCases }}\" GROUP BY {{ checkForSelect .FieldName .Name }}"
    Search{{ $.Model }}CategoriesFrom{{ .Name }}Query = "SELECT {{- range $index, $attr := $.Attrs }} {{- if $index }},{{- end }} {{ checkForSelect .FieldName .Name }} {{- end }} FROM {{ $.Schema }}\"{{ $.TableName.ToSnakeCases }}\" WHERE {{ checkUpdateOrDelete  .FieldName .Name }} IN ($1)"
    {{- end }}
    {{- end }}
    {{ .Model }}UpdateQuery = "UPDATE {{ .Schema }}\"{{ .TableName.ToSnakeCases }}\" SET {{- range $index, $attr := onlyNonKeys .Attrs }} {{- if $index }},{{- end }} {{ checkUpdateOrDelete $attr.FieldName $attr.Name }} = COALESCE(${{ add $index 1 }}, {{ checkUpdateOrDelete $attr.FieldName $attr.Name }}) {{- end }} WHERE {{- $count := len (onlyNonKeys .Attrs) -}} {{ range $index, $attr := onlyKeys .Attrs }} {{- if $index }}AND{{- end }} {{ checkUpdateOrDelete $attr.FieldName $attr.Name }} = ${{ add $index $count 1 }} {{- end }}"
    {{ .Model }}DeleteQuery = "DELETE FROM {{ .Schema }}\"{{ .TableName.ToSnakeCases }}\" WHERE {{- range $index, $attr := onlyKeys .Attrs }} {{- if $index }}AND{{- end }} \"{{ $attr.Name.ToSnakeCase  }}\" = ${{ add $index 1 }} {{- end }}"
    {{- if .Searcher }}
    {{ $.Model }}By{{- range $index, $attr := .Searcher }}{{- if $index }}And{{- end }}{{ .Name}}{{- end }}Query = "SELECT {{- range $index, $attr := $.Attrs }} {{- if $index }},{{- end }} {{ checkForSelect .FieldName .Name }} {{- end }} FROM {{ $.Schema }}\"{{ $.TableName.ToSnakeCases }}\" WHERE {{- $count := 0 -}} {{ range $index, $attr := .Searcher }} {{- if $index }} AND {{- end }} {{ checkUpdateOrDelete $attr.FieldName $attr.Name }} = ${{ add $index $count 1 }} {{- end }}"
    {{- end }}
)

type {{ .Model }} struct {
{{- range .Attrs }}
	{{ .Name }} {{ .Type }} {{ .Name.ToJsonTag }}
{{- end }}
}

{{- $receiver := .Model.ToCamelCase }}

func ({{ $receiver }} *{{ .Model }}) ScanRow(rows database.RowInterface) (interface{}, error) {
	var err error
	err = rows.Scan(
		{{- range .Attrs }}
		&{{ $receiver }}.{{ .Name }},
		{{- end }}
	)
	if err != nil {
		return nil, err
	}
	return *{{ $receiver }}, nil
}

func Create{{ .Model }}(
    query database.IQuery,
{{- range $index, $value := onlyNonKeys .Attrs }}
	{{ $value.Name.ToCamelCase }} {{ $value.Type }},
{{- end }}
) (*{{ .Model }}, error) {
	var err error
    var scanner scanners.{{ .KeyType }}
	var {{ $receiver }} *{{ .Model }}
	err = query.QueryOne(
	    {{ .Model }}CreateQuery,
	    &scanner,
		{{- range .Attrs }}
		{{- if  .Editable }}
		{{ .Name.ToCamelCase }},
		{{- end }}
		{{- if not .Key }}
		{{ .Name.ToCamelCase }},
		{{- end }}
		{{- end }}
	)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}
	{{ $receiver }} = &{{ .Model }}{
	    {{- range onlyKeys .Attrs }}
	    {{- if not .Editable }}
	    {{ .Name }}: scanner.Value,
	    {{- end }}
	    {{- end }}
	    {{- range onlyNonKeys .Attrs }}
	    {{ .Name }}: {{ .Name.ToCamelCase }},
	    {{- end }}
	}
	return {{ $receiver }}, nil
}

func ({{ $receiver }} *{{.Model}}) Update(
    query database.IQuery,
{{- range onlyNonKeys .Attrs }}
    {{ .Name.ToCamelCase }} {{.Type}},
{{- end }}
) error {
    var err error
    {{- if not .Attrs }}
    return nil
    {{- else }}
	err = query.Exec(
	    {{ .Model }}UpdateQuery,
        {{- range onlyNonKeys .Attrs }}
        {{ .Name.ToCamelCase }},
        {{- end }}
        {{- range onlyKeys .Attrs }}
        {{ $receiver }}.{{ .Name }},
        {{- end }}
    )
	{{- end }}
	if err != nil {
    	return err
	}
    {{- range onlyNonKeys .Attrs }}
    {{ $receiver }}.{{ .Name }} = {{ .Name.ToCamelCase }}
    {{- end }}
    return nil
}
{{ range .Attrs }}
{{- if .Key }}
func Get{{ $.Model }}By{{- .Name -}}(query database.IQuery, value {{ .Type }}) (*{{ $.Model }}, error) {
	var err error
	var {{ $receiver }} {{ $.Model }}
	err = query.QueryOne({{ $.Model }}By{{ .Name }}Query, &{{ $receiver }}, value)
	if err != nil {
		return  nil, err
	}
	return &{{ $receiver }}, nil
}
{{- end }}
{{- end }}

{{ range .Attrs }}
{{- if .KeyForMany }}
func Get{{ $.Model }}By{{- .Name -}}(query database.IQuery, value {{ .Type }}) ([]{{ $.Model }}, error) {
	var err error
	var rows []interface{}
	var list []{{ $.Model }}
	list = make([]{{ $.Model }}, 0)
	rows, err = query.QueryMany({{ $.Model }}By{{ .Name }}Query, &{{ $.Model }}{}, value)
	if err != nil {
		return list, err
	}
	for _, row := range rows {
		list = append(list, row.({{ $.Model }}))
	}
	return list, nil
}

{{- end }}
{{- end }}

{{- if .Searcher }}

func Get{{ $.Model }}By{{- range $index, $attr := .Searcher }}{{- if $index }}And{{- end }}{{ .Name}}{{- end }}(query database.IQuery,{{- range .Searcher }}
{{ .Name.ToCamelCase }} {{.Type}},
{{- end }}
) (*{{ $.Model }}, error) {
	var err error
	var {{ $receiver }} {{ $.Model }}
	err = query.QueryOne({{ $.Model }}By{{- range $index, $attr := .Searcher }}{{- if $index }}And{{- end }}{{ .Name}}{{- end }}Query, &{{ $receiver }},{{- range .Searcher }}
    {{ .Name.ToCamelCase }},
{{- end }}
)
	if err != nil {
		return  nil, err
	}
	return &{{ $receiver }}, nil
}
{{- end }}

{{ range .Attrs }}
{{- if .Searchable }}
func Get{{ $.Model }}By{{- .Name -}}(query database.IQuery, value {{ .Type }}) (*{{ $.Model }}, error) {
	var err error
	var {{ $receiver }} {{ $.Model }}
	err = query.QueryOne({{ $.Model }}By{{ .Name }}Query, &{{ $receiver }}, value)
	if err != nil {
		return  nil, err
	}
	return &{{ $receiver }}, nil
}
{{- end }}
{{- end }}

{{- if .AddPaging }}
func GetAll{{ .Model.Plural }}Paged(query database.IQuery, pageBase types.PageBase) (types.Page, error) {
	var err error
	var page types.Page
	page, err = query.QueryManyPaged(All{{ .Model }}Query, &{{ .Model }}{}, pageBase.PageNumber, pageBase.PageSize)
	if err != nil {
    		return  page, err
    }
	return page, nil
}
{{- end }}

{{ range .Attrs }}
{{- if .Discrete }}
func Get{{ $.Model }}CategoriesFrom{{- .Name -}}(query database.IQuery) ([]string, error) {
	var err error
	var rows []interface{}
	var {{ $receiver }} scanners.String
	var list []string
	rows, err = query.QueryMany(Get{{ $.Model }}CategoriesFrom{{ .Name }}Query, &{{ $receiver }})
	if err != nil {
		return  nil, err
	}
	list = make([]string, 0)
    for _, row := range rows {
		list = append(list, row.(string))
	}
	return list, nil
}


func Search{{ $.Model }}CategoriesFrom{{- .Name -}}(query database.IQuery, values types.StringArray, pageBase types.PageBase) (types.Page, error) {
	var err error
	var page types.Page
	var {{ $receiver }} {{ $.Model }}
	page, err = query.QueryManyPaged(Search{{ $.Model }}CategoriesFrom{{ .Name }}Query,
	 &{{ $receiver }},
	 pageBase.PageNumber,
	  pageBase.PageSize,
	  )
	if err != nil {
		return  page, err
	}
	return page, nil
}
{{- end }}
{{- end }}

func ({{ $receiver }} *{{.Model}}) Delete(query database.IQuery) error {
	var err error
	err = query.Exec(
	    {{ .Model }}DeleteQuery,
	    {{- range onlyKeys .Attrs }}
        {{ $receiver }}.{{ .Name }},
        {{- end }}
	)
	if err != nil {
		return err
	}
	return nil
}

func GetAll{{ .Model.Plural }}(query database.IQuery) ([]{{ .Model }}, error) {
	var err error
	var rows []interface{}
	var list []{{ .Model }}
	rows, err = query.QueryMany(All{{ .Model }}Query, &{{ .Model }}{})
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]{{ .Model }}, 0), nil
		}
		return nil, err
	}
	for _, row := range rows {
		list = append(list, row.({{ .Model }}))
	}
	return list, nil
}
