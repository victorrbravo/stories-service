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
	TaUserPolicyCreateQuery         = "INSERT INTO \"ta_user_policy\" (\"fk_user\",\"fk_policy\") VALUES ($1, $2) RETURNING (\"pk_user_policy\")"
	AllTaUserPolicyQuery            = "SELECT \"pk_user_policy\", \"fk_user\", \"fk_policy\" FROM \"ta_user_policy\""
	TaUserPolicyByPkUserPolicyQuery = "SELECT \"pk_user_policy\", \"fk_user\", \"fk_policy\" FROM \"ta_user_policy\" WHERE \"pk_user_policy\" = $1"
	TaUserPolicyUpdateQuery         = "UPDATE \"ta_user_policy\" SET \"fk_user\" = COALESCE($1, \"fk_user\"), \"fk_policy\" = COALESCE($2, \"fk_policy\") WHERE \"pk_user_policy\" = $3"
	TaUserPolicyDeleteQuery         = "DELETE FROM \"ta_user_policy\" WHERE \"pk_user_policy\" = $1"
	TaUserPolicyListOfFields        = " ta_user_policy.\"pk_user_policy\", ta_user_policy.\"fk_user\", ta_user_policy.\"fk_policy\""
)

type TaUserPolicy struct {
	PkUserPolicy types.Int `json:"pk_user_policy"`
	FkUser       types.Int `json:"fk_user"`
	FkPolicy     types.Int `json:"fk_policy"`
}

func (taUserPolicy TaUserPolicy) ListOfFields() string {
	return TaUserPolicyListOfFields
}

func (taUserPolicy TaUserPolicy) ScanRow(rows database.RowInterface) (interface{}, error) {
	var err error
	err = rows.Scan(
		&taUserPolicy.PkUserPolicy,
		&taUserPolicy.FkUser,
		&taUserPolicy.FkPolicy,
	)
	if err != nil {
		return nil, err
	}
	return taUserPolicy, nil
}

func CreateTaUserPolicy(
	query database.IQuery,
	fkUser types.Int,
	fkPolicy types.Int,
) (*TaUserPolicy, error) {
	var err error
	var scanner scanners.Int
	var taUserPolicy *TaUserPolicy
	_, err = query.QueryOne(
		TaUserPolicyCreateQuery,
		&scanner,
		fkUser,
		fkPolicy,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	taUserPolicy = &TaUserPolicy{
		PkUserPolicy: scanner.Value,
		FkUser:       fkUser,
		FkPolicy:     fkPolicy,
	}
	return taUserPolicy, nil
}

func (taUserPolicy *TaUserPolicy) UpdateByOne(
	query database.IQuery,
	options ...func() (string, any),
) error {
	var err error
	fkUser := types.NilInt
	fkPolicy := types.NilInt
	for _, option := range options {
		actionName, value := option()
		switch actionName {
		case "fkUser":
			fkUser = value.(types.Int)
		case "fkPolicy":
			fkPolicy = value.(types.Int)
		default:
			return errors.New("unsupported field")
		}
	}
	err = query.Exec(
		TaUserPolicyUpdateQuery,
		fkUser,
		fkPolicy,
		taUserPolicy.PkUserPolicy,
	)
	if err != nil {
		return err
	}
	return nil
}

func (taUserPolicy *TaUserPolicy) Update(
	query database.IQuery,
	fkUser types.Int,
	fkPolicy types.Int,
) error {
	var err error
	err = query.Exec(
		TaUserPolicyUpdateQuery,
		fkUser,
		fkPolicy,
		taUserPolicy.PkUserPolicy,
	)
	if err != nil {
		return err
	}
	taUserPolicy.FkUser = fkUser
	taUserPolicy.FkPolicy = fkPolicy
	return nil
}

func GetTaUserPolicyByPkUserPolicy(query database.IQuery, value types.Int) (*TaUserPolicy, error) {
	var err error
	var taUserPolicy TaUserPolicy
	var valueTaUserPolicy interface{}
	valueTaUserPolicy, err = query.QueryOne(TaUserPolicyByPkUserPolicyQuery, &taUserPolicy, value)
	if err != nil {
		return nil, err
	}
	taUserPolicy = valueTaUserPolicy.(TaUserPolicy)
	return &taUserPolicy, nil
}

func (taUserPolicy TaUserPolicy) GetID() types.Int {
	return taUserPolicy.PkUserPolicy
}

func (taUserPolicy TaUserPolicy) Getfk_user() types.Int {
	return taUserPolicy.FkUser
}

func (taUserPolicy TaUserPolicy) Getfk_policy() types.Int {
	return taUserPolicy.FkPolicy
}

func (taUserPolicy *TaUserPolicy) Delete(query database.IQuery) error {
	var err error
	err = query.Exec(
		TaUserPolicyDeleteQuery,
		taUserPolicy.PkUserPolicy,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetAllTaUserPolicies(query database.IQuery) ([]TaUserPolicy, error) {
	var err error
	var rows []interface{}
	var list []TaUserPolicy
	rows, err = query.QueryMany(AllTaUserPolicyQuery, &TaUserPolicy{})
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]TaUserPolicy, 0), nil
		}
		return nil, err
	}
	for _, row := range rows {
		list = append(list, row.(TaUserPolicy))
	}
	return list, nil
}
