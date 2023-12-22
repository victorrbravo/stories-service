package workflow_test

import (
	"context"
	"fmt"
	"log"
	"stories/database"
	db "stories/database"
	pgi "stories/generated/inputs"
	"stories/generated/models"
	"stories/pkg/workflow/common"
	input2 "stories/pkg/workflow/input"
	"stories/pkg/workflow/queries/pg"
	"stories/types"
	"testing"
)

const jsonFile = `
{
  "id": "%s",
  "name": "Workflow1",
  "variable": { "name": "basic", "queries": "pk_story,name,status" },
  "tracking_table_name": "ta_tracking",
  "tracking_keys": "fk_story,fk_status",
  "field_token": "pk_story",
  "table_name": "ta_story",
  "initial": "fk_status='ToDo'",
  "tasks": [
    "ToDo": {
      "info_field": "roles[1]",
      "time_field": "updated_at",
      "transition": {
        "pattern": "none",
        "links": [
          {
            "task": "Progress",
            "queries": "fk_status='Progress'"
          }
        ]
      }
    },
    "Progress": {
      "info_field": "roles[1]",
      "time_field": "updated_at",
      "transition": {
        "pattern": "or",
        "links": [
          {
            "task": "Done",
            "queries": "fk_status='Done'"
          }
        ]
      }
    },
    "Done": {
      "info_field": "roles[1]",
      "time_field": "updated_at",
      "transition": {
        "pattern": "or",
        "links": [
          {
            "task": "END",
            "queries": "END"
          }
        ]
      }
    }
  ]
}

`

func getPgQuery(config db.DatabaseConfig) (db.PostgresQuery, error) {
	var query db.PostgresQuery
	query.Initialize(&config)
	err := query.Connect()
	if err != nil {
		log.Println("Db Connect error")
		return query, err
	}
	return query, nil
}

func TestNextStatus(t *testing.T) {
	var err error
	ctx := context.Background()
	config := database.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "stories",
		Password: "stories",
	}
	jw := common.JobWorkflow[models.TaTracking]{}
	query, err := getPgQuery(config)
	if err != nil {
		t.Error(err)
		return
	}
	dataFile, err := getJsonDatafile(&query)
	if err != nil {
		t.Error(err)
		return
	}
	err = jw.Workflow.Load(ctx, dataFile)
	if err != nil {
		t.Error(err)
		return
	}
	var wf = jw.Workflow

	statuses, _ := wf.Statuses(ctx, false)
	fmt.Println("statuses ", statuses)

	status, err := wf.NextStatus("ToDo")
	fmt.Println("status:", status, "err", err)
	status, err = wf.NextStatus("Progress")
	fmt.Println("status:", status, "err", err)
	status, err = wf.NextStatus("Done")
	fmt.Println("status:", status, "err", err)
	status, err = wf.NextStatus("Standby")
	fmt.Println("no status:", status, "err", err)

}

// Test Read method
func TestGenerateStats(t *testing.T) {
	var err error
	ctx := context.Background()
	config := database.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "stories",
		Password: "stories",
	}
	jw := common.JobWorkflow[models.TaTracking]{}
	query, err := getPgQuery(config)
	if err != nil {
		t.Error(err)
		return
	}
	dataFile, err := getJsonDatafile(&query)
	if err != nil {
		t.Error(err)
		return
	}
	err = jw.Workflow.Load(ctx, dataFile)
	if err != nil {
		t.Error(err)
		return
	}
	var wf = jw.Workflow
	monad := pg.New[models.TaTracking](&wf, &query)

	monad = pg.Bind(monad, pg.TimeStats[models.TaTracking])
	monad = pg.Bind(monad, pg.RoleStats[models.TaTracking])
	monad = pg.Bind(monad, pg.Count[models.TaTracking])
	monad = pg.Bind(monad, pg.Total[models.TaTracking])
	monad = pg.Bind(monad, pg.Percentage[models.TaTracking])
	monad = pg.Bind(monad, pg.AllRows[models.TaTracking], "Done")
	monad = pg.Bind(monad, pg.Render[models.TaTracking])
	_, err = monad(nil)
	if err != nil {
		t.Error(err)
	}

	for key, task := range wf.Tasks {
		fmt.Println("===", key, "===")
		fmt.Println("Count", task.Data.Count)
		fmt.Println("TimeField", task.Data.TimeField)
		fmt.Println("Roles", task.Data.InfoField)
		fmt.Println("Percentage", task.Data.Percentage)
		fmt.Println("")
	}
	task := wf.Tasks["Done"]

	for _, data := range task.Data.Rows {
		record := data.(models.TaTracking)
		fmt.Println("fkStory", record.FkStory, "fkStatus", record.FkStatus)
	}
	fmt.Println("total", wf.Data.Total)

}

// Test Read method
func TestGenerateToken(t *testing.T) {
	var err error
	ctx := context.Background()
	config := database.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "stories",
		Password: "stories",
	}
	query, err := getPgQuery(config)
	if err != nil {
		t.Error(err)
		return
	}

	jw := common.JobWorkflow[models.TaTracking]{}
	dataFile, err := getJsonDatafile(&query)
	if err != nil {
		t.Error(err)
		return
	}
	err = jw.Workflow.Load(ctx, dataFile)
	if err != nil {
		t.Error(err)
	}
	jw.Workflow.CurrentIDInt = types.NewInt(9)
	wf := jw.Workflow
	monad := pg.New(&wf, &query)

	monad = pg.Bind(monad, pg.TimeField[models.TaTracking])
	monad = pg.Bind(monad, pg.RoleField[models.TaTracking])
	monad = pg.Bind(monad, pg.CountToken[models.TaTracking])
	monad = pg.Bind(monad, pg.Render[models.TaTracking])
	_, err = monad(nil)
	if err != nil {
		t.Error(err)
	}

	for key, task := range jw.Workflow.Tasks {
		fmt.Println("===", key, "===")
		fmt.Println("Count", task.Data.Count)
		fmt.Println("TimeField", task.Data.TimeField)
		fmt.Println("Roles", task.Data.InfoField)
		fmt.Println("Percentage", task.Data.Percentage)
		fmt.Println("")
	}
	fmt.Println("total", jw.Workflow.Data.Total)
}

func getJsonDatafile(query *db.PostgresQuery) ([]byte, error) {
	var err error
	var value []byte
	input := `
action: get_flow_by_PkFlow
pk_flow: 1
`
	order := []byte(input)
	q := query.NewQuery()
	err = input2.ProcessInput(order, &q)
	if err != nil {
		return value, err
	}
	if len(pgi.ValuesTaFlow) == 0 {
		return value, err
	}
	dataFile := pgi.ValuesTaFlow[0].Flow
	return dataFile, nil
}

func TestWhereExpressions(t *testing.T) {
	pg.CheckWhereExpression("status='Progress'")
}

func TestWhereClauses(t *testing.T) {
	var err error
	ctx := context.Background()
	config := database.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "postgres",
		Password: "stories",
	}
	query, err := getPgQuery(config)
	if err != nil {
		t.Error(err)
		return
	}
	jw := common.JobWorkflow[models.TaTracking]{}
	dataFile, err := getJsonDatafile(&query)
	if err != nil {
		t.Error(err)
		return
	}
	err = jw.Workflow.Load(ctx, dataFile)
	if err != nil {
		t.Error(err)
		return
	}
	wheres := pg.ClausesWheres(jw.Workflow)
	if len(wheres) == 0 {
		t.Logf("len wheres is zero")
	}
	fmt.Println("wheres", wheres)
}
