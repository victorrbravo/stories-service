package pg

import (
    "stories/generated/models"
    "stories/database"
	"stories/types"
)

type TaUser struct {
	PkUser types.Int `yaml:"pk_user"`
	Firstname types.String `yaml:"firstname"`
	Lastname types.String `yaml:"lastname"`
	Email types.String `yaml:"email"`
	Password types.String `yaml:"password"`
	AccountType types.String `yaml:"account_type"`
	Activated types.Bool `yaml:"activated"`
	Username types.String `yaml:"username"`
}

func  InsertTaUser(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaUser([]byte(order))
    if err != nil {
        return err
    }
    _, err = models.CreateTaUser(*query,
            handler.Firstname,
            handler.Lastname,
            handler.Email,
            handler.Password,
            handler.AccountType,
            handler.Activated,
            handler.Username,
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
            handler.Firstname,
            handler.Lastname,
            handler.Email,
            handler.Password,
            handler.AccountType,
            handler.Activated,
            handler.Username,
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
func  GetTaUserByUsername(order string, query *database.IQuery) ([]models.TaUser, error) {
    handler, err := UnmarshallTaUser([]byte(order))
    if err != nil {
        return nil, err
    }
    taUsers, err := models.GetTaUserByUsername(*query,
            handler.Username,
    )
    if err != nil {
        return taUsers, err
    }
    return taUsers, nil
}

