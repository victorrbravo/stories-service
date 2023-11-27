package pg


import (
    "gopkg.in/yaml.v3"
    "stories/generated/models"
)

{{range .Models}}
var Values{{.}}  []models.{{.}}
{{end}}

func ChooseProcessor(action string, order string, query *database.IQuery) error {
	switch action {
    {{range .Models}}
    			case "add_{{.Simplified}}":
    			    Insert{{.}}(order, query)
    			case "delete_{{.Simplified}}":
                    Delete{{.}}(order, query)
    			case "update_{{.Simplified}}":
                    Update{{.}}(order, query)
                case "all_{{.Simplified}}":
                    Values{{.}}, _ = GetAll{{.}}(query)
    {{end}}
	}
	return nil
}
{{range .Models}}

func Unmarshall{{.}}(order []byte) ({{.}}, error) {
	var {{ .ToCamelCase }} {{.}}
	var err error
	err = yaml.Unmarshal(order, &{{.ToCamelCase}})
	if err != nil {
		return {{.ToCamelCase}}, err
	}
	return {{.ToCamelCase}}, nil
}

{{end}}

