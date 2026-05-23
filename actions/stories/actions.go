package stories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	db "stories/database"
	"stories/generated/models"
	"stories/pkg/utils"
	"stories/pkg/workflow"
	"stories/pkg/workflow/common"
	"stories/pkg/workflow/queries/pg"
	"stories/scanners"
	"stories/types"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
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

const AddNewStepPolicy = `UPDATE ta_policy
SET condition = condition || jsonb_build_array(jsonb_build_object('ticket', $1::integer))
WHERE name = 'add_flow_step';
`
const DeleteNewStepPolicy = `UPDATE ta_policy
SET condition = (
    SELECT jsonb_agg(elem)
    FROM jsonb_array_elements(condition) AS elem
    WHERE elem->>'ticket' <> '205'
)
WHERE name = 'add_flow_step'
`

const SelectStatusTagsByStory = `
SELECT ta_tracking.fk_status, ta_tracking.tags, ta_tracking.updated_at
FROM ta_tracking
WHERE ta_tracking.fk_story = $1
`

// StatusInfo aggregates per-status data extracted from ta_tracking:
// the union of all tags seen for that status under a given story plus
// the earliest updated_at timestamp (so within a tag group, statuses
// can later be ordered ascending by time).
type StatusInfo struct {
	Tags      []string
	UpdatedAt time.Time
}

type ShowFlowInput struct {
	PkFlow         int                             `yaml:"flow" default:"1"`
	Ticket         types.Int                       `yaml:"ticket"`
	ShowEmptyTasks bool                            `yaml:"show_empty_tasks" default:"false"`
	Task           WorkflowType[models.TaTracking] `yaml:"task,omitempty"`
	Filter         types.String                    `yaml:"filter,omitempty"`
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
	yalmOrder := utils.FormatYaml(order)
	getFlow, err := UnmarshallGetFlow([]byte(yalmOrder), query)
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
	var diffStatuses []string
	statuses, err := wf.Statuses(ctx, false, true)
	if len(wf.Filter.IncludedWrapped) > 0 {
		diffStatuses = workflow.SliceDiff(statuses, wf.Filter.IncludedWrapped)
		statuses = wf.Filter.IncludedWrapped
	}
	log.Println("UpdateWorkflowDataInt", len(diffStatuses), len(wf.Filter.IncludedWrapped))
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
		return firstDate.Before(secondDate)
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

// TagsByStatus returns a map keyed by ta_tracking.fk_status whose value
// is a StatusInfo aggregating the rows of the given fk_story. Tags from
// every tracking row of the same status are unioned (deduplicated). The
// UpdatedAt field captures the earliest updated_at across those rows so
// callers can order statuses by time.
func TagsByStatus(query db.IQuery, fkStory types.Int) (map[string]StatusInfo, error) {
	result := make(map[string]StatusInfo)
	var scanner scanners.TupleStringStringArrayDateTime

	rows, err := query.QueryMany(SelectStatusTagsByStory, &scanner, fkStory)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, nil
		}
		log.Println("TagsByStatus QueryMany err", err)
		return result, err
	}

	// Per-status set used to deduplicate tags across multiple rows.
	seen := make(map[string]map[string]struct{})

	for _, row := range rows {
		tuple, ok := row.(scanners.TupleStringStringArrayDateTime)
		if !ok {
			log.Println("TagsByStatus unexpected row type")
			continue
		}
		status := tuple.First.String()
		if status == "" {
			continue
		}
		info := result[status]
		if seen[status] == nil {
			seen[status] = make(map[string]struct{})
		}
		for _, tag := range tuple.Second.Array {
			if _, dup := seen[status][tag]; dup {
				continue
			}
			seen[status][tag] = struct{}{}
			info.Tags = append(info.Tags, tag)
		}
		rowTime := tuple.Third.DateTime
		if !rowTime.IsZero() {
			if info.UpdatedAt.IsZero() || rowTime.Before(info.UpdatedAt) {
				info.UpdatedAt = rowTime
			}
		}
		result[status] = info
	}
	return result, nil
}

