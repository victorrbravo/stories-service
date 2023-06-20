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
SELECT {{.RoleField}} AS count
FROM  {{.TrackingTable}} where  {{.StatusField}} = {{.StatusValue}}
GROUP BY {{.RoleField}}
HAVING COUNT(*) = (
  SELECT MAX(count)
  FROM (
    SELECT {{.RoleField}}, COUNT(*) AS count
    FROM {{.TrackingTable}} tt where {{.StatusField}} = {{.StatusValue}} 
    GROUP BY {{.RoleField}}
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

func (pg *PgQuery) TimeField(ctx context.Context, wf *common.Workflow, task string) error {
	var count scanners.Interval
	newQuery := pg.query.NewQuery()
	current, ok := wf.Tasks[task]
	if ok == false {
		return errors.New(fmt.Sprintf("task doesn't exist", task))
	}
	nextTask := wf.NextTask(current.Name)
	if nextTask == nil {
		return errors.New(fmt.Sprintf("task has no a next task", task))
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
		Status:            task,
	}

	var outputBuffer bytes.Buffer
	err = tmpl.Execute(&outputBuffer, data)
	if err != nil {
		return err
	}
	newSQL = outputBuffer.String()
	err = newQuery.QueryOne(newSQL, &count)
	if err != nil {
		return err
	}
	fmt.Println("count max (interval)", count.Value.String())
	return nil

}

func extractKeys(tranckingKeys []string) (string, string) {
	linkStory := ""
	linkStatus := ""
	if len(tranckingKeys) > 0 {
		linkStory = tranckingKeys[0]
	}
	if len(tranckingKeys) > 1 {
		linkStatus = tranckingKeys[1]
	}
	return linkStory, linkStatus
}

func (pg *PgQuery) Count(ctx context.Context, wf *common.Workflow, task string) error {
	var count scanners.Int
	newQuery := pg.query.NewQuery()
	current := wf.Tasks[task]
	newSQL := fmt.Sprintf(SELECT_BASIC, wf.TableName, current.Transition.Links[0].Status,
		current.Transition.Links[0].Query)
	err := newQuery.QueryOne(newSQL, &count)
	if err != nil {
		return err
	}
	wf.Data.Count = count.Value
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
	for index, task := range wf.Tasks {
		err := query.Count(ctx, wf, index)
		if err != nil {
			return err
		}
		fmt.Println("task ", task.Name, "data count", wf.Data.Count)
		err = query.TimeField(ctx, wf, index)
		if err != nil {
			return err
		}
	}
	return nil
}
