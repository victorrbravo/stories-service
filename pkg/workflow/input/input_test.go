package input

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"stories/database"
	db "stories/database"
	pg "stories/generated/inputs"
	"stories/pkg/utils"
	"stories/scanners"
	"stories/types"
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
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("response:", response)

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
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &query)
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
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &query)
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
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &query)
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
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &query)
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
action: get_flow_by_PkFlow
pk_flow: 1
`
	order := []byte(input)
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &query)
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
name: Bucaramanga
summary: Check if the key field values were created autommatically
status: ToDo
`
	order := []byte(input)
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &query)
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
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &query)
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
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}
}

func TestProcessInputNewPayment(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: new_story_flow
name: Chacachinto
summary: activities
email: victorrbravo@gmail.com
content_type: Senda

`
	order := []byte(input)
	ctx := context.Background()
	ctxUsr := context.WithValue(ctx, "user_id", 2)
	_, err = ProcessInput(ctxUsr, order, &query)
	if err != nil {
		t.Error(err)
	}
}

func TestProcessInputNewStoryTIps(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: new_step_tips step_id: 2884
tips: La pizza La Veleña combina queso de Chitaga con bocadillo veleño#comma#,La pizza Gaucha incluye carne mechada con chimichurri y chorizo artesanal#comma#,El valor total de la pizza consumida en el video fue de 42.000 pesos#comma#,Ranking de pizzerías en Bucaramanga 1° Lucius#comma# 2° Gauchos y 3° Patagonia#comma#,La masa se describe como suave#comma# de buena calidad y fácil de doblar,`
	inputFormated := utils.FormatYaml(input)
	order := []byte(inputFormated)
	ctx := context.Background()
	newCtx := context.WithValue(ctx, "email", "victorrbravo@gmail.com")
	userCtx := context.WithValue(newCtx, "user_id", 2)
	_, err = ProcessInput(userCtx, order, &query)
	if err != nil {
		t.Error(err)
	}
}

func TestProcessAddStep(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
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
action: add_step
ticket: 10
status: Step2
image: imageBakery
amount: 25000
type: new_thing
`
	order := []byte(input)
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}
}

func TestProcessAddFlowStep(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_LOCAL_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_LOCAL_PASS,
	}
	err = os.Setenv("OPEN_API_KEY", OPEN_API_KEY)
	if err != nil {
		t.Fatalf("Failed to set env var: %v", err)
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input1 := `action: show_flow_data
ticket: 241
filter: '>jun7 and Padre'
by_day: false
page_number: 1
page_size: 12

`
	fmt.Println("input1", input1)
	order := []byte(input1)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}
	log.Println("response", response)
}

func TestProcessUpdateTracking(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
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
action: update_tracking
pk_tracking: 147
amount_spent: 31200.00
images:
   - 41026416-FD28-4357-A2FA-5C42DC257C77
   - ABC897B7-71DA-4C08-B7E8-B7F2F071D445
`
	fmt.Println("input", input)
	order := []byte(input)
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}
}

func TestProcessUpdateTrackingUpdated(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
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
	input1 := "action: updated_tracking pk_tracking: 147 updated_at: 2024-07-15T20:33:00Z"
	formmated := utils.FormatYaml(input1)
	input := `
action: update_tracking
pk_tracking: 147
updated_at: 2024-03-26T10:10:00Z
`
	fmt.Println("input", input)
	order := []byte(formmated)
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}
}

func TestProcessAddFlowStepSub(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := "action: new_step_sub step_id: 123 task:     status: mix tomatoes images: ABC897B7-71DA-4C08-B7E8-B7F2F071D445, amount: 2.5 summary: Buy tomatoes media_date: 2024-03-24T18:05:26Z type: step position: (10.2232,-72.2344)"

	formatInput := utils.FormatYaml(input)
	fmt.Println("formatInput", formatInput)
	order := []byte(formatInput)
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}
}

