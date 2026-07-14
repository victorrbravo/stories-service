package stories

import (
	"context"
	"encoding/json"
	"fmt"
	"stories/database"
	pg "stories/generated/inputs"
	"stories/generated/models"
	"stories/scanners"
	"stories/types"

	"gopkg.in/yaml.v3"
)

func FolderProcessor(ctx context.Context, action string, order string, query *database.IQuery) (string, error) {
	var err error
	var value string
	switch action {
	case "get_folders_by_user_id":
		value, err = GetFoldersByUserID[models.TaTracking](ctx, order, query)
	case "add_new_folder":
		value, err = AddNewFolder(ctx, order, query)
	case "get_folders":
		value, err = GetFoldersByUserID[models.TaTracking](ctx, order, query)
	case "add_story_to_folder":
		value, err = AddStoryToFolder(ctx, order, query)
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

func AddNewFolder(ctx context.Context, order string, query *database.IQuery) (string, error) {
	fmt.Println("...order:", order)
	handler, err := unmarshallTaFolder([]byte(order))
	if err != nil {
		return "", err
	}
	taFolder, err := models.CreateTaFolder(*query,
		handler.Name,
		handler.UserID,
		types.Now(),
		types.Now(),
	)
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(taFolder)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// addStoryToFolderInput accepts either the primary keys directly
// (fk_story + pk_folder) or the names (story_name + folder_name) that are
// resolved to their primary keys before linking.
type addStoryToFolderInput struct {
	FkStory    types.Int    `yaml:"fk_story"`
	PkFolder   types.Int    `yaml:"pk_folder"`
	StoryName  types.String `yaml:"story_name"`
	FolderName types.String `yaml:"folder_name"`
}

// AddStoryToFolder links a story to a folder by setting ta_story.fk_folder.
// It works with explicit ids (fk_story, pk_folder) or, when those are absent,
// by resolving story_name and folder_name to their ids first.
func AddStoryToFolder(ctx context.Context, order string, query *database.IQuery) (string, error) {
	fmt.Println("...order:", order)
	userID := types.NewInt(ctx.Value("user_id").(int))
	var input addStoryToFolderInput
	err := yaml.Unmarshal([]byte(order), &input)
	if err != nil {
		return "", err
	}

	pkStory := input.FkStory
	pkFolder := input.PkFolder

	// When ids are not provided, resolve them from the given names.
	if types.IsNil(pkStory) {
		if input.StoryName.IsEmpty() {
			return "", fmt.Errorf("add_story_to_folder requires fk_story or story_name")
		}
		pkStory, err = getPkStoryByName(input.StoryName, userID, query)
		if err != nil {
			return "", err
		}
	}
	if types.IsNil(pkFolder) {
		if input.FolderName.IsEmpty() {
			return "", fmt.Errorf("add_story_to_folder requires pk_folder or folder_name")
		}
		pkFolder, err = getPkFolderByName(input.FolderName, query)
		if err != nil {
			return "", err
		}
	}
	linkStory, err := models.GetTaStoryByPkStory(*query, pkStory)
	if err != nil {
		return "", err
	}
	err = linkStory.Update(*query,
		types.NilString(),
		types.NilString(),
		types.NilString(),
		types.NilInt,
		types.NilInt,
		types.NilString(),
		types.Now(),
		types.NilContentType,
		pkFolder)
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(*linkStory)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func getPkStoryByName(name types.String, userID types.Int, query *database.IQuery) (types.Int, error) {
	folders, err := models.GetTaFolderByUserIDAndName(*query, userID, name)
	if err != nil {
		return types.NilInt, err
	}
	if len(folders) == 0 {
		return types.NilInt, fmt.Errorf("getPkStoryByName not found for %s", name)
	}
	return folders[0].PkFolder, nil
}

func getPkFolderByName(name types.String, query *database.IQuery) (types.Int, error) {
	var value scanners.Int
	_, err := (*query).QueryOne("SELECT pk_folder FROM ta_folder WHERE name = $1", &value, name)
	if err != nil {
		return types.NilInt, err
	}
	return value.Value, nil
}

func unmarshallTaFolder(order []byte) (models.TaFolder, error) {
	// Decode via the generated input type, which carries the proper yaml tags
	// (models.TaFolder only has json tags, so yaml would never map user_id).
	handler, err := pg.UnmarshallTaFolder(order)
	if err != nil {
		return models.TaFolder{}, err
	}
	taFolder := models.TaFolder{
		PkFolder:  handler.PkFolder,
		Name:      handler.Name,
		UserID:    handler.UserID,
		CreatedAt: handler.CreatedAt,
		UpdatedAt: handler.UpdatedAt,
	}
	return taFolder, nil
}
