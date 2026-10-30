package tests

import (
    "testing"
{{- if .ImportsTime }}
    "time"
{{- end }}

    "github.com/victorrbravo/stories-service/database"
    "github.com/victorrbravo/stories-service/models"
{{- if .UsesTypes }}
    "github.com/victorrbravo/stories-service/types"
{{- end }}
)

{{- $variable := .Model.ToCamelCase }}
var {{ $variable }}TestID types.UUID
{{- range onlyNonKeys .Attrs }}
var {{ $variable }}Test{{ .Name }} {{ .Type }} = {{ .Random }}
{{- end  }}

func Test{{ .Model }}_Create(t *testing.T) {
    query := database.NewQuery()
    {{ $variable }}, err := models.Create{{ .Model }}(
        query,
        {{- range onlyNonKeys .Attrs }}
        {{ $variable }}Test{{ .Name }},
        {{- end  }}
    )
    if err != nil {
        t.Fatal(err)
    }
    if {{ $variable }} == nil {
        t.Fatal("could not create the `{{ .Model }}'")
    }
    {{ $variable }}TestID = {{ $variable }}.ID
}

func Test{{ .Model }}_GetByID(t *testing.T) {
    query := database.NewQuery()
    {{ $variable }}, err := models.Get{{ .Model }}ByID(query, {{ $variable }}TestID)
    if err != nil {
        t.Fatal(err)
    }
    if !types.UUIDEqual({{ $variable }}.ID, {{ $variable }}TestID) {
        t.Fatal("cannot query the object by it's id")
    }
}

func Test{{ .Model }}_Delete(t *testing.T) {
    query := database.NewQuery()
    {{ $variable }}, err := models.Get{{ .Model }}ByID(query, {{ $variable }}TestID)
    if err != nil {
        t.Fatal(err)
    }
    err = {{ $variable }}.Delete(query)
    if err != nil {
        t.Fatal(err)
    }
}