// Test for insert by Date
func TestProcessAddFlowStepOrderByDate(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_LOCAL_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_LOCAL_PASS,
	}
	var pgQuery db.PostgresQuery
	err = os.Setenv("OPEN_API_KEY", OPEN_API_KEY)
	if err != nil {
		t.Fatalf("Failed to set env var: %v", err)
	}

	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `action: add_flow_step
ticket: 244
status: Sports Broadcast1
videos: 
  - 035ACE95-B216-4BC9-9F52-FBB5014246EF
amount: 20000
summary: A vertical shot captures a TV displaying a sports game with players in yellow jerseys next to a wooden shelf filled with various items and visible wires. A foot in a black and blue striped slipper is seen in the upper right of the frame.
media_date: 2026-07-03T19:24:41Z
type: step
position: (-72.49, 7.8931000000000004)
tags: 
  - captures
  - wires
`
	order := []byte(input)
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}
}

func TestProcessDeleteFlowStep(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
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
action: delete_flow_step id: 174
`
	order := []byte(input)
	ctx := context.Background()
	value, err := ProcessInput(ctx, order, &query)
	if err != nil {

		t.Error(err)
	}
	fmt.Println("delete step value\n", value)
}

func TestProcessShowFlowGraph(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
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
action: show_flow_graph
ticket: 39
`
	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("response", response)
}

func TestProcessGeneratePresign(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
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
action: presign_image
name: coffe.jpg
`
	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("response", response)
}

func TestProcessGenerateMoveTaskInFlowx(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
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
action: move_task 
id_story: 43 
current: Biscuits
`
	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("response", response)
}

func TestProcessGenerateChangeStepReport(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
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
action: change_step_report 
id_story: 43 
name: Cream
report: true
`
	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("response", response)
}

func TestProcessGetStepsSchema(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
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
action: steps_schema
`
	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("response", response)
}

func TestProcessGetStoriesByUserId(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
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
action: get_story_by_user_id
user_id: 2
`
	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("response", response)
}

func TestProcessGetStoriesByEmail(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: get_stories_by_email
`
	input1 := utils.FormatYaml(input)
	order := []byte(input1)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "email", "victorrbravo@gmail.com")
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("response", response)
}

func TestProcessGetStoriesByEmailGrouped(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: get_stories_by_email_grouped
`
	input1 := utils.FormatYaml(input)
	order := []byte(input1)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "email", "victorrbravo@gmail.com")
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("response", response)
}

func TestProcessGePresignedImage(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
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
action: get_presigned_image
name: 762ED0C1-959C-400F-AABC-4E31C4BF5CD6.jpg
`
	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("response", response)
}

const DB_HOST = "doitall.work"
const DB_LOCAL_HOST = "localhost"
const DB_LOCAL_PASS = "stories"
const DB_NAME = "stories"
const DB_PASS = "3MncMarinoned9140**1020__"

const DB_USER = "stories"
const OPEN_API_KEY = "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"

// const GOOGLE_API_KEY = "AIzaSyB14bjScN9L_jLR2XXEUCe1Z98X2ICTII4"
const GOOGLE_API_KEY = "AQ.Ab8RN6KbZDSdsSehxGKkn4tHQqx38SnOgGPlczkYP-NKGBoaZw"
const AWS_ACCESS_KEY_ID = "AKIAYJKDBCUA4WH5EZOK"
const AWS_SECRET_ACCESS_KEY = "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"

func TestProcessShowFlowData(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_LOCAL_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_LOCAL_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()

	input := `
action: delete_user_by_email
email: victorrbravo7@gmail.com
`

	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response: |", response, "|")
}

func TestProcessGetPkTracking(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_LOCAL_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_LOCAL_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()

	input := `
action: get_tags_by_pk_tracking
pk_tracking: 3640`

	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response: |", response, "|")
}

func TestProcessOrderByDateStory(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_LOCAL_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_LOCAL_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()

	input := `
action: order_by_date_story
pk_story: 243`

	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response: |", response, "|")
}

