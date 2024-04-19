package pg

import (
    "stories/generated/models"
    "stories/database"
	"stories/types"
)

type TaUser struct {
	PkUser types.Int `yaml:"pk_user"`
	FirstName types.String `yaml:"first_name"`
	LastName types.String `yaml:"last_name"`
	Email types.String `yaml:"email"`
	Password types.String `yaml:"password"`
	AccountType types.String `yaml:"account_type"`
	Activated types.Bool `yaml:"activated"`
	Code types.String `yaml:"code"`
}

func  InsertTaUser(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaUser([]byte(order))
    if err != nil {
        return err
    }
    _, err = models.CreateTaUser(*query,
            handler.FirstName,
            handler.LastName,
            handler.Email,
            handler.Password,
            handler.AccountType,
            handler.Activated,
            handler.Code,
    )
    if err != nil {
        return err
    }
    return nil
}

func  DeleteTaUser(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaUser([]byte(order))
    if err != nil {
        return err
    }
    taUser, err := models.GetTaUserByPkUser(*query,
            handler.PkUser,
    )

    if err != nil {
        return err
    }
    err = taUser.Delete(*query)
    if err != nil {
        return err
    }
    return nil
}
func  GetTaUserByPkUser(order string, query *database.IQuery) (*models.TaUser, error) {
    handler, err := UnmarshallTaUser([]byte(order))
    if err != nil {
        return nil, err
    }
    taUser, err := models.GetTaUserByPkUser(*query,
            handler.PkUser,
    )
    if err != nil {
        return taUser, err
    }
    return taUser, nil
}


func  UpdateTaUser(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaUser([]byte(order))
    if err != nil {
        return err
    }
    taUser, err := models.GetTaUserByPkUser(*query,
            handler.PkUser,
    )
    if err != nil {
        return err
    }
    err = taUser.Update(*query,
            handler.FirstName,
            handler.LastName,
            handler.Email,
            handler.Password,
            handler.AccountType,
            handler.Activated,
            handler.Code,
    )
    if err != nil {
        return err
    }
    return nil
}


func  GetAllTaUser(query *database.IQuery) ([]models.TaUser, error) {

    taUser, err := models.GetAllTaUsers(*query)
    if err != nil {
        return taUser, err
    }
    return taUser, nil
}
func  GetTaUserByEmail(order string, query *database.IQuery) ([]models.TaUser, error) {
    handler, err := UnmarshallTaUser([]byte(order))
    if err != nil {
        return nil, err
    }
    taUsers, err := models.GetTaUserByEmail(*query,
            handler.Email,
    )
    if err != nil {
        return taUsers, err
    }
    return taUsers, nil
}

