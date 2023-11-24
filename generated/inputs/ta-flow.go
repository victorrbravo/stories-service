package pg

import (
	"stories/database"
	"stories/generated/models"
	"stories/types"
)

type TaFlow struct {
	PkFlow  types.Int    `yaml:"pk_flow"`
	Name    types.String `yaml:"name"`
	Flow    []byte       `yaml:"flow"`
	FkStory types.String `yaml:"fk_story"`
}

func InsertTaFlow(order string, query *database.IQuery) error {
	handler, err := UnmarshallTaFlow([]byte(order))
	if err != nil {
		return err
	}
	_, err = models.CreateTaFlow(*query,
		handler.Name,
		handler.Flow,
		handler.FkStory,
	)
	if err != nil {
		return err
	}
	return nil
}
