package workflow

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"stories/actions/payment"
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
	// Delete sentences for tasks

	DeleteTaskInFlow = `
WITH CTE AS (
    SELECT (flow->'tasks'->'{{.Name}}'->'transition'->'links'->0->'task')::text as task
FROM ta_flow WHERE pk_flow={{.PkFlow}}
)
UPDATE ta_flow SET flow=jsonb_set(flow, '{tasks, {{.Source}},transition,links,0,task }'::text[], CTE.task::jsonb,true)
FROM CTE WHERE pk_flow={{.PkFlow}};
WITH CTE AS (
    SELECT (flow->'tasks'->'{{.Name}}'->'transition'->'links'->0->'queries')::text as query
FROM ta_flow WHERE pk_flow={{.PkFlow}} 
)
UPDATE ta_flow SET flow=jsonb_set(flow, '{tasks, {{.Source}},transition,links,0,queries }'::text[], CTE.query::jsonb,true)
FROM CTE WHERE pk_flow={{.PkFlow}};
UPDATE ta_flow SET flow=flow #- '{tasks, {{.Name}} }'WHERE pk_flow={{.PkFlow}};
`

	EndDeleteTaskInFlow = `
UPDATE ta_flow SET flow=jsonb_set(flow, '{tasks, {{.Source}},transition,links,0,task }'::text[], '"END"'::jsonb,true)
WHERE pk_flow={{.PkFlow}};
UPDATE ta_flow SET flow=jsonb_set(flow, '{tasks, {{.Source}},transition,links,0,queries }'::text[], '"END"'::jsonb,true)
WHERE pk_flow={{.PkFlow}};
UPDATE ta_flow SET flow=flow #- '{tasks, {{.Name}} }'WHERE pk_flow={{.PkFlow}};

`

	StartDeleteTaskInFlow = `
WITH CTE AS (
    SELECT (flow->'tasks'->'{{.Name}}'->'transition'->'links'->0->'queries')::text as query
FROM ta_flow WHERE pk_flow={{.PkFlow}}
)
UPDATE ta_flow SET flow=jsonb_set(flow, '{initial}'::text[], CTE.query::jsonb,true)
FROM CTE WHERE pk_flow={{.PkFlow}};
UPDATE ta_flow SET flow=flow #- '{tasks, {{.Name}} }' WHERE pk_flow={{.PkFlow}};
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
	case "delete_task":
		err = DoDeleteTask(order, query)
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

func doCheckTask(task *AddTask, query *database.IQuery) error {
	var workflowType payment.WorkflowType
	workflowType.Query = query
	wf, err := workflowType.LoadWorkflow()
	if err != nil {
		return err
	}
	if task.Source.IsEmpty() && task.Target.IsEmpty() {
		return errors.New("source and target values are empty. At least one of them must have a value")
	}
	if task.Source.IsEmpty() {
		prev, err := wf.PreviousStatus(task.Target.String())
		if err != nil {
			return err
		}
		task.Source = types.NewString(prev)
	}
	if task.Target.IsEmpty() {
		next, err := wf.NextStatus(task.Source.String())
		if err != nil {
			return err
		}
		task.Target = types.NewString(next)
	}
	return nil
}

func DoAddTask(order string, query *database.IQuery) error {
	var currentInsert string
	task, err := UnmarshallAddTask([]byte(order))
	if err != nil {
		return err
	}
	fmt.Println("task struct flow", task.PkFlow, "name", task.Name)
	err = doCheckTask(&task, query)
	if err != nil {
		return err
	}
	switch task.Source.IsEmpty() {
	case true:
		currentInsert = StartInsertTaskInFlow
	default:
		if task.Target.IsEmpty() {
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

func DoDeleteTask(order string, query *database.IQuery) error {
	var currentDelete string
	task, err := UnmarshallAddTask([]byte(order))
	if err != nil {
		return err
	}
	fmt.Println("task struct flow", task.PkFlow, "name", task.Name)
	switch types.IsNil(task.Source) {
	case true:
		currentDelete = StartDeleteTaskInFlow
	default:
		if types.IsNil(task.Target) {
			currentDelete = EndDeleteTaskInFlow
		} else {
			currentDelete = DeleteTaskInFlow
		}
	}
	sqlDelete, err := utils.ParseTemplate(currentDelete, &task)
	if err != nil {
		return err
	}
	fmt.Println("sqlDelete", sqlDelete)
	err = (*query).Exec(sqlDelete)
	if err != nil {
		return err
	}
	return nil
}
