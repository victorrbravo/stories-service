package pg

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sort"
	"stories/constants"
	db "stories/database"
	"stories/date_parser"
	"stories/generated/models"
	"stories/pkg/workflow"
	"stories/pkg/workflow/common"
	"stories/scanners"
	"stories/types"
	"strings"
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
	ST_Distance(ST_Transform({{.Position}},3857), ST_Transform(LEAD({{.Position}}) over (ORDER BY {{.Key1Link}}, ARRAY_POSITION($1::text[], {{.Key2Link}})), 3857)) AS diff
FROM {{.TrackingTable}} JOIN {{.TableName}} ON {{.Key1Link}}={{.PkField}} WHERE {{.TrackingTable}}.{{.Key2Link}}=ANY($1::text[]) ORDER BY {{.Key1Link}}, ARRAY_POSITION($1::text[], {{.Key2Link}})  
)
SELECT 
CASE WHEN diff IS NULL THEN -1 ELSE diff END AS distance
FROM cte {{if ne .Filter ""}} WHERE {{.Filter}} AND {{.PkField}} = '{{.Token}}'{{end}};
`

	SelectTimeFieldToken = `
WITH cte AS (
SELECT {{.Key2Link}} ,{{.Key1Link}},{{.TrackingTable}}.{{.TimeField}},{{.PkField}},roles,images,
CASE
	WHEN EXTRACT (epoch from (LEAD({{.TrackingTable}}.{{.TimeField}}) over (ORDER BY {{.Key1Link}}, ARRAY_POSITION($1::text[], {{.Key2Link}})) - {{.TrackingTable}}.{{.TimeField}})) < 0 then INTERVAL '0 hour'
	ELSE LEAD({{.TrackingTable}}.{{.TimeField}}) over (ORDER BY {{.Key1Link}}, ARRAY_POSITION($1::text[], {{.Key2Link}})) - {{.TrackingTable}}.{{.TimeField}}
END AS diff
FROM {{.TrackingTable}} JOIN {{.TableName}} ON {{.Key1Link}}={{.PkField}} WHERE {{.TrackingTable}}.{{.Key2Link}}=ANY($1::text[]) ORDER BY {{.Key1Link}}, ARRAY_POSITION($1::text[], {{.Key2Link}})  
)
SELECT 
CASE WHEN diff IS NULL THEN '0s' ELSE diff END AS maxtime
FROM cte {{if ne .Filter ""}} WHERE {{.Filter}} AND {{.PkField}} = '{{.Token}}'{{end}};
`

	SelectSumTimeFieldToken = `
WITH cte AS (
SELECT fk_status ,fk_story,ta_tracking.updated_at,pk_story,roles,images,
CASE
	WHEN EXTRACT (epoch from (LEAD(ta_tracking.updated_at) over (ORDER BY fk_story, ARRAY_POSITION($1::text[], fk_status)) - ta_tracking.updated_at)) < 0 then INTERVAL '0 hour'
	ELSE LEAD(ta_tracking.updated_at) over (ORDER BY fk_story, ARRAY_POSITION($1::text[], fk_status)) - 
ta_tracking.updated_at
END AS diff
FROM ta_tracking JOIN ta_story ON fk_story=pk_story WHERE ta_tracking.fk_status=ANY($1::text[]) ORDER BY ARRAY_POSITION($1::text[], fk_status)  
)
SELECT 
SUM(CASE WHEN diff IS NULL THEN '0s' ELSE diff END) AS sumtime
FROM cte  WHERE pk_story = $2
`

	SelectSumPosFieldToken = `
WITH cte AS (
SELECT {{.Key2Link}} ,{{.Key1Link}}, {{.TrackingTable}}.{{.TimeField}},{{.PkField}},roles,images,
	ST_Distance(ST_Transform({{.Position}},3857), ST_Transform(LEAD({{.Position}}) over (ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.Statuses}}], {{.Key2Link}})), 3857)) AS diff
