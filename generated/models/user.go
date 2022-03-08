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
    UserCreateQuery = "INSERT INTO \"users\" (\"first_name\",\"last_name\",\"email\",\"account_type\") VALUES ($1, $2, $3, $4) RETURNING (\"id\")"
    AllUserQuery = "SELECT \"id\", \"first_name\", \"last_name\", \"email\", \"account_type\" FROM \"users\""
    UserByIDQuery = "SELECT \"id\", \"first_name\", \"last_name\", \"email\", \"account_type\" FROM \"users\" WHERE \"id\" = $1"
    UserByEmailQuery = "SELECT \"id\", \"first_name\", \"last_name\", \"email\", \"account_type\" FROM \"users\" WHERE \"email\" = $1"
    UserByAccountTypeQuery = "SELECT \"id\", \"first_name\", \"last_name\", \"email\", \"account_type\" FROM \"users\" WHERE \"account_type\" = $1"
    UserUpdateQuery = "UPDATE \"users\" SET \"first_name\" = COALESCE($1, \"first_name\"), \"last_name\" = COALESCE($2, \"last_name\"), \"email\" = COALESCE($3, \"email\"), \"account_type\" = COALESCE($4, \"account_type\") WHERE \"id\" = $5"
    UserDeleteQuery = "DELETE FROM \"users\" WHERE \"id\" = $1"
)

type User struct {
	ID types.Int `json:"id"`
	FirstName types.String `json:"first_name"`
	LastName types.String `json:"last_name"`
	Email types.String `json:"email"`
	AccountType types.String `json:"account_type"`
}

func (user *User) ScanRow(rows database.RowInterface) (interface{}, error) {
	var err error
	err = rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.AccountType,
	)
	if err != nil {
		return nil, err
	}
	return *user, nil
}

func CreateUser(
    query database.IQuery,
	firstName types.String,
	lastName types.String,
	email types.String,
	accountType types.String,
) (*User, error) {
	var err error
    var scanner scanners.Int
	var user *User
	err = query.QueryOne(
	    UserCreateQuery,
	    &scanner,
		firstName,
		lastName,
		email,
		accountType,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}
	user = &User{
	    ID: scanner.Value,
	    FirstName: firstName,
	    LastName: lastName,
	    Email: email,
	    AccountType: accountType,
	}
	return user, nil
}

func (user *User) Update(
    query database.IQuery,
    firstName types.String,
    lastName types.String,
    email types.String,
    accountType types.String,
) error {
    var err error
	err = query.Exec(
	    UserUpdateQuery,
        firstName,
        lastName,
        email,
        accountType,
        user.ID,
    )
	if err != nil {
    	return err
	}
    user.FirstName = firstName
    user.LastName = lastName
    user.Email = email
    user.AccountType = accountType
    return nil
}

func GetUserByID(query database.IQuery, value types.Int) (*User, error) {
	var err error
	var user User
	err = query.QueryOne(UserByIDQuery, &user, value)
	if err != nil {
		return  nil, err
	}
	return &user, nil
}




func GetUserByEmail(query database.IQuery, value types.String) (*User, error) {
	var err error
	var user User
	err = query.QueryOne(UserByEmailQuery, &user, value)
	if err != nil {
		return  nil, err
	}
	return &user, nil
}
func GetUserByAccountType(query database.IQuery, value types.String) (*User, error) {
	var err error
	var user User
	err = query.QueryOne(UserByAccountTypeQuery, &user, value)
	if err != nil {
		return  nil, err
	}
	return &user, nil
}



func (user *User) Delete(query database.IQuery) error {
	var err error
	err = query.Exec(
	    UserDeleteQuery,
        user.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetAllUsers(query database.IQuery) ([]User, error) {
	var err error
	var rows []interface{}
	var list []User
	rows, err = query.QueryMany(AllUserQuery, &User{})
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]User, 0), nil
		}
		return nil, err
	}
	for _, row := range rows {
		list = append(list, row.(User))
	}
	return list, nil
}
