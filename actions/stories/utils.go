package stories

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"regexp"
	db "stories/database"
	"stories/scanners"
	"strings"
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

func adaptStringBeforeYaml(input string) string {
	// Regular expression to match "ticket: <int>"
	re := regexp.MustCompile(`ticket:(\s|\t)*\d+`)

	reAction := regexp.MustCompile(`action:(\s|\t)*[a-zA-Z_]+`)

	// Remove the pattern from the string
	result := re.ReplaceAllString(input, "")

	result = reAction.ReplaceAllString(result, "")

	// Trim leading and trailing whitespace from the result
	result = strings.TrimSpace(result)

	// Add "task:" at the beginning of the string
	result = "task: \t" + result

	return result
}

func Unmarshall[T any](order []byte) (T, error) {
	var value T
	var err error
	err = yaml.Unmarshal(order, &value)
	if err != nil {
		return value, err
	}
	return value, nil
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
