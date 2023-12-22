package pg

import (
    "stories/generated/models"
    "stories/database"
	"stories/types"
)

type TaTracking struct {
	PkTracking types.Int `yaml:"pk_tracking"`
	FkStory types.Int `yaml:"fk_story"`
	FkStatus types.String `yaml:"fk_status"`
	Images types.StringArray `yaml:"images"`
	Roles types.StringArray `yaml:"roles"`
}

func  InsertTaTracking(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaTracking([]byte(order))
    if err != nil {
        return err
    }
    _, err = models.CreateTaTracking(*query,
            handler.FkStory,
            handler.FkStatus,
            handler.Images,
            handler.Roles,
    )
    if err != nil {
        return err
    }
    return nil
}

func  DeleteTaTracking(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaTracking([]byte(order))
    if err != nil {
        return err
    }
    taTracking, err := models.GetTaTrackingByPkTracking(*query,
            handler.PkTracking,
    )

    if err != nil {
        return err
    }
    err = taTracking.Delete(*query)
    if err != nil {
        return err
    }
    return nil
}
func  GetTaTrackingByPkTracking(order string, query *database.IQuery) (*models.TaTracking, error) {
    handler, err := UnmarshallTaTracking([]byte(order))
    if err != nil {
        return nil, err
    }
    taTracking, err := models.GetTaTrackingByPkTracking(*query,
            handler.PkTracking,
    )
    if err != nil {
        return taTracking, err
    }
    return taTracking, nil
}


func  UpdateTaTracking(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaTracking([]byte(order))
    if err != nil {
        return err
    }
    taTracking, err := models.GetTaTrackingByPkTracking(*query,
            handler.PkTracking,
    )
    if err != nil {
        return err
    }
    err = taTracking.Update(*query,
            handler.FkStory,
            handler.FkStatus,
            handler.Images,
            handler.Roles,
    )
    if err != nil {
        return err
    }
    return nil
}


func  GetAllTaTracking(query *database.IQuery) ([]models.TaTracking, error) {

    taTracking, err := models.GetAllTaTrackings(*query)
    if err != nil {
        return taTracking, err
    }
    return taTracking, nil
}