func TestProcessAddTagByPkTracking(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_LOCAL_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_LOCAL_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()

	input := `
action: add_tag_by_pk_tracking
pk_tracking: 3640
tag: Istana`

	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response: |", response, "|")
}

func TestProcessDeleteTagByPkTracking(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_LOCAL_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_LOCAL_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()

	input := `
action: delete_tag_by_pk_tracking
pk_tracking: 3640
tag: Istana Marcelo`

	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response: |", response, "|")
}

func TestProcessShowFlowDataCtgFunction(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_LOCAL_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_LOCAL_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()

	input := `
action: show_flow_data
ticket: 241
filter: 
by_day: true
page_number: 1
page_size: 12
`

	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response: |", response, "|")
}

func TestProcessCategoryShowFlowData(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_LOCAL_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_LOCAL_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()

	input := `
action: show_flow_data
ticket: 239
filter: ctg- bill,coffee,cash 
by_day: false
page_number: 1
page_size: 20
`

	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response: |", response, "|")
}

// Emails of the throwaway accounts TestProcessDeleteUserByEmail works on.
// delete_user_by_email erases an account for good, so the test never points
// it at real data: it creates both users itself and removes whatever is left
// of them afterwards.
const (
	deleteMeEmailA = "delete-me-a@example.invalid"
	deleteMeEmailB = "delete-me-b@example.invalid"
)

// seedDeleteMeUsers creates the two synthetic accounts. User A owns one row in
// every table the action touches; user B exists only to prove that another
// account's tracking row survives the wipe, with the fk_story_tips /
// fk_story_sub references to A's stories nulled instead of deleted.
const seedDeleteMeUsers = `
BEGIN;
INSERT INTO ta_role (name) VALUES ('delete-me-role') ON CONFLICT (name) DO NOTHING;
INSERT INTO ta_policy (resource, name, condition)
VALUES ('delete-me', 'delete-me-policy', '{}'::jsonb) ON CONFLICT (name, condition) DO NOTHING;

INSERT INTO ta_user (firstname, lastname, email, account_type, password, activated, username)
VALUES ('del', 'meA', 'delete-me-a@example.invalid', 'user', 'x', true, 'delete-me-a');
INSERT INTO ta_user (firstname, lastname, email, account_type, password, activated, username)
VALUES ('del', 'meB', 'delete-me-b@example.invalid', 'user', 'x', true, 'delete-me-b');

INSERT INTO ta_flow (name, flow, fk_user)
SELECT 'delete-me-flow-a', '{}'::jsonb, pk_user FROM ta_user WHERE email = 'delete-me-a@example.invalid';
INSERT INTO ta_folder (name, user_id)
SELECT 'delete-me-folder-a', pk_user FROM ta_user WHERE email = 'delete-me-a@example.invalid';

INSERT INTO ta_story (name, status, user_id, fk_flow, fk_folder)
SELECT 'delete-me-story-a1', 'Start', u.pk_user, f.pk_flow, d.pk_folder
FROM ta_user u, ta_flow f, ta_folder d
WHERE u.email = 'delete-me-a@example.invalid'
  AND f.name = 'delete-me-flow-a' AND d.name = 'delete-me-folder-a';
INSERT INTO ta_story (name, status, user_id)
SELECT 'delete-me-story-a2', 'Start', pk_user FROM ta_user WHERE email = 'delete-me-a@example.invalid';
INSERT INTO ta_story (name, status, user_id)
SELECT 'delete-me-story-b1', 'Start', pk_user FROM ta_user WHERE email = 'delete-me-b@example.invalid';

INSERT INTO ta_tracking (fk_status, fk_story, summary)
SELECT 'Start', pk_story, 'delete-me-a-row-1' FROM ta_story WHERE name = 'delete-me-story-a1';
INSERT INTO ta_tracking (fk_status, fk_story, summary)
SELECT 'Start', pk_story, 'delete-me-a-row-2' FROM ta_story WHERE name = 'delete-me-story-a2';
INSERT INTO ta_tracking (fk_status, fk_story, fk_story_tips, fk_story_sub, summary)
SELECT 'Start', b.pk_story, a1.pk_story, a2.pk_story, 'delete-me-b-row'
FROM ta_story b, ta_story a1, ta_story a2
WHERE b.name = 'delete-me-story-b1'
  AND a1.name = 'delete-me-story-a1' AND a2.name = 'delete-me-story-a2';

INSERT INTO ta_session (token, fk_user)
SELECT 'delete-me-token-1', pk_user FROM ta_user WHERE email = 'delete-me-a@example.invalid';
INSERT INTO ta_session (token, fk_user)
SELECT 'delete-me-token-2', pk_user FROM ta_user WHERE email = 'delete-me-a@example.invalid';

INSERT INTO ta_category (name, user_id)
SELECT 'delete-me-category', pk_user FROM ta_user WHERE email = 'delete-me-a@example.invalid';
INSERT INTO ta_permission (fk_from_user, category, permission)
SELECT pk_user, 'delete-me', '{}'::jsonb FROM ta_user WHERE email = 'delete-me-a@example.invalid';
INSERT INTO ta_user_policy (fk_user, fk_policy)
SELECT u.pk_user, p.pk_policy FROM ta_user u, ta_policy p
WHERE u.email = 'delete-me-a@example.invalid' AND p.name = 'delete-me-policy';
INSERT INTO ta_user_role (fk_user, fk_role)
SELECT u.pk_user, r.pk_role FROM ta_user u, ta_role r
WHERE u.email = 'delete-me-a@example.invalid' AND r.name = 'delete-me-role';
INSERT INTO ta_user_relationship (fk_follower_id, fk_followed_id)
SELECT a.pk_user, b.pk_user FROM ta_user a, ta_user b
WHERE a.email = 'delete-me-a@example.invalid' AND b.email = 'delete-me-b@example.invalid';
COMMIT;
`

