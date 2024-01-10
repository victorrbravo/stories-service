package pg


import (
    "gopkg.in/yaml.v3"
    "stories/database"
    "stories/generated/models"
)

{{range .Models}}
var Values{{.Name}}  []models.{{.Name}}
{{end}}

func MainProcessor(action string, order string, query *database.IQuery) (string, error) {
    var err error
    var value string
    {{range .Models}}
    if Values{{.Name}}  == nil {
        Values{{.Name}} = make([]models.{{.Name}}, 0)
    }
    {{end}}
	switch action {
    {{range .Models}}
    			case "add_{{.Name.Simplified}}":
    			    err = Insert{{.Name}}(order, query)
    			case "delete_{{.Name.Simplified}}":
                    err = Delete{{.Name}}(order, query)
    			case "update_{{.Name.Simplified}}":
                    err = Update{{.Name}}(order, query)
                case "all_{{.Name.Simplified}}":
                    Values{{.Name}}, _ = GetAll{{.Name}}(query)
                case "get_{{.Name.Simplified}}_by_{{.KeyField }}":
                    Value{{.Name}}, _ := Get{{.Name}}By{{.KeyField }}(order, query)
                    Values{{.Name}} = append(Values{{.Name}}, *Value{{.Name}})
    {{end}}
	}
    if err != nil  {
        return value, err
    }
	return value, nil
}
{{range .Models}}

func Unmarshall{{.Name}}(order []byte) ({{.Name}}, error) {
	var {{ .Name.ToCamelCase }} {{.Name}}
	var err error
	err = yaml.Unmarshal(order, &{{.Name.ToCamelCase}})
	if err != nil {
		return {{.Name.ToCamelCase}}, err
	}
	return {{.Name.ToCamelCase}}, nil
}

{{end}}

