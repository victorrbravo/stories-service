package pg

import (
    "stories/generated/models"
    "stories/database"
{{- if .UsesTypes }}
	"stories/types"
{{- end }}
)

type {{ .Model }} struct {
{{- range .Attrs }}
	{{ .Name }} {{ .Type }} {{ .Name.ToYamlTag }}
{{- end }}
}

func  Insert{{ .Model }}(order string, query *database.IQuery) error {
    handler, err := Unmarshall{{ .Model }}([]byte(order))
    if err != nil {
        return err
    }
    _, err = models.Create{{ .Model }}(*query,
        {{- range $index, $value := onlyNonKeys .Attrs }}
            handler.{{ $value.Name }},
        {{- end }}
    )
    if err != nil {
        return err
    }
    return nil
}

func  Delete{{ .Model }}(order string, query *database.IQuery) error {
    handler, err := Unmarshall{{ .Model }}([]byte(order))
    if err != nil {
        return err
    }
    {{- range $index, $value := onlyKeys .Attrs }}
    {{ $.Model.ToCamelCase }}, err := models.Get{{ $.Model }}By{{ $value.Name }}(*query,
            handler.{{ $value.Name }},
    )
    {{- end }}

    if err != nil {
        return err
    }
    err = {{ $.Model.ToCamelCase }}.Delete(*query)
    if err != nil {
        return err
    }
    return nil
}

{{- range $index, $value := onlyKeys .Attrs }}
func  Get{{ $.Model }}By{{ $value.Name }}(order string, query *database.IQuery) (*models.{{ $.Model }}, error) {
{{- end }}
    handler, err := Unmarshall{{ .Model }}([]byte(order))
    if err != nil {
        return nil, err
    }
    {{- range $index, $value := onlyKeys .Attrs }}
    {{ $.Model.ToCamelCase }}, err := models.Get{{ $.Model }}By{{ $value.Name }}(*query,
            handler.{{ $value.Name }},
    )
    {{- end }}
    if err != nil {
        return {{ $.Model.ToCamelCase }}, err
    }
    return {{ $.Model.ToCamelCase }}, nil
}


func  Update{{ .Model }}(order string, query *database.IQuery) error {
    handler, err := Unmarshall{{ .Model }}([]byte(order))
    if err != nil {
        return err
    }
    {{- range $index, $value := onlyKeys .Attrs }}
    {{ $.Model.ToCamelCase }}, err := models.Get{{ $.Model }}By{{ $value.Name }}(*query,
            handler.{{ $value.Name }},
    )
    {{- end }}
    if err != nil {
        return err
    }
    err = {{ $.Model.ToCamelCase }}.Update(*query,
    {{- range $index, $value := onlyNonKeys .Attrs }}
            handler.{{ $value.Name }},
        {{- end }}
    )
    if err != nil {
        return err
    }
    return nil
}


func  GetAll{{ .Model }}(query *database.IQuery) ([]models.{{ .Model }}, error) {

    {{ .Model.ToCamelCase }}, err := models.GetAll{{ .Model.Plural }}(*query)
    if err != nil {
        return {{ .Model.ToCamelCase }}, err
    }
    return {{ .Model.ToCamelCase }}, nil
}



{{- range $index, $value := onlySearchables .Attrs }}
func  Get{{ $.Model }}By{{ $value.Name }}(order string, query *database.IQuery) ([]models.{{ $.Model }}, error) {
    handler, err := Unmarshall{{ $.Model }}([]byte(order))
    if err != nil {
        return nil, err
    }
    {{ $.Model.ToCamelCasePlural }}, err := models.Get{{ $.Model }}By{{ $value.Name }}(*query,
            handler.{{ $value.Name }},
    )
    if err != nil {
        return {{ $.Model.ToCamelCasePlural }}, err
    }
    return {{ $.Model.ToCamelCasePlural }}, nil
}
{{- end }}

