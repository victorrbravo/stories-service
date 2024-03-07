package pg

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	db "stories/database"
	"stories/generated/models"
	"stories/pkg/workflow/common"
	"stories/scanners"
	"stories/types"
	"time"
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

	SelectTimeFieldToken = `
WITH cte AS (
SELECT {{.Key2Link}} ,{{.Key1Link}} , {{.TimeField}},{{.PkField}},roles,images,
CASE
	WHEN EXTRACT (epoch from (LEAD({{.TimeField}}) over (ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.Statuses}}], {{.Key2Link}})) - {{.TimeField}})) < 0 then INTERVAL '0 hour'
	ELSE LEAD({{.TimeField}}) over (ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.Statuses}}], {{.Key2Link}})) - {{.TimeField}}
END AS diff
FROM {{.TrackingTable}} JOIN {{.TableName}} ON {{.Key1Link}}={{.PkField}} ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.Statuses}}], {{.Key2Link}})  
)
SELECT 
CASE WHEN diff IS NULL THEN '0s' ELSE diff END AS maxtime
FROM cte {{if ne .Filter ""}} WHERE {{.Filter}} AND {{.PkField}} = '{{.Token}}'{{end}};
`

	SelectRoleToken = `
SELECT {{.InfoField}},{{.TimeField}},{{.AmountField}}  FROM {{.TrackingTable}}  JOIN  {{.TableName}} ON {{.Key1Link}}={{.Token}}                
{{if ne .Filter ""}} WHERE {{.Filter}} AND {{.PkField}} = '{{.Token}}'{{end}}
`
	SelectArrayToken = `
SELECT images, videos, roles  FROM {{.TrackingTable}}  JOIN  {{.TableName}} ON {{.Key1Link}}={{.Token}}                
{{if ne .Filter ""}} WHERE {{.Filter}} AND {{.PkField}} = '{{.Token}}'{{end}}
`

	SelectTimeToken = `
SELECT {{.TimeField}} FROM {{.TrackingTable}}  JOIN  {{.TableName}} ON {{.Key1Link}}={{.Token}}                
{{if ne .Filter ""}} WHERE {{.Filter}} AND {{.PkField}} = '{{.Token}}'{{end}}
`

	SelectFieldsToken = `
SELECT {{.Fields}} FROM {{.TrackingTable}}  JOIN  {{.TableName}} ON {{.Key1Link}}={{.Token}}                
{{if ne .Filter ""}} WHERE {{.Filter}} AND {{.PkField}} = '{{.Token}}'{{end}}
`
)

func getLinkFields[T db.Scanner](ctx context.Context, wf *common.Workflow[T]) (string, string) {
	return wf.ExtractKeys()
}

func TimeField[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {
	var duration scanners.String
	var outputBuffer bytes.Buffer
	newQuery := wf.Query
	return func(error, ...interface{}) (*common.Workflow[T], error) {
		for name, current := range wf.Tasks {
			if types.AsBool(current.Report) == false {
				continue
			}
			newSQL, err := buildSQL(wf, name, SelectTimeFieldToken)
			if err != nil {
				log.Println("Error TimeField err", err)
				continue
			}
			durationValue, err := (*newQuery).QueryOne(newSQL, &duration)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				log.Println("error query err", err, "newSQL", newSQL)
				continue
			}
			durationTime, ok := durationValue.(types.String)
			value := types.NewString(NA)
			if !errors.Is(err, sql.ErrNoRows) && ok == true {
				value = types.NewString(durationTime.String())
			}
			strValue := types.AsString(value)
			currentDuration, err := time.ParseDuration(strValue)
			if err == nil {
				log.Println("time diff |", fmt.Sprintf("%s", currentDuration))
			}
			current.Data.UpdateAtDiff = value
			outputBuffer.Reset()
		}
		return wf, nil
	}
	return nil
}

func RoleField[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {
	var role scanners.TupleThree
	newQuery := wf.Query
	return func(error, ...interface{}) (*common.Workflow[T], error) {
		for name, current := range wf.Tasks {
			if types.AsBool(current.Report) == false {
				continue
			}
			newSQL, err := buildSQL(wf, name, SelectRoleToken)
			if err != nil {
				log.Println("Error Role err", err)
				continue
			}
			_, err = (*newQuery).QueryOne(newSQL, &role)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					log.Println("Error Count SQL err", err)
				}
				continue
			}
			current.Data.InfoField = role.First
			current.Data.TimeField = role.Second
			current.Data.AmountField = role.Third
		}
		return wf, nil
	}
}

func ArrayField[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {
	var role scanners.TupleArrays
	newQuery := wf.Query
	return func(error, ...interface{}) (*common.Workflow[T], error) {
		for name, current := range wf.Tasks {
			if types.AsBool(current.Report) == false {
				continue
			}
			newSQL, err := buildSQL(wf, name, SelectArrayToken)
			if err != nil {
				log.Println("Error Role err", err)
				continue
			}
			_, err = (*newQuery).QueryOne(newSQL, &role)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					log.Println("Error Count SQL err", err)
				}
				continue
			}
			current.Data.Images = role.First
			current.Data.Videos = role.Second
			current.Data.Roles = role.Third
		}
		return wf, nil
	}
}

func VariableField[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {
	newQuery := wf.Query
	if len(parameters) == 0 {
		return func(error, ...interface{}) (*common.Workflow[T], error) {
			return wf, errors.New("error not specified 'task name' in parameters")
		}
	}
	nameTask := ""
	names := parameters[0].([]interface{})
	nameTask = names[0].(string)
	return func(error, ...interface{}) (*common.Workflow[T], error) {
		fields := models.TaTracking{}.ListOfFields()
		wf.Variable.Queries = fields
		newSQL, err := buildSQL(wf, nameTask, SelectFieldsToken)
		if err != nil {
			log.Println("Error Role err", err)
			return wf, err
		}

		var rows interface{}
		var scanner models.TaTracking
		rows, err = (*newQuery).QueryOne(newSQL, &scanner)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				log.Println("Error Count SQL err", err)
			}
			return wf, err
		}
		wf.Data.Rows = make([]interface{}, 1)
		wf.Data.Rows[0] = rows
		return wf, nil
	}
}

func CountToken[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {
	var count scanners.Int
	newQuery := wf.Query
	return func(error, ...interface{}) (*common.Workflow[T], error) {
		for name, current := range wf.Tasks {
			newSQL, err := buildSQL(wf, name, SelectCountToken)
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
