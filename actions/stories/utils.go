package stories

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	db "github.com/victorrbravo/stories-service/database"
	"github.com/victorrbravo/stories-service/generated/models"
	"github.com/victorrbravo/stories-service/pkg/utils"
	"github.com/victorrbravo/stories-service/scanners"
	"github.com/victorrbravo/stories-service/types"
	"strings"

	"gopkg.in/yaml.v3"
)

type WithKey interface {
	Key() int
}

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

func validateRole[T WithKey, U WithKey](ctx context.Context, action types.String, order string, query *db.IQuery) error {

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
	var value T
	condition := policies[0].Condition
	var storedPolicy []U
	err = json.Unmarshal(condition, &storedPolicy)
	if err != nil {
		log.Println("storedPolicy err", err)
		return err
	}
	yamlOrder := utils.FormatYaml(order)
	err = yaml.Unmarshal([]byte(yamlOrder), &value)
	if err != nil {
		log.Println("validateError", err)
		return err
	}
	fmt.Println("policies", storedPolicy)
	fmt.Println("value", value)
	for _, policy := range storedPolicy {
		if policy.Key() == value.Key() {
			return nil
		}
	}
	return fmt.Errorf("invalid role")
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
