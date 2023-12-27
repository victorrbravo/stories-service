package common

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"slices"
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

type JobWorkflow[T db.Scanner] struct {
	Workflow Workflow[T]
}

type WhereClause[T db.Scanner] func(wf Workflow[T]) types.String

func (w WhereClause[T]) MarshalJSON() ([]byte, error) {
	return []byte("\"\""), nil
}

type Exportable struct {
	Data     types.String   `json:"data"`
	UpdateAt types.DateTime `json:"update_at"`
	Synced   bool           `json:"synced"`
}

type Workflow[T db.Scanner] struct {
	ID                types.UUID        `json:"id"`
	Name              string            `json:"name"`
	TableName         string            `json:"table_name"`
	FieldToken        string            `json:"field_token"`
	FieldStatus       string            `json:"field_status,omitempty"`
	TrackingTableName types.String      `json:"tracking_table_name"`
	TrackingKeys      types.String      `json:"tracking_keys"`
	Variable          Variable          `json:"variable"`
	Tasks             map[string]*Task  `json:"tasks"`
	Data              Data              `json:"data,omitempty"`
	CurrentIDInt      types.Int         `json:"current_id_int,omitempty"`
	CurrentIDUUID     types.UUID        `json:"current_id_uuid,omitempty"`
	Initial           types.String      `json:"initial"`
	Query             *db.PostgresQuery `json:"omitempty"`
	WhereClause       WhereClause[T]    `json:"where_clause"`
	ValueScanner      T                 `json:"omitempty"`
	statuses          []string          `json:"omitempty"`
}

type Variable struct {
	Name   string `json:"name"`
	Fields string `json:"fields"`
}

type Task struct {
	Name       string     `json:"name"`
	TimeField  string     `json:"time_field"`
	InfoField  string     `json:"info_field"`
	Transition Transition `json:"transition"`
	Data       Data       `json:"data,omitempty"`
	Report     types.Bool `json:"report,omitempty"`
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

func (wf *Workflow[T]) AddTask(ctx context.Context, name string, linkTask *Task, status string) error {
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
func (wf *Workflow[T]) RemoveTask(ctx context.Context, name string) error {
	delete(wf.Tasks, name)
	return nil
}

// update a task in the workflow by name
func (wf *Workflow[T]) UpdateTask(ctx context.Context, name string, newTask *Task) error {
	for i, task := range wf.Tasks {
		if task.Name == name {
			wf.Tasks[i] = newTask
			return nil
		}
	}
	return nil
}

func (wf Workflow[T]) NextTask(name string) *Task {
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
func (wf *Workflow[T]) AddLink(ctx context.Context, taskName string, linkTask *Task, query string, status string) error {
	for _, task := range wf.Tasks {
		if task.Name == taskName {
			newLink := Link{Task: linkTask.Name, Query: query, Status: status}
			task.Transition.Links = append(task.Transition.Links, newLink)
			return nil
		}
	}
	return nil
}

func (wf *Workflow[T]) NextStatus(current string) (string, error) {
	var value string
	lenStatuses := len(wf.statuses)
	if lenStatuses == 0 {
		ctx := context.Background()
		_, err := wf.calculateStatuses(ctx, false)
		if err != nil {
			return value, err
		}
		lenStatuses = len(wf.statuses)
	}
	index := slices.Index(wf.statuses, current)
	if index == -1 {
		return value, fmt.Errorf("NextStatus can't find current %s", current)
	}
	if index == lenStatuses-1 {
		return value, nil
	}
	value = wf.statuses[index+1]
	return value, nil
}

func (wf *Workflow[T]) PreviousStatus(current string) (string, error) {
	var value string
	lenStatuses := len(wf.statuses)
	if lenStatuses == 0 {
		ctx := context.Background()
		_, err := wf.calculateStatuses(ctx, false)
		if err != nil {
			return value, err
		}
		lenStatuses = len(wf.statuses)
	}
	index := slices.Index(wf.statuses, current)
	if index == -1 {
		return value, fmt.Errorf("NextStatus can't find current %s", current)
	}
	if index == 0 {
		return value, nil
	}
	value = wf.statuses[index-1]
	return value, nil
}

func (wf *Workflow[T]) calculateStatuses(ctx context.Context, quotes bool) ([]string, error) {
	var err error
	var start string
	list := []string{}
	start, err = wf.CalculateInitialTask()
	if err != nil {
		return list, err
	}
	current := start
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
	wf.statuses = list
	return list, nil

}

func (wf *Workflow[T]) Statuses(ctx context.Context, quotes bool) ([]string, error) {
	if wf.statuses != nil && len(wf.statuses) == 0 {
		return wf.statuses, nil
	}
	return wf.calculateStatuses(ctx, quotes)
}

func (wf *Workflow[T]) CalculateInitialTask() (string, error) {
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
func (wf *Workflow[T]) Load(ctx context.Context, jsonFile []byte) error {
	err := json.Unmarshal(jsonFile, &wf)
	return err
}

func (wf Workflow[T]) Marshal(ctx context.Context) ([]byte, error) {
	value, err := json.Marshal(wf)
	return value, err
}

func (wf Workflow[T]) ExtractKeys(trackingKeys []string) (string, string) {
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
