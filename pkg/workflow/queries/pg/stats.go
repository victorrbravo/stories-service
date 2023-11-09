package pg

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	db "stories/database"
	"stories/pkg/workflow/common"
	"stories/pkg/workflow/queries"
	"stories/scanners"
	"stories/types"
	"strings"
	"text/template"
)

const (
	NA               = "n/a"
	SELECT_BASIC     = "SELECT count(*) FROM %s WHERE %s = '%s'"
	SELECT_TIMEFIELD = "SELECT %s FROM %s ON %s"
	STATS_TIMEFIELD  = `
WITH cte AS (
SELECT {{.Key2Link}} ,{{.Key1Link}} , {{.TimeField}} ,
CASE
	WHEN EXTRACT (epoch from (LEAD({{.TimeField}}) over (ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.Statuses}}], {{.Key2Link}})) - {{.TimeField}})) < 0 then INTERVAL '0 hour'
	ELSE LEAD({{.TimeField}}) over (ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.Statuses}}], {{.Key2Link}})) - {{.TimeField}}
END AS diff
FROM {{.TrackingTable}} tt  ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.Statuses}}], {{.Key2Link}})  
)
SELECT 
CASE WHEN AVG(diff) IS NULL THEN '0s' ELSE AVG(diff) END AS maxtime
FROM cte WHERE {{.Key2Link}}='{{.Status}}'  GROUP BY {{.Key2Link}} ;
`

	STATS_ROLEFIELD = `
SELECT {{.InfoField}} AS roles
FROM  {{.TrackingTable}} where  {{.Key2Link}} = '{{.Status}}'
GROUP BY {{.InfoField}}
HAVING COUNT(*) = (
  SELECT MAX(count)
  FROM (
    SELECT {{.InfoField}}, COUNT(*) AS count
    FROM {{.TrackingTable}} tt where {{.Key2Link}} = '{{.Status}}' 
    GROUP BY {{.InfoField}}
  ) subquery
) limit 1;
`
)

type Any interface{}

type Monad func(error) (*common.Workflow, error)

func Initialize(config *db.DatabaseConfig) Monad {
	return func(e error) (*common.Workflow, error) {
		return nil, e
	}
}

func New(v *common.Workflow, config *db.DatabaseConfig) Monad {
	return func(s error) (*common.Workflow, error) {
		var pgQuery db.PostgresQuery
		pgQuery.Initialize(config)
		err := pgQuery.Connect()
		if err != nil {
			log.Println("Db Connect error")
			return nil, err
		}
		v.Query = &pgQuery
		return v, s
	}
}

func Bind(m Monad, f func(*common.Workflow) Monad) Monad {
	return func(s error) (*common.Workflow, error) {
		a, e := m(s)
		if e != nil {
			return nil, e
		}
		return f(a)(e)
	}
}

func TimeStats(wf *common.Workflow) Monad {
	log.Println("time Stats")
	ctx := context.Background()
	return func(error) (*common.Workflow, error) {
		var duration scanners.Interval
		var outputBuffer bytes.Buffer
		newQuery := wf.Query.NewQuery()
		for name, current := range wf.Tasks {
			nextTask := wf.NextTask(name)
			if nextTask == nil {
				current.Data.TimeField = types.NewString("0s")
				// End of flow
				continue
			}
			tmpl, err := template.New("timeSQL").Parse(STATS_TIMEFIELD)
			trackingKeys := strings.Split(wf.TrackingKeys.String(), ",")
			linkStory, linkStatus := wf.ExtractKeys(trackingKeys)
			list, err := wf.Statuses(ctx)
			if err != nil {
				return wf, err
			}
			statuses := strings.Join(list, ",")
			fmt.Println("...statuses", statuses)
			data := queries.DataTemplate{
				Key1Link:      linkStory,
				Key2Link:      linkStatus,
				TimeField:     current.TimeField,
				TrackingTable: wf.TrackingTableName.String(),
				Statuses:      statuses,
				Status:        name,
			}
			err = tmpl.Execute(&outputBuffer, data)
			if err != nil {
				return wf, err
			}
			newSQL := outputBuffer.String()
			err = newQuery.QueryOne(newSQL, &duration)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					return wf, err
				}
			}
			value := types.NewString(NA)
			if !errors.Is(err, sql.ErrNoRows) {
				value = types.NewString(duration.Value.String())
			}
			current.Data.TimeField = value
			outputBuffer.Reset()
		}
		return wf, nil
	}
}

func RoleStats(wf *common.Workflow) Monad {
	log.Println("unmarshalling the file to json")
	ctx := context.Background()
	var roles scanners.String
	var outputBuffer bytes.Buffer
	newQuery := wf.Query.NewQuery()
	return func(error) (*common.Workflow, error) {
		for name, current := range wf.Tasks {
			tmpl, err := template.New("roleSQL").Parse(STATS_ROLEFIELD)
			trackingKeys := strings.Split(wf.TrackingKeys.String(), ",")
			linkStory, linkStatus := wf.ExtractKeys(trackingKeys)
			list, err := wf.Statuses(ctx)
			if err != nil {
				return wf, err
			}
			statuses := strings.Join(list, ",")
			data := queries.DataTemplate{Key1Link: linkStory,
				Key2Link:      linkStatus,
				TimeField:     current.TimeField,
				InfoField:     current.InfoField,
				TrackingTable: wf.TrackingTableName.String(),
				Statuses:      statuses,
				Status:        name,
			}
			err = tmpl.Execute(&outputBuffer, data)
			if err != nil {
				return wf, err
			}
			newSQL := outputBuffer.String()
			err = newQuery.QueryOne(newSQL, &roles)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					return wf, err
				}
			}
			if errors.Is(err, sql.ErrNoRows) {
				current.Data.InfoField = types.NewString(NA)
				return wf, nil
			}
			current.Data.InfoField = roles.Value
			outputBuffer.Reset()
		}
		return wf, nil
	}
}

func Count(wf *common.Workflow) Monad {
	var count scanners.Int
	newQuery := wf.Query.NewQuery()
	return func(error) (*common.Workflow, error) {
		for _, current := range wf.Tasks {
			newSQL := fmt.Sprintf(SELECT_BASIC, wf.TableName, current.Transition.Links[0].Status,
				current.Transition.Links[0].Query)
			err := newQuery.QueryOne(newSQL, &count)
			if err != nil {
				continue
			}
			current.Data.Count = count.Value
		}
		return wf, nil
	}
}

/// ******* STATS Query

func Percentage(wf *common.Workflow) Monad {
	total := 0
	for _, task := range wf.Tasks {
		total = total + types.AsInt(task.Data.Count)
	}
	return func(error) (*common.Workflow, error) {
		for _, task := range wf.Tasks {
			value := float64(types.AsInt(task.Data.Count)) / float64(total)
			task.Data.Percentage = value
			task.Data.Total = types.NewInt(total)
		}
		wf.Data.Total = types.NewInt(total)
		return wf, nil
	}
}
