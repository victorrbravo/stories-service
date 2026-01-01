package input

import (
	"context"
	"fmt"
	"os"
	"stories/database"
	db "stories/database"
	pg "stories/generated/inputs"
	"stories/pkg/utils"
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
action: new_step_tips step_id: 167
tips: tips1 for Oporto, tips2 for Burgos, tips3 from Barcelona, tips4 from Madrid
`
	inputFormated := utils.FormatYaml(input)
	order := []byte(inputFormated)
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &query)
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
	input1 := `action: add_flow_step ticket: 198 status: Tres media amount: 18876 summary: Bill media_date: 2025-11-12T01:59:14Z type: step  position: (-72.48997951061656, 7.893217399513308) tags: Tres,media,Bill `
	fmt.Println("input1", input1)
	order := []byte(input1)
	ctx := context.Background()
	_, err = ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}
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
		Host: DB_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_PASS,
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
	input := `
action: add_flow_step
ticket: 201
status: Metro market
images: 
 - ABC897B7-71DA-4C08-B7E8-B7F2F071D445
amount: 37300.00
summary: Market
media_date: 2025-11-25T23:02:26Z
type: step
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
const DB_NAME = "stories"
const DB_PASS = "3MncMarinoned9140**1020__"
const DB_USER = "stories"
const OPEN_API_KEY = "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"

func TestProcessShowFlowData(t *testing.T) {
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
action: show_flow_data 
ticket: 205
filter: >jan1
page_number: 1
page_size: 12
by_day: true
`

	order := []byte(input)
	ctx := context.Background()
	response, err := ProcessInput(ctx, order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response", response)
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
name: Tiramisu18
summary: tiramisu recipe for homemade cooking
email: solazver@gmail.com
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
