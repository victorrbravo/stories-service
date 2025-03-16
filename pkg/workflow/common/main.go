package common

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"stories/constants"
	db "stories/database"
	"stories/pkg/utils"
	"stories/types"
	"strings"
	"time"
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
	StartTask         = "Start"
)

const QueryStep = "{{.FieldName}}{{.Operator}}''{{.NameTask}}''"

const EqualOperator = "="

type QueryStepTemplate struct {
	FieldName string
	Operator  string
	NameTask  string
}

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

type Wrapped struct {
	Label    string        `json:"label"`
	Amount   types.String  `json:"amount"`
	Distance types.Float   `json:"distance"`
	Time     time.Duration `json:"time"`
	Comment  types.String  `json:"comment"`
}

type Data struct {
	Count        types.Int         `json:"count"` // checked
	TimeField    types.String      `json:"time_field"`
	InfoField    types.String      `json:"info_field"`
	AmountField  types.String      `json:"amount_field"`
	Roles        types.StringArray `json:"roles"`
	Images       types.StringArray `json:"images"`
	Videos       types.StringArray `json:"videos"`
	Percentage   float64           `json:"percentage"`
	UpdateAtDiff types.String      `json:"update_at_diff"`
	Summary      types.String      `json:"summary"`
	Amount       types.Float       `json:"amount"`
	Distance     types.Float       `json:"distance"`
	Total        types.Int         `json:"total"`
	FillValue    types.Int         `json:"fill_value"` // Measure progress in time, money or other type of value
	FillSelect   FillType          `json:"fill_select"`
	Rows         []interface{}     `json:"rows"`
	IDCurrent    types.Int         `json:"id_current"`
	ToExport     Exportable        `json:"to_export"`
	Stats        Values            `json:"stats"`
	ID           types.Int         `json:"id"`
	TipsIdField  types.Int         `json:"tips_id_field"`
	SubIdField   types.Int         `json:"fk_story_sub"`
	Position     types.Coordinates `json:"position"`
	Description  types.String      `json:"description"`
	DistanceDiff types.Float       `json:"distance_diff"`
	Tags         types.StringArray `json:"tags"`
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
	ID                types.UUID       `json:"id"`
	Name              string           `json:"name"`
	TableName         string           `json:"table_name"`
	FieldToken        string           `json:"field_token"`
	FieldStatus       string           `json:"field_status,omitempty"`
	TrackingTableName types.String     `json:"tracking_table_name"`
	TrackingKeys      types.String     `json:"tracking_keys"`
	TrackingID        types.String     `json:"tracking_id"`
	Variable          Variable         `json:"variable"`
	Tasks             map[string]*Task `json:"tasks"`
	Data              Data             `json:"data,omitempty"`
	CurrentIDInt      types.Int        `json:"current_id_int,omitempty"`
	CurrentIDUUID     types.UUID       `json:"current_id_uuid,omitempty"`
	Initial           types.String     `json:"initial"`
	Query             *db.IQuery       `json:"omitempty"`
	WhereClause       WhereClause[T]   `json:"where_clause"`
	ValueScanner      T                `json:"omitempty"`
	statuses          []string         `json:"omitempty"`
	JsonResponse      string           `json:"json_response"`
	Filter            WfFilter         `json:"filter,omitempty"`
	excludeReported   bool
}

type WfFilter struct {
	ExcludedWrapped []string `json:"excluded_wrapped"`
	IncludedWrapped []int    `json:"included_wrapped"`
	Sentence        string   `json:"sentence"`
}

type Variable struct {
	Name    string `json:"name"`
	Queries string `json:"queries"`
}

type Task struct {
	Name        string     `json:"name"`
	TimeField   string     `json:"time_field"`
	InfoField   string     `json:"info_field"`
	AmountField string     `json:"amount_field"`
	Transition  Transition `json:"transition"`
	Data        Data       `json:"data,omitempty"`
	Report      types.Bool `json:"report,omitempty"`
	VideosField string     `json:"videos_field" default:"videos"`
	ImagesField string     `json:"images_field" default:"images"`
	TipsIdField string     `json:"tips_id_field" default:"fk_story_tips"`
	SubIdField  string     `json:"sub_id_field" default:"fk_story_sub"`
	Position    string     `json:"position" default:"position"`
	Description string     `json:"description" default:"description"`
}

