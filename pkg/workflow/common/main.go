package common

import (
	"context"
	"encoding/json"
	"stories/types"
)

const (
	WorkflowName      = "worflowExample"
	QueryName         = "queries.FirstStation"
	NewUUID           = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	PatternSequencial = "sequencial"
	StationFirst      = "StationFirst"
)

type FillType string

const (
	FillTime  FillType = "Time"
	FillMoney FillType = "Money"
)

type Data struct {
	Count      types.Int
	TimeField  types.String
	InfoField  types.String
	Roles      []types.String
	Porcentage types.Float
	Total      types.Int
	FillValue  types.Int // Measaure progress in time, money or other type of value
	FillSelect FillType
	Rows       []interface{}
}

type JobWorkflow struct {
	Workflow Workflow
}

type Workflow struct {
	ID                types.UUID   `json:"id"`
	Name              string       `json:"name"`
	TableName         string       `json:"table_name"`
	FieldToken        string       `json:"field_token"`
	TrackingTableName types.String `json:"tracking_table_name"`
	tracking_keys     types.String `json:"tracking_keys"`
	Variable          Variable     `json:"variable"`
	Tasks             []*Task      `json:"tasks"`
	Data              Data         `json:"data"`
}

type Variable struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}

type Task struct {
	Name       string       `json:"name"`
	TimeField  string       `json:"time_field"`
	InfoField  string       `json:"info_field"`
	Transition Transition   `json:"transition"`
	data       types.String `json:"data"`
}

type Transition struct {
	Links   []Link `json:"links"`
	Pattern string `json:"pattern"`
}

type Link struct {
	LinkTask string       `json:"task"`
	Query    string       `json:"query"`
	Status   string       `json:"status,omitempty"`
	data     types.String `json:"data"`
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
	newLink := Link{LinkTask: linkTask.Name, Query: QueryName, Status: status}
	newTransition.Links = append(newTransition.Links, newLink)
	newTask.Transition = newTransition
	wf.Tasks = append(wf.Tasks, newTask)
	return nil
}

// remove a task from the workflow by name
func (wf *Workflow) RemoveTask(ctx context.Context, name string) error {
	for i, task := range wf.Tasks {
		if task.Name == name {
			wf.Tasks = append(wf.Tasks[:i], wf.Tasks[i+1:]...)
			return nil
		}
	}
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

// add a link to a task
func (wf *Workflow) AddLink(ctx context.Context, taskName string, linkTask *Task, query string, status string) error {
	for _, task := range wf.Tasks {
		if task.Name == taskName {
			newLink := Link{LinkTask: linkTask.Name, Query: query, Status: status}
			task.Transition.Links = append(task.Transition.Links, newLink)
			return nil
		}
	}
	return nil
}

// load a wokflow instance using a json file
func (wf *Workflow) Load(ctx context.Context, jsonFile []byte) error {
	err := json.Unmarshal(jsonFile, &wf)
	return err
}
