package stories

import (
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"sort"
	db "stories/database"
	"stories/generated/models"
	"stories/pkg/utils"
	"stories/pkg/workflow/common"
	"stories/pkg/workflow/queries/pg"
	"stories/types"
	"time"
)

const (
	NoStartChangeNameTaskInFlow = `
WITH CTE AS (
flow->'tasks'->'{{.Source}}'->'transition'->'links'->0->'queries' from ta_flow where pk_flow={{.PkFlow}} AS srcSentence	
)
UPDATE ta_flow SET flow=jsonb_set(flow, '{tasks, {{.Source}},transition,links,0,queries }'::text[], srcSentence,true)
UPDATE ta_flow SET flow=jsonb_set(flow, '{tasks, {{.Current}}, }'::text[], srcSentence,true)

`

	StartChangeNameTaskInFlow = `
WITH CTE AS (
flow->'tasks'->'{{.Source}}'->'transition'->'links'->0->'queries' from ta_flow where pk_flow={{.PkFlow}} AS srcSentence	
)
UPDATE ta_flow SET flow=jsonb_set(flow, '{tasks, {{.Source}},transition,links,0,queries }'::text[], srcSentence,true)
`

	InsertTaskInFlow = `
WITH CTE AS (
SELECT 
'{
"info_field": "{{.RoleField}}", 
"time_field": "{{.TimeField}}", 
"amount_field": "{{.AmountField}}",
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
"amount_field": "{{.AmountField}}",
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
"amount_field": "{{.AmountField}}",
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
   "tracking_id": "pk_tracking",
	"field_token": "pk_story",
  "field_status": "status",
  "table_name": "ta_story",
  "initial": "TRUE",
"tasks": {
	"Start" :{
		"report": false,
		"info_field": "images[1]", 
		"time_field": "updated_at", 
		"amount_field": "amount_spent",
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
"amount_field": "{{.AmountField}}",
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

const GetTaskPoints = `
WITH CTE AS (
SELECT flow -> 'tasks' -> '{{.Target}}' -> 'transition' AS tr1,
flow -> 'tasks' -> '{{.Current}}' -> 'transition' AS tr2,
flow -> 'tasks' -> '{{.Previous}}' -> 'transition' AS tr3
FROM ta_flow where pk_flow={{.PkFlow}}
)
SELECT tr1, tr2, tr3 FROM CTE
`
const UpdateAttributeFlow = `
UPDATE ta_flow SET flow=jsonb_set(flow, $1::text[], $2::jsonb,true) WHERE pk_flow=$3
`

type ShowFlowInput struct {
	PkFlow         int                             `yaml:"flow" default:"1"`
	Ticket         types.Int                       `yaml:"ticket"`
	ShowEmptyTasks bool                            `yaml:"show_empty_tasks" default:"false"`
	Task           WorkflowType[models.TaTracking] `yaml:"task,omitempty"`
	Filter         string                          `yaml:"filter,omitempty"`
	ByDay          bool                            `yaml:"by_day,omitempty"`
	PageNumber     types.Int                       `yaml:"page_number,omitempty"`
	PageSize       types.Int                       `yaml:"page_size,omitempty"`
}

type AddTask struct {
	PkFlow      int          `yaml:"flow" default:"1"`
	Name        types.String `yaml:"name" default:"Newtask"`
	Sentence    types.String `yaml:"sentence" default:"END"`
	Target      types.String `yaml:"target" default:"END"`
	RoleField   types.String `yaml:"role_field" default:"images[1]"`
	TimeField   types.String `yaml:"time_field" default:"updated_at"`
	AmountField types.String `yaml:"amount_field" default:"amount_spent"`
	Pattern     types.String `yaml:"pattern" default:"none"`
	Report      types.String `yaml:"report" default:"true"`
	Source      types.String `yaml:"source"`
	VideosField types.String `yaml:"videos_field" default:"videos"`
	ImagesField types.String `yaml:"images_field" default:"images"`
}

type WorkflowData struct {
	Name  string      `json:"name"`
	Data  common.Data `json:"data"`
	Tasks []TaskData  `json:"tasks"`
}
type TaskData struct {
	ID         types.Int         `json:"id"`
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

func WorkflowProcessor(ctx context.Context, action string, order string, query *db.IQuery) (string, error) {
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
	start1 := time.Now()
	taStory, err := models.GetTaStoryByPkStory(*query,
		getFlow.Ticket)
	if err != nil {
		return value, err
	}
	getFlow.PkFlow = types.AsInt(taStory.FkFlow)
	wf, err := UpdateWorkflow[T](getFlow, query, []interface{}{},
		pg.FilterField[T],
		pg.RoleField[T],
		pg.DistanceField[T],
		pg.TimeField[T],
		pg.CountToken[T],
		pg.ArrayField[T],
		pg.StatsField,
	)
	elapsed1 := time.Since(start1) // calculate duration
	fmt.Printf("Time spent1: %s\n", elapsed1)
	start2 := time.Now()
	if err != nil {
		return value, err
	}
	value, err = UpdateWorkflowData[T](wf)
	if err != nil {
		return value, err
	}
	elapsed2 := time.Since(start2) // calculate duration
	fmt.Printf("Time spent2: %s\n", elapsed2)
	return value, nil
}

func DoShowWorkflowDataInt[T db.Scanner](order string, query *db.IQuery) (WorkflowData, error) {
	wfData := WorkflowData{}
	getFlow, err := UnmarshallGetFlow([]byte(order), query)
	if err != nil {
		return wfData, err
	}
	taStory, err := models.GetTaStoryByPkStory(*query,
		getFlow.Ticket)
	if err != nil {
		return wfData, err
	}
	getFlow.PkFlow = types.AsInt(taStory.FkFlow)
	log.Println("DoShowWorkflowData getFlow pkFlow", getFlow.PkFlow)
	wf, err := UpdateWorkflow[T](getFlow, query, []interface{}{}, pg.TimeField[T], pg.RoleField[T],
		pg.CountToken[T], pg.ArrayField[T])
	if err != nil {
		return wfData, err
	}
	wfData, err = UpdateWorkflowDataInt[T](wf)
	if err != nil {
		return wfData, err
	}
	return wfData, nil
}

func UpdateWorkflowData[T db.Scanner](wf *common.Workflow[T]) (string, error) {
	var value string
	wfData, err := UpdateWorkflowDataInt[T](wf)
	if err != nil {
		return value, err
	}
	data, err := json.Marshal(wfData)
	if err != nil {
		return value, err
	}
	value = string(data)
	return value, nil
}

func UpdateWorkflowDataInt[T db.Scanner](wf *common.Workflow[T]) (WorkflowData, error) {
	ctx := context.Background()
	statuses, err := wf.Statuses(ctx, false, true)
	if len(wf.Filter.IncludedWrapped) > 0 {
		statuses = wf.Filter.IncludedWrapped
	}
	if err != nil {
		return WorkflowData{}, err
	}
	wfData := WorkflowData{Name: wf.Name}
	wfData.Data = wf.Data
	wfData.Tasks = make([]TaskData, 0)
	for _, status := range Paginate(statuses, wf.Page) {
		current := wf.Tasks[status]
		taskData := TaskData{
			Name:       current.Name,
			Transition: current.Transition,
			Data:       current.Data,
		}
		wfData.Tasks = append(wfData.Tasks, taskData)
	}
	sort.Slice(wfData.Tasks, func(i, j int) bool {
		firstDate, _ := time.Parse(time.RFC3339, types.AsString(wfData.Tasks[i].Data.TimeField))
		secondDate, _ := time.Parse(time.RFC3339, types.AsString(wfData.Tasks[j].Data.TimeField))
		return firstDate.After(secondDate)
	})
	return wfData, nil
}

func Paginate[T any](items []T, page types.Page) []T {
	start := 0
	end := len(items)

	if page.PageNumber > 0 && page.PageSize > 0 {
		start = (page.PageNumber - 1) * page.PageSize
		end = start + page.PageSize

		if start > len(items) {
			start = len(items)
		}
		if end > len(items) {
			end = len(items)
		}
	}
	return items[start:end]
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

func UpdateWorkflow[T db.Scanner](getFlow ShowFlowInput, query *db.IQuery,
	values []interface{},
	funcs ...func(*common.Workflow[T], ...interface{}) pg.Monad[T]) (*common.Workflow[T], error) {
	wf, err := GetWorkflow[T](getFlow.PkFlow, query)
	wf.Filter.Sentence = getFlow.Filter
	wf.ByDay = getFlow.ByDay
	wf.Page.PageNumber = types.AsInt(getFlow.PageNumber)
	wf.Page.PageSize = types.AsInt(getFlow.PageSize)
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
	taStory, err := models.GetTaStoryByPkStory(*query,
		getFlow.Ticket)
	if err != nil {
		return value, err
	}
	getFlow.PkFlow = types.AsInt(taStory.FkFlow)
	log.Println("DoShowWorkflowGraph getFlow pkFlow", getFlow.PkFlow)
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
	workflowType.PkFlow = task.PkFlow
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
	for _, task := range wf.Tasks {
		utils.SetDefaults[common.Task](task)
	}
	wf.CurrentIDInt = types.NewInt(pkFlow)
	return wf, nil
}

func DoAddTaskIntl[T db.Scanner](pkFlow int, task AddTask, query *db.IQuery) (*common.Workflow[T], error) {
	var currentInsert string
	log.Println("task struct flow", task.PkFlow, "name", task.Name)
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

func DoDeleteTaskInt[T db.Scanner](task AddTask, query *db.IQuery) error {
	var currentDelete string
	var err error
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

func DoDeleteTask[T db.Scanner](order string, query *db.IQuery) error {
	task, err := Unmarshall[AddTask]([]byte(order))
	if err != nil {
		return err
	}
	err = DoDeleteTaskInt[T](task, query)
	return err
}
