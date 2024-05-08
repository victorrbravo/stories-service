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
    TaUserCreateQuery = "INSERT INTO \"ta_user\" (\"firstname\",\"lastname\",\"email\",\"password\",\"account_type\",\"activated\",\"username\") VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING (\"pk_user\")"
    AllTaUserQuery = "SELECT \"pk_user\", \"firstname\", \"lastname\", \"email\", \"password\", \"account_type\", \"activated\", \"username\" FROM \"ta_user\""
    TaUserByPkUserQuery = "SELECT \"pk_user\", \"firstname\", \"lastname\", \"email\", \"password\", \"account_type\", \"activated\", \"username\" FROM \"ta_user\" WHERE \"pk_user\" = $1"
    TaUserByEmailQuery = "SELECT \"pk_user\", \"firstname\", \"lastname\", \"email\", \"password\", \"account_type\", \"activated\", \"username\" FROM \"ta_user\" WHERE \"email\" = $1"
    TaUserByUsernameQuery = "SELECT \"pk_user\", \"firstname\", \"lastname\", \"email\", \"password\", \"account_type\", \"activated\", \"username\" FROM \"ta_user\" WHERE \"username\" = $1"
    TaUserUpdateQuery = "UPDATE \"ta_user\" SET \"firstname\" = COALESCE($1, \"firstname\"), \"lastname\" = COALESCE($2, \"lastname\"), \"email\" = COALESCE($3, \"email\"), \"password\" = COALESCE($4, \"password\"), \"account_type\" = COALESCE($5, \"account_type\"), \"activated\" = COALESCE($6, \"activated\"), \"username\" = COALESCE($7, \"username\") WHERE \"pk_user\" = $8"
    TaUserDeleteQuery = "DELETE FROM \"ta_user\" WHERE \"pk_user\" = $1"
    TaUserListOfFields = " ta_user.\"pk_user\", ta_user.\"firstname\", ta_user.\"lastname\", ta_user.\"email\", ta_user.\"password\", ta_user.\"account_type\", ta_user.\"activated\", ta_user.\"username\""
)

type TaUser struct {
	PkUser types.Int `json:"pk_user"`
	Firstname types.String `json:"firstname"`
	Lastname types.String `json:"lastname"`
	Email types.String `json:"email"`
	Password types.String `json:"password"`
	AccountType types.String `json:"account_type"`
	Activated types.Bool `json:"activated"`
	Username types.String `json:"username"`
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
		&taUser.AccountType,
		&taUser.Activated,
		&taUser.Username,
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
	accountType types.String,
	activated types.Bool,
	username types.String,
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
		accountType,
		activated,
		username,
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
	    AccountType: accountType,
	    Activated: activated,
	    Username: username,
	}
	return taUser, nil
}

func (taUser *TaUser) Update(
    query database.IQuery,
    firstname types.String,
    lastname types.String,
    email types.String,
    password types.String,
    accountType types.String,
    activated types.Bool,
    username types.String,
) error {
    var err error
	err = query.Exec(
	    TaUserUpdateQuery,
        firstname,
        lastname,
        email,
        password,
        accountType,
        activated,
        username,
        taUser.PkUser,
    )
	if err != nil {
    	return err
	}
    taUser.Firstname = firstname
    taUser.Lastname = lastname
    taUser.Email = email
    taUser.Password = password
    taUser.AccountType = accountType
    taUser.Activated = activated
    taUser.Username = username
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
func GetTaUserByUsername(query database.IQuery, value types.String) ([]TaUser, error) {
	    var err error
    	var rows []interface{}
    	var list []TaUser
    	rows, err = query.QueryMany(TaUserByUsernameQuery, &TaUser{}, value)
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
