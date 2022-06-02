package main

import (
	"log"
	"stories/database"
	"stories/web"
)

//go:generate go run ../generators/queries/main.go
//go:generate go run ../generators/models/main.go
func main() {
	var server web.Server
	var err error
	var config = database.DatabaseConfig{Host: "localhost", Name: "panel", Port: 5432, User: "vbravo",
		Password: "vbravo"}
	var query database.PostgresQuery
	query.Initialize(&config)
	log.Println("Connect" +
		" Suceeded!")
	newQuery := query.NewQuery()
	server, err = web.NewServer()
	if err != nil {
		panic(err)
	}
	server.Listen("localhost", 8080, "/var/swagger", newQuery)
}
