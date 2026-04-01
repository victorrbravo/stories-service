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
	TaUserRelationshipCreateQuery                        = "INSERT INTO \"ta_user_relationship\" (\"fk_follower_id\",\"fk_followed_id\",\"status\",\"updated_at\") VALUES ($1, $2, $3, $4) RETURNING (\"pk_user_relationship\")"
	AllTaUserRelationshipQuery                           = "SELECT \"pk_user_relationship\", \"fk_follower_id\", \"fk_followed_id\", \"status\", \"updated_at\" FROM \"ta_user_relationship\""
	TaUserRelationshipByPkUserRelationshipQuery          = "SELECT \"pk_user_relationship\", \"fk_follower_id\", \"fk_followed_id\", \"status\", \"updated_at\" FROM \"ta_user_relationship\" WHERE \"pk_user_relationship\" = $1"
	TaUserRelationshipUpdateQuery                        = "UPDATE \"ta_user_relationship\" SET \"fk_follower_id\" = COALESCE($1, \"fk_follower_id\"), \"fk_followed_id\" = COALESCE($2, \"fk_followed_id\"), \"status\" = COALESCE($3, \"status\"), \"updated_at\" = COALESCE($4, \"updated_at\") WHERE \"pk_user_relationship\" = $5"
	TaUserRelationshipDeleteQuery                        = "DELETE FROM \"ta_user_relationship\" WHERE \"pk_user_relationship\" = $1"
	TaUserRelationshipByFkFollowerIdAndFkFollowedIdQuery = "SELECT \"pk_user_relationship\", \"fk_follower_id\", \"fk_followed_id\", \"status\", \"updated_at\" FROM \"ta_user_relationship\" WHERE \"fk_follower_id\" = $1 AND \"fk_followed_id\" = $2 "
	TaUserRelationshipListOfFields                       = " ta_user_relationship.\"pk_user_relationship\", ta_user_relationship.\"fk_follower_id\", ta_user_relationship.\"fk_followed_id\", ta_user_relationship.\"status\", ta_user_relationship.\"updated_at\""
)

type TaUserRelationship struct {
	PkUserRelationship types.Int      `json:"pk_user_relationship"`
	FkFollowerId       types.Int      `json:"fk_follower_id"`
	FkFollowedId       types.Int      `json:"fk_followed_id"`
	Status             types.String   `json:"status"`
	UpdatedAt          types.DateTime `json:"updated_at"`
}

func (taUserRelationship TaUserRelationship) ListOfFields() string {
	return TaUserRelationshipListOfFields
}

