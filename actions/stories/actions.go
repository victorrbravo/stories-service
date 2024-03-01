package stories

import (
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	db "stories/database"
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

type ShowFlowInput struct {
	PkFlow         int                             `yaml:"flow" default:"1"`
	Ticket         types.Int                       `yaml:"ticket"`
	ShowEmptyTasks bool                            `yaml:"show_empty_tasks" default:"false"`
	Task           WorkflowType[models.TaTracking] `yaml:"task,omitempty"`
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

type WorkflowData struct {
	Name  string      `json:"name"`
	Data  common.Data `json:"data"`
	Tasks []TaskData  `json:"tasks"`
}
type TaskData struct {
	Name       string            `json:"name"`
	Data       common.Data       `json:"data"`
	Transition common.Transition `json:"transition"`
}

func DoGetWorkflowSchema(order string, query *db.IQuery) (string, error) {
	var value = `
 {
	"actions" : {
		"add_task": {
			"description" : "create a new flow",
			"fields": [
				{ "name": "name", "type": "string", "validator" : "" },
				{ "name": "summary", "type": "string", "validator" : "" }
			]
		},
		"delete_task": {},
		"show_flow_graph": {},
		"show_flow_data": {},
		"show_task_data": {}
 	}
}
`

	return value, nil
}

func WorkflowProcessor(action string, order string, query *db.IQuery) (string, error) {
	var err error
	var value string
	switch action {
	case "workflow_schema":
		value, err = DoGetWorkflowSchema(order, query)
	case "add_task":
		_, err = DoAddTask[models.TaTracking](order, query)
	case "delete_task":
		err = DoDeleteTask[models.TaTracking](order, query)
	case "show_flow_graph":
		value, err = DoShowWorkflowGraph[models.TaTracking](order, query)
	case "show_flow_data":
		value, err = DoShowWorkflowData[models.TaTracking](order, query)
	case "show_task_data":
		value, err = DoShowWorkflowTaskData[models.TaTracking](order, query)
	}
	if err != nil {
		return value, err
	}
	return value, nil
}

func DoShowWorkflowTaskData[T db.Scanner](order string, query *db.IQuery) (string, error) {
	var value string
	getFlow, err := UnmarshallGetFlow([]byte(order), query)
	if err != nil {
		return value, err
	}
	taStory, err := models.GetTaStoryByPkStory(*query,
		getFlow.Ticket)
	if err != nil {
		return value, err
	}
	getFlow.PkFlow = types.AsInt(taStory.FkFlow)
	log.Println("getFlow pkFlow", getFlow.PkFlow)
	wf, err := UpdateWorkflow[T](getFlow, query, []interface{}{getFlow.Task.String()}, pg.VariableField[T])
	if err != nil {
		return value, err
	}
	data, err := json.Marshal(wf.Data.Rows)
	if err != nil {
		return value, err
	}
	value = string(data)
	fmt.Println("wf", wf.ID)
	return value, nil

}

func DoShowWorkflowData[T db.Scanner](order string, query *db.IQuery) (string, error) {
	var value string
	getFlow, err := UnmarshallGetFlow([]byte(order), query)
	if err != nil {
		return value, err
	}
	taStory, err := models.GetTaStoryByPkStory(*query,
		getFlow.Ticket)
	if err != nil {
		return value, err
	}
	getFlow.PkFlow = types.AsInt(taStory.FkFlow)
	log.Println("DoShowWorkflowData getFlow pkFlow", getFlow.PkFlow)
	wf, err := UpdateWorkflow[T](getFlow, query, []interface{}{}, pg.TimeField[T], pg.RoleField[T],
		pg.CountToken[T])
	if err != nil {
		return value, err
	}
	value, err = updateWorkflowData[T](wf)
	if err != nil {
		return value, err
	}
	return value, nil
}

func updateWorkflowData[T db.Scanner](wf *common.Workflow[T]) (string, error) {
	var err error
	var value string
	ctx := context.Background()
	statuses, _ := wf.Statuses(ctx, false)
	wfData := WorkflowData{Name: wf.Name}
	wfData.Data = wf.Data
	wfData.Tasks = make([]TaskData, 0)
	for _, status := range statuses {
		current := wf.Tasks[status]
		if types.AsBool(current.Report) == false {
			continue
		}
		taskData := TaskData{
			Name:       current.Name,
			Transition: current.Transition,
			Data:       current.Data,
		}
		wfData.Tasks = append(wfData.Tasks, taskData)
	}
	data, err := json.Marshal(wfData)
	if err != nil {
		return value, err
	}
	value = string(data)
	return value, nil
}

func UnmarshallGetFlow(order []byte, query *db.IQuery) (ShowFlowInput, error) {
	var getFlow ShowFlowInput
	var err error
	getFlow.Task.Query = query
	err = yaml.Unmarshal(order, &getFlow)
	if err != nil {
		return getFlow, err
	}
	utils.SetDefaults(&getFlow)
	return getFlow, nil
}

func UpdateWorkflow[T db.Scanner](getFlow ShowFlowInput, query *db.IQuery, values []interface{}, funcs ...func(*common.Workflow[T], ...interface{}) pg.Monad[T]) (*common.Workflow[T], error) {
	wf, err := GetWorkflow[T](getFlow.PkFlow, query)
	if err != nil {
		return wf, err
	}
	wf.CurrentIDInt = getFlow.Ticket
	monad := pg.New[T](wf, query)
	for _, f := range funcs {
		monad = pg.Bind[T](monad, f, values...)
	}
	_, err = monad(nil)
	if err != nil {
		log.Print("update workflow monad error", err)
		return wf, err
	}
	return wf, nil
}

func DoShowWorkflowGraph[T db.Scanner](order string, query *db.IQuery) (string, error) {
	var value string
	getFlow, err := UnmarshallGetFlow([]byte(order), query)
	if err != nil {
		return value, err
	}
	log.Println("getFlow", getFlow)
	wf, err := UpdateWorkflow[T](getFlow, query, []interface{}{}, pg.TimeField[T], pg.RoleField[T],
		pg.CountToken[T], pg.Render[T])
	if err != nil {
		return value, err
	}
	value = wf.JsonResponse
	return value, nil
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

func doCheckAddTask[T db.Scanner](pkFlow int, task *AddTask, query *db.IQuery) (*common.Workflow[T], error) {
	var workflowType WorkflowType[T]
	workflowType.PkFlow = pkFlow
	workflowType.Query = query
	wf, err := workflowType.LoadWorkflow(workflowType.PkFlow)
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

func doCheckDelTask[T db.Scanner](task *AddTask, query *db.IQuery) error {
	var workflowType WorkflowType[T]
	workflowType.Query = query
	workflowType.PkFlow = 1
	wf, err := workflowType.LoadWorkflow(workflowType.PkFlow)
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

func GetWorkflow[T db.Scanner](pkFlow int, query *db.IQuery) (*common.Workflow[T], error) {
	var workflowType WorkflowType[T]
	workflowType.Query = query
	workflowType.PkFlow = pkFlow
	wf, err := workflowType.LoadWorkflow(workflowType.PkFlow)
	if err != nil {
		return wf, err
	}
	return wf, nil
}

func DoAddTaskIntl[T db.Scanner](pkFlow int, task AddTask, query *db.IQuery) (*common.Workflow[T], error) {
	var currentInsert string
	fmt.Println("task struct flow", task.PkFlow, "name", task.Name)
	wf, err := doCheckAddTask[T](pkFlow, &task, query)
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

func DoAddTask[T db.Scanner](order string, query *db.IQuery) (*common.Workflow[T], error) {
	task, err := UnmarshallAddTask([]byte(order))
	if err != nil {
		return nil, err
	}
	return DoAddTaskIntl[T](task.PkFlow, task, query)

}

func DoDeleteTask[T db.Scanner](order string, query *db.IQuery) error {
	var currentDelete string
	task, err := UnmarshallAddTask([]byte(order))
	if err != nil {
		return err
	}
	err = doCheckDelTask[T](&task, query)
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