// OrderTasksByTags rebuilds wf.Tasks so that every task is preserved,
// but their iteration order reflects the priority list provided in
// wantedTags.
//
// The natural traversal is established by walking the workflow chain
// starting at common.StartTask ("Start") and following
// Task -> Transition -> Links[0] -> Task on each hop. Any orphan task
// (not reachable from "Start") is appended afterwards so nothing is
// lost.
//
// The final order is:
//  1. The "Start" task itself.
//  2. For each tag in wantedTags (in order), the tasks whose entry in
//     tagsMap contains that tag, sorted by UpdatedAt ascending within
//     the tag bucket.
//  3. Any remaining tasks, in chain order.
//
// Example: with wantedTags = []string{"bill", "coffee"} the result
// starts with "Start", then all tasks whose tagsMap value contains
// "bill" (oldest UpdatedAt first), then all tasks whose tagsMap value
// contains "coffee" (oldest first), then the rest.
//
// tagsMap is keyed by task name (matching wf.Tasks keys / fk_status
// values returned by TagsByStatus) and its value carries both the union
// of tags and the earliest UpdatedAt for that status. An empty/nil
// wantedTags is treated as "no priority" — every task keeps the chain
// order.
//
// The same map-rebuild pattern used by pg.OrderingTaskField is applied
// here: a fresh map is assigned to wf.Tasks so the insertion order
// reflects the priority walk.
func OrderTasksByTags[T db.Scanner](
	wf *common.Workflow[T],
	tagsMap map[string]StatusInfo,
	wantedTags []string,
) *common.Workflow[T] {
	if wf == nil || len(wf.Tasks) == 0 {
		return wf
	}

	// 1. Walk the chain from "Start" to establish the natural chain order.
	chainOrder := make([]string, 0, len(wf.Tasks))
	visited := make(map[string]bool, len(wf.Tasks))
	current := common.StartTask
	if _, ok := wf.Tasks[current]; !ok {
		log.Println("OrderTasksByTags: 'Start' task not found in wf.Tasks")
	}
	for {
		task, ok := wf.Tasks[current]
		if !ok || visited[current] {
			break
		}
		visited[current] = true
		chainOrder = append(chainOrder, current)

		if len(task.Transition.Links) == 0 {
			break
		}
		next := task.Transition.Links[0].Task
		if next == "" || next == common.EndTask {
			break
		}
		current = next
	}
	// Append any orphan tasks (not reachable from "Start") so all
	// tasks are preserved in the rebuilt map.
	for name := range wf.Tasks {
		if !visited[name] {
			chainOrder = append(chainOrder, name)
			visited[name] = true
		}
	}

	// hasTag returns true when tagsMap[taskName].Tags contains the
	// given tag. The comparison is case-insensitive (strings.EqualFold),
	// so e.g. wantedTags=["Bill"] matches a stored tag "bill" or "BILL".
	hasTag := func(taskName, tag string) bool {
		info, ok := tagsMap[taskName]
		if !ok {
			return false
		}
		for _, t := range info.Tags {
			if strings.EqualFold(t, tag) {
				return true
			}
		}
		return false
	}

	// timeOf returns the UpdatedAt timestamp registered for a task in
	// tagsMap. Tasks not present in tagsMap (e.g. "Start") get the zero
	// value, which sorts ahead of every real timestamp.
	timeOf := func(taskName string) time.Time {
		if info, ok := tagsMap[taskName]; ok {
			return info.UpdatedAt
		}
		return time.Time{}
	}

	ordered := make(map[string]*common.Task, len(wf.Tasks))
	added := make(map[string]bool, len(wf.Tasks))
	// priorityOrder records the final task sequence (without END).
	// We later use it to rewire each task's Transition.Links[0] so
	// that wf.Statuses() walks the chain in this same order.
	priorityOrder := make([]string, 0, len(wf.Tasks))

	addToOrdered := func(name string) {
		if added[name] {
			return
		}
		if task, ok := wf.Tasks[name]; ok {
			ordered[name] = task
			added[name] = true
			priorityOrder = append(priorityOrder, name)
		}
	}

	// 2. "Start" goes first.
	addToOrdered(common.StartTask)

	// 3. For each wanted tag (in priority order), gather every task
	//    whose tagsMap entry contains that tag, sort the bucket by
	//    UpdatedAt ascending, and append in that order.
	for _, tag := range wantedTags {
		bucket := make([]string, 0)
		for _, name := range chainOrder {
			if !added[name] && hasTag(name, tag) {
				bucket = append(bucket, name)
			}
		}
		sort.SliceStable(bucket, func(i, j int) bool {
			return timeOf(bucket[i]).Before(timeOf(bucket[j]))
		})
		for _, name := range bucket {
			addToOrdered(name)
		}
	}

	// 4. Append every remaining task in chain order.
	for _, name := range chainOrder {
		if !added[name] {
			addToOrdered(name)
		}
	}

	// 5. Rewire each task's Transition.Links[0] so that the chain
	//    walk (Task -> Transition -> Links[0] -> Task) follows
	//    priorityOrder. The last task in the order links to
	//    common.EndTask, which is the sentinel used by
	//    Workflow.calculateStatuses to terminate iteration.
	for i, name := range priorityOrder {
		task := ordered[name]
		if task == nil {
			continue
		}
		var nextName string
		if i < len(priorityOrder)-1 {
			nextName = priorityOrder[i+1]
		} else {
			nextName = common.EndTask
		}
		nextQuery := ""
		if nextName != common.EndTask {
			nextQuery = wf.Sentence(nextName)
		}
		if len(task.Transition.Links) == 0 {
			task.Transition.Links = []common.Link{{
				Task:  nextName,
				Query: nextQuery,
			}}
		} else {
			task.Transition.Links[0].Task = nextName
			task.Transition.Links[0].Query = nextQuery
		}
	}

	// Invalidate any cached statuses so subsequent calls to
	// wf.Statuses(ctx, ..., recalculate=true) actually re-walk the
	// rewired chain.
	wf.CacheStatuses = common.CacheStatuses{}

	wf.Tasks = ordered
	return wf
}

func UpdateWorkflow[T db.Scanner](getFlow ShowFlowInput, query *db.IQuery,
	values []interface{},
	funcs ...func(*common.Workflow[T], ...interface{}) pg.Monad[T]) (*common.Workflow[T], error) {
	wf, err := GetWorkflow[T](getFlow.PkFlow, query)
	if err != nil {
		return nil, err
	}
	tagsMap, err := TagsByStatus(*query, getFlow.Ticket)
	if err != nil {
		log.Println("FilterField err", err)
		return nil, err
	}
	log.Println("tagMaps:", tagsMap)
	nwWf := OrderTasksByTags(wf, tagsMap, []string{"bill", "coffee"})
	ctx := context.Background()
	newStatuses, err := nwWf.Statuses(ctx, false, true)
	log.Println("newStatuses:", newStatuses)
	wf.Filter.Sentence = getFlow.Filter.String()
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
