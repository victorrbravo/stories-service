package workflow

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"stories/database"
	"stories/pkg/utils"
	"stories/types"
)

const (
	InsertTaskInFlow = `
WITH CTE AS (
SELECT 
'{
"info_field": "{{.RoleField}}", 
"time_field": "{{.TimeField}}", 
"report": {{.Report}},
"transition": {
"links": [
	{
	"task": ' || (flow->'tasks'->'{{.Source}}'->'transition'->'links'->0->'task')::text || ', 
	"queries": ' || (flow->'tasks'->'{{.Source}}'->'transition'->'links'->0->'queries')::text || ', 
	"pattern": "{{.Pattern}}"
	}
  ]
}
}' AS newtask
FROM ta_flow
WHERE pk_flow={{.PkFlow}}
)
UPDATE ta_flow SET flow=jsonb_set(flow::jsonb, '{tasks, {{.Name}} }'::text[], CTE.newtask::jsonb,true) FROM CTE
WHERE pk_flow={{.PkFlow}};
UPDATE ta_flow SET flow=jsonb_set(flow, '{tasks, {{.Source}},transition,links,0,task }', '"{{.Name}}"',true)
WHERE pk_flow={{.PkFlow}};
UPDATE ta_flow SET flow=jsonb_set(flow, '{tasks, {{.Source}},transition,links,0,queries }', '"{{.Sentence}}"',true)
WHERE pk_flow={{.PkFlow}};
`
	EndInsertTaskInFlow = `
WITH CTE AS (
SELECT 
'{
"info_field": "{{.RoleField}}", 
"time_field": "{{.TimeField}}", 
"report": {{.Report}},
"transition": {
"links": [
	{
	"task": "END", 
	"queries": "END", 
	"pattern": "{{.Pattern}}"
	}
  ]
}
}' AS newtask
FROM ta_flow
WHERE pk_flow={{.PkFlow}}
)
UPDATE ta_flow SET flow=jsonb_set(flow::jsonb, '{tasks, {{.Name}} }'::text[], CTE.newtask::jsonb,true) FROM CTE
WHERE pk_flow={{.PkFlow}};
UPDATE ta_flow SET flow=jsonb_set(flow, '{tasks, {{.Source}},transition,links,0,task }', '"{{.Name}}"',true)
WHERE pk_flow={{.PkFlow}};
UPDATE ta_flow SET flow=jsonb_set(flow, '{tasks, {{.Source}},transition,links,0,queries }', '"{{.Sentence}}"',true)
WHERE pk_flow={{.PkFlow}};
`
	StartInsertTaskInFlow = `
WITH CTE AS (
SELECT 
'{
"info_field": "{{.RoleField}}", 
"time_field": "{{.TimeField}}", 
"report": {{.Report}},
"transition": {
"links": [
	{
	"task": "{{.Target}}", 
	"queries": ' || (flow->'initial')::text || ', 
	"pattern": "{{.Pattern}}"
	}
  ]
}
}' AS newtask
FROM ta_flow
WHERE pk_flow={{.PkFlow}}
)
UPDATE ta_flow SET flow=jsonb_set(flow::jsonb, '{tasks, {{.Name}} }'::text[], CTE.newtask::jsonb,true) FROM CTE
WHERE pk_flow={{.PkFlow}};
UPDATE ta_flow SET flow=jsonb_set(flow, '{initial}', '"{{.Sentence}}"',true)
WHERE pk_flow={{.PkFlow}};
`
	DeleteTaskInFlow = `
UPDATE ta_flow SET flow=flow #- '{tasks, {{.Name}} }'
WHERE pk_flow={{.PkFlow}}
`
	EmptyWorkflow = `
{
  "id": "3f12f5a5-57ac-4861-8771-f85de203f863",
  "name": "Workflow1",
  "variable": { "name": "basic", "queries": "pk_story,name,status" },
  "tracking_table_name": "ta_tracking",
  "tracking_keys": "fk_story,fk_status",
  "field_token": "pk_story",
  "field_status": "status",
  "table_name": "ta_story",
  "initial": "TRUE",
"tasks": {
	"Start" :{
		"report": false,
		"info_field": "roles[1]", 
		"time_field": "updated_at", 
	"transition": {
"links": [
	{
	"task": "END", 
	"queries": "END", 
	"pattern": "none"
	}
  ]
}	}
}
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
  ]
}
}
`
)

type AddTask struct {
	PkFlow    int          `yaml:"flow" default:"1"`
	Name      types.String `yaml:"name" default:"newtask"`
	Sentence  types.String `yaml:"sentence"`
	Target    types.String `yaml:"target"`
	RoleField types.String `yaml:"role_field" default:"roles[1]"`
	TimeField types.String `yaml:"time_field" default:"updated_at"`
	Pattern   types.String `yaml:"pattern" default:"none"`
	Report    types.String `yaml:"report" default:"true"`
	Source    types.String `yaml:"source"`
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

func UnmarshallAddTask(order []byte) (AddTask, error) {
	var addTask AddTask
	var err error
	err = yaml.Unmarshal(order, &addTask)
	if err != nil {
		return addTask, err
	}
	utils.SetDefaults(&addTask)
	return addTask, nil
}

func DoAddTask(order string, query *database.IQuery) error {
	var currentInsert string
	task, err := UnmarshallAddTask([]byte(order))
	if err != nil {
		return err
	}
	fmt.Println("task struct flow", task.PkFlow, "name", task.Name)
	switch types.IsNil(task.Source) {
	case true:
		currentInsert = StartInsertTaskInFlow
	default:
		if types.IsNil(task.Target) {
			currentInsert = EndInsertTaskInFlow
		} else {
			currentInsert = InsertTaskInFlow
		}
	}
	sqlInsert, err := utils.ParseTemplate(currentInsert, &task)
	if err != nil {
		return err
	}
	fmt.Println("sqlInsert", sqlInsert)
	err = (*query).Exec(sqlInsert)
	if err != nil {
		return err
	}
	return nil
}
