package pg

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
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
	WHITE            = "#FFFFFF"
	NA               = "nil"
	SELECT_BASIC     = "SELECT count(*) FROM %s JOIN %s ON %s=%s WHERE %s"
	SelectCountStats = `
SELECT count(*) FROM {{.TrackingTable}}  JOIN  {{.TableName}} ON {{.Key1Link}}={{.PkField}}                
{{if ne .Filter ""}} WHERE {{.Filter}} {{end}}
`

	SelectCountToken = `
SELECT count(*) FROM {{.TrackingTable}}  JOIN  {{.TableName}} ON {{.Key1Link}}={{.PkField}}                
{{if ne .Filter ""}} WHERE {{.Filter}} AND {{.PkField}} = '{{.Token}}' AND {{.PkStatus}} = '{{.Status}}'{{end}}
`
	SelectRowsStats = `
SELECT {{.Fields}} FROM {{.TrackingTable}}  JOIN  {{.TableName}} ON {{.Key1Link}}={{.Token}}                
{{if ne .Filter ""}} WHERE {{.Filter}} {{end}}
`
	STATS_TIMEFIELD = `
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
FROM cte WHERE {{.Filter}}  GROUP BY {{.Key2Link}} ;
`

	STATS_ROLEFIELD = `
SELECT {{.InfoField}} AS roles
FROM  {{.TrackingTable}} WHERE {{.Filter}}
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

type Monad[T db.Scanner] func(error, ...interface{}) (*common.Workflow[T], error)

func New[T db.Scanner](wf *common.Workflow[T], q *db.IQuery) Monad[T] {
	return func(s error, v ...interface{}) (*common.Workflow[T], error) {
		wf.Query = q
		if wf.WhereClause == nil {
			wf.WhereClause = CalculateWhere[T]
		}
		return wf, s
	}
}

func Bind[T db.Scanner](m Monad[T], f func(*common.Workflow[T], ...interface{}) Monad[T], values ...interface{}) Monad[T] {
	return func(s error, v ...interface{}) (*common.Workflow[T], error) {
		a, e := m(s)
		if e != nil {
			return nil, e
		}
		return f(a, values)(e)
	}
}

func TimeStats[T db.Scanner](wf *common.Workflow[T], names ...interface{}) Monad[T] {
	log.Println("time Stats")
	return func(error, ...interface{}) (*common.Workflow[T], error) {
		var duration scanners.Interval

		newQuery := wf.Query
		for name, current := range wf.Tasks {
			if !types.IsNil(current.Report) && current.Report.ToBool() == false {
				continue
			}
			nextTask := wf.NextTask(name)
			if nextTask == nil {
				current.Data.TimeField = types.NewString("0s")
				// End of flow
				continue
			}
			newSQL, err := buildSQL(wf, name, STATS_TIMEFIELD)
			if err != nil {
				return wf, err
			}
			durationValue, err := (*newQuery).QueryOne(newSQL, &duration)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					log.Println("Error TimeStats on sql task", name)
					current.Data.TimeField = types.NewString(NA)
					// FIXME: error catch here
					continue
				}
			}
			durationTime, ok := durationValue.(*scanners.Interval)
			log.Println("ok", ok)
			value := types.NewString(NA)
			if !errors.Is(err, sql.ErrNoRows) {
				value = types.NewString(durationTime.Value.String())
			}
			current.Data.TimeField = value

		}
		return wf, nil
	}
}

func RoleStats[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {
	log.Println("unmarshalling the file to json")
	var roles scanners.String
	newQuery := wf.Query
	return func(error, ...interface{}) (*common.Workflow[T], error) {
		for name, current := range wf.Tasks {
			if !types.IsNil(current.Report) && current.Report.ToBool() == false {
				continue
			}
			newSQL, err := buildSQL(wf, name, STATS_ROLEFIELD)
			if err != nil {
				log.Println("Error RoleStats err", err)
				continue
			}
			rolesValue, err := (*newQuery).QueryOne(newSQL, &roles)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					log.Println("Error query one err", err)
					continue
				}
			}
			if errors.Is(err, sql.ErrNoRows) || rolesValue == nil {
				current.Data.InfoField = types.NewString(NA)
				continue
			}
			rolesString := rolesValue.(types.String)
			current.Data.InfoField = rolesString
		}
		return wf, nil
	}
}

func Count[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {
	var count scanners.Int
	newQuery := wf.Query
	return func(error, ...interface{}) (*common.Workflow[T], error) {
		for name, current := range wf.Tasks {
			newSQL, err := buildSQL(wf, name, SelectCountStats)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					log.Println("Error Count err", err)
				}
				continue
			}
			countValue, err := (*newQuery).QueryOne(newSQL, &count)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					log.Println("Error Count SQL err", err)
				}
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

func AllRows[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {
	newQuery := wf.Query
	var lookTask string
	return func(error, ...interface{}) (*common.Workflow[T], error) {
		lookTask = firstParameter(parameters)
		for name, current := range wf.Tasks {
			if len(lookTask) > 0 && strings.Compare(name, lookTask) != 0 {
				continue
			}
			newSQL, err := buildSQL[T](wf, name, SelectRowsStats)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					log.Println("Error Count err", err)
				}
				continue
			}
			current.Data.Rows, err = (*newQuery).QueryMany(newSQL, wf.ValueScanner)
			if err != nil {
				log.Println("Error Rows SQL err", err)
			}
		}
		return wf, nil
	}
}

func firstParameter(parameters []interface{}) string {
	var parSlice []interface{}
	var ok bool
	var value string
	if len(parameters) > 0 {
		parSlice, ok = parameters[0].([]interface{})
		if !ok {
			return ""
		}
		if len(parSlice) > 0 {
			value, ok = parSlice[0].(string)
			if !ok {
				return ""
			}
		}
	}
	return value
}

func Total[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {
	return func(error, ...interface{}) (*common.Workflow[T], error) {
		total := 0
		if types.IsNil(wf.Data.Total) {
			for _, task := range wf.Tasks {
				if types.IsNil(task.Report) || task.Report.ToBool() == true {
					total = total + types.AsInt(task.Data.Count)
				}
			}
			wf.Data.Total = types.NewInt(total)
		}
		return wf, nil
	}
}

func Percentage[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {

	return func(error, ...interface{}) (*common.Workflow[T], error) {
		for _, task := range wf.Tasks {
			var value float64
			if types.IsNil(task.Report) || task.Report.ToBool() == true {
				total := float64(types.AsInt(wf.Data.Total))
				if total > 0 {
					value = float64(types.AsInt(task.Data.Count)) / total
				}
			}
			task.Data.Percentage = value
		}
		return wf, nil
	}
}

func buildSQL[T db.Scanner](wf *common.Workflow[T], name string, sqlString string) (string, error) {
	var outputBuffer bytes.Buffer
	ctx := context.Background()
	tmpl, err := template.New("timeSQL").Parse(sqlString)
	linkStory, linkStatus := wf.ExtractKeys()
	list, err := wf.Statuses(ctx, true)
	if err != nil {
		return "", err
	}
	statuses := strings.Join(list, ",")
	queryValue, err := ReducesFilter(*wf, name)
	if err != nil {
		return queryValue, err
	}
	current := wf.Tasks[name]
	data := queries.DataTemplate{
		Key1Link:      linkStory,
		Key2Link:      linkStatus,
		TableName:     wf.TableName,
		TimeField:     current.TimeField,
		InfoField:     current.InfoField,
		AmountField:   current.AmountField,
		TrackingTable: wf.TrackingTableName.String(),
		SummaryField:  SummaryFieldName,
		Statuses:      statuses,
		Status:        name,
		Filter:        queryValue,
		Token:         wf.CurrentIDInt.String(),
		Fields:        wf.Variable.Queries,
		PkField:       wf.FieldToken,
		PkStatus:      wf.FieldStatus,
		PkTracking:    checkPkTracking(wf),
		VideosField:   current.VideosField,
		ImagesField:   current.ImagesField,
	}
	err = tmpl.Execute(&outputBuffer, data)
	if err != nil {
		return "", err
	}
	newSQL := outputBuffer.String()
	return newSQL, nil
}

func checkPkTracking[T db.Scanner](wf *common.Workflow[T]) string {
	if wf.TrackingID.IsEmpty() {
		wf.TrackingID = types.NewString(strings.Replace(wf.TrackingTableName.String(),
			"ta_",
			"pk_",
			1,
		))
	}
	return wf.TrackingID.String()
}
