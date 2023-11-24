package pg


import (
    "gopkg.in/yaml.v3"
    "stories/generated/models"
)

func ChooseProcessor(action string, order string, query *database.IQuery) error {
	switch action {
    {{range .Models}}
    			case "add_{{.Simplified}}":
    			Insert{{.}}(order, query)
    {{end}}
	}
	return nil
}
{{range .Models}}

func Unmarshall{{.}}(order []byte) (models.{{.}}, error) {
	var {{ .ToCamelCase }} models.{{.}}
	var err error
	err = yaml.Unmarshal(order, &{{.ToCamelCase}})
	if err != nil {
		return {{.ToCamelCase}}, err
	}
	return {{.ToCamelCase}}, nil
}

{{end}}

