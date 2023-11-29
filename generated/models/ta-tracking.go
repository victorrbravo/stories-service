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
	TaTrackingCreateQuery     = "INSERT INTO \"ta_tracking\" (\"fk_story\") VALUES ($1) RETURNING (\"fk_status\")"
	AllTaTrackingQuery        = "SELECT \"fk_status\", \"fk_story\" FROM \"ta_tracking\""
	TaTrackingByFkStatusQuery = "SELECT \"fk_status\", \"fk_story\" FROM \"ta_tracking\" WHERE \"fk_status\" = $1"
	TaTrackingUpdateQuery     = "UPDATE \"ta_tracking\" SET \"fk_story\" = COALESCE($1, \"fk_story\") WHERE \"fk_status\" = $2"
	TaTrackingDeleteQuery     = "DELETE FROM \"ta_tracking\" WHERE \"fk_status\" = $1"
)

type TaTracking struct {
	FkStatus types.String `json:"fk_status"`
	FkStory  types.Int    `json:"fk_story"`
}

func (taTracking TaTracking) ScanRow(rows database.RowInterface) (interface{}, error) {
	var err error
	err = rows.Scan(
		&taTracking.FkStatus,
		&taTracking.FkStory,
	)
	if err != nil {
		return nil, err
	}
	return taTracking, nil
}

func CreateTaTracking(
	query database.IQuery,
	fkStory types.Int,
) (*TaTracking, error) {
	var err error
	var scanner scanners.String
	var taTracking *TaTracking
	err = query.QueryOne(
		TaTrackingCreateQuery,
		&scanner,
		fkStory,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}
	taTracking = &TaTracking{
		FkStatus: scanner.Value,
		FkStory:  fkStory,
	}
	return taTracking, nil
}

func (taTracking *TaTracking) Update(
	query database.IQuery,
	fkStory types.Int,
) error {
	var err error
	err = query.Exec(
		TaTrackingUpdateQuery,
		fkStory,
		taTracking.FkStatus,
	)
	if err != nil {
		return err
	}
	taTracking.FkStory = fkStory
	return nil
}

func GetTaTrackingByFkStatus(query database.IQuery, value types.String) (*TaTracking, error) {
	var err error
	var taTracking TaTracking
	err = query.QueryOne(TaTrackingByFkStatusQuery, &taTracking, value)
	if err != nil {
		return nil, err
	}
	return &taTracking, nil
}

func (taTracking TaTracking) GetStatus() types.String {
	return taTracking.FkStatus
}

func (taTracking TaTracking) GetStory() types.Int {
	return taTracking.FkStory
}

func (taTracking *TaTracking) Delete(query database.IQuery) error {
	var err error
	err = query.Exec(
		TaTrackingDeleteQuery,
		taTracking.FkStatus,
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
