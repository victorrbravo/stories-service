package pg

import (
	"bytes"
	"context"
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
	WHEN EXTRACT (epoch from (LEAD({{.TimeField}}) over (ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.Statuses}}], {{.Key2Link}})) - {{.TimeField}})) < 0 then INTERVAL '0 hour'
	ELSE LEAD({{.TimeField}}) over (ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.Statuses}}], {{.Key2Link}})) - {{.TimeField}}
END AS diff
FROM {{.TrackingTableName}} tt  ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.Statuses}}], {{.Key2Link}})  
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

type DataTemplate struct {
	Key1Link          string
	Key2Link          string
	TimeField         string
	InfoField         string
	TrackingTableName string
	Statuses          string
	Status            string
}

type StatsQuery struct {
	query *db.PostgresQuery
}

func (query *StatsQuery) Initialize(config *db.DatabaseConfig) error {
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(config)
	err := pgQuery.Connect()
	if err != nil {
		log.Println("Db Connect error")
		return err
	}
	query.query = &pgQuery
	return nil
}

func (query *StatsQuery) RoleField(ctx context.Context, wf *common.Workflow, id string) error {
	var roles scanners.String
	newQuery := query.query.NewQuery()
	current, ok := wf.Tasks[id]
	if ok == false {
		return fmt.Errorf("task doesn't exist %s", id)
	}
	newSQL := ""
	tmpl, err := template.New("roleSQL").Parse(STATS_ROLEFIELD)
	trackingKeys := strings.Split(wf.TrackingKeys.String(), ",")
	linkStory, linkStatus := extractKeys(trackingKeys)
	list, err := wf.Statuses(ctx)
	if err != nil {
		return err
	}
	statuses := strings.Join(list, ",")
	data := DataTemplate{Key1Link: linkStory,
		Key2Link:          linkStatus,
		TimeField:         current.TimeField,
		InfoField:         current.InfoField,
		TrackingTableName: wf.TrackingTableName.String(),
		Statuses:          statuses,
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

func (query *StatsQuery) TimeField(ctx context.Context, wf *common.Workflow, id string) error {
	var duration scanners.Interval
	var outputBuffer bytes.Buffer
	newQuery := query.query.NewQuery()
	current, ok := wf.Tasks[id]
	if ok == false {
		return fmt.Errorf("id doesn't exist %s", id)
	}
	nextTask := wf.NextTask(current.Name)
	if nextTask == nil {
		current.Data.TimeField = types.NewString("0s")
		// End of flow
		return nil
	}
	newSQL := ""
	tmpl, err := template.New("timeSQL").Parse(STATS_TIMEFIELD)
	trackingKeys := strings.Split(wf.TrackingKeys.String(), ",")
	linkStory, linkStatus := extractKeys(trackingKeys)
	list, err := wf.Statuses(ctx)
	if err != nil {
		return err
	}
	statuses := strings.Join(list, ",")
	fmt.Println("...statuses", statuses)
	data := DataTemplate{
		Key1Link:          linkStory,
		Key2Link:          linkStatus,
		TimeField:         current.TimeField,
		TrackingTableName: wf.TrackingTableName.String(),
		Statuses:          statuses,
		Status:            id,
	}
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

func (query *StatsQuery) Percentage(ctx context.Context, wf *common.Workflow, id string) error {
	total := 0
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
func (query *StatsQuery) Count(ctx context.Context, wf *common.Workflow, id string) error {
	var count scanners.Int
	newQuery := query.query.NewQuery()
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

func (query *StatsQuery) Generate(ctx context.Context, wf *common.Workflow) error {
	for name, _ := range wf.Tasks {
		err := query.Count(ctx, wf, name)
		if err != nil {
			return err
		}
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
