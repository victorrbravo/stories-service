package pg

import (
    "github.com/victorrbravo/stories-service/generated/models"
    "github.com/victorrbravo/stories-service/database"
	"github.com/victorrbravo/stories-service/types"
)

type TaCategory struct {
	PkCategory types.Int `yaml:"pk_category"`
	Name types.String `yaml:"name"`
	FkMainCategory types.Int `yaml:"fk_main_category"`
}

func  InsertTaCategory(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaCategory([]byte(order))
    if err != nil {
        return err
    }
    _, err = models.CreateTaCategory(*query,
            handler.Name,
            handler.FkMainCategory,
    )
    if err != nil {
        return err
    }
    return nil
}

func  DeleteTaCategory(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaCategory([]byte(order))
    if err != nil {
        return err
    }
    taCategory, err := models.GetTaCategoryByPkCategory(*query,
            handler.PkCategory,
    )

    if err != nil {
        return err
    }
    err = taCategory.Delete(*query)
    if err != nil {
        return err
    }
    return nil
}
func  GetTaCategoryByPkCategory(order string, query *database.IQuery) (*models.TaCategory, error) {
    handler, err := UnmarshallTaCategory([]byte(order))
    if err != nil {
        return nil, err
    }
    taCategory, err := models.GetTaCategoryByPkCategory(*query,
            handler.PkCategory,
    )
    if err != nil {
        return taCategory, err
    }
    return taCategory, nil
}


func  UpdateTaCategory(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaCategory([]byte(order))
    if err != nil {
        return err
    }
    taCategory, err := models.GetTaCategoryByPkCategory(*query,
            handler.PkCategory,
    )
    if err != nil {
        return err
    }
    err = taCategory.Update(*query,
            handler.Name,
            handler.FkMainCategory,
    )
    if err != nil {
        return err
    }
    return nil
}


func  GetAllTaCategory(query *database.IQuery) ([]models.TaCategory, error) {

    taCategory, err := models.GetAllTaCategories(*query)
    if err != nil {
        return taCategory, err
    }
    return taCategory, nil
}

