package stories

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"stories/database"
	"stories/generated/models"
	"stories/pkg/utils"
	"stories/pkg/workflow/common"
	"stories/pkg/workflow/queries/pg"
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
		"info_field": "images[1]", 
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

type GetFlow struct {
	PkFlow         int       `yaml:"flow" default:"1"`
	Ticket         types.Int `yaml:"ticket"`
	ShowEmptyTasks bool      `yaml:"show_empty_tasks" default:"false"`
}

type AddTask struct {
	PkFlow    int          `yaml:"flow" default:"1"`
	Name      types.String `yaml:"name" default:"Newtask"`
	Sentence  types.String `yaml:"sentence" default:"END"`
	Target    types.String `yaml:"target" default:"END"`
	RoleField types.String `yaml:"role_field" default:"images[1]"`
	TimeField types.String `yaml:"time_field" default:"updated_at"`
	Pattern   types.String `yaml:"pattern" default:"none"`
	Report    types.String `yaml:"report" default:"true"`
	Source    types.String `yaml:"source"`
}

func WorkflowProcessor(action string, order string, query *database.IQuery) error {
	var err error
	switch action {
	case "add_task":
		_, err = DoAddTask(order, query)
	case "delete_task":
		err = DoDeleteTask(order, query)
	case "show_flow":
		err = DoGetWorkflow(order, query)
	}
	if err != nil {
		return err
	}
	return nil
}

func UnmarshallGetFlow(order []byte) (GetFlow, error) {
	var getFlow GetFlow
	var err error
	err = yaml.Unmarshal(order, &getFlow)
	if err != nil {
		return getFlow, err
	}
	utils.SetDefaults(&getFlow)
	return getFlow, nil
}

func DoGetWorkflow(order string, query *database.IQuery) error {
	getFlow, err := UnmarshallGetFlow([]byte(order))
	if err != nil {
		return err
	}
	log.Println("getFlow", getFlow)
	wf, err := GetWorkflow(getFlow.PkFlow, query)
	if err != nil {
		return err
	}
	wf.CurrentIDInt = getFlow.Ticket
	monad := pg.New(wf, query)
	monad = pg.Bind(monad, pg.TimeField[models.TaTracking])
	monad = pg.Bind(monad, pg.RoleField[models.TaTracking])
	monad = pg.Bind(monad, pg.CountToken[models.TaTracking])
	monad = pg.Bind(monad, pg.Render[models.TaTracking])
	_, err = monad(nil)
	if err != nil {
		log.Print("monad error", err)
		return err
	}
	log.Print("wf rendered successfully")
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

func doCheckAddTask(task *AddTask, query *database.IQuery) (*common.Workflow[models.TaTracking], error) {
	var workflowType WorkflowType
	workflowType.Query = query
	wf, err := workflowType.LoadWorkflow()
	if err != nil {
		return wf, err
	}
	if task.Source.IsEmpty() && task.Target.IsEmpty() {
		return wf, fmt.Errorf("target and source values can't be both empty")
	}
	if err = wf.Contains(task.Name.String()); err == nil {
		return wf, fmt.Errorf("task called '%v' already exists on workflow", task.Name.String())
	}
	if err = wf.Contains(task.Source.String()); err != nil {
		return wf, err
	}
	if err = wf.Contains(task.Target.String()); err != nil {
		return wf, err
	}
	if task.Source.IsEmpty() {
		prev, err := wf.PreviousStatus(task.Target.String())
		if err != nil {
			return wf, err
		}
		task.Source = types.NewString(prev)
	}
	if task.Target.IsEmpty() {
		next, err := wf.NextStatus(task.Source.String())
		if err != nil {
			return wf, err
		}
		task.Target = types.NewString(next)
	}
	return wf, nil
}

func doCheckDelTask(task *AddTask, query *database.IQuery) error {
	var workflowType WorkflowType
	workflowType.Query = query
	wf, err := workflowType.LoadWorkflow()
	if err != nil {
		return err
	}
	current := task.Name.String()
	if err = wf.Contains(current); err != nil {
		return err
	}
	prev, err := wf.PreviousStatus(current)
	if err != nil {
		return err
	}
	task.Source = types.NewString(prev)
	next, err := wf.NextStatus(current)
	if err != nil {
		return err
	}
	task.Target = types.NewString(next)
	return nil
}

func GetWorkflow(pkFlow int, query *database.IQuery) (*common.Workflow[models.TaTracking], error) {
	var workflowType WorkflowType
	workflowType.Query = query
	wf, err := workflowType.LoadWorkflow()
	if err != nil {
		return wf, err
	}
	return wf, nil
}

func DoAddTaskIntl(task AddTask, query *database.IQuery) (*common.Workflow[models.TaTracking], error) {
	var currentInsert string
	fmt.Println("task struct flow", task.PkFlow, "name", task.Name)
	wf, err := doCheckAddTask(&task, query)
	if err != nil {
		return wf, err
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
		return wf, err
	}
	fmt.Println("sqlInsert", sqlInsert)
	err = (*query).Exec(sqlInsert)
	if err != nil {
		return wf, err
	}
	return wf, nil
}

func DoAddTask(order string, query *database.IQuery) (*common.Workflow[models.TaTracking], error) {
	task, err := UnmarshallAddTask([]byte(order))
	if err != nil {
		return nil, err
	}
	return DoAddTaskIntl(task, query)

}

func DoDeleteTask(order string, query *database.IQuery) error {
	var currentDelete string
	task, err := UnmarshallAddTask([]byte(order))
	if err != nil {
		return err
	}
	err = doCheckDelTask(&task, query)
	if err != nil {
		return err
	}
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
