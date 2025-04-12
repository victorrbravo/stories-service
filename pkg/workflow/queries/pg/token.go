package pg

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"stories/constants"
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
	WHEN EXTRACT (epoch from (LEAD({{.TrackingTable}}.{{.TimeField}}) over (ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.Statuses}}], {{.Key2Link}})) - {{.TrackingTable}}.{{.TimeField}})) < 0 then INTERVAL '0 hour'
	ELSE LEAD({{.TrackingTable}}.{{.TimeField}}) over (ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.Statuses}}], {{.Key2Link}})) - {{.TrackingTable}}.{{.TimeField}}
END AS diff
FROM {{.TrackingTable}} tt  WHERE {{.Key1Link}}='{{.Token}}' ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.Statuses}}], {{.Key2Link}})  
)
SELECT 
CASE WHEN diff IS NULL THEN '0s' ELSE diff END AS maxtime
FROM cte WHERE {{.Key2Link}}='{{.Status}}' AND  {{.Key1Link}}='{{.Token}}'  ;
`

	SelectPositionFieldToken = `
WITH cte AS (
SELECT {{.Key2Link}} ,{{.Key1Link}}, {{.TrackingTable}}.{{.TimeField}},{{.PkField}},roles,images,
	ST_Distance(ST_Transform({{.Position}},3857), ST_Transform(LEAD({{.Position}}) over (ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.Statuses}}], {{.Key2Link}})), 3857)) AS diff
FROM {{.TrackingTable}} JOIN {{.TableName}} ON {{.Key1Link}}={{.PkField}} ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.Statuses}}], {{.Key2Link}})  
)
SELECT 
CASE WHEN diff IS NULL THEN -1 ELSE diff END AS distance
FROM cte {{if ne .Filter ""}} WHERE {{.Filter}} AND {{.PkField}} = '{{.Token}}'{{end}};
`
	SelectTimeFieldToken = `
WITH cte AS (
SELECT {{.Key2Link}} ,{{.Key1Link}},{{.TrackingTable}}.{{.TimeField}},{{.PkField}},roles,images,
CASE
	WHEN EXTRACT (epoch from (LEAD({{.TrackingTable}}.{{.TimeField}}) over (ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.Statuses}}], {{.Key2Link}})) - {{.TrackingTable}}.{{.TimeField}})) < 0 then INTERVAL '0 hour'
	ELSE LEAD({{.TrackingTable}}.{{.TimeField}}) over (ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.Statuses}}], {{.Key2Link}})) - {{.TrackingTable}}.{{.TimeField}}
END AS diff
FROM {{.TrackingTable}} JOIN {{.TableName}} ON {{.Key1Link}}={{.PkField}} ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.Statuses}}], {{.Key2Link}})  
)
SELECT 
CASE WHEN diff IS NULL THEN '0s' ELSE diff END AS maxtime
FROM cte {{if ne .Filter ""}} WHERE {{.Filter}} AND {{.PkField}} = '{{.Token}}'{{end}};
`

	SelectRoleToken = `
SELECT {{.InfoField}},{{.TrackingTable}}.{{.TimeField}},{{.AmountField}}, {{.TrackingTable}}.{{.SummaryField}}, {{.PkTracking}}, {{.VideosField}}, {{.ImagesField}}, {{.TipsIdField}}, {{.SubIdField}}, ST_AsBinary({{.Position}}), {{.Description}},{{.TrackingTable}}."tags"  FROM {{.TrackingTable}}  JOIN  {{.TableName}} ON {{.Key1Link}}={{.Token}}                
{{if ne .Filter ""}} WHERE {{.Filter}} AND {{.PkField}} = '{{.Token}}'{{end}}
`
	SelectArrayToken = `
SELECT images, videos, roles  FROM {{.TrackingTable}}  JOIN  {{.TableName}} ON {{.Key1Link}}={{.Token}}                
{{if ne .Filter ""}} WHERE {{.Filter}} AND {{.PkField}} = '{{.Token}}'{{end}}
`

	SelectTimeToken = `
SELECT {{.TrackingTable}}.{{.TimeField}} FROM {{.TrackingTable}}  JOIN  {{.TableName}} ON {{.Key1Link}}={{.Token}}                
{{if ne .Filter ""}} WHERE {{.Filter}} AND {{.PkField}} = '{{.Token}}'{{end}}
`

	SelectFieldsToken = `
SELECT {{.Fields}} FROM {{.TrackingTable}}  JOIN  {{.TableName}} ON {{.Key1Link}}={{.Token}}                
{{if ne .Filter ""}} WHERE {{.Filter}} AND {{.PkField}} = '{{.Token}}'{{end}}
`
	SelectFilterForReportToken = `
