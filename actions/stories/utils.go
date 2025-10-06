package stories

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	db "stories/database"
	"stories/generated/models"
	"stories/scanners"
	"stories/types"
	"strings"

	"gopkg.in/yaml.v3"
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
	result = "task:     " + result

	return result
}

func validateRole[T any](ctx context.Context, action types.String, query *db.IQuery) error {

	type Permission struct {
		Permises []string `json:"permises"`
		Type     string   `json:"type"`
	}

	policies, err := models.GetTaPolicyByName(*query, action)
	if err != nil {
		return err
	}
	if policies == nil {
		log.Println("No policies found")
		return nil
	}
	condition := policies[0].Condition
	var storedPolicy map[string]Permission
	err = json.Unmarshal(condition, &storedPolicy)
	fmt.Println("policies", storedPolicy)
	return nil
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
