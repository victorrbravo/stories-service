package query

import (
	"context"
	"errors"
	"fmt"
	"log"
	db "stories/database"
	"stories/pkg/workflow/common"
	"stories/scanners"
)

const (
	SELECT_BASIC     = "SELECT count(*) FROM %s WHERE %s = '%s'"
	SELECT_TIMEFIELD = "SELECT %s FROM %s JOIN %s ON %s"
)

type PgQuery struct {
	query *db.PostgresQuery
}

func (pg *PgQuery) Initialize(config *db.DatabaseConfig) error {
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

func (pg *PgQuery) TimeField(wf *common.Workflow, task int) error {
	var count scanners.Int
	if task >= len(wf.Tasks) {
		return errors.New("len task invalid")
	}
	newQuery := pg.query.NewQuery()
	current := wf.Tasks[task]
	newSQL := fmt.Sprintf(SELECT_TIMEFIELD, current.Transition.Links[0].Status,
		current.Transition.Links[0].Query)
	err := newQuery.QueryOne(newSQL, &count)
	if err != nil {
		return err
	}
	wf.Data.Count = count.Value
	return nil

}

func (pg *PgQuery) Count(wf *common.Workflow, task int) error {
	var count scanners.Int
	if task >= len(wf.Tasks) {
		return errors.New("len task invalid")
	}
	newQuery := pg.query.NewQuery()
	current := wf.Tasks[task]
	newSQL := fmt.Sprintf(SELECT_BASIC, wf.TableName, current.Transition.Links[0].Status,
		current.Transition.Links[0].Query)
	err := newQuery.QueryOne(newSQL, &count)
	if err != nil {
		return err
	}
	wf.Data.Count = count.Value
	return nil
}

func (query *PgQuery) Generate(ctx context.Context, wf *common.Workflow) error {
	fmt.Println("len tasks", len(wf.Tasks))
	for index, task := range wf.Tasks {
		err := query.Count(wf, index)
		if err != nil {
			return err
		}
		fmt.Println("task ", task.Name, "data count", wf.Data.Count)
	}
	return nil
}
