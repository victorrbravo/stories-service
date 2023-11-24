package common

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	db "stories/database"
	"stories/types"
	"strings"
)

const (
	WorkflowName      = "worflowExample"
	QueryName         = "queries.FirstStation"
	NewUUID           = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	PatternSequencial = "sequencial"
	StationFirst      = "StationFirst"
	EndTask           = "END"
	ARROW_TASK        = " -> "
	OR_OPERATOR       = " OR "
	AndOperator       = " AND "
	TablePrefix       = "ta_"
)

type FillType string

const (
	FillTime  FillType = "Time"
	FillMoney FillType = "Money"
)

type Values struct {
	Value types.Int
	Min   types.Int
	Max   types.Int
}

type WhereTask struct {
	Name  string
	Where string
}

func (obj *Values) Assign(v types.Int) types.Int {
	obj.Value = v
	if types.IsNil(obj.Max) || types.AsInt(v) > types.AsInt(obj.Max) {
		obj.Max = v
	}
	if types.IsNil(obj.Min) || types.AsInt(v) < types.AsInt(obj.Min) {
		obj.Min = v
	}
	return v
}

type Data struct {
	Count      types.Int    // checked
	TimeField  types.String // checked
	InfoField  types.String
	Roles      []types.String // checked
	Percentage float64
	Total      types.Int
	FillValue  types.Int // Measure progress in time, money or other type of value
	FillSelect FillType
	Rows       []interface{}
	IDCurrent  types.Int
	ToExport   Exportable
	Stats      Values
}

type JobWorkflow struct {
	Workflow Workflow
}

type WhereClause func(wf Workflow) types.String

type Exportable struct {
	Data     types.String   `json:"data"`
	UpdateAt types.DateTime `json:"update_at"`
	Synced   bool           `json:"synced"`
}

type Workflow struct {
	ID                types.UUID        `json:"id"`
	Name              string            `json:"name"`
	TableName         string            `json:"table_name"`
	FieldToken        string            `json:"field_token"`
	TrackingTableName types.String      `json:"tracking_table_name"`
	TrackingKeys      types.String      `json:"tracking_keys"`
	Variable          Variable          `json:"variable"`
	Tasks             map[string]*Task  `json:"tasks"`
	initialTask       types.String      `json:"initial_task,omitempty"`
	Data              Data              `json:"data,omitempty"`
	CurrentIDInt      types.Int         `json:"current_id_int,omitempty"`
	CurrentIDUUID     types.UUID        `json:"current_id_uuid,omitempty"`
	Initial           types.String      `json:"initial"`
	Query             *db.PostgresQuery `json:"omitempty"`
	WhereClause       WhereClause       `json:"where_clause"`
}

type Variable struct {
	Name  string `json:"name"`
	Query string `json:"queries"`
}

type Task struct {
	Name       string     `json:"name"`
	TimeField  string     `json:"time_field"`
	InfoField  string     `json:"info_field"`
	Transition Transition `json:"transition"`
	Data       Data       `json:"data,omitempty"`
}

type Transition struct {
	Links   []Link `json:"links"`
	Pattern string `json:"pattern"`
}

type Link struct {
	Task   string `json:"task"`
	Query  string `json:"queries"`
	Status string `json:"status,omitempty"`
}

func Create(ctx context.Context, name string) (*Workflow, error) {
	wf := Workflow{}
	return &wf, nil
}

func CreateExample(ctx context.Context, name string) (*Workflow, error) {
	var err error
	wf := Workflow{}
	wf.Name = WorkflowName
	wf.ID, err = types.NewUUID(NewUUID)
	if err != nil {
		return nil, err
	}
	wf.Name = name
	err = wf.AddTask(ctx, StationFirst, nil, "")
	return &wf, err
}

func (wf *Workflow) AddTask(ctx context.Context, name string, linkTask *Task, status string) error {
	newTask := new(Task)
	newTask.Name = name
	newTransition := Transition{Pattern: PatternSequencial}
	newLink := Link{Task: linkTask.Name, Query: QueryName, Status: status}
	newTransition.Links = append(newTransition.Links, newLink)
	newTask.Transition = newTransition
	wf.Tasks[name] = newTask
	return nil
}

