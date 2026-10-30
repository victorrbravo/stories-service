package pg

import (
    "github.com/victorrbravo/stories-service/generated/models"
    "github.com/victorrbravo/stories-service/database"
	"github.com/victorrbravo/stories-service/types"
)

type TaTracking struct {
	PkTracking types.Int `yaml:"pk_tracking"`
	FkStory types.Int `yaml:"fk_story"`
	FkStatus types.String `yaml:"fk_status"`
	Summary types.String `yaml:"summary"`
	AmountSpent types.String `yaml:"amount_spent"`
	Images types.StringArray `yaml:"images"`
	Videos types.StringArray `yaml:"videos"`
	Roles types.StringArray `yaml:"roles"`
	UpdatedAt types.DateTime `yaml:"updated_at"`
	FkStoryTips types.Int `yaml:"fk_story_tips"`
	FkStorySub types.Int `yaml:"fk_story_sub"`
	Position types.Coordinates `yaml:"position"`
	MetaOrder types.String `yaml:"meta_order"`
	Tags types.StringArray `yaml:"tags"`
}

func  InsertTaTracking(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaTracking([]byte(order))
    if err != nil {
        return err
    }
    _, err = models.CreateTaTracking(*query,
            handler.FkStory,
            handler.FkStatus,
            handler.Summary,
            handler.AmountSpent,
            handler.Images,
            handler.Videos,
            handler.Roles,
            handler.UpdatedAt,
            handler.FkStoryTips,
            handler.FkStorySub,
            handler.Position,
            handler.MetaOrder,
            handler.Tags,
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
            handler.Summary,
            handler.AmountSpent,
            handler.Images,
            handler.Videos,
            handler.Roles,
            handler.UpdatedAt,
            handler.FkStoryTips,
            handler.FkStorySub,
            handler.Position,
            handler.MetaOrder,
            handler.Tags,
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
func  GetTaTrackingByFkStoryTips(order string, query *database.IQuery) ([]models.TaTracking, error) {
    handler, err := UnmarshallTaTracking([]byte(order))
    if err != nil {
        return nil, err
    }
    taTrackings, err := models.GetTaTrackingByFkStoryTips(*query,
            handler.FkStoryTips,
    )
    if err != nil {
        return taTrackings, err
    }
    return taTrackings, nil
}
func  GetTaTrackingByFkStorySub(order string, query *database.IQuery) ([]models.TaTracking, error) {
    handler, err := UnmarshallTaTracking([]byte(order))
    if err != nil {
        return nil, err
    }
    taTrackings, err := models.GetTaTrackingByFkStorySub(*query,
            handler.FkStorySub,
    )
    if err != nil {
        return taTrackings, err
    }
    return taTrackings, nil
}
func  GetTaTrackingByPosition(order string, query *database.IQuery) ([]models.TaTracking, error) {
    handler, err := UnmarshallTaTracking([]byte(order))
    if err != nil {
        return nil, err
    }
    taTrackings, err := models.GetTaTrackingByPosition(*query,
            handler.Position,
    )
    if err != nil {
        return taTrackings, err
    }
    return taTrackings, nil
}