FROM {{.TrackingTable}} JOIN {{.TableName}} ON {{.Key1Link}}={{.PkField}} ORDER BY {{.Key1Link}}, ARRAY_POSITION(ARRAY[{{.Statuses}}], {{.Key2Link}})  
)
SELECT 
SUM(CASE WHEN diff IS NULL THEN -1 ELSE diff END) AS sumdistance
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

	SelectFilterUsingOptions = `
SELECT ta_tracking.fk_status FROM ta_tracking JOIN  ta_story ON ta_tracking.fk_story=ta_story.pk_story
WHERE ( (($2::text[] IS NULL) AND ($6::text[] IS NULL)) OR (EXISTS ( SELECT 1 FROM unnest(tags) AS tag
    JOIN unnest($2::text[]) AS input_word ON similarity(tag, input_word) > 0.3
    ) OR (($6::text[] IS NOT NULL) AND (lower_text_array(tags::text[]) && lower_text_array($6::text[]))))) AND ($3::text[] IS NULL OR  (NOT (lower_text_array(tags::text[]) && lower_text_array($3::text[]))))  AND ta_story.pk_story=$1
	AND  ($4::text IS NULL OR ta_tracking.updated_at < $4::timestamptz) AND ($5::text IS NULL OR ta_tracking.updated_at > $5::timestamptz)
AND ($7::text IS NULL OR ta_tracking.pk_tracking IN (
    SELECT pk_tracking FROM ta_tracking
    WHERE ta_tracking.fk_story = $1
    ORDER BY ta_tracking.amount_spent DESC NULLS LAST
    LIMIT $7::int
))
ORDER BY ta_tracking.updated_at;
`
	// Same as SelectFilterUsingOptions but, instead of ordering by
	// updated_at, it orders the rows by ta_tracking.fk_status following the
	// explicit order of the status list passed in $8 (i.e. wf.GetStatuses()).
	// array_position returns the 1-based index of each fk_status inside the
	// list, so the result preserves that ordering. Statuses missing from the
	// list (NULL position) are pushed to the end.
	SelectFilterUsingOptionsOrdered = `
SELECT ta_tracking.fk_status FROM ta_tracking JOIN  ta_story ON ta_tracking.fk_story=ta_story.pk_story
WHERE ( (($2::text[] IS NULL) AND ($6::text[] IS NULL)) OR (EXISTS ( SELECT 1 FROM unnest(tags) AS tag
    JOIN unnest($2::text[]) AS input_word ON similarity(tag, input_word) > 0.3
    ) OR (($6::text[] IS NOT NULL) AND (lower_text_array(tags::text[]) && lower_text_array($6::text[]))))) AND ($3::text[] IS NULL OR  (NOT (lower_text_array(tags::text[]) && lower_text_array($3::text[]))))  AND ta_story.pk_story=$1
	AND  ($4::text IS NULL OR ta_tracking.updated_at < $4::timestamptz) AND ($5::text IS NULL OR ta_tracking.updated_at > $5::timestamptz)
AND ($7::text IS NULL OR ta_tracking.pk_tracking IN (
    SELECT pk_tracking FROM ta_tracking
    WHERE ta_tracking.fk_story = $1
    ORDER BY ta_tracking.amount_spent DESC NULLS LAST
    LIMIT $7::int
))
ORDER BY array_position($8::text[], ta_tracking.fk_status) NULLS LAST;
`
	SumAmountFieldTaTracking = `SELECT SUM(ta_tracking.amount_spent)::money::numeric FROM ta_tracking WHERE ta_tracking.fk_status=ANY($1::text[])
AND ta_tracking.fk_story=$2
`
	SumAmountFieldForWeekTaTracking = `SELECT SUM(ta_tracking.amount_spent)::money::numeric FROM ta_tracking WHERE 
ta_tracking.fk_story=$2 AND ta_tracking.fk_status=ANY($3::text[])   AND updated_at >= date_trunc('week', $1::date) 
  AND updated_at < date_trunc('week', $1::date) + INTERVAL '7 days';
`
	SumAmountFieldForDayTaTracking = `SELECT SUM(ta_tracking.amount_spent)::money::numeric FROM ta_tracking WHERE 
ta_tracking.fk_story=$2 AND ta_tracking.fk_status=ANY($3::text[]) AND updated_at >= $1::timestamp
  AND updated_at < $1::timestamp + INTERVAL '1 day';
`

	CountFieldTaTracking = `SELECT COUNT(*) FROM ta_tracking WHERE ta_tracking.fk_status=ANY($1::text[])
AND ta_tracking.fk_story=$2
`
)

