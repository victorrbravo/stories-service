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

	SelectRoleToken = `
SELECT {{.InfoField}} FROM {{.TrackingTable}}  JOIN  {{.TableName}} ON {{.Key1Link}}={{.Token}}                
{{if ne .Filter ""}} WHERE {{.Filter}} AND {{.PkField}} = '{{.Token}}{{end}}'
`

	SELECT_ROLEFIELD_TOKEN = `
SELECT {{.InfoField}}  FROM {{.TrackingTable}} tt WHERE {{.Key2Link}} = '{{.Status}}' AND {{.Key1Link}} = {{.Token}};
`
)

func getLinkFields[T db.Scanner](ctx context.Context, wf *common.Workflow[T]) (string, string) {
	return wf.ExtractKeys()
}

func TimeField[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {
	var duration scanners.String
	var outputBuffer bytes.Buffer
	ctx := context.Background()
	newQuery := wf.Query.NewQuery()
	return func(error, ...interface{}) (*common.Workflow[T], error) {
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
			durationValue, err := newQuery.QueryOne(newSQL, &duration)

			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				log.Println("error query err", err)
				continue
			}
			durationTime, ok := durationValue.(types.String)
			log.Println("ok", ok)
			value := types.NewString(NA)
			if !errors.Is(err, sql.ErrNoRows) && ok == true {
				value = types.NewString(durationTime.String())
			}
			current.Data.TimeField = value
			outputBuffer.Reset()
		}
		return wf, nil
	}
	return nil
}

func RoleField[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {
	var role scanners.String
	newQuery := wf.Query.NewQuery()
	return func(error, ...interface{}) (*common.Workflow[T], error) {
		for name, current := range wf.Tasks {
			newSQL, err := buildSQL(wf, name, SelectRoleToken)
			if err != nil {
				log.Println("Error Role err", err)
				continue
			}
			roleValue, err := newQuery.QueryOne(newSQL, &role)
			if err != nil {
				log.Println("Error Count SQL err", err)
				continue
			}
			roleString := roleValue.(types.String)
			current.Data.InfoField = roleString
		}
		return wf, nil
	}
}
func CountToken[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {
	var count scanners.Int
	newQuery := wf.Query.NewQuery()
	return func(error, ...interface{}) (*common.Workflow[T], error) {
		for name, current := range wf.Tasks {
			newSQL, err := buildSQL(wf, name, SelectCountToken)
			if err != nil {
				log.Println("Error Count err", err)
				continue
			}
			countValue, err := newQuery.QueryOne(newSQL, &count)
			if err != nil {
				log.Println("Error Count SQL err", err)
				continue
			}
			countInt := countValue.(types.Int)
			current.Data.Count = countInt
			if types.IsNil(current.Report) || current.Report.ToBool() == true {
				wf.Data.Total = types.NewInt(types.AsInt(wf.Data.Total) + types.AsInt(count.Value))
			}
			wf.Data.Stats.Assign(count.Value)
		}
		return wf, nil
	}
}
