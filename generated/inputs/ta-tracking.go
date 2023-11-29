package pg

import (
    "stories/generated/models"
    "stories/database"
	"stories/types"
)

type TaTracking struct {
	FkStatus types.String `yaml:"fk_status"`
	FkStory types.Int `yaml:"fk_story"`
}

func  InsertTaTracking(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaTracking([]byte(order))
    if err != nil {
        return err
    }
    _, err = models.CreateTaTracking(*query,
            handler.FkStory,
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
    taTracking, err := models.GetTaTrackingByFkStatus(*query,
            handler.FkStatus,
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
func  GetTaTrackingByFkStatus(order string, query *database.IQuery) (*models.TaTracking, error) {
    handler, err := UnmarshallTaTracking([]byte(order))
    if err != nil {
        return nil, err
    }
    taTracking, err := models.GetTaTrackingByFkStatus(*query,
            handler.FkStatus,
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
    taTracking, err := models.GetTaTrackingByFkStatus(*query,
            handler.FkStatus,
    )
    if err != nil {
        return err
    }
    err = taTracking.Update(*query,
            handler.FkStory,
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

