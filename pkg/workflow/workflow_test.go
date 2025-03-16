package workflow_test

import (
	"context"
	"fmt"
	"stories/actions/stories"
	"stories/database"
	"stories/generated/models"
	"stories/pkg/workflow"
	"stories/pkg/workflow/queries/pg"
	"stories/types"
	"testing"
)

func TestNextStatus(t *testing.T) {
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
	q := query.NewQuery()
	wf, err := stories.GetWorkflow[models.TaTracking](1, &q)
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

const DB_HOST = "104.251.218.55"
const DB_NAME = "stories"
const DB_PASS = "3MncMarinoned9140**1020__"
const DB_USER = "stories"

func TestIncludedWrapped(t *testing.T) {
	var err error
	ctx := context.Background()
	config := database.DatabaseConfig{
		Host: DB_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_PASS,
	}
	query, err := workflow.GetPgQuery(config)
	if err != nil {
		t.Error(err)
		return
	}
	q := query.NewQuery()
	wf, err := stories.GetWorkflow[models.TaTracking](175, &q)
	monad := pg.New(wf, &q)
	monad = pg.Bind(monad, pg.RoleField[models.TaTracking])
	_, err = monad(nil)
	statuses, _ := wf.Statuses(ctx, false)
	fmt.Println("statuses ", statuses)
}

// Filter Function
func TestFilterIncludedWrapped(t *testing.T) {
	var err error
	ctx := context.Background()
	config := database.DatabaseConfig{
		Host: DB_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_PASS,
	}
	query, err := workflow.GetPgQuery(config)
	if err != nil {
		t.Error(err)
		return
	}
	q := query.NewQuery()
	wf, err := stories.GetWorkflow[models.TaTracking](175, &q)
	wf.Filter.Sentence = "restaurant"
	monad := pg.New(wf, &q)
	monad = pg.Bind(monad, pg.RoleField[models.TaTracking])
	monad = pg.Bind(monad, pg.FilterField[models.TaTracking])

	_, err = monad(nil)
	statuses, _ := wf.Statuses(ctx, false)
	fmt.Println("statuses ", statuses)
}

// Test Read method
func TestGenerateStats(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "stories",
		Password: "stories",
	}
	query, err := workflow.GetPgQuery(config)
	if err != nil {
		t.Error(err)
		return
	}
	q := query.NewQuery()
	wf, err := stories.GetWorkflow[models.TaTracking](174, &q)
	if err != nil {
		t.Error(err)
		return
	}
	monad := pg.New[models.TaTracking](wf, &q)

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
	config := database.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "stories",
		Password: "stories",
	}
	query, err := workflow.GetPgQuery(config)
	if err != nil {
		t.Error(err)
		return
	}

	q := query.NewQuery()
	wf, err := stories.GetWorkflow[models.TaTracking](1, &q)
	wf.CurrentIDInt = types.NewInt(11)
	monad := pg.New(wf, &q)
	monad = pg.Bind(monad, pg.TimeField[models.TaTracking])
	monad = pg.Bind(monad, pg.RoleField[models.TaTracking])
	monad = pg.Bind(monad, pg.CountToken[models.TaTracking])
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

func TestWhereExpressions(t *testing.T) {
	pg.CheckWhereExpression("status='Progress'")
}

func TestWhereClauses(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "postgres",
		Password: "stories",
	}
	query, err := workflow.GetPgQuery(config)
	if err != nil {
		t.Error(err)
		return
	}
	q := query.NewQuery()
	wf, err := stories.GetWorkflow[models.TaTracking](1, &q)
	wheres := pg.ClausesWheres[models.TaTracking](*wf)
	if len(wheres) == 0 {
		t.Logf("len wheres is zero")
	}
	fmt.Println("wheres", wheres)
}