// purgeDeleteMeUsers removes every seeded row, children first. It runs before
// the seed as well, so a run that died halfway through does not poison the
// next one.
const purgeDeleteMeUsers = `
BEGIN;
DELETE FROM ta_tracking WHERE summary LIKE 'delete-me-%'
   OR fk_story IN (SELECT pk_story FROM ta_story WHERE name LIKE 'delete-me-story-%');
DELETE FROM ta_story WHERE name LIKE 'delete-me-story-%';
DELETE FROM ta_folder WHERE name = 'delete-me-folder-a';
DELETE FROM ta_flow WHERE name = 'delete-me-flow-a';
DELETE FROM ta_category WHERE name = 'delete-me-category';
DELETE FROM ta_permission WHERE category = 'delete-me';
DELETE FROM ta_user_policy WHERE fk_user IN (SELECT pk_user FROM ta_user WHERE email LIKE 'delete-me-%@example.invalid');
DELETE FROM ta_user_role WHERE fk_user IN (SELECT pk_user FROM ta_user WHERE email LIKE 'delete-me-%@example.invalid');
DELETE FROM ta_user_relationship
WHERE fk_follower_id IN (SELECT pk_user FROM ta_user WHERE email LIKE 'delete-me-%@example.invalid')
   OR fk_followed_id IN (SELECT pk_user FROM ta_user WHERE email LIKE 'delete-me-%@example.invalid');
DELETE FROM ta_session WHERE token LIKE 'delete-me-token%';
DELETE FROM ta_user WHERE email LIKE 'delete-me-%@example.invalid';
DELETE FROM ta_policy WHERE name = 'delete-me-policy';
DELETE FROM ta_role WHERE name = 'delete-me-role';
COMMIT;
`