// remove a task from the workflow by name
func (wf *Workflow) RemoveTask(ctx context.Context, name string) error {
	delete(wf.Tasks, name)
	return nil
}

// update a task in the workflow by name
func (wf *Workflow) UpdateTask(ctx context.Context, name string, newTask *Task) error {
	for i, task := range wf.Tasks {
		if task.Name == name {
			wf.Tasks[i] = newTask
			return nil
		}
	}
	return nil
}

func (wf Workflow) NextTask(name string) *Task {
	ctx := context.Background()
	task, ok := wf.Tasks[name]
	if !ok {
		log.Println("error task", name)
		return nil
	}
	nextName := task.Transition.Links[0].Task
	if strings.Compare(nextName, EndTask) == 0 {
		return nil
	}
	nextTask, ok := wf.Tasks[nextName]
	if !ok {
		log.Println(ctx, "error Next task", nextName)
		return nil
	}
	return nextTask
}

// add a link to a task
func (wf *Workflow) AddLink(ctx context.Context, taskName string, linkTask *Task, query string, status string) error {
	for _, task := range wf.Tasks {
		if task.Name == taskName {
			newLink := Link{Task: linkTask.Name, Query: query, Status: status}
			task.Transition.Links = append(task.Transition.Links, newLink)
			return nil
		}
	}
	return nil
}

func (wf *Workflow) Statuses(ctx context.Context, quotes bool) ([]string, error) {
	var err error
	var start string
	list := []string{}
	if types.IsNil(wf.initialTask) {
		start, err = wf.CalculateInitialTask()
		if err != nil {
			return list, err
		}
		wf.initialTask = types.NewString(start)
	}
	current := wf.initialTask.String()
	title := current
	if quotes {
		title = "'" + current + "'"
	}
	list = append(list, title)
	lenMax := len(wf.Tasks)
	i := 0
	for {
		task, ok := wf.Tasks[current]
		if !ok {
			return list, fmt.Errorf("Incorrect value for task name %s", current)
		}
		current = task.Transition.Links[0].Task
		if len(current) == 0 || strings.Compare(current, EndTask) == 0 {
			break
		}
		title = current
		if quotes {
			title = "'" + current + "'"
		}
		list = append(list, title)
		if i > lenMax {
			return list, fmt.Errorf("Error iterating through tasks")
		}
		i++
	}
	return list, nil
}

func (wf *Workflow) CalculateInitialTask() (string, error) {
	tasksVisited := map[string]bool{}
	value := ""
	for key, task := range wf.Tasks {
		task.Name = key
		next := task.Transition.Links[0].Task
		tasksVisited[next] = true
	}
	for key, _ := range wf.Tasks {
		_, ok := tasksVisited[key]
		if !ok {
			return key, nil
		}
	}
	return value, fmt.Errorf("fail calculating the start task. No start task")
}

// load a wokflow instance using a json file
func (wf *Workflow) Load(ctx context.Context, jsonFile []byte) error {
	err := json.Unmarshal(jsonFile, &wf)
	return err
}

func (wf Workflow) ExtractKeys(trackingKeys []string) (string, string) {
	linkStory := ""
	linkStatus := ""
	if len(trackingKeys) > 0 {
		linkStory = trackingKeys[0]
	}
	if len(trackingKeys) > 1 {
		linkStatus = trackingKeys[1]
	}
	return linkStory, linkStatus
}

func (wf Workflow) ReducesFilter(name string) string {
	wheres := wf.ClausesWheres()
	value := ""
	for _, tuple := range wheres {
		if strings.Compare(tuple.Name, name) == 0 {
			value = tuple.Where
		}
	}
	return value
}

func (wf Workflow) ClausesWheres() map[int]WhereTask {
	ctx := context.Background()
	values := make(map[int]WhereTask)
	statuses, err := wf.Statuses(ctx, false)
	if err != nil {
		log.Println("Error WhereClauses Statuses")
		return values
	}
	index := 0
	str := ""
	for _, name := range statuses {
		current, _ := wf.Tasks[name]
		if index == 0 {
			values[index] = WhereTask{name, wf.Initial.String()}
		} else {
			values[index] = WhereTask{name, str}
		}
		str = current.Transition.Links[0].Query
		if strings.Compare(str, EndTask) != 0 {
			index++
		}
	}
	return values
}
