package payment

import (
	"log"
	db "stories/database"
	"stories/scanners"
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
	var value scanners.String
	var dataFile []byte
	q := query.NewQuery()

	_, err = q.QueryOne("SELECT flow FROM ta_flow WHERE pk_flow=1", &value)
	if err != nil {
		return dataFile, err
	}
	dataFile = []byte(value.Value.String())
	return dataFile, nil
}
