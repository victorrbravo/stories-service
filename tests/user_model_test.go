package tests

import (
	"database/sql"
	"stories/types"
	"testing"
)
import "stories/database"
import "stories/generated/models"

func TestModelUser(t *testing.T) {
	var config = database.DatabaseConfig{Host: "localhost", Name: "panel", Port: 5432, User: "vbravo",
		Password: "vbravo"}
	var query database.PostgresQuery
	var users []models.User
	query.Initialize(&config)
	err := query.Connect()
	t.Log("Connect err", err)
	newQuery := query.NewQuery()
	t.Log("Get All Users")
	users, err = models.GetAllUsers(newQuery)
	if err != nil {
		t.Logf("Error GetAllUsers %v", err)
		return
	}
	t.Logf("Users: %v", users)
}

func TestModelUserByID(t *testing.T) {
	var config = database.DatabaseConfig{Host: "localhost", Name: "panel", Port: 5432, User: "vbravo",
		Password: "vbravo"}
	var query database.PostgresQuery
	var user *models.User
	query.Initialize(&config)
	err := query.Connect()
	t.Log("Connect err", err)
	newQuery := query.NewQuery()
	t.Log("Get User By ID")
	user, err = models.GetUserByID(newQuery, types.NewInt(2))
	if err != nil {
		if err == sql.ErrNoRows {
			t.Logf("Error No Rows GetUserByID %v", err)
			return
		}
		t.Errorf("Error f %v", err)
		return
	}
	t.Logf("User: %v", user)
}