func DistanceField[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {
	var distance scanners.Float
	newQuery := wf.Query
	ctx := context.Background()
	return func(error, ...interface{}) (*common.Workflow[T], error) {
		distanceSum := 0.0
		nTasks := len(wf.Tasks)
		start1 := time.Now()
		n := 0
		statuses, err := wf.Statuses(ctx, false, false)
		if len(wf.Filter.IncludedWrapped) > 0 {
			statuses = wf.Filter.IncludedWrapped
		}
		if err != nil {
			return wf, err
		}
		pagStatuses, before, after := Paginate(statuses, wf.Page)
		log.Println("before", before)
		log.Println("after", after)
		for _, name := range pagStatuses {
			current := wf.Tasks[name]
			paginatedStatuses, _, _ := Paginate(statuses, wf.Page)
			if types.AsBool(current.Report) == false {
				n++
				continue
			}
			newSQL, err := buildSQL(wf, name, SelectPositionFieldToken)
			if err != nil {
				log.Println("Error PositionField err", err)
				n++
				continue
			}
			_, err = (*newQuery).QueryOne(newSQL, &distance, types.NewStringArray(paginatedStatuses))
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				log.Println("error position query err", err, "newSQL", newSQL)
				n++
				continue
			}
			if types.IsNil(distance.Value) {
				current.Data.DistanceDiff = types.NewFloat(-1)
			} else {
				current.Data.DistanceDiff = distance.Value
				if n >= 0 && n < (nTasks-1) {
					distanceSum = distanceSum + types.AsFloat(distance.Value)
				}
			}
			n++
		}
		elapsed1 := time.Since(start1)
		log.Printf("distance spent %s", elapsed1)
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
		if len(wf.Filter.IncludedWrapped) > 0 {
			statuses = wf.Filter.IncludedWrapped
		}
		if err != nil {
			return wf, err
		}
		pagStatuses, before, after := Paginate(statuses, wf.Page)
		log.Println("before", before)
		log.Println("after", after)
		for _, name := range pagStatuses {
			current := wf.Tasks[name]
			if types.AsBool(current.Report) == false {
				n++
				continue
			}
			paginatedStatuses, _, _ := Paginate(statuses, wf.Page)
			newSQL, err := buildSQL(wf, name, SelectTimeFieldToken)
			if err != nil {
				log.Println("Error TimeField err", err)
				n++
				continue
			}
			durationValue, err := (*newQuery).QueryOne(newSQL, &duration, types.NewStringArray(paginatedStatuses))
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
		amountSumParcial := 0.0
		ctx := context.Background()
		var previousTask *common.Task
		var current *common.Task

		loc, err := time.LoadLocation(constants.TimezoneBogota)
		if err != nil {
			fmt.Errorf("invalid timezone: %v", err)
			return wf, nil
		}
		start1 := time.Now()
		statuses, err := wf.Statuses(ctx, false, false)
		allStatuses, err := wf.AllStatuses(ctx, false)
		if len(wf.Filter.IncludedWrapped) > 0 {
			statuses = wf.Filter.IncludedWrapped
		}
		if err != nil {
			return wf, err
		}
		pagStatuses, before, after := Paginate(statuses, wf.Page)
		log.Println("before", before)
		log.Println("after", after)
		for _, name := range pagStatuses {
			current = wf.Tasks[name]
			amountField, err := ParseDollarAmount(current.Data.AmountField.String())
			if err != nil {
				log.Println("task name:", name, "Error AmountField err", err)
				continue
			}
			if previousTask != nil && CheckPeriodDates(previousTask.Data.TimeField, current.Data.TimeField, constants.TimezoneBogota, wf.ByDay) {
				firstDate, _ := time.Parse(time.RFC3339, types.AsString(previousTask.Data.TimeField))
				firstDate = firstDate.In(loc)
				wrapped := new(common.Wrapped)
				wrapped.Label = firstDate.Format("Mon, Jan 02")
				log.Println("StatsField* name", name, "*wrapped.amount", amountSumParcial)
				calculateWrappedAmount(wf, firstDate, wrapped, statuses)
				previousTask.Transition.Wrapped = wrapped
				amountSumParcial = amountField
			} else {
				amountSumParcial = amountSumParcial + amountField
			}
			previousTask = current
		}
		elapsed1 := time.Since(start1)
		log.Printf("StatsField elapsed %s", elapsed1)
		newQuery := wf.Query
		var value scanners.Float
		var valueAll scanners.Float

		_, err = (*newQuery).QueryOne(SumAmountFieldTaTracking, &value, types.NewStringArray(statuses), wf.CurrentIDInt)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Println("Error Filter StatsField err", err)
			return wf, err
		}
		_, err = (*newQuery).QueryOne(SumAmountFieldTaTracking, &valueAll, types.NewStringArray(allStatuses), wf.CurrentIDInt)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Println("Error All Statuses StatsField err", err)
			return wf, err
		}
		var valueFloat, valuePercFloat float64
		valueFloat = types.AsFloat(value.Value)
		if types.AsFloat(valueAll.Value) > 0 {
			valuePercFloat = 100 * types.AsFloat(value.Value) / types.AsFloat(valueAll.Value)
		}
		wf.Data.AmountField = types.NewString(FormatToMoney(valueFloat, true))
		wf.Data.AmountFieldPerc = types.NewString(FormatToMoney(valuePercFloat, false))

		var duration scanners.String
		durationValue, err := (*newQuery).QueryOne(SelectSumTimeFieldToken, &duration, types.NewStringArray(statuses), wf.CurrentIDInt)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Println("Error StatsField err", err)
			return wf, err
		}
		durationTime, ok := durationValue.(types.String)
		if !ok {
			log.Println("Error StatsField err", err)
		}
		valueTime := types.NewString(durationTime.String())

		strValue := types.AsString(valueTime)
		currentDuration, err := ParsePostgresInterval(strValue)
		daysDuration, err := ConvertToDaysHours(currentDuration.String())
		if err != nil {
			log.Println("Error StatsField converting hours err", err)
		}
		log.Println("StatsField duration", daysDuration)
		wf.Data.UpdateAtDiff = types.NewString(daysDuration)
		wf.Data.Count = types.NewInt(len(statuses))
		wf.Data.Total = types.NewInt(len(statuses))
		return wf, nil
	}
}

