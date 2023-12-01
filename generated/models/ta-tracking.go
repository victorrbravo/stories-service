package models


import (
	"database/sql"
	"stories/database"
	"stories/scanners"
	"stories/types"
)
/*
 * Generated queries (these can be used perfectly instead of writing them by hand
 */
const (
    TaTrackingCreateQuery = "INSERT INTO \"ta_tracking\" (\"fk_story\",\"fk_status\",\"images\") VALUES ($1, $2, $3) RETURNING (\"pk_tracking\")"
    AllTaTrackingQuery = "SELECT \"pk_tracking\", \"fk_story\", \"fk_status\", \"images\" FROM \"ta_tracking\""
    TaTrackingByPkTrackingQuery = "SELECT \"pk_tracking\", \"fk_story\", \"fk_status\", \"images\" FROM \"ta_tracking\" WHERE \"pk_tracking\" = $1"
    TaTrackingUpdateQuery = "UPDATE \"ta_tracking\" SET \"fk_story\" = COALESCE($1, \"fk_story\"), \"fk_status\" = COALESCE($2, \"fk_status\"), \"images\" = COALESCE($3, \"images\") WHERE \"pk_tracking\" = $4"
    TaTrackingDeleteQuery = "DELETE FROM \"ta_tracking\" WHERE \"pk_tracking\" = $1"
)

type TaTracking struct {
	PkTracking types.Int `json:"pk_tracking"`
	FkStory types.Int `json:"fk_story"`
	FkStatus types.String `json:"fk_status"`
	Images types.StringArray `json:"images"`
}

func (taTracking TaTracking) ScanRow(rows database.RowInterface) (interface{}, error) {
	var err error
	err = rows.Scan(
		&taTracking.PkTracking,
		&taTracking.FkStory,
		&taTracking.FkStatus,
		&taTracking.Images,
	)
	if err != nil {
		return nil, err
	}
	return taTracking, nil
}

func CreateTaTracking(
    query database.IQuery,
	fkStory types.Int,
	fkStatus types.String,
	images types.StringArray,
) (*TaTracking, error) {
	var err error
    var scanner scanners.Int
	var taTracking *TaTracking
	_, err = query.QueryOne(
	    TaTrackingCreateQuery,
	    &scanner,
		fkStory,
		fkStatus,
		images,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}
	taTracking = &TaTracking{
	    PkTracking: scanner.Value,
	    FkStory: fkStory,
	    FkStatus: fkStatus,
	    Images: images,
	}
	return taTracking, nil
}

func (taTracking *TaTracking) Update(
    query database.IQuery,
    fkStory types.Int,
    fkStatus types.String,
    images types.StringArray,
) error {
    var err error
	err = query.Exec(
	    TaTrackingUpdateQuery,
        fkStory,
        fkStatus,
        images,
        taTracking.PkTracking,
    )
	if err != nil {
    	return err
	}
    taTracking.FkStory = fkStory
    taTracking.FkStatus = fkStatus
    taTracking.Images = images
    return nil
}

func GetTaTrackingByPkTracking(query database.IQuery, value types.Int) (*TaTracking, error) {
	var err error
	var taTracking TaTracking
	var valueTaTracking interface{}
	valueTaTracking, err = query.QueryOne(TaTrackingByPkTrackingQuery, &taTracking, value)
	if err != nil {
		return  nil, err
	}
    taTracking = valueTaTracking.(TaTracking)
	return &taTracking, nil
}








func (taTracking TaTracking) GetPkTracking() types.Int {
	return taTracking.PkTracking}


func (taTracking TaTracking) GetFkStory() types.Int {
	return taTracking.FkStory}


func (taTracking TaTracking) GetFkStatus() types.String {
	return taTracking.FkStatus}


func (taTracking TaTracking) GetImages() types.StringArray {
	return taTracking.Images}


func (taTracking *TaTracking) Delete(query database.IQuery) error {
	var err error
	err = query.Exec(
	    TaTrackingDeleteQuery,
        taTracking.PkTracking,
	)
	if err != nil {
		return err
	}
	return nil
}


func GetAllTaTrackings(query database.IQuery) ([]TaTracking, error) {
	var err error
	var rows []interface{}
	var list []TaTracking
	rows, err = query.QueryMany(AllTaTrackingQuery, &TaTracking{})
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]TaTracking, 0), nil
		}
		return nil, err
	}
	for _, row := range rows {
		list = append(list, row.(TaTracking))
	}
	return list, nil
}
