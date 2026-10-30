package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/victorrbravo/stories-service/database"
	"github.com/victorrbravo/stories-service/web"
)

//go:generate go run ../../generators/queries/main.go
//go:generate go run ../../generators/models/main.go
//go:generate go run ../../generators/inputs/main.go
func main() {
	var server web.Server
	var err error
	// Load secrets and settings from .env when present; real env vars take precedence
	if err = godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}
	port, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		log.Fatalf("invalid DB_PORT: %v", err)
	}
	config := database.DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Name:     getEnv("DB_NAME", "stories"),
		Port:     port,
		User:     getEnv("DB_USER", "stories"),
		Password: getEnv("DB_PASS", "stories"),
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

func getEnv(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
