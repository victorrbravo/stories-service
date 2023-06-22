package query

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	db "stories/database"
	"stories/pkg/workflow/common"
	"stories/scanners"
	"stories/types"
	"strings"
	"text/template"
)

const (
	SELECT_BASIC     = "SELECT count(*) FROM %s WHERE %s = '%s'"
	SELECT_TIMEFIELD = "SELECT %s FROM %s ON %s"
	STATS_TIMEFIELD  = `
WITH cte AS (
SELECT {{.Key2Link}} ,{{.Key1Link}} , {{.TimeField}} ,
CASE
	WHEN EXTRACT (epoch from (LEAD({{.TimeField}}) over (ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.StatusesList}}], {{.Key2Link}})) - {{.TimeField}})) < 0 then INTERVAL '0 hour'
	ELSE LEAD({{.TimeField}}) over (ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.StatusesList}}], {{.Key2Link}})) - {{.TimeField}}
END AS diff
FROM {{.TrackingTableName}} tt  ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.StatusesList}}], {{.Key2Link}})  
)
SELECT 
MAX(diff) AS maxtime
FROM cte WHERE {{.Key2Link}}='{{.Status}}'  GROUP BY {{.Key2Link}} ;
`

	STATS_ROLEFIELD = `
SELECT {{.InfoField}} AS roles
FROM  {{.TrackingTableName}} where  {{.Key2Link}} = '{{.Status}}'
GROUP BY {{.InfoField}}
HAVING COUNT(*) = (
  SELECT MAX(count)
  FROM (
    SELECT {{.InfoField}}, COUNT(*) AS count
    FROM {{.TrackingTableName}} tt where {{.Key2Link}} = '{{.Status}}' 
    GROUP BY {{.InfoField}}
  ) subquery
) limit 1;
`
)

type PgQuery struct {
	query *db.PostgresQuery
}

func (pg *PgQuery) Initialize(config *db.DatabaseConfig) error {
	var query db.PostgresQuery
	query.Initialize(config)
	err := query.Connect()
	if err != nil {
		log.Println("Db Connect error")
		return err
	}
	pg.query = &query
	return nil
}

func GenerateJoins(tables ...string) string {
	value := strings.Join(tables, " JOIN ")
	return value
}

func GenerateOnJoins(tracking string, filter string, task *common.Task) string {
	values := strings.Split(filter, ",")
	statement := ""
	if len(values) > 0 {
		value := values[0]
		statement = statement + fmt.Sprintf("%s.%s = %s.%s ", tracking, value, "1")
	}
	if len(values) > 1 {
		value := values[1]
		statement = statement + " AND "
		statement = statement + fmt.Sprintf("%s.%s = '%s' ", tracking, value, task.Transition.Links[0].Query)
	}
	return statement
}

func (pg *PgQuery) RoleField(ctx context.Context, wf *common.Workflow, id string) error {
	var roles scanners.String
	newQuery := pg.query.NewQuery()
	current, ok := wf.Tasks[id]
	if ok == false {
		return errors.New(fmt.Sprintf("task doesn't exist", id))
	}
	newSQL := ""
	tmpl, err := template.New("roleSQL").Parse(STATS_ROLEFIELD)
	trackingKeys := strings.Split(wf.TrackingKeys.String(), ",")
	linkStory, linkStatus := extractKeys(trackingKeys)
	data := struct {
		Key1Link          string
		Key2Link          string
		TimeField         string
		InfoField         string
		TrackingTableName string
		StatusesList      string
		Status            string
	}{
		Key1Link:          linkStory,
		Key2Link:          linkStatus,
		TimeField:         current.TimeField,
		InfoField:         current.InfoField,
		TrackingTableName: wf.TrackingTableName.String(),
		StatusesList:      pg.StatusesList(ctx, wf),
		Status:            id,
	}
	var outputBuffer bytes.Buffer
	err = tmpl.Execute(&outputBuffer, data)
	if err != nil {
		return err
	}
	newSQL = outputBuffer.String()
	err = newQuery.QueryOne(newSQL, &roles)
	if err != nil {
		return err
	}
	current.Data.InfoField = roles.Value
	return nil
}