func TestProcessDeleteUserByEmail(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_LOCAL_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_LOCAL_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Fatal(err)
	}
	query := pgQuery.NewQuery()

	runScript(t, &query, purgeDeleteMeUsers)
	runScript(t, &query, seedDeleteMeUsers)
	t.Cleanup(func() { runScript(t, &query, purgeDeleteMeUsers) })

	input := `
action: delete_user_by_email
email: ` + deleteMeEmailA + `
`

	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("value data response: |", response, "|")

	var result struct {
		Email   string         `json:"email"`
		PkUser  int            `json:"pk_user"`
		Deleted map[string]int `json:"deleted"`
	}
	err = json.Unmarshal([]byte(response), &result)
	if err != nil {
		t.Fatalf("response is not the expected json: %v", err)
	}
	if result.Email != deleteMeEmailA {
		t.Errorf("expected email %q, got %q", deleteMeEmailA, result.Email)
	}
	if result.PkUser == 0 {
		t.Error("expected the pk_user of the deleted account")
	}

	// One row per table was seeded, except tracking (two rows of user A) and
	// sessions (two tokens). tracking_unlinked is user B's row, kept but with
	// its references to A's stories nulled.
	expected := map[string]int{
		"tracking": 2, "tracking_unlinked": 1, "stories": 2, "folders": 1,
		"flows": 1, "categories": 1, "permissions": 1, "user_policies": 1,
		"user_roles": 1, "relationships": 1, "sessions": 2, "user": 1,
	}
	for name, want := range expected {
		if got := result.Deleted[name]; got != want {
			t.Errorf("deleted[%s] = %d, expected %d", name, got, want)
		}
	}

	// Nothing of user A survives, in any of the tables.
	leftovers := []struct {
		what      string
		statement string
		values    []interface{}
	}{
		{"users", "SELECT count(*) FROM ta_user WHERE email = $1", []interface{}{deleteMeEmailA}},
		{"stories", "SELECT count(*) FROM ta_story WHERE user_id IN (SELECT pk_user FROM ta_user WHERE email = $1)", []interface{}{deleteMeEmailA}},
		{"flows", "SELECT count(*) FROM ta_flow WHERE name = 'delete-me-flow-a'", nil},
		{"folders", "SELECT count(*) FROM ta_folder WHERE name = 'delete-me-folder-a'", nil},
		{"sessions", "SELECT count(*) FROM ta_session WHERE token LIKE 'delete-me-token%'", nil},
		{"tracking", "SELECT count(*) FROM ta_tracking WHERE summary IN ('delete-me-a-row-1', 'delete-me-a-row-2')", nil},
	}
	for _, leftover := range leftovers {
		if got := countRows(t, &query, leftover.statement, leftover.values...); got != 0 {
			t.Errorf("%s: %d rows of %s survived the delete", leftover.what, got, deleteMeEmailA)
		}
	}

	// User B keeps its own row and its account, with the two references to
	// A's stories nulled rather than deleted.
	kept := countRows(t, &query,
		"SELECT count(*) FROM ta_tracking WHERE summary = 'delete-me-b-row' AND fk_story_tips IS NULL AND fk_story_sub IS NULL")
	if kept != 1 {
		t.Errorf("expected user B's tracking row to survive with nulled references, got %d", kept)
	}
	if got := countRows(t, &query, "SELECT count(*) FROM ta_user WHERE email = $1", deleteMeEmailB); got != 1 {
		t.Errorf("expected user B to survive, got %d rows", got)
	}

	// A second run reports the account is gone instead of a phantom wipe.
	_, err = ProcessInput(ctx, order, &query)
	if err == nil {
		t.Error("expected an error when the email no longer exists")
	}
	fmt.Println("second call err:", err)
}

// runScript executes a multi-statement SQL script. lib/pq only allows that
// when no bind parameters are used, which is why the scripts above inline
// their (constant) values.
func runScript(t *testing.T, query *database.IQuery, script string) {
	t.Helper()
	err := (*query).Exec(script)
	if err != nil {
		t.Fatalf("script failed: %v", err)
	}
}

func countRows(t *testing.T, query *database.IQuery, statement string, values ...interface{}) int {
	t.Helper()
	var counter scanners.Int
	_, err := (*query).QueryOne(statement, &counter, values...)
	if err != nil {
		t.Fatalf("count failed: %v\n%s", err, statement)
	}
	return types.AsInt(counter.Value)
}

