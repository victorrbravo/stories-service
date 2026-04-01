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
	TaSessionCreateQuery      = "INSERT INTO \"ta_session\" (\"token\",\"fk_user\",\"started_at\",\"invalid\") VALUES ($1, $2, $3, $4) RETURNING (\"pk_session\")"
	AllTaSessionQuery         = "SELECT \"pk_session\", \"token\", \"fk_user\", \"started_at\", \"invalid\" FROM \"ta_session\""
	TaSessionByPkSessionQuery = "SELECT \"pk_session\", \"token\", \"fk_user\", \"started_at\", \"invalid\" FROM \"ta_session\" WHERE \"pk_session\" = $1"
	TaSessionByTokenQuery     = "SELECT \"pk_session\", \"token\", \"fk_user\", \"started_at\", \"invalid\" FROM \"ta_session\" WHERE \"token\" = $1"
	TaSessionByInvalidQuery   = "SELECT \"pk_session\", \"token\", \"fk_user\", \"started_at\", \"invalid\" FROM \"ta_session\" WHERE \"invalid\" = $1"
	TaSessionUpdateQuery      = "UPDATE \"ta_session\" SET \"token\" = COALESCE($1, \"token\"), \"fk_user\" = COALESCE($2, \"fk_user\"), \"started_at\" = COALESCE($3, \"started_at\"), \"invalid\" = COALESCE($4, \"invalid\") WHERE \"pk_session\" = $5"
	TaSessionDeleteQuery      = "DELETE FROM \"ta_session\" WHERE \"pk_session\" = $1"
	TaSessionListOfFields     = " ta_session.\"pk_session\", ta_session.\"token\", ta_session.\"fk_user\", ta_session.\"started_at\", ta_session.\"invalid\""
)

type TaSession struct {
	PkSession types.Int      `json:"pk_session"`
	Token     types.String   `json:"token"`
	FkUser    types.Int      `json:"fk_user"`
	StartedAt types.DateTime `json:"started_at"`
	Invalid   types.Bool     `json:"invalid"`
}

func (taSession TaSession) ListOfFields() string {
	return TaSessionListOfFields
}

func (taSession TaSession) ScanRow(rows database.RowInterface) (interface{}, error) {
	var err error
	err = rows.Scan(
		&taSession.PkSession,
		&taSession.Token,
		&taSession.FkUser,
		&taSession.StartedAt,
		&taSession.Invalid,
	)
	if err != nil {
		return nil, err
	}
	return taSession, nil
}

func CreateTaSession(
	query database.IQuery,
	token types.String,
	fkUser types.Int,
	startedAt types.DateTime,
	invalid types.Bool,
) (*TaSession, error) {
	var err error
	var scanner scanners.Int
	var taSession *TaSession
	_, err = query.QueryOne(
		TaSessionCreateQuery,
		&scanner,
		token,
		fkUser,
		startedAt,
		invalid,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	taSession = &TaSession{
		PkSession: scanner.Value,
		Token:     token,
		FkUser:    fkUser,
		StartedAt: startedAt,
		Invalid:   invalid,
	}
	return taSession, nil
}

func (taSession *TaSession) UpdateByOne(
	query database.IQuery,
	options ...func() (string, any),
) error {
	var err error
	token := types.NilString()
	fkUser := types.NilInt
	startedAt := types.NilDateTime
	invalid := types.NilBool
	for _, option := range options {
		actionName, value := option()
		switch actionName {
		case "token":
			token = value.(types.String)
		case "fkUser":
			fkUser = value.(types.Int)
		case "startedAt":
			startedAt = value.(types.DateTime)
		case "invalid":
			invalid = value.(types.Bool)
		default:
			return errors.New("unsupported field")
		}
	}
	err = query.Exec(
		TaSessionUpdateQuery,
		token,
		fkUser,
		startedAt,
		invalid,
		taSession.PkSession,
	)
	if err != nil {
		return err
	}
	return nil
}

func (taSession *TaSession) Update(
	query database.IQuery,
	token types.String,
	fkUser types.Int,
	startedAt types.DateTime,
	invalid types.Bool,
) error {
	var err error
	err = query.Exec(
		TaSessionUpdateQuery,
		token,
		fkUser,
		startedAt,
		invalid,
		taSession.PkSession,
	)
	if err != nil {
		return err
	}
	taSession.Token = token
	taSession.FkUser = fkUser
	taSession.StartedAt = startedAt
	taSession.Invalid = invalid
	return nil
}

func GetTaSessionByPkSession(query database.IQuery, value types.Int) (*TaSession, error) {
	var err error
	var taSession TaSession
	var valueTaSession interface{}
	valueTaSession, err = query.QueryOne(TaSessionByPkSessionQuery, &taSession, value)
	if err != nil {
		return nil, err
	}
	taSession = valueTaSession.(TaSession)
	return &taSession, nil
}

func GetTaSessionByToken(query database.IQuery, value types.String) ([]TaSession, error) {
	var err error
	var rows []interface{}
	var list []TaSession
	rows, err = query.QueryMany(TaSessionByTokenQuery, &TaSession{}, value)
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]TaSession, 0), nil
		}
		return nil, err
	}
	for _, row := range rows {
		list = append(list, row.(TaSession))
	}
	return list, nil
}
func GetTaSessionByInvalid(query database.IQuery, value types.Bool) ([]TaSession, error) {
	var err error
	var rows []interface{}
	var list []TaSession
	rows, err = query.QueryMany(TaSessionByInvalidQuery, &TaSession{}, value)
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]TaSession, 0), nil
		}
		return nil, err
	}
	for _, row := range rows {
		list = append(list, row.(TaSession))
	}
	return list, nil
}

func (taSession TaSession) GetPkSession() types.Int {
	return taSession.PkSession
}

func (taSession TaSession) GetToken() types.String {
	return taSession.Token
}

func (taSession TaSession) GetFkUser() types.Int {
	return taSession.FkUser
}

func (taSession TaSession) GetStartedAt() types.DateTime {
	return taSession.StartedAt
}

func (taSession TaSession) GetInvalid() types.Bool {
	return taSession.Invalid
}

func (taSession *TaSession) Delete(query database.IQuery) error {
	var err error
	err = query.Exec(
		TaSessionDeleteQuery,
		taSession.PkSession,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetAllTaSessions(query database.IQuery) ([]TaSession, error) {
	var err error
	var rows []interface{}
	var list []TaSession
	rows, err = query.QueryMany(AllTaSessionQuery, &TaSession{})
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]TaSession, 0), nil
		}
		return nil, err
	}
	for _, row := range rows {
		list = append(list, row.(TaSession))
	}
	return list, nil
}
