package models


import (
	"database/sql"
    "errors"
	"github.com/victorrbravo/stories-service/database"
	"github.com/victorrbravo/stories-service/scanners"
	"github.com/victorrbravo/stories-service/types"
)
/*
 * Generated queries (these can be used perfectly instead of writing them by hand
 */
const (
    TaStoryCreateQuery = "INSERT INTO \"ta_story\" (\"name\",\"status\",\"summary\",\"user_id\",\"fk_flow\",\"story_type\",\"updated_at\",\"content_type\",\"fk_folder\") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING (\"pk_story\")"
    AllTaStoryQuery = "SELECT \"pk_story\", \"name\", \"status\", \"summary\", \"user_id\", \"fk_flow\", \"story_type\", \"updated_at\", \"content_type\", \"fk_folder\" FROM \"ta_story\" ORDER BY \"updated_at asc\" "
    TaStoryByPkStoryQuery = "SELECT \"pk_story\", \"name\", \"status\", \"summary\", \"user_id\", \"fk_flow\", \"story_type\", \"updated_at\", \"content_type\", \"fk_folder\" FROM \"ta_story\" WHERE \"pk_story\" = $1"
    TaStoryByUserIDQuery = "SELECT \"pk_story\", \"name\", \"status\", \"summary\", \"user_id\", \"fk_flow\", \"story_type\", \"updated_at\", \"content_type\", \"fk_folder\" FROM \"ta_story\" WHERE \"user_id\" = $1"
    TaStoryByFkFlowQuery = "SELECT \"pk_story\", \"name\", \"status\", \"summary\", \"user_id\", \"fk_flow\", \"story_type\", \"updated_at\", \"content_type\", \"fk_folder\" FROM \"ta_story\" WHERE \"fk_flow\" = $1"
    TaStoryByFkFolderQuery = "SELECT \"pk_story\", \"name\", \"status\", \"summary\", \"user_id\", \"fk_flow\", \"story_type\", \"updated_at\", \"content_type\", \"fk_folder\" FROM \"ta_story\" WHERE \"fk_folder\" = $1"
    TaStoryUpdateQuery = "UPDATE \"ta_story\" SET \"name\" = COALESCE($1, \"name\"), \"status\" = COALESCE($2, \"status\"), \"summary\" = COALESCE($3, \"summary\"), \"user_id\" = COALESCE($4, \"user_id\"), \"fk_flow\" = COALESCE($5, \"fk_flow\"), \"story_type\" = COALESCE($6, \"story_type\"), \"updated_at\" = COALESCE($7, \"updated_at\"), \"content_type\" = COALESCE($8, \"content_type\"), \"fk_folder\" = COALESCE($9, \"fk_folder\") WHERE \"pk_story\" = $10"
    TaStoryDeleteQuery = "DELETE FROM \"ta_story\" WHERE \"pk_story\" = $1"
    TaStoryByStoryTypeAndUserIDQuery = "SELECT \"pk_story\", \"name\", \"status\", \"summary\", \"user_id\", \"fk_flow\", \"story_type\", \"updated_at\", \"content_type\", \"fk_folder\" FROM \"ta_story\" WHERE \"story_type\" = $1 AND \"user_id\" = $2 ORDER BY updated_at ASC"
    TaStoryByNameAndUserIDQuery = "SELECT \"pk_story\", \"name\", \"status\", \"summary\", \"user_id\", \"fk_flow\", \"story_type\", \"updated_at\", \"content_type\", \"fk_folder\" FROM \"ta_story\" WHERE \"name\" = $1 AND \"user_id\" = $2 ORDER BY updated_at ASC"
    TaStoryListOfFields = " ta_story.\"pk_story\", ta_story.\"name\", ta_story.\"status\", ta_story.\"summary\", ta_story.\"user_id\", ta_story.\"fk_flow\", ta_story.\"story_type\", ta_story.\"updated_at\", ta_story.\"content_type\", ta_story.\"fk_folder\""
)

type TaStory struct {
	PkStory types.Int `json:"pk_story"`
	Name types.String `json:"name"`
	Status types.String `json:"status"`
	Summary types.String `json:"summary"`
	UserID types.Int `json:"user_id"`
	FkFlow types.Int `json:"fk_flow"`
	StoryType types.String `json:"story_type"`
	UpdatedAt types.DateTime `json:"updated_at"`
	ContentType types.ContentType `json:"content_type"`
	FkFolder types.Int `json:"fk_folder"`
}

func (taStory TaStory) ListOfFields() string {
	return TaStoryListOfFields
}


