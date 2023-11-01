package workflow_test

import (
	"context"
	"fmt"
	"github.com/twinj/uuid"
	"stories/database"
	"stories/pkg/workflow/common"
	"stories/pkg/workflow/query/pg"
	"stories/types"
	"testing"
)

const jsonFile = `
{
  "id": "%s",
  "name": "Workflow1",
  "variable": { "name": "basic", "query": "pk_story,name,status" },
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
            "query": "ToDo",
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
            "query": "Progress",
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
            "query": "Done",
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
	var pgQuery pg.StatsQuery
	ctx := context.Background()
	config := database.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "postgres",
		Password: "stories",
	}
	err = pgQuery.Initialize(&config)
	if err != nil {
		t.Error(err)
	}
	jw := common.JobWorkflow{}
	newJson := fmt.Sprintf(jsonFile, uuid.NewV1().String())
	dataFile := []byte(newJson)
	err = jw.Workflow.Load(ctx, dataFile)
	if err != nil {
		t.Error(err)
	}
	err = pgQuery.Generate(ctx, &jw.Workflow)
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
	pgQuery := pg.TokenQuery{}
	ctx := context.Background()
	config := database.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "postgres",
		Password: "stories",
	}
	pgQuery.Token = types.NewString("1")
	err = pgQuery.Initialize(&config)
	if err != nil {
		t.Error(err)
	}
	jw := common.JobWorkflow{}
	newJson := fmt.Sprintf(jsonFile, uuid.NewV1().String())
	dataFile := []byte(newJson)
	err = jw.Workflow.Load(ctx, dataFile)
	if err != nil {
		t.Error(err)
	}
	err = pgQuery.Generate(ctx, &jw.Workflow)
	for key, task := range jw.Workflow.Tasks {
		fmt.Println("===", key, "===")
		fmt.Println("Count", task.Data.Count)
		fmt.Println("TimeFields", task.Data.TimeField)
		fmt.Println("Roles", task.Data.InfoField)
		fmt.Println("Percentage", task.Data.Percentage)
		fmt.Println("")
	}
	fmt.Println("total", jw.Workflow.Data.Total)
}