func TestProcessAddUserPolicy(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: add_user_policy
fk_user: 2
fk_policy: 5 
`

	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response", response)

}

func TestProcessAddPolicy(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: add_policy
resource: flow
name: add_flow_step
condition: '{"193": { "permises": ["read", "write"], "type": "flow" } }'
`

	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response", response)

}
func TestProcessGetStoriesFollowing(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: get_following_stories
`
	order := []byte(input)
	ctx := context.WithValue(context.Background(), "email", "victorrbravo@gmail.com")

	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response", response)
}

func TestProcessFollowUser(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: follow_user
from_username: v2 
to_username: v1
`
	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response", response)
}

func TestProcessAcceptedUser(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: accept_follower
from_username: v2 
to_username: v1
`
	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response", response)
}

func TestProcessUnFollowUser(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: unfollow_user
from_username: v2 
to_username: v1
`
	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response", response)
}

func TestProcessShowAllCategories(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: all_category
`
	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response", response)
}

func TestProcessShowFlowDeleteTips(t *testing.T) {

}

func TestProcessShowFlowTips(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
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
action: show_step_tips 
tips_id: 20
`
	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response", response)
}

func TestProcessChangeStepField(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
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
action: change_step_field
step_id: 6
position: (70.27255,-82.3224)
`
	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response", response)
}

func TestProcessSuggestShortName(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_PASS,
	}
	err = os.Setenv("OPEN_API_KEY", OPEN_API_KEY)
	if err != nil {
		t.Fatalf("Failed to set env var: %v", err)
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: suggest_short_name_for_image
image_id: 5895050E-3FE2-48AF-B191-80BF2AEF8E63.jpg
`
	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("value data response", response)
}

func TestProcessSuggestVideoShortWordsName(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_LOCAL_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_LOCAL_PASS,
	}
	err = os.Setenv("GOOGLE_API_KEY", GOOGLE_API_KEY)
	if err != nil {
		t.Fatalf("Failed to set env var: %v", err)
	}
	err = os.Setenv("AWS_ACCESS_KEY_ID", AWS_ACCESS_KEY_ID)
	if err != nil {
		t.Fatalf("Failed to set env var: %v", err)
	}
	err = os.Setenv("AWS_SECRET_ACCESS_KEY", AWS_SECRET_ACCESS_KEY)
	if err != nil {
		t.Fatalf("Failed to set env var: %v", err)
	}

	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: suggest_short_name_words_for_video 
image_id: 2A544CD7-7175-4599-80A0-60279915DA49
`
	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("value data response", response)
}

func TestProcessAskAboutCategory(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_PASS,
	}
	err = os.Setenv("OPEN_API_KEY", OPEN_API_KEY)
	if err != nil {
		t.Fatalf("Failed to set env var: %v", err)
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: select_the_category 
question: Given these words '%s' which category does it most likely belong to '%s'?. The response must be only the name of category
`
	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("value data response", response)
}

func TestProcessShowFlowTaskData(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: show_task_data
ticket: 11
task: Step4
`
	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response", response)
}

func TestProcessDoAddTaskToWorkflowNormal(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
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
action: add_task
name: Step2
sentence: fk_status=''Step2''
source: Step1
`
	order := []byte(input)
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &query)
	if err != nil {
		t.Errorf("error reported %v", err)
	}

}

func TestProcessDoAddTaskToWorkflowNoSource(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
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
action: add_task
name: Filter
sentence: roles[1]::int>26000
target: Start
`
	order := []byte(input)
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &query)
	if err != nil {
		t.Errorf("Error reported to add task %v", err.Error())
	}

}

func TestProcessDoAddTaskToWorkflowNoTarget(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
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
action: add_task
name: Payment2
sentence: fk_status=''Payment2''
source: Payment1
`
	order := []byte(input)
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &query)

}

