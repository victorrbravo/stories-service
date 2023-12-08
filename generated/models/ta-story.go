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
    TaStoryCreateQuery = "INSERT INTO \"ta_story\" (\"name\",\"status\",\"summary\") VALUES ($1, $2, $3) RETURNING (\"pk_story\")"
    AllTaStoryQuery = "SELECT \"pk_story\", \"name\", \"status\", \"summary\" FROM \"ta_story\""
    TaStoryByPkStoryQuery = "SELECT \"pk_story\", \"name\", \"status\", \"summary\" FROM \"ta_story\" WHERE \"pk_story\" = $1"
    TaStoryUpdateQuery = "UPDATE \"ta_story\" SET \"name\" = COALESCE($1, \"name\"), \"status\" = COALESCE($2, \"status\"), \"summary\" = COALESCE($3, \"summary\") WHERE \"pk_story\" = $4"
    TaStoryDeleteQuery = "DELETE FROM \"ta_story\" WHERE \"pk_story\" = $1"
)

type TaStory struct {
	PkStory types.Int `json:"pk_story"`
	Name types.String `json:"name"`
	Status types.String `json:"status"`
	Summary types.String `json:"summary"`
}

func (taStory TaStory) ScanRow(rows database.RowInterface) (interface{}, error) {
	var err error
	err = rows.Scan(
		&taStory.PkStory,
		&taStory.Name,
		&taStory.Status,
		&taStory.Summary,
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
	}
	return taStory, nil
}

func (taStory *TaStory) Update(
    query database.IQuery,
    name types.String,
    status types.String,
    summary types.String,
) error {
    var err error
	err = query.Exec(
	    TaStoryUpdateQuery,
        name,
        status,
        summary,
        taStory.PkStory,
    )
	if err != nil {
    	return err
	}
    taStory.Name = name
    taStory.Status = status
    taStory.Summary = summary
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








func (taStory TaStory) GetPkStory() types.Int {
	return taStory.PkStory}


func (taStory TaStory) GetName() types.String {
	return taStory.Name}


func (taStory TaStory) GetStatus() types.String {
	return taStory.Status}


func (taStory TaStory) GetSummary() types.String {
	return taStory.Summary}


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
