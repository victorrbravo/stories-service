package pg

import (
	"gopkg.in/yaml.v3"
	"stories/database"
	"stories/generated/models"
)

var ValuesTaFlow []models.TaFlow

var ValuesTaUser []models.TaUser

func ChooseProcessor(action string, order string, query *database.IQuery) error {

	if ValuesTaFlow == nil {
		ValuesTaFlow = make([]models.TaFlow, 0)
	}

	if ValuesTaUser == nil {
		ValuesTaUser = make([]models.TaUser, 0)
	}

	switch action {

	case "add_flow":
		InsertTaFlow(order, query)
	case "delete_flow":
		DeleteTaFlow(order, query)
	case "update_flow":
		UpdateTaFlow(order, query)
	case "all_flow":
		ValuesTaFlow, _ = GetAllTaFlow(query)
	case "get_flow_by_PkFlow":
		ValueTaFlow, _ := GetTaFlowByPkFlow(order, query)
		ValuesTaFlow = append(ValuesTaFlow, *ValueTaFlow)

	case "add_user":
		InsertTaUser(order, query)
	case "delete_user":
		DeleteTaUser(order, query)
	case "update_user":
		UpdateTaUser(order, query)
	case "all_user":
		ValuesTaUser, _ = GetAllTaUser(query)
	case "get_user_by_PkUser":
		ValueTaUser, _ := GetTaUserByPkUser(order, query)
		ValuesTaUser = append(ValuesTaUser, *ValueTaUser)

	}
	return nil
}

func UnmarshallTaFlow(order []byte) (TaFlow, error) {
	var taFlow TaFlow
	var err error
	err = yaml.Unmarshal(order, &taFlow)
	if err != nil {
		return taFlow, err
	}
	return taFlow, nil
}

func UnmarshallTaUser(order []byte) (TaUser, error) {
	var taUser TaUser
	var err error
	err = yaml.Unmarshal(order, &taUser)
	if err != nil {
		return taUser, err
	}
	return taUser, nil
}
