package pg

import (
	"stories/database"
	"stories/generated/models"
	"stories/types"
)

type TaUserPolicy struct {
	PkUserPolicy types.Int `yaml:"pk_user_policy"`
	FkUser       types.Int `yaml:"fk_user"`
	FkPolicy     types.Int `yaml:"fk_policy"`
}

func InsertTaUserPolicy(order string, query *database.IQuery) error {
	handler, err := UnmarshallTaUserPolicy([]byte(order))
	if err != nil {
		return err
	}
	_, err = models.CreateTaUserPolicy(*query,
		handler.FkUser,
		handler.FkPolicy,
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTaUserPolicy(order string, query *database.IQuery) error {
	handler, err := UnmarshallTaUserPolicy([]byte(order))
	if err != nil {
		return err
	}
	taUserPolicy, err := models.GetTaUserPolicyByPkUserPolicy(*query,
		handler.PkUserPolicy,
	)

	if err != nil {
		return err
	}
	err = taUserPolicy.Delete(*query)
	if err != nil {
		return err
	}
	return nil
}
func GetTaUserPolicyByPkUserPolicy(order string, query *database.IQuery) (*models.TaUserPolicy, error) {
	handler, err := UnmarshallTaUserPolicy([]byte(order))
	if err != nil {
		return nil, err
	}
	taUserPolicy, err := models.GetTaUserPolicyByPkUserPolicy(*query,
		handler.PkUserPolicy,
	)
	if err != nil {
		return taUserPolicy, err
	}
	return taUserPolicy, nil
}

func UpdateTaUserPolicy(order string, query *database.IQuery) error {
	handler, err := UnmarshallTaUserPolicy([]byte(order))
	if err != nil {
		return err
	}
	taUserPolicy, err := models.GetTaUserPolicyByPkUserPolicy(*query,
		handler.PkUserPolicy,
	)
	if err != nil {
		return err
	}
	err = taUserPolicy.Update(*query,
		handler.FkUser,
		handler.FkPolicy,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetAllTaUserPolicy(query *database.IQuery) ([]models.TaUserPolicy, error) {

	taUserPolicy, err := models.GetAllTaUserPolicies(*query)
	if err != nil {
		return taUserPolicy, err
	}
	return taUserPolicy, nil
}
