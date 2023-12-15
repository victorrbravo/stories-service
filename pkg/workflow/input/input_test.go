package input

import (
	"context"
	"fmt"
	"stories/database"
	db "stories/database"
	pg "stories/generated/inputs"
	"stories/generated/models"
	"stories/pkg/workflow/common"

	"testing"
)

func TestProcessInputCreate(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "postgres",
		Password: "stories",
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
firstname: Omaira
lastname: sole
account: user
password: secret1
email: solazver3@gmail.com
action: add_user
`
	order := []byte(input)
	err = ProcessInput(order, &query)
	if err != nil {
		t.Error(err)
	}

}

func TestProcessInputDeleteUser(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "postgres",
		Password: "stories",
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: delete_user
pk_user: 39
`
	order := []byte(input)
	err = ProcessInput(order, &query)
	if err != nil {
		t.Error(err)
	}

}

func TestProcessInputUpdateUser(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "postgres",
		Password: "stories",
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: update_user
pk_user: 42
firstname: Hecmair
`
	order := []byte(input)
	err = ProcessInput(order, &query)
	if err != nil {
		t.Error(err)
	}

}

func TestProcessInputGetAll(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "postgres",
		Password: "stories",
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: all_user
`
	order := []byte(input)
	err = ProcessInput(order, &query)
	if err != nil {
		t.Error(err)
	}
	for index, value := range pg.ValuesTaUser {
		fmt.Println(index+1, "Firstname", value.Firstname, "Lastname", value.Lastname, "Email", value.Email)
	}

}

func TestProcessInputGetAllFLow(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "postgres",
		Password: "stories",
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: all_flow
`
	order := []byte(input)
	err = ProcessInput(order, &query)
	if err != nil {
		t.Error(err)
	}
	for index, value := range pg.ValuesTaFlow {
		fmt.Println(index+1, "Flow", string(value.Flow))
	}

}

func TestProcessInputGetByPkFlow(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "postgres",
		Password: "stories",
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: get_flow_by_PkFlow
pk_flow: 1
`
	order := []byte(input)
	err = ProcessInput(order, &query)
	if err != nil {
		t.Error(err)
	}
	if len(pg.ValuesTaFlow) == 0 {
		t.Error("Value Ta Flow not found")
		return
	}
	fmt.Println("Flow name:", pg.ValuesTaFlow[0].Name)
	fmt.Println("Flow string:", string(pg.ValuesTaFlow[0].Flow))

}

func TestProcessInputCreateStory(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "stories",
		Password: "stories",
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: add_story
name: checkKey
summary: Check if the key field values were created autommatically
status: ToDo
`
	order := []byte(input)
	err = ProcessInput(order, &query)
	if err != nil {
		t.Error(err)
	}
}

func TestProcessInputCreateTracking(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "postgres",
		Password: "stories",
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: add_tracking
fk_story: 10
fk_status: Progress
images: [image3]
`
	order := []byte(input)
	err = ProcessInput(order, &query)
	if err != nil {
		t.Error(err)
	}
}

func TestProcessInputDeleteTracking(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "postgres",
		Password: "stories",
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: delete_tracking
pk_tracking: 14
`
	order := []byte(input)
	err = ProcessInput(order, &query)
	if err != nil {
		t.Error(err)
	}
}

func TestProcessInputAddTicket(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "postgres",
		Password: "stories",
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: add_ticket
name: add field project1
summary: in1 tickets should be a field that reference to projects
`
	order := []byte(input)
	err = ProcessInput(order, &query)
	if err != nil {
		t.Error(err)
	}
}

func TestProcessChangeStatus(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "postgres",
		Password: "stories",
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: change_status
ticket: 15
status: Progress
`
	order := []byte(input)
	datafile, err := getJsonDatafile(&pgQuery)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()
	wf := common.Workflow[models.TaTracking]{}
	err = wf.Load(ctx, datafile)
	statuses, _ := wf.Statuses(ctx, false)
	fmt.Println("Statuses", statuses)
	if err != nil {
		t.Error(err)
		return
	}
	err = ProcessInput(order, &query)
	if err != nil {
		t.Error(err)
	}
}

func TestProcessDeleteTicket(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "postgres",
		Password: "stories",
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: delete_ticket
ticket: 14
`
	order := []byte(input)
	err = ProcessInput(order, &query)
	if err != nil {
		t.Error(err)
	}
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
	err = ProcessInput(order, &q)
	if err != nil {
		return value, err
	}
	if len(pg.ValuesTaFlow) == 0 {
		return value, err
	}
	dataFile := pg.ValuesTaFlow[0].Flow
	return dataFile, nil
}
