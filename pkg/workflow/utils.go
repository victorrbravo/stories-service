package workflow

import (
	"log"
	db "stories/database"
	pgi "stories/generated/inputs"
	input2 "stories/pkg/workflow/input"
)

func GetPgQuery(config db.DatabaseConfig) (db.PostgresQuery, error) {
	var query db.PostgresQuery
	query.Initialize(&config)
	err := query.Connect()
	if err != nil {
		log.Println("Db Connect error")
		return query, err
	}
	return query, nil
}

func GetJsonDatafile(query *db.PostgresQuery) ([]byte, error) {
	var err error
	var value []byte
	input := `
action: get_flow_by_PkFlow
pk_flow: 1
`
	order := []byte(input)
	q := query.NewQuery()
	err = input2.ProcessInput(order, &q)
	if err != nil {
		return value, err
	}
	if len(pgi.ValuesTaFlow) == 0 {
		return value, err
	}
	dataFile := pgi.ValuesTaFlow[0].Flow
	return dataFile, nil
}
