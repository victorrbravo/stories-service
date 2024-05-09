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
	TaTrackingCreateQuery          = "INSERT INTO \"ta_tracking\" (\"fk_story\",\"fk_status\",\"summary\",\"amount_spent\",\"images\",\"videos\",\"roles\",\"updated_at\",\"fk_story_tips\") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING (\"pk_tracking\")"
	AllTaTrackingQuery             = "SELECT \"pk_tracking\", \"fk_story\", \"fk_status\", \"summary\", \"amount_spent\", \"images\", \"videos\", \"roles\", \"updated_at\", \"fk_story_tips\" FROM \"ta_tracking\""
	TaTrackingByPkTrackingQuery    = "SELECT \"pk_tracking\", \"fk_story\", \"fk_status\", \"summary\", \"amount_spent\", \"images\", \"videos\", \"roles\", \"updated_at\", \"fk_story_tips\" FROM \"ta_tracking\" WHERE \"pk_tracking\" = $1"
	TaTrackingByFkStoryQuery       = "SELECT \"pk_tracking\", \"fk_story\", \"fk_status\", \"summary\", \"amount_spent\", \"images\", \"videos\", \"roles\", \"updated_at\", \"fk_story_tips\" FROM \"ta_tracking\" WHERE \"fk_story\" = $1"
	TaTrackingByFkStoryDeleteQuery = "DELETE FROM \"ta_tracking\" WHERE \"fk_story\" = $1"
	TaTrackingByFkStoryTipsQuery   = "SELECT \"pk_tracking\", \"fk_story\", \"fk_status\", \"summary\", \"amount_spent\", \"images\", \"videos\", \"roles\", \"updated_at\", \"fk_story_tips\" FROM \"ta_tracking\" WHERE \"fk_story_tips\" = $1"
	TaTrackingUpdateQuery          = "UPDATE \"ta_tracking\" SET \"fk_story\" = COALESCE($1, \"fk_story\"), \"fk_status\" = COALESCE($2, \"fk_status\"), \"summary\" = COALESCE($3, \"summary\"), \"amount_spent\" = COALESCE($4, \"amount_spent\"), \"images\" = COALESCE($5, \"images\"), \"videos\" = COALESCE($6, \"videos\"), \"roles\" = COALESCE($7, \"roles\"), \"updated_at\" = COALESCE($8, \"updated_at\"), \"fk_story_tips\" = COALESCE($9, \"fk_story_tips\") WHERE \"pk_tracking\" = $10"
	TaTrackingDeleteQuery          = "DELETE FROM \"ta_tracking\" WHERE \"pk_tracking\" = $1"
	TaTrackingListOfFields         = " ta_tracking.\"pk_tracking\", ta_tracking.\"fk_story\", ta_tracking.\"fk_status\", ta_tracking.\"summary\", ta_tracking.\"amount_spent\", ta_tracking.\"images\", ta_tracking.\"videos\", ta_tracking.\"roles\", ta_tracking.\"updated_at\", ta_tracking.\"fk_story_tips\""
)

type TaTracking struct {
	PkTracking  types.Int         `json:"pk_tracking"`
	FkStory     types.Int         `json:"fk_story"`
	FkStatus    types.String      `json:"fk_status"`
	Summary     types.String      `json:"summary"`
	AmountSpent types.String      `json:"amount_spent"`
	Images      types.StringArray `json:"images"`
	Videos      types.StringArray `json:"videos"`
	Roles       types.StringArray `json:"roles"`
	UpdatedAt   types.DateTime    `json:"updated_at"`
	FkStoryTips types.Int         `json:"fk_story_tips"`
}

func (taTracking TaTracking) ListOfFields() string {
	return TaTrackingListOfFields
}

