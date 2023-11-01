package pg

import (
	"bytes"
	"context"
	"database/sql"
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
	SELECT_WHERE_TOKEN = `
SELECT COUNT(*) FROM {{.TableName}} WHERE {{.PkField}} = '{{.ID}}' AND {{.PkStatus}} = '{{.Status}}';
`
	SELECT_TIMEFIELD_TOKEN = `º
SELECT {{.TimeField}}  FROM {{.TrackingTableName}} tt WHERE {{.Key2Link}} = '{{.Status}}' AND {{.Key1Link}} = {{.Token}};
`

	SELECT_ROLEFIELD_TOKEN = `
SELECT {{.InfoField}}  FROM {{.TrackingTableName}} tt WHERE {{.Key2Link}} = '{{.Status}}' AND {{.Key1Link}} = {{.Token}};
`
)

type TokenQuery struct {
	Token types.String
	query *db.PostgresQuery
}

func (query *TokenQuery) Initialize(config *db.DatabaseConfig) error {
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

func (query *TokenQuery) TimeField(ctx context.Context, wf *common.Workflow, statusValue string) error {
	var duration scanners.String
	var outputBuffer bytes.Buffer
	newQuery := query.query.NewQuery()
	current, ok := wf.Tasks[statusValue]
	if ok == false {
		return fmt.Errorf("id doesn't exist %s", statusValue)
	}
	tmpl, err := template.New("timeSQL").Parse(SELECT_TIMEFIELD_TOKEN)
	trackingKeys := strings.Split(wf.TrackingKeys.String(), ",")
	linkStory, linkStatus := wf.ExtractKeys(trackingKeys)
	data := DataTemplate{
		Key1Link:          linkStory,
		Key2Link:          linkStatus,
		TimeField:         current.TimeField,
		TrackingTableName: wf.TrackingTableName.String(),
		Status:            statusValue,
		Token:             query.Token.String(),
	}
	err = tmpl.Execute(&outputBuffer, data)
	if err != nil {
		return err
	}
	newSQL := outputBuffer.String()
	err = newQuery.QueryOne(newSQL, &duration)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	current.Data.TimeField = duration.Value
	return nil
}

func (query *TokenQuery) RoleField(ctx context.Context, wf *common.Workflow, statusValue string) error {
	var role scanners.String
	var outputBuffer bytes.Buffer
	newQuery := query.query.NewQuery()
	current, ok := wf.Tasks[statusValue]
	if ok == false {
		return fmt.Errorf("statusValue doesn't exist %s", statusValue)
	}
	tmpl, err := template.New("timeSQL").Parse(SELECT_ROLEFIELD_TOKEN)
	trackingKeys := strings.Split(wf.TrackingKeys.String(), ",")
	linkStory, linkStatus := wf.ExtractKeys(trackingKeys)
	data := DataTemplate{
		Key1Link:          linkStory,
		Key2Link:          linkStatus,
		InfoField:         current.InfoField,
		TrackingTableName: wf.TrackingTableName.String(),
		Status:            statusValue,
		Token:             query.Token.String(),
	}
	err = tmpl.Execute(&outputBuffer, data)
	if err != nil {
		return err
	}
	newSQL := outputBuffer.String()
	err = newQuery.QueryOne(newSQL, &role)
	fmt.Println("Role newSQL", newSQL)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	current.Data.InfoField = role.Value
	return nil

}

func (query *TokenQuery) Percentage(ctx context.Context, wf *common.Workflow, id string) error {
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
func (query *TokenQuery) Count(ctx context.Context, wf *common.Workflow, id string) error {
	var count scanners.Int
	var outputBuffer bytes.Buffer
	newQuery := query.query.NewQuery()
	current, ok := wf.Tasks[id]
	if !ok { // task is not found
		return fmt.Errorf("Task not found: %s", id)
	}
	tmpl, err := template.New("roleSQL").Parse(SELECT_WHERE_TOKEN)
	data := DataTemplate{
		TableName: wf.TableName,
		Status:    id,
		PkStatus:  current.Transition.Links[0].Status,
	}
	err = tmpl.Execute(&outputBuffer, data)
	if err != nil {
		return err
	}
	newSQL := outputBuffer.String()
	fmt.Println("newSQL", newSQL)
	err = newQuery.QueryOne(newSQL, &count)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	current.Data.Count = count.Value
	return nil
}

func (query *TokenQuery) Generate(ctx context.Context, wf *common.Workflow) error {
	if types.IsNil(query.Token) {
		return fmt.Errorf("token not set")
	}
	for name, _ := range wf.Tasks {
		err := query.TimeField(ctx, wf, name)
		if err != nil {
			return err
		}
		err = query.RoleField(ctx, wf, name)
		if err != nil {
			return err
		}
	}
	return nil
}