func FilterField[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {
	var statusScanned scanners.String
	return func(error, ...interface{}) (*common.Workflow[T], error) {
		log.Println("Executing FilterField...start")
		wfId := wf.CurrentIDInt
		newQuery := wf.Query
		start1 := time.Now()
		if wf.Filter.Sentence == "" {
			log.Println("FilterField Sentence empty filter empty")
			wf.Filter.IncludedWrapped = []string{}
			return wf, nil
		}
		p, err := date_parser.ExtractTerms(wf.Filter.Sentence)
		if err != nil {
			log.Println("FilterField err", err)
			return wf, nil
		}

		andsWords := types.NewStringArray(p.AndTerms)
		andsLiteralWords := types.NewStringArray(p.AndLiteralTerms)
		notWords := types.NewStringArray(p.AndNotTerms)
		stringAndWords := andsWords.String()
		if strings.Compare(stringAndWords, "{}") == 0 {
			andsWords = types.NewStringArray([]string{wf.Filter.Sentence})

		}
		if len(andsWords.Array) == 0 {
			andsWords = types.NilStringArray
		}
		if len(notWords.Array) == 0 {
			notWords = types.NilStringArray
		}
		if len(andsLiteralWords.Array) == 0 {
			andsLiteralWords = types.NilStringArray
		}
		lessThan := types.NilString()
		greaterThan := types.NilString()

		if p.LessThan != "" {
			l, err := date_parser.GenerateSQLWhereDate(p.LessThan, "<")
			if err != nil {
				log.Println("LessFilterField err", err)
				return wf, nil
			}
			lessThan = types.NewString(l)
		}
		if p.GreaterThan != "" {
			g, err := date_parser.GenerateSQLWhereDate(p.GreaterThan, ">")
			if err != nil {
				log.Println("Greater FilterField err", err)
				return wf, nil
			}
			greaterThan = types.NewString(g)
		}

		// Pull a function-style term out of the parsed sentence, e.g.
		// "max- amount,2" means: keep only the 2 rows with the largest
		// ta_tracking.amount_spent for this story. The first param picks
		// the field (currently only "amount" -> amount_spent is wired up)
		// and the second param is the N. When no such function is
		// present we send NULL and the SQL clause becomes a no-op.
		maxAmountLimit := types.NilString()
		for _, fn := range p.Functions {
			if fn.Name != "max" || len(fn.Params) < 2 {
				continue
			}
			if fn.Params[0] != "amount" {
				log.Println("FilterField max function: unsupported field", fn.Params[0])
				continue
			}
			maxAmountLimit = types.NewString(fn.Params[1])
			log.Println("FilterField applying max- amount limit:", fn.Params[1])
			break
		}

		var rows []interface{}
		if len(wf.OrderedByTags) > 0 {
			orderedStatuses := types.NewStringArray(wf.GetStatuses())
			rows, err = (*newQuery).QueryMany(SelectFilterUsingOptionsOrdered, &statusScanned, wfId, andsWords, notWords, lessThan, greaterThan, andsLiteralWords, maxAmountLimit, orderedStatuses)
		} else {
			rows, err = (*newQuery).QueryMany(SelectFilterUsingOptions, &statusScanned, wfId, andsWords, notWords, lessThan, greaterThan, andsLiteralWords, maxAmountLimit)
		}
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				log.Println("FilterField Error Count SQL err", err)
				return wf, err
			}
		}
		log.Println("FilterField rows rows", rows)
		for _, row := range rows {
			valueStr := row.(types.String)
			wf.Filter.IncludedWrapped = append(wf.Filter.IncludedWrapped, valueStr.String())
		}
		elapsed := time.Since(start1)
		log.Printf("FilterField elapsed %s", elapsed)
		log.Println("FilterField rows IncludedWrapped", wf.Filter.IncludedWrapped)
		return wf, nil
	}
}