func TestProcessDoDeleteTaskToWorkflowNormal(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
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
action: delete_task
name: Step10
pk_flow: 1
`
	order := []byte(input)
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &query)
	if err != nil {
		t.Errorf("error reported '%v'", err.Error())
	}
}

func TestProcessDoDeleteTaskToWorkflowNoTarget(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
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
action: delete_task
name: Step2
pk_flow: 1
`
	order := []byte(input)
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &query)

}

func TestProcessDoDeleteTaskToWorkflowNoSource(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
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
action: delete_task
name: Start
target: Payment1
pk_flow: 1
`
	order := []byte(input)
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &query)

}

func TestProcessNewStoryFlow(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: new_story_flow
name: February 2026
summary: February days
email: victorrbravo@gmail.com
content_type: Senda
`
	order := []byte(input)
	ctx := context.Background()
	newCtx := context.WithValue(ctx, "email", "victorrbravo@gmail.com")
	userCtx := context.WithValue(newCtx, "user_id", 2)
	_, err = ProcessInput(userCtx, order, &query)
	if err != nil {
		t.Error(err)
	}
}

func TestProcessUpdatePolicy(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: update_policy
pk_policy: 12
condition: '[{"ticket": 210}, {"ticket": 205}, {"ticket": 193}, {"ticket": 198}, {"ticket": 201}]'
`
	order := []byte(input)
	ctx := context.Background()
	newCtx := context.WithValue(ctx, "email", "victorrbravo@gmail.com")
	userCtx := context.WithValue(newCtx, "user_id", 2)
	_, err = ProcessInput(userCtx, order, &query)
	if err != nil {
		t.Error(err)
	}
}

func TestProcessDeleteTicket(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
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
action: delete_ticket
ticket: 43
`
	order := []byte(input)
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}
}

func TestProcessDeleteTips(t *testing.T) {
	var err error
	config := db.DatabaseConfig{
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
action: delete_step_tips
step_id: 123
`
	order := []byte(input)
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}
}

func TestProcessAddNewFolder(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_LOCAL_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_LOCAL_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: add_new_folder
name: Cities
user_id: 2
`
	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response: |", response, "|")
}

func TestProcessGetFolders(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_LOCAL_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_LOCAL_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: get_folders
user_id: 2
`
	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response: |", response, "|")
}

func TestProcessAddStoryToFolderByIds(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_LOCAL_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_LOCAL_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: add_story_to_folder
fk_story: 244
pk_folder: 2
`
	order := []byte(input)
	ctx := context.Background()
	ctxUsr := context.WithValue(ctx, "user_id", 2)
	response, err := ProcessInput(ctxUsr, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response: |", response, "|")
}

func TestProcessAddStoryToFolderByNames(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_LOCAL_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_LOCAL_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: add_story_to_folder
story_name: July 2026
folder_name: Cities
`
	order := []byte(input)
	ctx := context.Background()
	ctxUsr := context.WithValue(ctx, "user_id", 2)
	response, err := ProcessInput(ctxUsr, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response: |", response, "|")
}

func TestProcessUpdateFolderOfStory(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_LOCAL_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_LOCAL_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: update_folder_of_story
pk_story: 243
pk_folder: 2
remove: false
`
	order := []byte(input)
	ctx := context.Background()
	ctxUsr := context.WithValue(ctx, "user_id", 2)
	response, err := ProcessInput(ctxUsr, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response: |", response, "|")
}

func TestProcessUpdateFolderOfStoryRemove(t *testing.T) {
	var err error
	config := database.DatabaseConfig{
		Host: DB_LOCAL_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_LOCAL_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err = pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()
	input := `
action: update_folder_of_story
pk_story: 244
remove: true
`
	order := []byte(input)
	ctx := context.Background()
	ctxUsr := context.WithValue(ctx, "user_id", 2)
	response, err := ProcessInput(ctxUsr, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response: |", response, "|")
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
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &q)
	if err != nil {
		return value, err
	}
	if len(pg.ValuesTaFlow) == 0 {
		return value, err
	}
	dataFile := pg.ValuesTaFlow[0].Flow
	return dataFile, nil
}
