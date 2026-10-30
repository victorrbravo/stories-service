package stories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/victorrbravo/stories-service/database"
	pg "github.com/victorrbravo/stories-service/generated/inputs"
	"github.com/victorrbravo/stories-service/generated/models"
	"github.com/victorrbravo/stories-service/scanners"
	"github.com/victorrbravo/stories-service/types"

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
	case "update_folder_of_story":
		value, err = UpdateFolderOfStory(ctx, order, query)
	case "delete_user_by_email":
		value, err = DeleteUserByEmail(ctx, order, query)
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
	valueID, ok := ctx.Value("user_id").(int)
	if !ok {
		return "", errors.New("User ID required")
	}
	userID := types.NewInt(valueID)
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

// updateFolderOfStoryInput carries the ids of the story and folder to link,
// plus a remove flag. When remove is true the story's fk_folder is set to null
// (pk_folder is ignored); otherwise fk_folder is set to pk_folder.
type updateFolderOfStoryInput struct {
	PkStory  types.Int `yaml:"pk_story"`
	PkFolder types.Int `yaml:"pk_folder"`
	Remove   bool      `yaml:"remove"`
}

// UpdateFolderOfStory updates ta_story.fk_folder for the given story. It follows
// the same approach as AddStoryToFolder, but when remove is true it clears the
// association by setting fk_folder to null. The generated Update uses COALESCE,
// so nulling a column requires an explicit statement instead.
func UpdateFolderOfStory(ctx context.Context, order string, query *database.IQuery) (string, error) {
	fmt.Println("...order:", order)
	var input updateFolderOfStoryInput
	err := yaml.Unmarshal([]byte(order), &input)
	if err != nil {
		return "", err
	}

	if types.IsNil(input.PkStory) {
		return "", fmt.Errorf("update_folder_of_story requires pk_story")
	}

	linkStory, err := models.GetTaStoryByPkStory(*query, input.PkStory)
	if err != nil {
		return "", err
	}

	if input.Remove {
		err = (*query).Exec(
			"UPDATE \"ta_story\" SET \"fk_folder\" = NULL, \"updated_at\" = $1 WHERE \"pk_story\" = $2",
			types.Now(),
			input.PkStory,
		)
		if err != nil {
			return "", err
		}
		linkStory.FkFolder = types.NilInt
		linkStory.UpdatedAt = types.Now()
	} else {
		if types.IsNil(input.PkFolder) {
			return "", fmt.Errorf("update_folder_of_story requires pk_folder when remove is false")
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
			input.PkFolder)
		if err != nil {
			return "", err
		}
	}

	data, err := json.Marshal(*linkStory)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// deleteUserByEmailInput carries the email of the account to erase.
type deleteUserByEmailInput struct {
	Email types.String `yaml:"email"`
}

// deletedRows is the per-table row count of one delete_user_by_email run.
type deletedRows struct {
	Tracking         int `json:"tracking"`
	TrackingUnlinked int `json:"tracking_unlinked"`
	Stories          int `json:"stories"`
	Folders          int `json:"folders"`
	Flows            int `json:"flows"`
	Categories       int `json:"categories"`
	Permissions      int `json:"permissions"`
	UserPolicies     int `json:"user_policies"`
	UserRoles        int `json:"user_roles"`
	Relationships    int `json:"relationships"`
	Sessions         int `json:"sessions"`
	User             int `json:"user"`
}

type deleteUserByEmailResult struct {
	Email   types.String `json:"email"`
	PkUser  types.Int    `json:"pk_user"`
	Deleted deletedRows  `json:"deleted"`
}

// DeleteUserByEmailQuery erases every row a user owns and the user itself.
//
// It is a single statement on purpose: IQuery hands out a plain *Query whose
// Commit/Rollback are no-ops, so a sequence of Exec calls could leave the
// account half deleted. One statement is atomic, and PostgreSQL fires the
// foreign-key triggers only once it finishes, so parent and child rows may be
// removed together inside it.
//
// Every data-modifying CTE reads the same snapshot, so user_stories still
// lists the stories after deleted_stories removed them. Tracking rows of
// *other* users that point at a deleted story through fk_story_tips /
// fk_story_sub are not deleted — those two columns are nulled instead, in a
// single UPDATE (two CTEs touching the same row would be undefined) that skips
// the rows deleted_tracking already claims.
//
// ta_user_role and ta_user_relationship cascade on delete, but they are
// deleted explicitly so the counts below are complete.
const DeleteUserByEmailQuery = `
WITH target AS (
    SELECT pk_user FROM ta_user WHERE lower(email) = lower($1::text)
), user_stories AS (
    SELECT pk_story FROM ta_story WHERE user_id IN (SELECT pk_user FROM target)
), unlinked_tracking AS (
    UPDATE ta_tracking SET
        fk_story_tips = CASE WHEN fk_story_tips IN (SELECT pk_story FROM user_stories)
                             THEN NULL ELSE fk_story_tips END,
        fk_story_sub  = CASE WHEN fk_story_sub IN (SELECT pk_story FROM user_stories)
                             THEN NULL ELSE fk_story_sub END
    WHERE (fk_story IS NULL OR fk_story NOT IN (SELECT pk_story FROM user_stories))
      AND (fk_story_tips IN (SELECT pk_story FROM user_stories)
           OR fk_story_sub IN (SELECT pk_story FROM user_stories))
    RETURNING 1
), deleted_tracking AS (
    DELETE FROM ta_tracking WHERE fk_story IN (SELECT pk_story FROM user_stories) RETURNING 1
), deleted_stories AS (
    DELETE FROM ta_story WHERE user_id IN (SELECT pk_user FROM target) RETURNING 1
), deleted_folders AS (
    DELETE FROM ta_folder WHERE user_id IN (SELECT pk_user FROM target) RETURNING 1
), deleted_flows AS (
    DELETE FROM ta_flow WHERE fk_user IN (SELECT pk_user FROM target) RETURNING 1
), deleted_categories AS (
    DELETE FROM ta_category WHERE user_id IN (SELECT pk_user FROM target) RETURNING 1
), deleted_permissions AS (
    DELETE FROM ta_permission WHERE fk_from_user IN (SELECT pk_user FROM target) RETURNING 1
), deleted_user_policies AS (
    DELETE FROM ta_user_policy WHERE fk_user IN (SELECT pk_user FROM target) RETURNING 1
), deleted_user_roles AS (
    DELETE FROM ta_user_role WHERE fk_user IN (SELECT pk_user FROM target) RETURNING 1
), deleted_relationships AS (
    DELETE FROM ta_user_relationship
    WHERE fk_follower_id IN (SELECT pk_user FROM target)
       OR fk_followed_id IN (SELECT pk_user FROM target)
    RETURNING 1
), deleted_sessions AS (
    DELETE FROM ta_session WHERE fk_user IN (SELECT pk_user FROM target) RETURNING 1
), deleted_user AS (
    DELETE FROM ta_user WHERE pk_user IN (SELECT pk_user FROM target) RETURNING 1
)
SELECT json_build_object(
    'tracking', (SELECT count(*) FROM deleted_tracking),
    'tracking_unlinked', (SELECT count(*) FROM unlinked_tracking),
    'stories', (SELECT count(*) FROM deleted_stories),
    'folders', (SELECT count(*) FROM deleted_folders),
    'flows', (SELECT count(*) FROM deleted_flows),
    'categories', (SELECT count(*) FROM deleted_categories),
    'permissions', (SELECT count(*) FROM deleted_permissions),
    'user_policies', (SELECT count(*) FROM deleted_user_policies),
    'user_roles', (SELECT count(*) FROM deleted_user_roles),
    'relationships', (SELECT count(*) FROM deleted_relationships),
    'sessions', (SELECT count(*) FROM deleted_sessions),
    'user', (SELECT count(*) FROM deleted_user)
)::text
`

// DeleteUserByEmail wipes the account identified by email: its tracking rows,
// stories, flows, folders, sessions and the ta_user row, plus the categories,
// permissions, policies, roles and relationships that reference the user — a
// ta_user row cannot be removed while any of those still point at it.
//
// The deletion is permanent and has no undo. The response reports how many
// rows each table lost, e.g.
//
//	{"email":"someone@example.com","pk_user":42,
//	 "deleted":{"tracking":118,"stories":4,...,"user":1}}
func DeleteUserByEmail(ctx context.Context, order string, query *database.IQuery) (string, error) {
	fmt.Println("...order:", order)
	var input deleteUserByEmailInput
	err := yaml.Unmarshal([]byte(order), &input)
	if err != nil {
		return "", err
	}
	if input.Email.IsEmpty() {
		return "", fmt.Errorf("delete_user_by_email requires email")
	}

	users, err := models.GetTaUserByEmail(*query, input.Email)
	if err != nil {
		return "", err
	}
	if len(users) == 0 {
		return "", fmt.Errorf("delete_user_by_email no user found for %s", input.Email)
	}

	var counters scanners.String
	_, err = (*query).QueryOne(DeleteUserByEmailQuery, &counters, input.Email)
	if err != nil {
		return "", err
	}

	result := deleteUserByEmailResult{Email: users[0].Email, PkUser: users[0].PkUser}
	err = json.Unmarshal([]byte(counters.Value.String()), &result.Deleted)
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func getPkStoryByName(name types.String, userID types.Int, query *database.IQuery) (types.Int, error) {
	names, err := models.GetTaStoryByNameAndUserID(*query, name, userID)
	if err != nil {
		return types.NilInt, err
	}
	if len(names) == 0 {
		return types.NilInt, fmt.Errorf("getPkStoryByName not found for %s", name)
	}
	return names[0].PkStory, nil
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
