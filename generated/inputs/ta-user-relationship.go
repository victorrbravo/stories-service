package pg

import (
    "github.com/victorrbravo/stories-service/generated/models"
    "github.com/victorrbravo/stories-service/database"
	"github.com/victorrbravo/stories-service/types"
)

type TaUserRelationship struct {
	PkUserRelationship types.Int `yaml:"pk_user_relationship"`
	FkFollowerId types.Int `yaml:"fk_follower_id"`
	FkFollowedId types.Int `yaml:"fk_followed_id"`
	Status types.String `yaml:"status"`
	UpdatedAt types.DateTime `yaml:"updated_at"`
}

func  InsertTaUserRelationship(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaUserRelationship([]byte(order))
    if err != nil {
        return err
    }
    _, err = models.CreateTaUserRelationship(*query,
            handler.FkFollowerId,
            handler.FkFollowedId,
            handler.Status,
            handler.UpdatedAt,
    )
    if err != nil {
        return err
    }
    return nil
}

func  DeleteTaUserRelationship(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaUserRelationship([]byte(order))
    if err != nil {
        return err
    }
    taUserRelationship, err := models.GetTaUserRelationshipByPkUserRelationship(*query,
            handler.PkUserRelationship,
    )

    if err != nil {
        return err
    }
    err = taUserRelationship.Delete(*query)
    if err != nil {
        return err
    }
    return nil
}
func  GetTaUserRelationshipByPkUserRelationship(order string, query *database.IQuery) (*models.TaUserRelationship, error) {
    handler, err := UnmarshallTaUserRelationship([]byte(order))
    if err != nil {
        return nil, err
    }
    taUserRelationship, err := models.GetTaUserRelationshipByPkUserRelationship(*query,
            handler.PkUserRelationship,
    )
    if err != nil {
        return taUserRelationship, err
    }
    return taUserRelationship, nil
}


func  UpdateTaUserRelationship(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaUserRelationship([]byte(order))
    if err != nil {
        return err
    }
    taUserRelationship, err := models.GetTaUserRelationshipByPkUserRelationship(*query,
            handler.PkUserRelationship,
    )
    if err != nil {
        return err
    }
    err = taUserRelationship.Update(*query,
            handler.FkFollowerId,
            handler.FkFollowedId,
            handler.Status,
            handler.UpdatedAt,
    )
    if err != nil {
        return err
    }
    return nil
}


func  GetAllTaUserRelationship(query *database.IQuery) ([]models.TaUserRelationship, error) {

    taUserRelationship, err := models.GetAllTaUserRelationships(*query)
    if err != nil {
        return taUserRelationship, err
    }
    return taUserRelationship, nil
}