func (pg *PgQuery) TimeField(ctx context.Context, wf *common.Workflow, id string) error {
	var duration scanners.Interval
	newQuery := pg.query.NewQuery()
	current, ok := wf.Tasks[id]
	if ok == false {
		return errors.New(fmt.Sprintf("id doesn't exist", id))
	}
	nextTask := wf.NextTask(current.Name)
	if nextTask == nil {
		return errors.New(fmt.Sprintf("id has no a next id", id))
	}
	newSQL := ""
	tmpl, err := template.New("timeSQL").Parse(STATS_TIMEFIELD)
	trackingKeys := strings.Split(wf.TrackingKeys.String(), ",")
	linkStory, linkStatus := extractKeys(trackingKeys)
	data := struct {
		Key1Link          string
		Key2Link          string
		TimeField         string
		TrackingTableName string
		StatusesList      string
		Status            string
	}{
		Key1Link:          linkStory,
		Key2Link:          linkStatus,
		TimeField:         current.TimeField,
		TrackingTableName: wf.TrackingTableName.String(),
		StatusesList:      pg.StatusesList(ctx, wf),
		Status:            id,
	}

	var outputBuffer bytes.Buffer
	err = tmpl.Execute(&outputBuffer, data)
	if err != nil {
		return err
	}
	newSQL = outputBuffer.String()
	err = newQuery.QueryOne(newSQL, &duration)
	if err != nil {
		return err
	}
	current.Data.TimeField = types.NewString(duration.Value.String())
	return nil

}

func extractKeys(trackingKeys []string) (string, string) {
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

func (pg *PgQuery) Percentage(ctx context.Context, wf *common.Workflow, id string) error {
	total := 0
	if !types.IsNil(wf.Data.Total) {
		total = types.AsInt(wf.Data.Total)
		return nil
	}
	for _, task := range wf.Tasks {
		total = total + types.AsInt(task.Data.Count)
	}

	for _, task := range wf.Tasks {
		value := float64(types.AsInt(task.Data.Count)) / float64(total)
		task.Data.Percentage = value
		task.Data.Total = types.NewInt(total)
	}
	wf.Data.Total = types.NewInt(total)
	return nil
}
func (pg *PgQuery) Count(ctx context.Context, wf *common.Workflow, id string) error {
	var count scanners.Int
	newQuery := pg.query.NewQuery()
	current, ok := wf.Tasks[id]
	if !ok { // task is not found
		return fmt.Errorf("Task not found: %s", id)
	}
	newSQL := fmt.Sprintf(SELECT_BASIC, wf.TableName, current.Transition.Links[0].Status,
		current.Transition.Links[0].Query)
	err := newQuery.QueryOne(newSQL, &count)
	if err != nil {
		return err
	}
	current.Data.Count = count.Value
	return nil
}

func (query *PgQuery) StatusesList(ctx context.Context, wf *common.Workflow) string {
	tasksList := ""
	if true {
		return "'ToDo','Progress','Done'"
	}
	for key, _ := range wf.Tasks {
		tasksList = tasksList + "'" + key + "'"
		tasksList = tasksList + ","
	}
	tasksList = strings.TrimRight(tasksList, ", ")
	return tasksList
}

func (query *PgQuery) Generate(ctx context.Context, wf *common.Workflow) error {
	for name, task := range wf.Tasks {
		err := query.Count(ctx, wf, name)
		if err != nil {
			return err
		}
		fmt.Println("task ", task.Name, "data count", wf.Data.Count)
		err = query.RoleField(ctx, wf, name)
		if err != nil {
			return err
		}
		err = query.TimeField(ctx, wf, name)
		if err != nil {
			return err
		}
	}
	for name, _ := range wf.Tasks {
		err := query.Percentage(ctx, wf, name)
		if err != nil {
			return err
		}
	}
	return nil
}
