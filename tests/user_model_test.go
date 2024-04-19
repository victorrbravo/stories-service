package tests

import (
	"database/sql"
	"testing"

	"stories/types"
)

import (
	"stories/database"
	"stories/generated/models"
)

func TestModelUser(t *testing.T) {
	config := database.DatabaseConfig{
		Host: "localhost", Name: "panel", Port: 5432, User: "vbravo",
		Password: "vbravo",
	}
	var query database.PostgresQuery
	var users []models.TaUser
	query.Initialize(&config)
	err := query.Connect()
	t.Log("Connect err", err)
	newQuery := query.NewQuery()
	t.Log("Get All Users")
	users, err = models.GetAllTaUsers(newQuery)
	if err != nil {
		t.Logf("Error GetAllUsers %v", err)
		return
	}
	t.Logf("Users: %v", users)
}

func TestUserGraph(t *testing.T) {
}

func TestModelAddUser(t *testing.T) {
	config := database.DatabaseConfig{
		Host: "localhost", Name: "panel", Port: 5432, User: "vbravo",
		Password: "vbravo",
	}
	var query database.PostgresQuery
	var user *models.TaUser
	query.Initialize(&config)
	err := query.Connect()
	t.Log("Connect err", err)
	newQuery := query.NewQuery()
	t.Log("Get Add User")
	user, err = models.CreateTaUser(newQuery,
		types.NewString("Marcelo"),
		types.NewString("Bravo"),
		types.NewString("victor2@email.com"),
		types.NewString("pass"),
		types.NewString("Admin"),
		types.NewBool(false),
		types.NilString())
	if err != nil {
		t.Logf("Error GetAllUsers %v", err)
		return
	}
	t.Logf("Users: %v", user)
}

func TestModelUserByID(t *testing.T) {
	config := database.DatabaseConfig{
		Host: "localhost", Name: "panel", Port: 5432, User: "vbravo",
		Password: "vbravo",
	}
	var query database.PostgresQuery
	var user *models.TaUser
	query.Initialize(&config)
	err := query.Connect()
	t.Log("Connect err", err)
	newQuery := query.NewQuery()
	t.Log("Get User By ID")
	user, err = models.GetTaUserByPkUser(newQuery, types.NewInt(1))
	if err != nil {
		if err == sql.ErrNoRows {
			t.Logf("Error No Rows GetUserByID %v", err)
			return
		}
		t.Errorf("Error f %v", err)
		return
	}
	t.Logf("User : %v", user)
}
