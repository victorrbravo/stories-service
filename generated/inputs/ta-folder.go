package pg

import (
    "github.com/victorrbravo/stories-service/generated/models"
    "github.com/victorrbravo/stories-service/database"
	"github.com/victorrbravo/stories-service/types"
)

type TaFolder struct {
	PkFolder types.Int `yaml:"pk_folder"`
	Name types.String `yaml:"name"`
	UserID types.Int `yaml:"user_id"`
	CreatedAt types.DateTime `yaml:"created_at"`
	UpdatedAt types.DateTime `yaml:"updated_at"`
}

func  InsertTaFolder(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaFolder([]byte(order))
    if err != nil {
        return err
    }
    _, err = models.CreateTaFolder(*query,
            handler.Name,
            handler.UserID,
            handler.CreatedAt,
            handler.UpdatedAt,
    )
    if err != nil {
        return err
    }
    return nil
}

func  DeleteTaFolder(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaFolder([]byte(order))
    if err != nil {
        return err
    }
    taFolder, err := models.GetTaFolderByPkFolder(*query,
            handler.PkFolder,
    )

    if err != nil {
        return err
    }
    err = taFolder.Delete(*query)
    if err != nil {
        return err
    }
    return nil
}
func  GetTaFolderByPkFolder(order string, query *database.IQuery) (*models.TaFolder, error) {
    handler, err := UnmarshallTaFolder([]byte(order))
    if err != nil {
        return nil, err
    }
    taFolder, err := models.GetTaFolderByPkFolder(*query,
            handler.PkFolder,
    )
    if err != nil {
        return taFolder, err
    }
    return taFolder, nil
}


func  UpdateTaFolder(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaFolder([]byte(order))
    if err != nil {
        return err
    }
    taFolder, err := models.GetTaFolderByPkFolder(*query,
            handler.PkFolder,
    )
    if err != nil {
        return err
    }
    err = taFolder.Update(*query,
            handler.Name,
            handler.UserID,
            handler.CreatedAt,
            handler.UpdatedAt,
    )
    if err != nil {
        return err
    }
    return nil
}


func  GetAllTaFolder(query *database.IQuery) ([]models.TaFolder, error) {

    taFolder, err := models.GetAllTaFolders(*query)
    if err != nil {
        return taFolder, err
    }
    return taFolder, nil
}
func  GetTaFolderByName(order string, query *database.IQuery) ([]models.TaFolder, error) {
    handler, err := UnmarshallTaFolder([]byte(order))
    if err != nil {
        return nil, err
    }
    taFolders, err := models.GetTaFolderByName(*query,
            handler.Name,
    )
    if err != nil {
        return taFolders, err
    }
    return taFolders, nil
}
func  GetTaFolderByUserID(order string, query *database.IQuery) ([]models.TaFolder, error) {
    handler, err := UnmarshallTaFolder([]byte(order))
    if err != nil {
        return nil, err
    }
    taFolders, err := models.GetTaFolderByUserID(*query,
            handler.UserID,
    )
    if err != nil {
        return taFolders, err
    }
    return taFolders, nil
}

