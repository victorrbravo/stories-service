package pg

import (
	"stories/database"
	"stories/generated/models"
	"stories/types"
)

type TaUser struct {
	PkUser    types.Int    `yaml:"pk_user"`
	Firstname types.String `yaml:"firstname"`
	Lastname  types.String `yaml:"lastname"`
	Email     types.String `yaml:"email"`
	Password  types.String `yaml:"password"`
	Account   types.String `yaml:"account"`
}

func InsertTaUser(order string, query *database.IQuery) error {
	handler, err := UnmarshallTaUser([]byte(order))
	if err != nil {
		return err
	}
	_, err = models.CreateTaUser(*query,
		handler.Firstname,
		handler.Lastname,
		handler.Email,
		handler.Password,
		handler.Account,
	)
	if err != nil {
		return err
	}
	return nil
}
