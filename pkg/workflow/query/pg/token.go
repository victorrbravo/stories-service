package pg

import (
	"bytes"
	"context"
	"errors"
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
	TOKEN_BASIC = ""
)

type TokenQuery struct {
	query *db.PostgresQuery
}

func (pg *TokenQuery) Initialize(config *db.DatabaseConfig) error {
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

func (pg *TokenQuery) RoleField(ctx context.Context, wf *common.Workflow, id string) error {
	var roles scanners.String
	newQuery := pg.query.NewQuery()
	current, ok := wf.Tasks[id]
	if ok == false {
		return errors.New(fmt.Sprintf("task doesn't exist", id))
	}
	newSQL := ""
	tmpl, err := template.New("roleSQL").Parse(STATS_ROLEFIELD)
	trackingKeys := strings.Split(wf.TrackingKeys.String(), ",")
	linkStory, linkStatus := extractKeys(trackingKeys)
	data := struct {
		Key1Link          string
		Key2Link          string
		TimeField         string
		InfoField         string
		TrackingTableName string
		StatusesList      string
		Status            string
	}{
		Key1Link:          linkStory,
		Key2Link:          linkStatus,
		TimeField:         current.TimeField,
		InfoField:         current.InfoField,
		TrackingTableName: wf.TrackingTableName.String(),
		StatusesList:      pg.StatusesList(ctx, wf),
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

func (pg *TokenQuery) TimeField(ctx context.Context, wf *common.Workflow, id string) error {
	var duration scanners.Interval
	newQuery := pg.query.NewQuery()
	current, ok := wf.Tasks[id]
	if ok == false {
		return errors.New(fmt.Sprintf("id doesn't exist", id))
	}
	nextTask := wf.NextTask(current.Name)
	if nextTask == nil {
		return errors.New(fmt.Sprintf("id has no a next id", id))
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
		Status:            id,
	}

	var outputBuffer bytes.Buffer
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

func (pg *TokenQuery) Percentage(ctx context.Context, wf *common.Workflow, id string) error {
	total := 0
	if !types.IsNil(wf.Data.Total) {
		total = types.AsInt(wf.Data.Total)
		return nil
	}
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
func (pg *TokenQuery) Count(ctx context.Context, wf *common.Workflow, id string) error {
	var count scanners.Int
	newQuery := pg.query.NewQuery()
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

func (query *TokenQuery) StatusesList(ctx context.Context, wf *common.Workflow) string {
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

func (query *TokenQuery) Generate(ctx context.Context, wf *common.Workflow) error {
	for name, task := range wf.Tasks {
		err := query.Count(ctx, wf, name)
		if err != nil {
			return err
		}
		fmt.Println("task ", task.Name, "data count", wf.Data.Count)
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
