package pg

import (
    "stories/generated/models"
    "stories/database"
	"stories/types"
)

type TaStory struct {
	PkStory types.Int `yaml:"pk_story"`
	Name types.String `yaml:"name"`
	Status types.String `yaml:"status"`
	Summary types.String `yaml:"summary"`
	UserID types.Int `yaml:"user_id"`
	FkFlow types.Int `yaml:"fk_flow"`
	StoryType types.String `yaml:"story_type"`
	UpdatedAt types.DateTime `yaml:"updated_at"`
	ContentType types.ContentType `yaml:"content_type"`
}

func  InsertTaStory(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaStory([]byte(order))
    if err != nil {
        return err
    }
    _, err = models.CreateTaStory(*query,
            handler.Name,
            handler.Status,
            handler.Summary,
            handler.UserID,
            handler.FkFlow,
            handler.StoryType,
            handler.UpdatedAt,
            handler.ContentType,
    )
    if err != nil {
        return err
    }
    return nil
}

func  DeleteTaStory(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaStory([]byte(order))
    if err != nil {
        return err
    }
    taStory, err := models.GetTaStoryByPkStory(*query,
            handler.PkStory,
    )

    if err != nil {
        return err
    }
    err = taStory.Delete(*query)
    if err != nil {
        return err
    }
    return nil
}
func  GetTaStoryByPkStory(order string, query *database.IQuery) (*models.TaStory, error) {
    handler, err := UnmarshallTaStory([]byte(order))
    if err != nil {
        return nil, err
    }
    taStory, err := models.GetTaStoryByPkStory(*query,
            handler.PkStory,
    )
    if err != nil {
        return taStory, err
    }
    return taStory, nil
}


func  UpdateTaStory(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaStory([]byte(order))
    if err != nil {
        return err
    }
    taStory, err := models.GetTaStoryByPkStory(*query,
            handler.PkStory,
    )
    if err != nil {
        return err
    }
    err = taStory.Update(*query,
            handler.Name,
            handler.Status,
            handler.Summary,
            handler.UserID,
            handler.FkFlow,
            handler.StoryType,
            handler.UpdatedAt,
            handler.ContentType,
    )
    if err != nil {
        return err
    }
    return nil
}


func  GetAllTaStory(query *database.IQuery) ([]models.TaStory, error) {

    taStory, err := models.GetAllTaStories(*query)
    if err != nil {
        return taStory, err
    }
    return taStory, nil
}
func  GetTaStoryByUserID(order string, query *database.IQuery) ([]models.TaStory, error) {
    handler, err := UnmarshallTaStory([]byte(order))
    if err != nil {
        return nil, err
    }
    taStories, err := models.GetTaStoryByUserID(*query,
            handler.UserID,
    )
    if err != nil {
        return taStories, err
    }
    return taStories, nil
}
func  GetTaStoryByFkFlow(order string, query *database.IQuery) ([]models.TaStory, error) {
    handler, err := UnmarshallTaStory([]byte(order))
    if err != nil {
        return nil, err
    }
    taStories, err := models.GetTaStoryByFkFlow(*query,
            handler.FkFlow,
    )
    if err != nil {
        return taStories, err
    }
    return taStories, nil
}

