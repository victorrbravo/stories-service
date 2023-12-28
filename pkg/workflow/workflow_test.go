package workflow_test

import (
	"context"
	"fmt"
	"stories/database"
	"stories/generated/models"
	"stories/pkg/workflow"
	"stories/pkg/workflow/common"
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

func TestNextStatus(t *testing.T) {
	var err error
	ctx := context.Background()
	config := database.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "stories",
		Password: "stories",
	}
	jw := common.JobWorkflow[models.TaTracking]{}
	query, err := workflow.GetPgQuery(config)
	if err != nil {
		t.Error(err)
		return
	}
	dataFile, err := workflow.GetJsonDatafile(&query)
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
	query, err := workflow.GetPgQuery(config)
	if err != nil {
		t.Error(err)
		return
	}
	dataFile, err := workflow.GetJsonDatafile(&query)
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
	query, err := workflow.GetPgQuery(config)
	if err != nil {
		t.Error(err)
		return
	}

	jw := common.JobWorkflow[models.TaTracking]{}
	dataFile, err := workflow.GetJsonDatafile(&query)
	if err != nil {
		t.Error(err)
		return
	}
	err = jw.Workflow.Load(ctx, dataFile)
	if err != nil {
		t.Error(err)
	}
	jw.Workflow.CurrentIDInt = types.NewInt(10)
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
	query, err := workflow.GetPgQuery(config)
	if err != nil {
		t.Error(err)
		return
	}
	jw := common.JobWorkflow[models.TaTracking]{}
	dataFile, err := workflow.GetJsonDatafile(&query)
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
