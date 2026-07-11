package stories

import (
	"context"
	"encoding/json"
	"fmt"
	"stories/database"
	"stories/generated/models"

	"gopkg.in/yaml.v3"
)

func FolderProcessor(ctx context.Context, action string, order string, query *database.IQuery) (string, error) {
	var err error
	var value string
	switch action {
	case "get_folders_by_user_id":
		value, err = GetFoldersByUserID[models.TaTracking](ctx, order, query)
	}
	if err != nil {
		return value, err
	}
	return value, nil

}

func GetFoldersByUserID[T database.Scanner](ctx context.Context, order string, query *database.IQuery) (string, error) {
	fmt.Println("...order:", order)
	var data []byte
	handler, err := unmarshallTaFolder([]byte(order))
	if err != nil {
		return "", err
	}
	taFolder, err := models.GetTaFolderByUserID(*query, handler.UserID)
	data, err = json.Marshal(taFolder)
	if err != nil {
		return "", err
	}
	value := string(data)
	if err != nil {
		return value, err
	}
	return value, nil

}

func unmarshallTaFolder(order []byte) (models.TaFolder, error) {
	var taFolder models.TaFolder
	var err error
	err = yaml.Unmarshal(order, &taFolder)
	if err != nil {
		return taFolder, err
	}
	return taFolder, nil
}
