package input

import (
	"fmt"
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
	response, err := ProcessInput(order, &query)
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
	_, err = ProcessInput(order, &query)
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
	_, err = ProcessInput(order, &query)
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
	_, err = ProcessInput(order, &query)
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
	_, err = ProcessInput(order, &query)
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
	_, err = ProcessInput(order, &query)
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
	_, err = ProcessInput(order, &query)
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
	_, err = ProcessInput(order, &query)
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
	_, err = ProcessInput(order, &query)
	if err != nil {
		t.Error(err)
	}
}

func TestProcessInputNewPayment(t *testing.T) {
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
action: new_story_flow
name: victor_bravo
summary: an afternoon in Merida
`
	order := []byte(input)
	_, err = ProcessInput(order, &query)
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
	_, err = ProcessInput(order, &query)
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
	_, err = ProcessInput(order, &query)
	if err != nil {
		t.Error(err)
	}
}

func TestProcessAddFlowStep(t *testing.T) {
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
action: add_flow_step
ticket: 44
status: StepNewWithMeta3
images: 
 - ABC897B7-71DA-4C08-B7E8-B7F2F071D445
amount: 3500.00
summary: Last Step
media_date: 2024-03-24T18:05:26Z
type: step
position: (10.2232,-72.2344)
`
	fmt.Println("input", input)
	input1 := `action: add_flow_step ticket: 36 status: Preparing videos: 5FC7B644-1F56-4984-9154-63E78846BD4E.mov, amount: 0 summary: Preparing media_date: 2024-06-10T20:52:42Z type: step position:  (-72.4900176935573, 7.893132814213494)`
	order := []byte(input1)
	_, err = ProcessInput(order, &query)
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
	_, err = ProcessInput(order, &query)
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
	_, err = ProcessInput(order, &query)
	if err != nil {
		t.Error(err)
	}
}

// Test for insert by Date
func TestProcessAddFlowStepOrderByDate(t *testing.T) {
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
action: add_flow_step_ordered
ticket: 44
status: San Pio
images: 
 - ABC897B7-71DA-4C08-B7E8-B7F2F071D445
amount: 1270.00
summary: In San Pio
media_date: 2024-03-20T02:02:26Z
type: step
`
	order := []byte(input)
	_, err = ProcessInput(order, &query)
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
	value, err := ProcessInput(order, &query)
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
	response, err := ProcessInput(order, &query)
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
	response, err := ProcessInput(order, &query)
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
	response, err := ProcessInput(order, &query)
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
	response, err := ProcessInput(order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("response", response)
}

func TestProcessGetStoriesByEmail(t *testing.T) {
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
action: get_stories_by_email email: victorrbravo@gmail.com
`
	input1 := utils.FormatYaml(input)
	order := []byte(input1)
	response, err := ProcessInput(order, &query)
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
	response, err := ProcessInput(order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("response", response)
}

func TestProcessShowFlowData(t *testing.T) {
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
action: show_flow_data
ticket: 49
`
	order := []byte(input)
	response, err := ProcessInput(order, &query)
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
	response, err := ProcessInput(order, &query)
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
	response, err := ProcessInput(order, &query)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("value data response", response)
}

func TestProcessSuggestShortName(t *testing.T) {
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
action: suggest_short_name_for_image
image_id: 5895050E-3FE2-48AF-B191-80BF2AEF8E63.jpg
`
	order := []byte(input)
	response, err := ProcessInput(order, &query)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("value data response", response)
}

func TestProcessShowFlowTaskData(t *testing.T) {
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
action: show_task_data
ticket: 11
task: Step4
`
	order := []byte(input)
	response, err := ProcessInput(order, &query)
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
	_, err = ProcessInput(order, &query)
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
	_, err = ProcessInput(order, &query)
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
	_, err = ProcessInput(order, &query)

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
	_, err = ProcessInput(order, &query)
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
	_, err = ProcessInput(order, &query)

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
	_, err = ProcessInput(order, &query)

}

func TestProcessNewStoryFlow(t *testing.T) {
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
action: new_story_flow
name: Tiramisu6
summary: tiramisu recipe for homemade cooking
email: solazver@gmail.com
`
	order := []byte(input)
	_, err = ProcessInput(order, &query)
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
	_, err = ProcessInput(order, &query)
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
step_id: 167
`
	order := []byte(input)
	_, err = ProcessInput(order, &query)
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
	_, err = ProcessInput(order, &q)
	if err != nil {
		return value, err
	}
	if len(pg.ValuesTaFlow) == 0 {
		return value, err
	}
	dataFile := pg.ValuesTaFlow[0].Flow
	return dataFile, nil
}
