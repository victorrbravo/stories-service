package pg

import (
	"gopkg.in/yaml.v3"
	"stories/database"
	"stories/generated/models"
)

func ChooseProcessor(action string, order string, query *database.IQuery) error {
	switch action {

	case "add_flow":
		InsertTaFlow(order, query)

	case "add_user":
		InsertTaUser(order, query)

	}
	return nil
}

func UnmarshallTaFlow(order []byte) (models.TaFlow, error) {
	var taFlow models.TaFlow
	var err error
	err = yaml.Unmarshal(order, &taFlow)
	if err != nil {
		return taFlow, err
	}
	return taFlow, nil
}

func UnmarshallTaUser(order []byte) (models.TaUser, error) {
	var taUser models.TaUser
	var err error
	err = yaml.Unmarshal(order, &taUser)
	if err != nil {
		return taUser, err
	}
	return taUser, nil
}
