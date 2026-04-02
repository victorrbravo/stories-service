package pg

import (
	"stories/database"
	"stories/generated/models"
	"stories/types"
)

type TaPolicy struct {
	PkPolicy  types.Int    `yaml:"pk_policy"`
	Resource  types.String `yaml:"resource"`
	Name      types.String `yaml:"name"`
	Condition types.Bytes  `yaml:"condition"`
}

func InsertTaPolicy(order string, query *database.IQuery) error {
	handler, err := UnmarshallTaPolicy([]byte(order))
	if err != nil {
		return err
	}
	_, err = models.CreateTaPolicy(*query,
		handler.Resource,
		handler.Name,
		handler.Condition,
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTaPolicy(order string, query *database.IQuery) error {
	handler, err := UnmarshallTaPolicy([]byte(order))
	if err != nil {
		return err
	}
	taPolicy, err := models.GetTaPolicyByPkPolicy(*query,
		handler.PkPolicy,
	)

	if err != nil {
		return err
	}
	err = taPolicy.Delete(*query)
	if err != nil {
		return err
	}
	return nil
}
func GetTaPolicyByPkPolicy(order string, query *database.IQuery) (*models.TaPolicy, error) {
	handler, err := UnmarshallTaPolicy([]byte(order))
	if err != nil {
		return nil, err
	}
	taPolicy, err := models.GetTaPolicyByPkPolicy(*query,
		handler.PkPolicy,
	)
	if err != nil {
		return taPolicy, err
	}
	return taPolicy, nil
}

func UpdateTaPolicy(order string, query *database.IQuery) error {
	handler, err := UnmarshallTaPolicy([]byte(order))
	if err != nil {
		return err
	}
	taPolicy, err := models.GetTaPolicyByPkPolicy(*query,
		handler.PkPolicy,
	)
	if err != nil {
		return err
	}
	err = taPolicy.Update(*query,
		handler.Resource,
		handler.Name,
		handler.Condition,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetAllTaPolicy(query *database.IQuery) ([]models.TaPolicy, error) {

	taPolicy, err := models.GetAllTaPolicies(*query)
	if err != nil {
		return taPolicy, err
	}
	return taPolicy, nil
}
