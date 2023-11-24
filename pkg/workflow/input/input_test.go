package input

import (
	"stories/database"
	db "stories/database"
	"testing"
)

func TestProcessInput(t *testing.T) {
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
firstname: Solazver
lastname: sole
account: user
password: secret1
email: solazver@gmail.com
action: add_user
`
	order := []byte(input)
	err = ProcessInput(order, &query)
	if err != nil {
		t.Error(err)
	}

}
