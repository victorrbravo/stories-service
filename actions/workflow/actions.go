package workflow

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"reflect"
	"stories/database"
	"stories/types"
	"strings"
	"text/template"
)

const (
	InsertTaskInFlow = `
UPDATE ta_flow SET flow=jsonb_set(flow, '{tasks, {{.Name}} }', '{{.NewTask}}',true)
WHERE pk_flow={{.KeyValue}}
`
	DeleteTaskInFlow = `
UPDATE ta_flow SET flow=flow #- '{tasks, {{.Name}} }'
WHERE pk_flow={{.KeyValue}}
`
	emptyWorkflow = `
{
  "id": "3f12f5a5-57ac-4861-8771-f85de203f863",
  "name": "Workflow1",
  "variable": { "name": "basic", "queries": "pk_story,name,status" },
  "tracking_table_name": "ta_tracking",
  "tracking_keys": "fk_story,fk_status",
  "field_token": "pk_story",
  "table_name": "ta_story",
  "initial": "fk_status='start'",
  "tasks": []
}
`
	tmplTask = `
{
"info_field": "{{.RoleField}}", 
"time_field": "{{.TimeField}}", 
"transition": {
"links": [
	{
	"task": "{{.Target}}", 
	"queries": "{{.Sentence}}", 
	"pattern": "{{.Pattern}}"
	}
}
`
)

type AddTask struct {
	Name      types.String `yaml:"name" default:"newtask"`
	Sentence  types.String `yaml:"sentence"`
	Target    types.String `yaml:"target"`
	RoleField types.String `yaml:"role_field" default:"roles[1]"`
	TimeField types.String `yaml:"time_field" default:"updated_at"`
	Pattern   types.String `yaml:"pattern" default:"none"`
}

func Processor(action string, order string, query *database.IQuery) error {
	var err error
	switch action {
	case "add_task":
		err = DoAddTask(order, query)
	}
	if err != nil {
		return err
	}
	return nil
}

func setDefaults(p *AddTask) {
	// Iterate over the fields of the Person struct using reflection
	// and set the default value for each field if the field is not provided
	// by the caller of the constructor function.
	for i := 0; i < reflect.TypeOf(*p).NumField(); i++ {
		field := reflect.TypeOf(*p).Field(i)
		fieldName := field.Name
		valueOf := reflect.ValueOf(p).Elem()
		valueField := valueOf.FieldByName(fieldName).Interface()
		if value, ok := field.Tag.Lookup("default"); ok {
			currentType := reflect.TypeOf(valueField).String()
			switch currentType {
			case "types.String":
				currentValue := valueField.(types.String)
				currentString := currentValue.String()
				if strings.Compare(currentString, "nil") == 0 {
					valueOf.FieldByName(field.Name).Set(reflect.ValueOf(types.NewString(value)))
				}

			}
		}
	}
}

func UnmarshallAddTask(order []byte) (AddTask, error) {
	var addTask AddTask
	var err error
	err = yaml.Unmarshal(order, &addTask)
	if err != nil {
		return addTask, err
	}
	setDefaults(&addTask)
	return addTask, nil
}

func DoAddTask(order string, query *database.IQuery) error {
	task, err := UnmarshallAddTask([]byte(order))
	var outputBuffer bytes.Buffer
	if err != nil {
		return err
	}
	tmpl, err := template.New("doAddTask").Parse(tmplTask)
	if err != nil {
		return err
	}

	tmpl.Execute(&outputBuffer, task)

	if err != nil {
		return err
	}
	newTask := outputBuffer.String()
	fmt.Println("newTask", newTask)
	return nil
}
