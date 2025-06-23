package pg


import (
    {{if .HasJson}}"encoding/json"{{end}}
    "gopkg.in/yaml.v3"
    "stories/database"
    "stories/generated/models"
)

{{range .Models}}
var Values{{.Name}}  []models.{{.Name}}
{{end}}


func MainProcessor(ctx context.Context, action string, order string, query *database.IQuery) (string, error) {
    var err error
    var value string
    var data []byte
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
                    data, err = json.Marshal(Values{{.Name}})
                    if err != nil {
                        return "", err
                    }
                    value = string(data)
                case "get_{{.Name.Simplified}}_by_{{.KeyField }}":
                    Value{{.Name}}, _ := Get{{.Name}}By{{.KeyField }}(order, query)
                    Values{{.Name}} = append(Values{{.Name}}, *Value{{.Name}})
                {{if len .Searchables}}
                case "get_{{.Name.Simplified}}_by_{{ (index .Searchables 0).Name.Simplified}}":
                    Value{{.Name.Plural}}, errSearch := Get{{.Name}}By{{ (index .Searchables 0).Name}}(order, query)
                    if errSearch != nil  {
                            return value, errSearch
                    }
                    data, err = json.Marshal( Value{{.Name.Plural}} )
                    value = string(data){{end}}{{end}}
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

