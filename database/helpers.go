package database

import "fmt"

func PrintHello() {
	fmt.Println("Hello database")
}

// NewQuery creates a database IQuery using environment config if available or defaults.
func NewQuery() IQuery {
	config := DatabaseConfig{
		Host:     "localhost",
		Name:     "stories",
		Port:     5432,
		User:     "postgres",
		Password: "stories",
	}
	if cfg, err := InitializeConfig(); err == nil && cfg != nil {
		if cfg.DatabaseConfig.Host != "" {
			config = cfg.DatabaseConfig
		}
	}

	var pg PostgresQuery
	pg.Initialize(&config)
	if err := pg.Connect(); err != nil {
		panic(err)
	}
	return pg.NewQuery()
}
