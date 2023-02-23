package workflow

import (
	"context"
	"stories/types"
)

const (
	WorkflowName      = "worflowExample"
	Task1Name         = "Task1Name"
	QueryName         = "queries.FirstStation"
	NewUUID           = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	PatternSequencial = "sequencial"
	StationFirst      = "StationFirst"
)

type Workflow struct {
	ID                types.UUID   `json:"id"`
	Name              string       `json:"name"`
	TableName         string       `json:"table_name"`
	FieldToken        string       `json:"field_token"`
	TrackingTableName types.String `json:"tracking_table_name"`
	Variable          Variable     `json:"variable"`
	Tasks             []*Task      `json:"task"`
}

type Variable struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}

type Task struct {
	Name       string     `json:"name"`
	TimeField  string     `json:"time_field"`
	InfoField  string     `json:"info_field"`
	Transition Transition `json:"transition"`
}

type Transition struct {
	Links   []Link `json:"link"`
	Pattern string `json:"pattern"`
}

type Link struct {
	LinkTask *Task  `json:"link_task"`
	Query    string `json:"query"`
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
	err = wf.AddTask(ctx, StationFirst, nil)
	return &wf, err
}

func (wf *Workflow) AddTask(ctx context.Context, name string, linkTask *Task) error {
	newTask := new(Task)
	newTask.Name = name
	newTransition := Transition{Pattern: PatternSequencial}
	newLink := Link{linkTask, QueryName}
	newTransition.Links = append(newTransition.Links, newLink)
	newTask.Transition = newTransition
	wf.Tasks = append(wf.Tasks, newTask)
	return nil
}
