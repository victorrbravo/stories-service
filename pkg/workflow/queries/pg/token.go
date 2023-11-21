package pg

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"stories/pkg/workflow/common"
	"stories/pkg/workflow/queries"
	"stories/scanners"
	"strings"
	"text/template"
)

const (
	SELECT_WHERE_TOKEN = `
SELECT COUNT(*) FROM {{.TableName}} WHERE {{.PkField}} = '{{.Token}}' AND {{.PkStatus}} = '{{.Status}}';
`
	SELECT_TIMEFIELD_TOKEN = `
WITH cte AS (
SELECT {{.Key2Link}} ,{{.Key1Link}} , {{.TimeField}} ,
CASE
	WHEN EXTRACT (epoch from (LEAD({{.TimeField}}) over (ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.Statuses}}], {{.Key2Link}})) - {{.TimeField}})) < 0 then INTERVAL '0 hour'
	ELSE LEAD({{.TimeField}}) over (ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.Statuses}}], {{.Key2Link}})) - {{.TimeField}}
END AS diff
FROM {{.TrackingTable}} tt  WHERE {{.Key1Link}}='{{.Token}}' ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.Statuses}}], {{.Key2Link}})  
)
SELECT 
CASE WHEN diff IS NULL THEN '0s' ELSE diff END AS maxtime
FROM cte WHERE {{.Key2Link}}='{{.Status}}' AND  {{.Key1Link}}='{{.Token}}'  ;

`

	SELECT_ROLEFIELD_TOKEN = `
SELECT {{.InfoField}}  FROM {{.TrackingTable}} tt WHERE {{.Key2Link}} = '{{.Status}}' AND {{.Key1Link}} = {{.Token}};
`
)

func getLinkFields(ctx context.Context, wf *common.Workflow) (string, string) {
	trackingKeys := strings.Split(wf.TrackingKeys.String(), ",")
	return wf.ExtractKeys(trackingKeys)
}

func TimeField(wf *common.Workflow) Monad {
	var duration scanners.String
	var outputBuffer bytes.Buffer
	ctx := context.Background()
	newQuery := wf.Query.NewQuery()
	return func(error) (*common.Workflow, error) {
		for name, current := range wf.Tasks {
			tmpl, err := template.New("timeSQL").Parse(SELECT_TIMEFIELD_TOKEN)
			linkStory, linkStatus := getLinkFields(ctx, wf)
			list, err := wf.Statuses(ctx, true)
			if err != nil {
				log.Println("Error statuses err", err)
				continue
			}
			statuses := strings.Join(list, ",")
			fmt.Println("...statuses", statuses)
			data := queries.DataTemplate{
				Key1Link:      linkStory,
				Key2Link:      linkStatus,
				TimeField:     current.TimeField,
				TrackingTable: wf.TrackingTableName.String(),
				Status:        name,
				Statuses:      statuses,
				Token:         wf.CurrentIDInt.String(),
			}
			err = tmpl.Execute(&outputBuffer, data)
			if err != nil {
				log.Println("error template execution err", err)
			}
			newSQL := outputBuffer.String()
			err = newQuery.QueryOne(newSQL, &duration)
			if err != nil && err != sql.ErrNoRows {
				log.Println("error query err", err)
				continue
			}
			current.Data.TimeField = duration.Value
			outputBuffer.Reset()
		}
		return wf, nil
	}
	return nil
}

func RoleField(wf *common.Workflow) Monad {
	var role scanners.String
	var outputBuffer bytes.Buffer
	newQuery := wf.Query.NewQuery()
	ctx := context.Background()
	return func(error) (*common.Workflow, error) {
		for name, current := range wf.Tasks {
			tmpl, err := template.New("timeSQL").Parse(SELECT_ROLEFIELD_TOKEN)
			linkStory, linkStatus := getLinkFields(ctx, wf)
			data := queries.DataTemplate{
				Key1Link:      linkStory,
				Key2Link:      linkStatus,
				InfoField:     current.InfoField,
				TrackingTable: wf.TrackingTableName.String(),
				Status:        name,
				Token:         wf.CurrentIDInt.String(),
			}
			err = tmpl.Execute(&outputBuffer, data)
			if err != nil {
				log.Println("Error template execute err", err)
				continue
			}
			newSQL := outputBuffer.String()
			err = newQuery.QueryOne(newSQL, &role)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				log.Println("Error Query one err", err)
				continue
			}
			current.Data.InfoField = role.Value
			outputBuffer.Reset()
		}
		return wf, nil
	}
}

func CountToken(wf *common.Workflow) Monad {
	var count scanners.Int
	var outputBuffer bytes.Buffer
	newQuery := wf.Query.NewQuery()
	return func(error) (*common.Workflow, error) {
		for name, current := range wf.Tasks {

			tmpl, err := template.New("roleSQL").Parse(SELECT_WHERE_TOKEN)
			data := queries.DataTemplate{
				TableName: wf.TableName,
				Status:    name,
				Token:     wf.CurrentIDInt.String(),
				PkStatus:  current.Transition.Links[0].Status,
				PkField:   wf.FieldToken,
			}
			err = tmpl.Execute(&outputBuffer, data)
			if err != nil {
				log.Println("Error execute")
				continue
			}
			newSQL := outputBuffer.String()
			err = newQuery.QueryOne(newSQL, &count)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				log.Println("Error query token err", err)
				continue
			}
			current.Data.Count = count.Value
			wf.Data.Stats.Assign(count.Value)
			outputBuffer.Reset()
		}
		return wf, nil
	}
}
