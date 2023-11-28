package workflow_test

import (
	"context"
	"fmt"
	"github.com/twinj/uuid"
	"log"
	"stories/database"
	db "stories/database"
	pgi "stories/generated/inputs"
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
  "tasks": {
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
  }
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

// Test Read method
func TestGenerateStats(t *testing.T) {
	var err error
	ctx := context.Background()
	config := database.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "postgres",
		Password: "stories",
	}
	jw := common.JobWorkflow{}

	input := `
action: get_flow_by_PkFlow
pk_flow: 1
`

	order := []byte(input)
	query, err := getPgQuery(config)
	if err != nil {
		t.Error("error connecting db")
		return
	}
	q := query.NewQuery()
	err = input2.ProcessInput(order, &q)
	if err != nil {
		t.Error(err)
	}
	if len(pgi.ValuesTaFlow) == 0 {
		t.Error("Values TaFlow is zero 0")
		return
	}

	uuidValue := uuid.NewV1().String()
	println(uuidValue)
	dataFile := pgi.ValuesTaFlow[0].Flow
	err = jw.Workflow.Load(ctx, dataFile)
	if err != nil {
		t.Error(err)
		return
	}
	newJson, err := jw.Workflow.Marshal(ctx)
	if err != nil {
		fmt.Printf("marshal err %s\n", err)
		t.Error(err)
	}
	fmt.Println("newJson Marshal:", string(newJson))
	var wf common.Workflow = jw.Workflow
	monad := pg.New(&wf, &query)

	monad = pg.Bind(monad, pg.TimeStats)
	monad = pg.Bind(monad, pg.RoleStats)
	monad = pg.Bind(monad, pg.Count)
	monad = pg.Bind(monad, pg.Percentage)
	monad = pg.Bind(monad, pg.Render)
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

// Test Read method
func TestGenerateToken(t *testing.T) {
	var err error
	ctx := context.Background()
	config := database.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "postgres",
		Password: "stories",
	}
	jw := common.JobWorkflow{}
	newJson := fmt.Sprintf(jsonFile, uuid.NewV1().String())
	dataFile := []byte(newJson)
	err = jw.Workflow.Load(ctx, dataFile)
	if err != nil {
		t.Error(err)
	}
	jw.Workflow.CurrentIDInt = types.NewInt(1)
	query, err := getPgQuery(config)
	if err != nil {
		t.Error("error connecting db")
		return
	}
	monad := pg.New(&jw.Workflow, &query)

	monad = pg.Bind(monad, pg.TimeField)
	monad = pg.Bind(monad, pg.RoleField)
	monad = pg.Bind(monad, pg.CountToken)
	monad = pg.Bind(monad, pg.Render)
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

	jw := common.JobWorkflow{}
	newJson := fmt.Sprintf(jsonFile, uuid.NewV1().String())
	dataFile := []byte(newJson)
	err = jw.Workflow.Load(ctx, dataFile)
	if err != nil {
		t.Error(err)
	}
	wheres := jw.Workflow.ClausesWheres()
	if len(wheres) == 0 {
		t.Logf("len wheres is zero")
	}
	fmt.Println("wheres", wheres)
}
