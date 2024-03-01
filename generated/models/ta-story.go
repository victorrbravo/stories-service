package models


import (
	"database/sql"
    "errors"
	"stories/database"
	"stories/scanners"
	"stories/types"
)
/*
 * Generated queries (these can be used perfectly instead of writing them by hand
 */
const (
    TaStoryCreateQuery = "INSERT INTO \"ta_story\" (\"name\",\"status\",\"summary\",\"user_id\",\"fk_flow\") VALUES ($1, $2, $3, $4, $5) RETURNING (\"pk_story\")"
    AllTaStoryQuery = "SELECT \"pk_story\", \"name\", \"status\", \"summary\", \"user_id\", \"fk_flow\" FROM \"ta_story\""
    TaStoryByPkStoryQuery = "SELECT \"pk_story\", \"name\", \"status\", \"summary\", \"user_id\", \"fk_flow\" FROM \"ta_story\" WHERE \"pk_story\" = $1"
    TaStoryByUserIDQuery = "SELECT \"pk_story\", \"name\", \"status\", \"summary\", \"user_id\", \"fk_flow\" FROM \"ta_story\" WHERE \"user_id\" = $1"
    TaStoryByFkFlowQuery = "SELECT \"pk_story\", \"name\", \"status\", \"summary\", \"user_id\", \"fk_flow\" FROM \"ta_story\" WHERE \"fk_flow\" = $1"
    TaStoryUpdateQuery = "UPDATE \"ta_story\" SET \"name\" = COALESCE($1, \"name\"), \"status\" = COALESCE($2, \"status\"), \"summary\" = COALESCE($3, \"summary\"), \"user_id\" = COALESCE($4, \"user_id\"), \"fk_flow\" = COALESCE($5, \"fk_flow\") WHERE \"pk_story\" = $6"
    TaStoryDeleteQuery = "DELETE FROM \"ta_story\" WHERE \"pk_story\" = $1"
    TaStoryListOfFields = " ta_story.\"pk_story\", ta_story.\"name\", ta_story.\"status\", ta_story.\"summary\", ta_story.\"user_id\", ta_story.\"fk_flow\""
)

type TaStory struct {
	PkStory types.Int `json:"pk_story"`
	Name types.String `json:"name"`
	Status types.String `json:"status"`
	Summary types.String `json:"summary"`
	UserID types.Int `json:"user_id"`
	FkFlow types.Int `json:"fk_flow"`
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
	}
	return taStory, nil
}

func (taStory *TaStory) Update(
    query database.IQuery,
    name types.String,
    status types.String,
    summary types.String,
    userID types.Int,
    fkFlow types.Int,
) error {
    var err error
	err = query.Exec(
	    TaStoryUpdateQuery,
        name,
        status,
        summary,
        userID,
        fkFlow,
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
