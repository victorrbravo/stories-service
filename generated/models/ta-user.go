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
    TaUserCreateQuery = "INSERT INTO \"ta_user\" (\"first_name\",\"last_name\",\"email\",\"password\",\"account_type\",\"activated\") VALUES ($1, $2, $3, $4, $5, $6) RETURNING (\"pk_user\")"
    AllTaUserQuery = "SELECT \"pk_user\", \"first_name\", \"last_name\", \"email\", \"password\", \"account_type\", \"activated\" FROM \"ta_user\""
    TaUserByPkUserQuery = "SELECT \"pk_user\", \"first_name\", \"last_name\", \"email\", \"password\", \"account_type\", \"activated\" FROM \"ta_user\" WHERE \"pk_user\" = $1"
    TaUserByEmailQuery = "SELECT \"pk_user\", \"first_name\", \"last_name\", \"email\", \"password\", \"account_type\", \"activated\" FROM \"ta_user\" WHERE \"email\" = $1"
    TaUserUpdateQuery = "UPDATE \"ta_user\" SET \"first_name\" = COALESCE($1, \"first_name\"), \"last_name\" = COALESCE($2, \"last_name\"), \"email\" = COALESCE($3, \"email\"), \"password\" = COALESCE($4, \"password\"), \"account_type\" = COALESCE($5, \"account_type\"), \"activated\" = COALESCE($6, \"activated\") WHERE \"pk_user\" = $7"
    TaUserDeleteQuery = "DELETE FROM \"ta_user\" WHERE \"pk_user\" = $1"
    TaUserListOfFields = " ta_user.\"pk_user\", ta_user.\"first_name\", ta_user.\"last_name\", ta_user.\"email\", ta_user.\"password\", ta_user.\"account_type\", ta_user.\"activated\""
)

type TaUser struct {
	PkUser types.Int `json:"pk_user"`
	FirstName types.String `json:"first_name"`
	LastName types.String `json:"last_name"`
	Email types.String `json:"email"`
	Password types.String `json:"password"`
	AccountType types.String `json:"account_type"`
	Activated types.Bool `json:"activated"`
}

func (taUser TaUser) ListOfFields() string {
	return TaUserListOfFields
}


func (taUser TaUser) ScanRow(rows database.RowInterface) (interface{}, error) {
	var err error
	err = rows.Scan(
		&taUser.PkUser,
		&taUser.FirstName,
		&taUser.LastName,
		&taUser.Email,
		&taUser.Password,
		&taUser.AccountType,
		&taUser.Activated,
	)
	if err != nil {
		return nil, err
	}
	return taUser, nil
}

func CreateTaUser(
    query database.IQuery,
	firstName types.String,
	lastName types.String,
	email types.String,
	password types.String,
	accountType types.String,
	activated types.Bool,
) (*TaUser, error) {
	var err error
    var scanner scanners.Int
	var taUser *TaUser
	_, err = query.QueryOne(
	    TaUserCreateQuery,
	    &scanner,
		firstName,
		lastName,
		email,
		password,
		accountType,
		activated,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	taUser = &TaUser{
	    PkUser: scanner.Value,
	    FirstName: firstName,
	    LastName: lastName,
	    Email: email,
	    Password: password,
	    AccountType: accountType,
	    Activated: activated,
	}
	return taUser, nil
}

func (taUser *TaUser) Update(
    query database.IQuery,
    firstName types.String,
    lastName types.String,
    email types.String,
    password types.String,
    accountType types.String,
    activated types.Bool,
) error {
    var err error
	err = query.Exec(
	    TaUserUpdateQuery,
        firstName,
        lastName,
        email,
        password,
        accountType,
        activated,
        taUser.PkUser,
    )
	if err != nil {
    	return err
	}
    taUser.FirstName = firstName
    taUser.LastName = lastName
    taUser.Email = email
    taUser.Password = password
    taUser.AccountType = accountType
    taUser.Activated = activated
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




func GetTaUserByEmail(query database.IQuery, value types.String) ([]TaUser, error) {
	    var err error
    	var rows []interface{}
    	var list []TaUser
    	rows, err = query.QueryMany(TaUserByEmailQuery, &TaUser{}, value)
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




func (taUser TaUser) GetID() types.Int {
	return taUser.PkUser}


func (taUser TaUser) GetFullName() types.String {
	return taUser.FirstName}

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
