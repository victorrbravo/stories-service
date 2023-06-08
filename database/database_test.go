package database

import (
	"testing"

	"stories/generated/queries"
	"stories/types"
)

type User struct {
	FirstName types.String `json:"firstname"`
	LastName  types.String `json:"lastname"`
}

func (user *User) ScanRow(rows RowInterface) (interface{}, error) {
	var err error
	// We scan into thw newly allocated object
	err = rows.Scan(
		&user.FirstName,
		&user.LastName,
	)
	if err != nil {
		return User{}, err
	}
	return *user, nil
}

func TestDatabaseConnect(t *testing.T) {
	config := DatabaseConfig{
		Host: "localhost", Name: "panel", Port: 5432, User: "vbravo",
		Password: "vbravo",
	}
	var query PostgresQuery
	query.Initialize(&config)
	err := query.Connect()
	t.Log("Connect err", err)
	var user User
	// Connect to the database
	newQuery := query.NewQuery()
	newQuery.QueryMany(queries.GetUsers, &user)
	t.Logf("User: Nombre: %s Apellido: %s", user.FirstName, user.LastName)
}
