package pg

import (
    "stories/generated/models"
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
    user, err := models.Get{{ $.Model }}By{{ $value.Name }}(*query,
            handler.{{ $value.Name }},
    )
    {{- end }}

    if err != nil {
        return err
    }
    err = user.Delete(*query)
    if err != nil {
        return err
    }
    return nil
}


func  Update{{ .Model }}(order string, query *database.IQuery) error {
    handler, err := Unmarshall{{ .Model }}([]byte(order))
    if err != nil {
        return err
    }
    {{- range $index, $value := onlyKeys .Attrs }}
    user, err := models.Get{{ $.Model }}By{{ $value.Name }}(*query,
            handler.{{ $value.Name }},
    )
    {{- end }}
    if err != nil {
        return err
    }
    err = user.Update(*query,
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

    users, err := models.GetAll{{ .Model.Plural }}(*query)
    if err != nil {
        return users, err
    }
    return users, nil
}