func OrderingTaskField[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {
	return func(error, ...interface{}) (*common.Workflow[T], error) {
		ctx := context.Background()
		timeStart := time.Now()
		statuses, err := wf.Statuses(ctx, false, false)
		if err != nil {
			return wf, err
		}
		var newTasks map[string]*common.Task = make(map[string]*common.Task)
		for _, status := range statuses {
			newTasks[status] = wf.Tasks[status]
		}
		wf.Tasks = newTasks
		elapsed := time.Since(timeStart)
		log.Printf("OrderingField elapsed: %v", elapsed)
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
		var diffStatuses []string
		n := 0
		statuses, err := wf.Statuses(ctx, false, false)
		if err != nil {
			return wf, err
		}
		if len(wf.Filter.IncludedWrapped) > 0 {
			diffStatuses = workflow.SliceDiff(statuses, wf.Filter.IncludedWrapped)
			statuses = wf.Filter.IncludedWrapped
		}
		log.Println("diffStatuses: ", diffStatuses)
		var previousTask *common.Task
		loc, err := time.LoadLocation(constants.TimezoneBogota)
		if err != nil {
			fmt.Errorf("invalid timezone: %v", err)
			return wf, nil
		}
		timeStart := time.Now()

		pagStatuses, before, after := Paginate(statuses, wf.Page)
		log.Println("before", before)
		log.Println("after", after)
		for _, name := range pagStatuses {
			current, ok := wf.Tasks[name]
			if !ok {
				log.Println("Statuses in paginated value not found name:", name)
				continue
			}
			newSQL, errBuild := buildSQL(wf, name, SelectRoleToken)
			if errBuild != nil {
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
			if previousTask != nil && CheckPeriodDates(previousTask.Data.TimeField, current.Data.TimeField, constants.TimezoneBogota, wf.ByDay) {
				firstDate, _ := time.Parse(time.RFC3339, types.AsString(previousTask.Data.TimeField))
				firstDate = firstDate.In(loc)
				wrapped := new(common.Wrapped)
				wrapped.Label = firstDate.Format("Mon, Jan 02")
				calculateWrappedAmount(wf, firstDate, wrapped, statuses)
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
		var current *common.Task
		if after != "" {
			current = wf.Tasks[after]
		}
		if after == "" || (previousTask != nil && current != nil && CheckPeriodDates(previousTask.Data.TimeField, current.Data.TimeField, constants.TimezoneBogota, wf.ByDay)) {
			firstDate, _ := time.Parse(time.RFC3339, types.AsString(previousTask.Data.TimeField))
			firstDate = firstDate.In(loc)
			wrapped := new(common.Wrapped)
			wrapped.Label = firstDate.Format("Mon, Jan 02")
			calculateWrappedAmount(wf, firstDate, wrapped, statuses)
			previousTask.Transition.Wrapped = wrapped
			amountSumParcial = 0.0
		}
		elapsed := time.Since(timeStart)
		log.Printf("RoleField elapsed: %v", elapsed)
		wf.Data.AmountField = types.NewString(FormatToMoney(amountSum, true))

		return wf, nil
	}
}

func calculateWrappedAmount[T db.Scanner](wf *common.Workflow[T], firstDate time.Time, wrapped *common.Wrapped, statuses []string) {
	var valueAmount scanners.Float
	newQuery := wf.Query
	var currentSQL = SumAmountFieldForWeekTaTracking
	if wf.ByDay {
		currentSQL = SumAmountFieldForDayTaTracking
	}

	toCalculateDateString := truncateBogotaMidnightAsUTC(firstDate)
	log.Println("CalculateDate", toCalculateDateString)
	_, err := (*newQuery).QueryOne(currentSQL, &valueAmount, toCalculateDateString,
		wf.CurrentIDInt, types.NewStringArray(statuses))
	if err != nil {
		log.Println("SumAmountField Error SumAmountFieldForWeekTaTracking", err)
	}
	wrapped.Amount = types.NewString(FormatToMoney(types.AsFloat(valueAmount.Value), true))
}

func truncateBogotaMidnightAsUTC(t time.Time) string {
	// Load Bogotá location
	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		panic(err)
	}

	// Convert to Bogotá time and truncate to midnight
	tInBogota := t.In(loc)
	tMidnightBogota := time.Date(tInBogota.Year(), tInBogota.Month(), tInBogota.Day(), 0, 0, 0, 0, loc)

	// Convert Bogotá midnight to UTC
	tMidnightUTC := tMidnightBogota.UTC()

	// Format in desired layout
	return tMidnightUTC.Format("2006-01-02 15:04:05")
}

func Paginate[T any](items []T, page types.Page) ([]T, T, T) {
	start := 0
	end := len(items)
	var before T
	var after T

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
	if start > 0 {
		before = items[start-1]
	}
	if end < len(items)-1 {
		after = items[end+1]
	}
	return items[start:end], before, after
}

func ArrayField[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {
	var role scanners.TupleArrays
	newQuery := wf.Query
	ctx := context.Background()
	return func(error, ...interface{}) (*common.Workflow[T], error) {
		start1 := time.Now()
		statuses, err := wf.Statuses(ctx, false, false)
		if len(wf.Filter.IncludedWrapped) > 0 {
			statuses = wf.Filter.IncludedWrapped
		}
		if err != nil {
			return wf, err
		}
		pagStatuses, before, after := Paginate(statuses, wf.Page)
		log.Println("before", before)
		log.Println("after", after)
		for _, name := range pagStatuses {
			current := wf.Tasks[name]
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
		elapsed := time.Since(start1)
		log.Printf("ArrayField elapsed: %v", elapsed)
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
	ctx := context.Background()

	return func(error, ...interface{}) (*common.Workflow[T], error) {
		start1 := time.Now()
		statuses, err := wf.Statuses(ctx, false, false)
		if len(wf.Filter.IncludedWrapped) > 0 {
			statuses = wf.Filter.IncludedWrapped
		}
		if err != nil {
			return wf, err
		}
		pagStatuses, before, after := Paginate(statuses, wf.Page)
		log.Println("before", before)
		log.Println("after", after)

		for _, name := range pagStatuses {
			current := wf.Tasks[name]
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
			wf.Data.Stats.Assign(count.Value)
		}
		newQuery := wf.Query
		var value scanners.Int
		_, err = (*newQuery).QueryOne(CountFieldTaTracking, &value, types.NewStringArray(statuses), wf.CurrentIDInt)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Println("Count Error StatsField err", err)
			return wf, err
		}
		wf.Data.Total = value.Value
		wf.Data.Stats.Assign(value.Value)
		elapsed := time.Since(start1)
		log.Println("Count Item elapsed:", elapsed)
		return wf, nil
	}
}

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
// The same map-rebuild pattern used by OrderingTaskField is applied
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

// TasksOrderField is a workflow Monad (same shape as FilterField) that
// reorders wf.Tasks based on a "ctg-" function found in
// wf.Filter.Sentence.
//
// The sentence is parsed via date_parser.ExtractTerms (same call used
// by FilterField above). For every parsed Function whose Name is "ctg"
// the Params slice is taken as the wantedTags priority list — e.g. the
// sentence "ctg- bill,coffee" produces wantedTags ==
// []string{"bill", "coffee"}.
//
// When at least one ctg- function is found, the per-status tags +
// updated_at map is fetched with TagsByStatus(*wf.Query, wf.CurrentIDInt)
// and OrderTasksByTags is invoked. Absence of a ctg- function (or an
// empty Filter.Sentence) makes this a no-op.
func TasksOrderField[T db.Scanner](wf *common.Workflow[T], parameters ...interface{}) Monad[T] {
	return func(error, ...interface{}) (*common.Workflow[T], error) {
		log.Println("Executing TasksOrderField...start")
		start1 := time.Now()

		if wf.Filter.Sentence == "" {
			log.Println("TasksOrderField empty Filter.Sentence, skipping")
			return wf, nil
		}

		p, err := date_parser.ExtractTerms(wf.Filter.Sentence)
		if err != nil {
			log.Println("TasksOrderField parse err", err)
			return wf, nil
		}

		// Pull a function-style term out of the parsed sentence, e.g.
		// "ctg- bill,coffee" means: reorder tasks so that those whose
		// tags contain "bill" come first (sorted by updated_at
		// ascending), then those tagged "coffee", then the rest. When
		// no ctg- function is present this is a no-op.
		var wantedTags []string
		for _, fn := range p.Functions {
			if fn.Name != "ctg" || len(fn.Params) == 0 {
				continue
			}
			wantedTags = fn.Params
			log.Println("TasksOrderField applying ctg- tags:", wantedTags)
			break
		}
		if len(wantedTags) == 0 {
			log.Println("TasksOrderField no ctg- function in sentence, skipping")
			return wf, nil
		}

		if wf.Query == nil {
			log.Println("TasksOrderField wf.Query is nil, skipping")
			return wf, nil
		}
		tagsMap, err := TagsByStatus(*wf.Query, wf.CurrentIDInt)
		if err != nil {
			log.Println("TasksOrderField TagsByStatus err", err)
			return wf, err
		}
		log.Println("TasksOrderField tagsMap:", tagsMap)

		OrderTasksByTags(wf, tagsMap, wantedTags)

		// Force statuses to be re-walked against the rewired chain.
		ctx := context.Background()
		statuses, _ := wf.Statuses(ctx, false, true)
		wf.OrderedByTags = wantedTags
		log.Println("TasksOrderField is ordered:", statuses)
		elapsed := time.Since(start1)
		log.Printf("TasksOrderField elapsed %s", elapsed)
		return wf, nil
	}
}