func (taTracking TaTracking) ScanRow(rows database.RowInterface) (interface{}, error) {
	var err error
	err = rows.Scan(
		&taTracking.PkTracking,
		&taTracking.FkStory,
		&taTracking.FkStatus,
		&taTracking.Summary,
		&taTracking.AmountSpent,
		&taTracking.Images,
		&taTracking.Videos,
		&taTracking.Roles,
		&taTracking.UpdatedAt,
		&taTracking.FkStoryTips,
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
	summary types.String,
	amountSpent types.String,
	images types.StringArray,
	videos types.StringArray,
	roles types.StringArray,
	updatedAt types.DateTime,
	fkStoryTips types.Int,
) (*TaTracking, error) {
	var err error
	var scanner scanners.Int
	var taTracking *TaTracking
	_, err = query.QueryOne(
		TaTrackingCreateQuery,
		&scanner,
		fkStory,
		fkStatus,
		summary,
		amountSpent,
		images,
		videos,
		roles,
		updatedAt,
		fkStoryTips,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	taTracking = &TaTracking{
		PkTracking:  scanner.Value,
		FkStory:     fkStory,
		FkStatus:    fkStatus,
		Summary:     summary,
		AmountSpent: amountSpent,
		Images:      images,
		Videos:      videos,
		Roles:       roles,
		UpdatedAt:   updatedAt,
		FkStoryTips: fkStoryTips,
	}
	return taTracking, nil
}

func (taTracking *TaTracking) UpdateByOne(query database.IQuery, options ...func() (string, any)) error {
	fkStory := types.NilInt
	fkStatus := types.NilString()
	summary := types.NilString()
	amountSpent := types.NilString()
	images := types.NilStringArray
	videos := types.NilStringArray
	roles := types.NilStringArray
	updatedAt := types.NilDateTime
	fkStoryTips := types.NilInt
	for _, option := range options {
		name, value := option()
		switch name {
		case "fkStory":
			fkStory = value.(types.Int)
			taTracking.FkStory = fkStory
		case "fkStatus":
			fkStatus = value.(types.String)
			taTracking.FkStatus = fkStatus
		case "summary":
			summary = value.(types.String)
			taTracking.Summary = summary
		case "amountSpent":
			amountSpent = value.(types.String)
			taTracking.AmountSpent = amountSpent
		case "images":
			images = value.(types.StringArray)
			taTracking.Images = images
		case "videos":
			videos = value.(types.StringArray)
			taTracking.Videos = videos
		case "roles":
			roles = value.(types.StringArray)
			taTracking.Roles = roles
		case "updatedAt":
			updatedAt = value.(types.DateTime)
			taTracking.UpdatedAt = updatedAt
		case "fkStoryTips":
			fkStoryTips = value.(types.Int)
			taTracking.FkStoryTips = fkStoryTips
		default:
			return errors.New("unsupported field")
		}
	}
	err := query.Exec(
		TaTrackingUpdateQuery,
		fkStory,
		fkStatus,
		summary,
		amountSpent,
		images,
		videos,
		roles,
		updatedAt,
		fkStoryTips,
		taTracking.PkTracking,
	)
	if err != nil {
		return err
	}
	return nil
}

func (taTracking *TaTracking) Update(
	query database.IQuery,
	fkStory types.Int,
	fkStatus types.String,
	summary types.String,
	amountSpent types.String,
	images types.StringArray,
	videos types.StringArray,
	roles types.StringArray,
	updatedAt types.DateTime,
	fkStoryTips types.Int,
) error {
	var err error
	err = query.Exec(
		TaTrackingUpdateQuery,
		fkStory,
		fkStatus,
		summary,
		amountSpent,
		images,
		videos,
		roles,
		updatedAt,
		fkStoryTips,
		taTracking.PkTracking,
	)
	if err != nil {
		return err
	}
	taTracking.FkStory = fkStory
	taTracking.FkStatus = fkStatus
	taTracking.Summary = summary
	taTracking.AmountSpent = amountSpent
	taTracking.Images = images
	taTracking.Videos = videos
	taTracking.Roles = roles
	taTracking.UpdatedAt = updatedAt
	taTracking.FkStoryTips = fkStoryTips
	return nil
}

func GetTaTrackingByPkTracking(query database.IQuery, value types.Int) (*TaTracking, error) {
	var err error
	var taTracking TaTracking
	var valueTaTracking interface{}
	valueTaTracking, err = query.QueryOne(TaTrackingByPkTrackingQuery, &taTracking, value)
	if err != nil {
		return nil, err
	}
	taTracking = valueTaTracking.(TaTracking)
	return &taTracking, nil
}

func GetTaTrackingByFkStory(query database.IQuery, value types.Int) ([]TaTracking, error) {
	var err error
	var rows []interface{}
	var list []TaTracking
	list = make([]TaTracking, 0)
	rows, err = query.QueryMany(TaTrackingByFkStoryQuery, &TaTracking{}, value)
	if err != nil {
		return list, err
	}
	for _, row := range rows {
		list = append(list, row.(TaTracking))
	}
	return list, nil
}

func DeleteTaTrackingByFkStory(query database.IQuery, value types.Int) error {
	var err error
	err = query.Exec(TaTrackingByFkStoryDeleteQuery, value)
	if err != nil {
		return err
	}
	return nil
}

func GetTaTrackingByFkStoryTips(query database.IQuery, value types.Int) ([]TaTracking, error) {
	var err error
	var rows []interface{}
	var list []TaTracking
	rows, err = query.QueryMany(TaTrackingByFkStoryTipsQuery, &TaTracking{}, value)
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

func (taTracking TaTracking) GetPkTracking() types.Int {
	return taTracking.PkTracking
}

func (taTracking TaTracking) GetFkStory() types.Int {
	return taTracking.FkStory
}

func (taTracking TaTracking) GetFkStatus() types.String {
	return taTracking.FkStatus
}

func (taTracking TaTracking) GetSummary() types.String {
	return taTracking.Summary
}

func (taTracking TaTracking) GetAmountSpent() types.String {
	return taTracking.AmountSpent
}

func (taTracking TaTracking) GetImages() types.StringArray {
	return taTracking.Images
}

func (taTracking TaTracking) GetVideos() types.StringArray {
	return taTracking.Videos
}

func (taTracking TaTracking) GetRoles() types.StringArray {
	return taTracking.Roles
}

func (taTracking TaTracking) GetUpdatedAt() types.DateTime {
	return taTracking.UpdatedAt
}

func (taTracking TaTracking) GetFkStoryTips() types.Int {
	return taTracking.FkStoryTips
}

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
