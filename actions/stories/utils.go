package stories

import (
	"fmt"
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

func GetJsonDatafile(pkFlow int, query *db.IQuery) ([]byte, error) {
	var err error
	var value scanners.String
	var dataFile []byte
	newSQL := fmt.Sprintf("SELECT flow FROM ta_flow WHERE pk_flow=%d", pkFlow)
	_, err = (*query).QueryOne(newSQL, &value)
	if err != nil {
		return dataFile, err
	}
	dataFile = []byte(value.Value.String())
	return dataFile, nil
}
