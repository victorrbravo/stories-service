package pg

import (
	"stories/database"
	"stories/generated/models"
	"stories/types"
)

type TaSession struct {
	PkSession types.Int      `yaml:"pk_session"`
	Token     types.String   `yaml:"token"`
	FkUser    types.Int      `yaml:"fk_user"`
	StartedAt types.DateTime `yaml:"started_at"`
	Invalid   types.Bool     `yaml:"invalid"`
}

func InsertTaSession(order string, query *database.IQuery) error {
	handler, err := UnmarshallTaSession([]byte(order))
	if err != nil {
		return err
	}
	_, err = models.CreateTaSession(*query,
		handler.Token,
		handler.FkUser,
		handler.StartedAt,
		handler.Invalid,
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTaSession(order string, query *database.IQuery) error {
	handler, err := UnmarshallTaSession([]byte(order))
	if err != nil {
		return err
	}
	taSession, err := models.GetTaSessionByPkSession(*query,
		handler.PkSession,
	)

	if err != nil {
		return err
	}
	err = taSession.Delete(*query)
	if err != nil {
		return err
	}
	return nil
}
func GetTaSessionByPkSession(order string, query *database.IQuery) (*models.TaSession, error) {
	handler, err := UnmarshallTaSession([]byte(order))
	if err != nil {
		return nil, err
	}
	taSession, err := models.GetTaSessionByPkSession(*query,
		handler.PkSession,
	)
	if err != nil {
		return taSession, err
	}
	return taSession, nil
}

func UpdateTaSession(order string, query *database.IQuery) error {
	handler, err := UnmarshallTaSession([]byte(order))
	if err != nil {
		return err
	}
	taSession, err := models.GetTaSessionByPkSession(*query,
		handler.PkSession,
	)
	if err != nil {
		return err
	}
	err = taSession.Update(*query,
		handler.Token,
		handler.FkUser,
		handler.StartedAt,
		handler.Invalid,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetAllTaSession(query *database.IQuery) ([]models.TaSession, error) {

	taSession, err := models.GetAllTaSessions(*query)
	if err != nil {
		return taSession, err
	}
	return taSession, nil
}
func GetTaSessionByToken(order string, query *database.IQuery) ([]models.TaSession, error) {
	handler, err := UnmarshallTaSession([]byte(order))
	if err != nil {
		return nil, err
	}
	taSessions, err := models.GetTaSessionByToken(*query,
		handler.Token,
	)
	if err != nil {
		return taSessions, err
	}
	return taSessions, nil
}
func GetTaSessionByInvalid(order string, query *database.IQuery) ([]models.TaSession, error) {
	handler, err := UnmarshallTaSession([]byte(order))
	if err != nil {
		return nil, err
	}
	taSessions, err := models.GetTaSessionByInvalid(*query,
		handler.Invalid,
	)
	if err != nil {
		return taSessions, err
	}
	return taSessions, nil
}