type Transition struct {
	Links   []Link   `json:"links"`
	Pattern string   `json:"pattern"`
	Wrapped *Wrapped `json:"wrapped"`
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
func (wf *Workflow[T]) ContainsAll(current string) error {
	wf.ExcludeReport(true)
	err := wf.Contains(current)
	wf.ExcludeReport(false)
	return err
}

func (wf *Workflow[T]) Contains(current string) error {
	var err error
	var statuses []string
	if len(current) == 0 || strings.Compare(current, constants.GolangNil) == 0 || strings.Compare(current, constants.EndTaskName) == 0 {
		return nil
	}
	ctx := context.Background()
	statuses, err = wf.Statuses(ctx, false)
	if err != nil {
		return err
	}
	index := slices.Index(statuses, current)
	if index == -1 {
		return fmt.Errorf("error: value '%s' is not found in list '%v'", current, statuses)
	}
	return nil
}

func (wf *Workflow[T]) NextStatus(current string) (string, error) {
	var value string
	if len(current) == 0 || strings.Compare(current, constants.GolangNil) == 0 {
		return value, nil
	}
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
	if len(current) == 0 || strings.Compare(current, constants.GolangNil) == 0 {
		return value, nil
	}
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
	if strings.Compare(title, StartTask) != 0 {
		list = append(list, title)
	}
	lenMax := len(wf.Tasks)
	i := 0
	for {
		task, ok := wf.Tasks[current]
		if !ok {
			return list, fmt.Errorf("Incorrect value for task name %s", current)
		}
		current, err = wf.findNextTask(task)
		log.Println("filterField findNextTask current", current, " err", err)
		if err != nil {
			return list, err
		}
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
	if len(list) == 0 {
		list = append(list, StartTask)
	}
	wf.statuses = list
	return list, nil

}

func containsIncluded[T db.Scanner](n int, wf *Workflow[T]) bool {
	if len(wf.Filter.IncludedWrapped) == 0 || wf.Filter.Sentence == "" {
		return true
	}
	log.Println("containsIncluded wf.Filter.IncludedWrapped:", wf.Filter.IncludedWrapped)
	log.Println("containsIncluded: n", n)
	for _, num := range wf.Filter.IncludedWrapped {
		if num == n {
			log.Println("containsIncluded...true n:", n, "wf.Filter.IncludedWrapped", wf.Filter.IncludedWrapped)
			return true
		}
	}
	log.Println("containsIncluded...false n:", n, "wf.Filter.IncludedWrapped", wf.Filter.IncludedWrapped)
	return false
}

func (wf *Workflow[T]) findNextTask(task *Task) (string, error) {
	nextName := ""
	var next *Task
	var ok bool
	if task == nil {
		return nextName, fmt.Errorf("findNextTask, Task is nil")
	}
	next = task
	for {
		nextName = next.Transition.Links[0].Task
		if len(nextName) == 0 || strings.Compare(nextName, EndTask) == 0 {
			return nextName, nil
		}
		next, ok = wf.Tasks[nextName]
		if !ok {

			return "", fmt.Errorf("findNextTask Incorrect value for task name %s", nextName)
		}
		idCurrent := types.AsInt(next.Data.ID)
		log.Println("findNextTask name", nextName, "idCurrent", idCurrent)
		if wf.excludeReported == false && (types.AsBool(next.Report) == false || !containsIncluded(idCurrent, wf)) {
			continue
		}
		break
	}

	return nextName, nil
}
func (wf *Workflow[T]) ExcludeReport(exclude bool) {
	wf.excludeReported = exclude
}

func (wf *Workflow[T]) Sentence(name string) string {
	_, linkStatus := wf.ExtractKeys()
	queryStepTmpl := QueryStepTemplate{
		NameTask:  name,
		Operator:  EqualOperator,
		FieldName: linkStatus,
	}
	sentence, err := utils.ParseTemplate(QueryStep, &queryStepTmpl)
	if err != nil {
		log.Print("error building sentence err:%v", err)
		return ""
	}
	return sentence
}

func (wf *Workflow[T]) Last(ctx context.Context) (*Task, error) {
	var task *Task
	statuses, err := wf.Statuses(ctx, false)
	if err != nil {
		return task, err
	}
	nLast := len(statuses)
	if nLast == 0 {
		return task, fmt.Errorf("Last function error nLast is equal to 0")
	}
	name := statuses[nLast-1]
	task, ok := wf.Tasks[name]
	if ok == false {
		return task, fmt.Errorf("task '%s' not found in statuses", name)
	}
	return task, nil
}

func (wf *Workflow[T]) Statuses(ctx context.Context, quotes bool) ([]string, error) {
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

func (wf Workflow[T]) ExtractKeys() (string, string) {
	linkStory := ""
	linkStatus := ""
	trackingKeys := strings.Split(wf.TrackingKeys.String(), ",")
	if len(trackingKeys) > 0 {
		linkStory = trackingKeys[0]
	}
	if len(trackingKeys) > 1 {
		linkStatus = trackingKeys[1]
	}
	return linkStory, linkStatus
}
