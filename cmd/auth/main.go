package main

import (
	"log"

	"github.com/victorrbravo/stories-service/database"
	"github.com/victorrbravo/stories-service/web"
)

//go:generate go run ../../generators/queries/main.go
//go:generate go run ../../generators/models/main.go
//go:generate go run ../../generators/inputs/main.go
func main() {
	var server web.Server
	var err error
	config := database.DatabaseConfig{
		Host: "localhost", Name: "stories", Port: 5432, User: "postgres",
		Password: "stories",
	}
	var query database.PostgresQuery
	result := query.Initialize(&config)
	err = query.Connect()
	log.Println("Connect Suceeded! result", result)
	newQuery := query.NewQuery()
	server, err = web.NewServer()
	if err != nil {
		panic(err)
	}
	server.Listen("localhost", 8080, "/var/swagger", newQuery)
}
