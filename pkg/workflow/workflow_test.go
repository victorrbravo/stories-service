package workflow_test

import (
	"context"
	"fmt"
	"github.com/twinj/uuid"
	"stories/database"
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
  "tasks": {
    "ToDo": {
      "info_field": "roles[1]",
      "time_field": "updated_at",
      "transition": {
        "pattern": "none",
        "links": [
          {
            "task": "Progress",
            "queries": "ToDo",
            "status": "status"
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
            "queries": "Progress",
            "status": "status"
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
            "queries": "Done",
            "status": "status"
          }
        ]
      }
    }
  }
}

`

// Test Read method
func TestGenerateStats(t *testing.T) {
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
	monad := pg.New(&jw.Workflow, &config)

	monad = pg.Bind(monad, pg.TimeStats)
	monad = pg.Bind(monad, pg.RoleStats)
	monad = pg.Bind(monad, pg.Count)
	monad = pg.Bind(monad, pg.Percentage)
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
	monad := pg.New(&jw.Workflow, &config)

	monad = pg.Bind(monad, pg.TimeField)
	monad = pg.Bind(monad, pg.RoleField)
	monad = pg.Bind(monad, pg.CountToken)
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
