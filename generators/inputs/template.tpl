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