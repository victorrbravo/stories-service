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
    TaUserCreateQuery = "INSERT INTO \"ta_user\" (\"firstname\",\"lastname\",\"email\",\"password\",\"account\") VALUES ($1, $2, $3, $4, $5) RETURNING (\"pk_user\")"
    AllTaUserQuery = "SELECT \"pk_user\", \"firstname\", \"lastname\", \"email\", \"password\", \"account\" FROM \"ta_user\""
    TaUserByPkUserQuery = "SELECT \"pk_user\", \"firstname\", \"lastname\", \"email\", \"password\", \"account\" FROM \"ta_user\" WHERE \"pk_user\" = $1"
    TaUserByEmailQuery = "SELECT \"pk_user\", \"firstname\", \"lastname\", \"email\", \"password\", \"account\" FROM \"ta_user\" WHERE \"email\" = $1"
    TaUserByAccountQuery = "SELECT \"pk_user\", \"firstname\", \"lastname\", \"email\", \"password\", \"account\" FROM \"ta_user\" WHERE \"account\" = $1"
    TaUserUpdateQuery = "UPDATE \"ta_user\" SET \"firstname\" = COALESCE($1, \"firstname\"), \"lastname\" = COALESCE($2, \"lastname\"), \"email\" = COALESCE($3, \"email\"), \"password\" = COALESCE($4, \"password\"), \"account\" = COALESCE($5, \"account\") WHERE \"pk_user\" = $6"
    TaUserDeleteQuery = "DELETE FROM \"ta_user\" WHERE \"pk_user\" = $1"
    TaUserListOfFields = " ta_user.\"pk_user\", ta_user.\"firstname\", ta_user.\"lastname\", ta_user.\"email\", ta_user.\"password\", ta_user.\"account\""
)

type TaUser struct {
	PkUser types.Int `json:"pk_user"`
	Firstname types.String `json:"firstname"`
	Lastname types.String `json:"lastname"`
	Email types.String `json:"email"`
	Password types.String `json:"password"`
	Account types.String `json:"account"`
}

func (taUser TaUser) ListOfFields() string {
	return TaUserListOfFields
}


func (taUser TaUser) ScanRow(rows database.RowInterface) (interface{}, error) {
	var err error
	err = rows.Scan(
		&taUser.PkUser,
		&taUser.Firstname,
		&taUser.Lastname,
		&taUser.Email,
		&taUser.Password,
		&taUser.Account,
	)
	if err != nil {
		return nil, err
	}
	return taUser, nil
}

func CreateTaUser(
    query database.IQuery,
	firstname types.String,
	lastname types.String,
	email types.String,
	password types.String,
	account types.String,
) (*TaUser, error) {
	var err error
    var scanner scanners.Int
	var taUser *TaUser
	_, err = query.QueryOne(
	    TaUserCreateQuery,
	    &scanner,
		firstname,
		lastname,
		email,
		password,
		account,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	taUser = &TaUser{
	    PkUser: scanner.Value,
	    Firstname: firstname,
	    Lastname: lastname,
	    Email: email,
	    Password: password,
	    Account: account,
	}
	return taUser, nil
}

func (taUser *TaUser) Update(
    query database.IQuery,
    firstname types.String,
    lastname types.String,
    email types.String,
    password types.String,
    account types.String,
) error {
    var err error
	err = query.Exec(
	    TaUserUpdateQuery,
        firstname,
        lastname,
        email,
        password,
        account,
        taUser.PkUser,
    )
	if err != nil {
    	return err
	}
    taUser.Firstname = firstname
    taUser.Lastname = lastname
    taUser.Email = email
    taUser.Password = password
    taUser.Account = account
    return nil
}

func GetTaUserByPkUser(query database.IQuery, value types.Int) (*TaUser, error) {
	var err error
	var taUser TaUser
	var valueTaUser interface{}
	valueTaUser, err = query.QueryOne(TaUserByPkUserQuery, &taUser, value)
	if err != nil {
		return  nil, err
	}
    taUser = valueTaUser.(TaUser)
	return &taUser, nil
}




func (taUser TaUser) GetTaUserByEmail(query database.IQuery, value types.String) (interface{}, error) {
	var err error
	var valueTaUser interface{}
	valueTaUser, err = query.QueryOne(TaUserByEmailQuery, &taUser, value)
	if err != nil {
		return  nil, err
	}
	taUser = valueTaUser.(TaUser)
	return &taUser, nil
}
func (taUser TaUser) GetTaUserByAccount(query database.IQuery, value types.String) (interface{}, error) {
	var err error
	var valueTaUser interface{}
	valueTaUser, err = query.QueryOne(TaUserByAccountQuery, &taUser, value)
	if err != nil {
		return  nil, err
	}
	taUser = valueTaUser.(TaUser)
	return &taUser, nil
}




func (taUser TaUser) GetID() types.Int {
	return taUser.PkUser}


func (taUser TaUser) GetFullName() types.String {
	return taUser.Firstname}

func (taUser TaUser) GetEmail() types.String {
	return taUser.Email}

func (taUser TaUser) GetPassword() types.String {
	return taUser.Password}


func (taUser *TaUser) Delete(query database.IQuery) error {
	var err error
	err = query.Exec(
	    TaUserDeleteQuery,
        taUser.PkUser,
	)
	if err != nil {
		return err
	}
	return nil
}


func GetAllTaUsers(query database.IQuery) ([]TaUser, error) {
	var err error
	var rows []interface{}
	var list []TaUser
	rows, err = query.QueryMany(AllTaUserQuery, &TaUser{})
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]TaUser, 0), nil
		}
		return nil, err
	}
	for _, row := range rows {
		list = append(list, row.(TaUser))
	}
	return list, nil
}