func (taUserRelationship TaUserRelationship) ScanRow(rows database.RowInterface) (interface{}, error) {
	var err error
	err = rows.Scan(
		&taUserRelationship.PkUserRelationship,
		&taUserRelationship.FkFollowerId,
		&taUserRelationship.FkFollowedId,
		&taUserRelationship.Status,
		&taUserRelationship.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return taUserRelationship, nil
}

func CreateTaUserRelationship(
	query database.IQuery,
	fkFollowerId types.Int,
	fkFollowedId types.Int,
	status types.String,
	updatedAt types.DateTime,
) (*TaUserRelationship, error) {
	var err error
	var scanner scanners.Int
	var taUserRelationship *TaUserRelationship
	_, err = query.QueryOne(
		TaUserRelationshipCreateQuery,
		&scanner,
		fkFollowerId,
		fkFollowedId,
		status,
		updatedAt,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	taUserRelationship = &TaUserRelationship{
		PkUserRelationship: scanner.Value,
		FkFollowerId:       fkFollowerId,
		FkFollowedId:       fkFollowedId,
		Status:             status,
		UpdatedAt:          updatedAt,
	}
	return taUserRelationship, nil
}

func (taUserRelationship *TaUserRelationship) UpdateByOne(
	query database.IQuery,
	options ...func() (string, any),
) error {
	var err error
	fkFollowerId := types.NilInt
	fkFollowedId := types.NilInt
	status := types.NilString()
	updatedAt := types.NilDateTime
	for _, option := range options {
		actionName, value := option()
		switch actionName {
		case "fkFollowerId":
			fkFollowerId = value.(types.Int)
		case "fkFollowedId":
			fkFollowedId = value.(types.Int)
		case "status":
			status = value.(types.String)
		case "updatedAt":
			updatedAt = value.(types.DateTime)
		default:
			return errors.New("unsupported field")
		}
	}
	err = query.Exec(
		TaUserRelationshipUpdateQuery,
		fkFollowerId,
		fkFollowedId,
		status,
		updatedAt,
		taUserRelationship.PkUserRelationship,
	)
	if err != nil {
		return err
	}
	return nil
}

func (taUserRelationship *TaUserRelationship) Update(
	query database.IQuery,
	fkFollowerId types.Int,
	fkFollowedId types.Int,
	status types.String,
	updatedAt types.DateTime,
) error {
	var err error
	err = query.Exec(
		TaUserRelationshipUpdateQuery,
		fkFollowerId,
		fkFollowedId,
		status,
		updatedAt,
		taUserRelationship.PkUserRelationship,
	)
	if err != nil {
		return err
	}
	taUserRelationship.FkFollowerId = fkFollowerId
	taUserRelationship.FkFollowedId = fkFollowedId
	taUserRelationship.Status = status
	taUserRelationship.UpdatedAt = updatedAt
	return nil
}

func GetTaUserRelationshipByPkUserRelationship(query database.IQuery, value types.Int) (*TaUserRelationship, error) {
	var err error
	var taUserRelationship TaUserRelationship
	var valueTaUserRelationship interface{}
	valueTaUserRelationship, err = query.QueryOne(TaUserRelationshipByPkUserRelationshipQuery, &taUserRelationship, value)
	if err != nil {
		return nil, err
	}
	taUserRelationship = valueTaUserRelationship.(TaUserRelationship)
	return &taUserRelationship, nil
}

func GetTaUserRelationshipByFkFollowerIdAndFkFollowedId(query database.IQuery,
	fkFollowerId types.Int,
	fkFollowedId types.Int,
) ([]TaUserRelationship, error) {
	var err error
	var rows []interface{}
	var list []TaUserRelationship
	rows, err = query.QueryMany(TaUserRelationshipByFkFollowerIdAndFkFollowedIdQuery, &TaUserRelationship{},
		fkFollowerId,
		fkFollowedId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]TaUserRelationship, 0), nil
		}
		return nil, err
	}
	for _, row := range rows {
		list = append(list, row.(TaUserRelationship))
	}
	return list, nil
}

func GetTaUserRelationshipsByFkFollowerIdAndFkFollowedIdPaged(query database.IQuery,
	fkFollowerId types.Int,
	fkFollowedId types.Int,
	pageBase types.PageBase) (types.Page, error) {
	var err error
	var page types.Page
	page, err = query.QueryManyPaged(TaUserRelationshipByFkFollowerIdAndFkFollowedIdQuery, &TaUserRelationship{},
		pageBase.PageNumber,
		pageBase.PageSize,
		fkFollowerId,
		fkFollowedId,
	)
	if err != nil {
		return page, err
	}
	return page, nil
}

func (taUserRelationship TaUserRelationship) GetID() types.Int {
	return taUserRelationship.PkUserRelationship
}

func (taUserRelationship TaUserRelationship) Getfollower_id() types.Int {
	return taUserRelationship.FkFollowerId
}

func (taUserRelationship TaUserRelationship) Getfollowed_id() types.Int {
	return taUserRelationship.FkFollowedId
}

func (taUserRelationship TaUserRelationship) Getstatus() types.String {
	return taUserRelationship.Status
}

func (taUserRelationship TaUserRelationship) GetUpdatedAt() types.DateTime {
	return taUserRelationship.UpdatedAt
}

func (taUserRelationship *TaUserRelationship) Delete(query database.IQuery) error {
	var err error
	err = query.Exec(
		TaUserRelationshipDeleteQuery,
		taUserRelationship.PkUserRelationship,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetAllTaUserRelationships(query database.IQuery) ([]TaUserRelationship, error) {
	var err error
	var rows []interface{}
	var list []TaUserRelationship
	rows, err = query.QueryMany(AllTaUserRelationshipQuery, &TaUserRelationship{})
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]TaUserRelationship, 0), nil
		}
		return nil, err
	}
	for _, row := range rows {
		list = append(list, row.(TaUserRelationship))
	}
	return list, nil
}