func (taStory TaStory) ScanRow(rows database.RowInterface) (interface{}, error) {
	var err error
	err = rows.Scan(
		&taStory.PkStory,
		&taStory.Name,
		&taStory.Status,
		&taStory.Summary,
		&taStory.UserID,
		&taStory.FkFlow,
		&taStory.StoryType,
		&taStory.UpdatedAt,
		&taStory.ContentType,
		&taStory.FkFolder,
	)
	if err != nil {
		return nil, err
	}
	return taStory, nil
}

func CreateTaStory(
    query database.IQuery,
	name types.String,
	status types.String,
	summary types.String,
	userID types.Int,
	fkFlow types.Int,
	storyType types.String,
	updatedAt types.DateTime,
	contentType types.ContentType,
	fkFolder types.Int,
) (*TaStory, error) {
	var err error
    var scanner scanners.Int
	var taStory *TaStory
	_, err = query.QueryOne(
	    TaStoryCreateQuery,
	    &scanner,
		name,
		status,
		summary,
		userID,
		fkFlow,
		storyType,
		updatedAt,
		contentType,
		fkFolder,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	taStory = &TaStory{
	    PkStory: scanner.Value,
	    Name: name,
	    Status: status,
	    Summary: summary,
	    UserID: userID,
	    FkFlow: fkFlow,
	    StoryType: storyType,
	    UpdatedAt: updatedAt,
	    ContentType: contentType,
	    FkFolder: fkFolder,
	}
	return taStory, nil
}

func (taStory *TaStory) UpdateByOne(
    query database.IQuery,
    options ...func() (string, any),
) error {
var err error
    name := types.NilString()
    status := types.NilString()
    summary := types.NilString()
    userID := types.NilInt
    fkFlow := types.NilInt
    storyType := types.NilString()
    updatedAt := types.NilDateTime
    contentType := types.NilContentType
    fkFolder := types.NilInt
    for _, option := range options {
    		actionName, value := option()
    		switch actionName {
            case "name":
            name = value.(types.String)
            case "status":
            status = value.(types.String)
            case "summary":
            summary = value.(types.String)
            case "userID":
            userID = value.(types.Int)
            case "fkFlow":
            fkFlow = value.(types.Int)
            case "storyType":
            storyType = value.(types.String)
            case "updatedAt":
            updatedAt = value.(types.DateTime)
            case "contentType":
            contentType = value.(types.ContentType)
            case "fkFolder":
            fkFolder = value.(types.Int)
		default:
			return errors.New("unsupported field")
		}
	}
	err = query.Exec(
	    TaStoryUpdateQuery,
        name,
        status,
        summary,
        userID,
        fkFlow,
        storyType,
        updatedAt,
        contentType,
        fkFolder,
        taStory.PkStory,
    )
	if err != nil {
    	return err
	}
    return nil
}


func (taStory *TaStory) Update(
    query database.IQuery,
    name types.String,
    status types.String,
    summary types.String,
    userID types.Int,
    fkFlow types.Int,
    storyType types.String,
    updatedAt types.DateTime,
    contentType types.ContentType,
    fkFolder types.Int,
) error {
    var err error
	err = query.Exec(
	    TaStoryUpdateQuery,
        name,
        status,
        summary,
        userID,
        fkFlow,
        storyType,
        updatedAt,
        contentType,
        fkFolder,
        taStory.PkStory,
    )
	if err != nil {
    	return err
	}
    taStory.Name = name
    taStory.Status = status
    taStory.Summary = summary
    taStory.UserID = userID
    taStory.FkFlow = fkFlow
    taStory.StoryType = storyType
    taStory.UpdatedAt = updatedAt
    taStory.ContentType = contentType
    taStory.FkFolder = fkFolder
    return nil
}

func GetTaStoryByPkStory(query database.IQuery, value types.Int) (*TaStory, error) {
	var err error
	var taStory TaStory
	var valueTaStory interface{}
	valueTaStory, err = query.QueryOne(TaStoryByPkStoryQuery, &taStory, value)
	if err != nil {
		return  nil, err
	}
    taStory = valueTaStory.(TaStory)
	return &taStory, nil
}



func GetTaStoryByStoryTypeAndUserID(query database.IQuery,
storyType types.String,
userID types.Int,
) ([]TaStory, error) {
	    var err error
        var rows []interface{}
        var list []TaStory
	rows, err = query.QueryMany(TaStoryByStoryTypeAndUserIDQuery, &TaStory{},
    storyType,
    userID,
)
	if err != nil {
        if err == sql.ErrNoRows {
        	return make([]TaStory, 0), nil
        }
        return nil, err
       }
       for _, row := range rows {
        	list = append(list, row.(TaStory))
        }
       return list, nil
}

func GetTaStoriesByStoryTypeAndUserIDPaged(query database.IQuery,
storyType types.String,
userID types.Int,
pageBase types.PageBase) (types.Page, error) {
	var err error
	var page types.Page
	page, err = query.QueryManyPaged(TaStoryByStoryTypeAndUserIDQuery,  &TaStory{},
	pageBase.PageNumber,
	pageBase.PageSize,
        storyType,
        userID,
	)
	if err != nil {
    		return  page, err
    }
	return page, nil
}

func GetTaStoryByNameAndUserID(query database.IQuery,
name types.String,
userID types.Int,
) ([]TaStory, error) {
	    var err error
        var rows []interface{}
        var list []TaStory
	rows, err = query.QueryMany(TaStoryByNameAndUserIDQuery, &TaStory{},
    name,
    userID,
)
	if err != nil {
        if err == sql.ErrNoRows {
        	return make([]TaStory, 0), nil
        }
        return nil, err
       }
       for _, row := range rows {
        	list = append(list, row.(TaStory))
        }
       return list, nil
}

func GetTaStoriesByNameAndUserIDPaged(query database.IQuery,
name types.String,
userID types.Int,
pageBase types.PageBase) (types.Page, error) {
	var err error
	var page types.Page
	page, err = query.QueryManyPaged(TaStoryByNameAndUserIDQuery,  &TaStory{},
	pageBase.PageNumber,
	pageBase.PageSize,
        name,
        userID,
	)
	if err != nil {
    		return  page, err
    }
	return page, nil
}


func GetTaStoryByUserID(query database.IQuery, value types.Int) ([]TaStory, error) {
	    var err error
    	var rows []interface{}
    	var list []TaStory
    	rows, err = query.QueryMany(TaStoryByUserIDQuery, &TaStory{}, value)
    	if err != nil {
    		if err == sql.ErrNoRows {
    			return make([]TaStory, 0), nil
    		}
    		return nil, err
    	}
    	for _, row := range rows {
    		list = append(list, row.(TaStory))
    	}
    	return list, nil
}
func GetTaStoryByFkFlow(query database.IQuery, value types.Int) ([]TaStory, error) {
	    var err error
    	var rows []interface{}
    	var list []TaStory
    	rows, err = query.QueryMany(TaStoryByFkFlowQuery, &TaStory{}, value)
    	if err != nil {
    		if err == sql.ErrNoRows {
    			return make([]TaStory, 0), nil
    		}
    		return nil, err
    	}
    	for _, row := range rows {
    		list = append(list, row.(TaStory))
    	}
    	return list, nil
}
func GetTaStoryByFkFolder(query database.IQuery, value types.Int) ([]TaStory, error) {
	    var err error
    	var rows []interface{}
    	var list []TaStory
    	rows, err = query.QueryMany(TaStoryByFkFolderQuery, &TaStory{}, value)
    	if err != nil {
    		if err == sql.ErrNoRows {
    			return make([]TaStory, 0), nil
    		}
    		return nil, err
    	}
    	for _, row := range rows {
    		list = append(list, row.(TaStory))
    	}
    	return list, nil
}




func (taStory TaStory) GetPkStory() types.Int {
	return taStory.PkStory}


func (taStory TaStory) GetName() types.String {
	return taStory.Name}


func (taStory TaStory) GetStatus() types.String {
	return taStory.Status}


func (taStory TaStory) GetSummary() types.String {
	return taStory.Summary}


func (taStory TaStory) GetUserID() types.Int {
	return taStory.UserID}


func (taStory TaStory) GetFkFlow() types.Int {
	return taStory.FkFlow}


func (taStory TaStory) GetFkFolder() types.Int {
	return taStory.FkFolder}


func (taStory *TaStory) Delete(query database.IQuery) error {
	var err error
	err = query.Exec(
	    TaStoryDeleteQuery,
        taStory.PkStory,
	)
	if err != nil {
		return err
	}
	return nil
}


func GetAllTaStories(query database.IQuery) ([]TaStory, error) {
	var err error
	var rows []interface{}
	var list []TaStory
	rows, err = query.QueryMany(AllTaStoryQuery, &TaStory{})
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]TaStory, 0), nil
		}
		return nil, err
	}
	for _, row := range rows {
		list = append(list, row.(TaStory))
	}
	return list, nil
}