SELECT ta_tracking.pk_tracking FROM ta_tracking  
    JOIN  ta_story ON ta_tracking.fk_story=ta_story.pk_story 
                               WHERE EXISTS ( SELECT 1 FROM UNNEST(tags) AS item 
                            WHERE similarity(item, $2) > 0.3) AND ta_story.fk_flow=$1
`
)

func getLinkFields[T db.Scanner](ctx context.Context, wf *common.Workflow[T]) (string, string) {
	return wf.ExtractKeys()
}
func DistanceField[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {
	var distance scanners.Float
	newQuery := wf.Query
	ctx := context.Background()
	return func(error, ...interface{}) (*common.Workflow[T], error) {
		distanceSum := 0.0
		nTasks := len(wf.Tasks)
		n := 0
		statuses, err := wf.Statuses(ctx, false, false)
		if err != nil {
			return wf, err
		}
		for _, name := range statuses {
			current := wf.Tasks[name]
			idCurrent := types.AsInt(current.Data.ID)
			if types.AsBool(current.Report) == false || !common.ContainsIncluded(idCurrent, wf) {
				n++
				continue
			}
			newSQL, err := buildSQL(wf, name, SelectPositionFieldToken)
			if err != nil {
				log.Println("Error PositionField err", err)
				n++
				continue
			}
			_, err = (*newQuery).QueryOne(newSQL, &distance)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				log.Println("error position query err", err, "newSQL", newSQL)
				n++
				continue
			}
			if types.IsNil(distance.Value) {
				current.Data.DistanceDiff = types.NewFloat(-1)
			} else {
				current.Data.DistanceDiff = distance.Value
				if n > 0 && n < (nTasks-1) {
					distanceSum = distanceSum + types.AsFloat(distance.Value)
				}
			}
			n++
		}
		wf.Data.DistanceDiff = types.NewFloat(distanceSum)
		return wf, nil
	}
	return nil
}

func TimeField[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {
	var duration scanners.String
	var outputBuffer bytes.Buffer
	newQuery := wf.Query
	ctx := context.Background()
	return func(error, ...interface{}) (*common.Workflow[T], error) {
		timeSum := 0 * time.Second
		nTasks := len(wf.Tasks)
		n := 0
		statuses, err := wf.Statuses(ctx, false, false)
		if err != nil {
			return wf, err
		}
		for _, name := range statuses {
			current := wf.Tasks[name]
			if types.AsBool(current.Report) == false {
				n++
				continue
			}
			newSQL, err := buildSQL(wf, name, SelectTimeFieldToken)
			if err != nil {
				log.Println("Error TimeField err", err)
				n++
				continue
			}
			durationValue, err := (*newQuery).QueryOne(newSQL, &duration)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				log.Println("error query err", err, "newSQL", newSQL)
				n++
				continue
			}
			durationTime, ok := durationValue.(types.String)
			value := types.NewString(NA)
			if !errors.Is(err, sql.ErrNoRows) && ok == true {
				value = types.NewString(durationTime.String())
			}
			strValue := types.AsString(value)
			currentDuration, err := ParsePostgresInterval(strValue)
			if err == nil && (n < (nTasks - 1)) && (n > 0) {
				timeSum = timeSum + currentDuration
			}
			current.Data.UpdateAtDiff = value
			outputBuffer.Reset()
			n++
		}
		daysDuration, err := ConvertToDaysHours(timeSum.String())
		if err == nil {
			wf.Data.UpdateAtDiff = types.NewString(daysDuration)
		} else {
			wf.Data.UpdateAtDiff = types.NewString(timeSum.String())
		}
		return wf, nil
	}
	return nil
}

// Calculate sum of AmountField and calculate IncludeWrapped

func StatsField[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {

	return func(error, ...interface{}) (*common.Workflow[T], error) {
		amountSum := 0.0
		amountSumParcial := 0.0
		ctx := context.Background()
		var previousTask *common.Task
		var current *common.Task

		loc, err := time.LoadLocation(constants.TimezoneBogota)
		if err != nil {
			fmt.Errorf("invalid timezone: %v", err)
			return wf, nil
		}
		statuses, err := wf.Statuses(ctx, false, false)
		if err != nil {
			return wf, err
		}
		log.Println("StatsField statuses", statuses)
		for _, name := range statuses {
			current = wf.Tasks[name]
			amountField, err := ParseDollarAmount(current.Data.AmountField.String())
			log.Println("StatsField* name", name, "amount", amountField)
			if err == nil {
				amountSum = amountSum + amountField
			}
			if previousTask != nil && CheckPeriodDates(previousTask.Data.TimeField, current.Data.TimeField, constants.TimezoneBogota) {
				firstDate, _ := time.Parse(time.RFC3339, types.AsString(previousTask.Data.TimeField))
				firstDate = firstDate.In(loc)
				wrapped := new(common.Wrapped)
				wrapped.Label = firstDate.Format("Jan 02")
				log.Println("StatsField* name", name, "*wrapped.amount", amountSumParcial)
				wrapped.Amount = types.NewString(FormatToMoney(amountSumParcial))
				previousTask.Transition.Wrapped = wrapped
				amountSumParcial = amountField
			} else {
				amountSumParcial = amountSumParcial + amountField
			}
			previousTask = current
		}
		wf.Data.AmountField = types.NewString(FormatToMoney(amountSum))
		return wf, nil
	}
}

func FilterField[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {
	var number scanners.Int

	return func(error, ...interface{}) (*common.Workflow[T], error) {
		log.Println("Executing FilterField...start")
		wfId := wf.CurrentIDInt
		newQuery := wf.Query
		if wf.Filter.Sentence == "" {
			log.Println("FilterField Sentence empty filter empty")
			wf.Filter.IncludedWrapped = []int{}
			return wf, nil
		}
		rows, err := (*newQuery).QueryMany(SelectFilterForReportToken, &number, wfId, wf.Filter.Sentence)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				log.Println("FilterField Error Count SQL err", err)
				return wf, err
			}
		}
		log.Println("FilterField rows rows", rows)
		for _, row := range rows {
			valueInt := types.AsInt(row.(types.Int))
			wf.Filter.IncludedWrapped = append(wf.Filter.IncludedWrapped, valueInt)
		}
		log.Println("FilterField rows IncludedWrapped", wf.Filter.IncludedWrapped)
		return wf, nil
	}
}

func RoleField[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {
	var role scanners.TupleEleven
	newQuery := wf.Query
	ctx := context.Background()
	return func(error, ...interface{}) (*common.Workflow[T], error) {
		amountSum := 0.0
		amountSumParcial := 0.0
		n := 0
		statuses, err := wf.Statuses(ctx, false, false)
		if err != nil {
			return wf, err
		}
		var previousTask *common.Task
		var current *common.Task
		loc, err := time.LoadLocation(constants.TimezoneBogota)
		if err != nil {
			fmt.Errorf("invalid timezone: %v", err)
			return wf, nil
		}
		timeStart := time.Now()
		for _, name := range statuses {
			current = wf.Tasks[name]
			newSQL, err := buildSQL(wf, name, SelectRoleToken)
			if err != nil {
				log.Println("RoleField Error Role err", err)
				n++
				continue
			}
			_, err = (*newQuery).QueryOne(newSQL, &role)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					log.Println("Error Count SQL err", err)
				}
				n++
				continue
			}
			current.Data.InfoField = role.First
			current.Data.TimeField = role.Second

			current.Data.AmountField = role.Third
			current.Data.Summary = role.Four
			current.Data.ID = role.ID
			current.Data.Images = role.Fifth
			current.Data.Videos = role.Sixth
			current.Data.TipsIdField = role.Seventh
			current.Data.SubIdField = role.Eight
			current.Data.Position = role.Ninth
			current.Data.Description = role.Tenth
			current.Data.Tags = role.Eleventh
			if previousTask != nil && CheckPeriodDates(previousTask.Data.TimeField, current.Data.TimeField, constants.TimezoneBogota) {
				firstDate, _ := time.Parse(time.RFC3339, types.AsString(previousTask.Data.TimeField))
				firstDate = firstDate.In(loc)
				wrapped := new(common.Wrapped)
				wrapped.Label = firstDate.Format("Jan 02")
				wrapped.Amount = types.NewString(FormatToMoney(amountSumParcial))
				previousTask.Transition.Wrapped = wrapped
				amountSumParcial = 0.0
			}
			amountField, err := ParseDollarAmount(role.Third.String())
			if err == nil {
				amountSum = amountSum + amountField
				amountSumParcial = amountSumParcial + amountField
			} else {
				fmt.Println("_error role.Third.String()", role.Third.String(), "amountField", amountField, "err", err)
			}
			n++
			previousTask = current
		}
		if previousTask != nil {
			firstDate, _ := time.Parse(time.RFC3339, types.AsString(previousTask.Data.TimeField))
			firstDate = firstDate.In(loc)
			wrapped := new(common.Wrapped)
			wrapped.Label = firstDate.Format("Jan 02")
			wrapped.Amount = types.NewString(FormatToMoney(amountSumParcial))
			previousTask.Transition.Wrapped = wrapped
		}
		elapsed := time.Since(timeStart)
		log.Println("RoleField elapsed:", elapsed)
		wf.Data.AmountField = types.NewString(FormatToMoney(amountSum))
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
