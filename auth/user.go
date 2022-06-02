package auth

import (
	"database/sql"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"stories/crypto"
	"stories/database"
	"stories/generated/queries"
	"stories/helpers"
	"stories/scanners"
	"stories/types"
)

type AccountType string

const (
	Public  AccountType = "Public"
	Private             = "Private"
)

type User struct {
	ID                   types.UUID                 `json:"id"`
	FirstName            types.String               `json:"firstname"`
	LastName             types.String               `json:"lastname"`
	Email                types.String               `json:"email"`
	AccountType          AccountType                `json:"account_type,omitempty"`
	CreatedAt            time.Time                  `json:"create_at"`
	Password             crypto.Password            `json:"-"`
	Activated            types.Bool                 `json:"activated"`
	ProfessionalCategory types.ProfessionalCategory `json:"professional_category"`
	Message              types.String               `json:"message"`
	PreRegistration      bool                       `json:"pre_registration"`
	authToken            crypto.Token
}

type BasicWebUser struct {
	ID        string       `json:"id"`
	FirstName types.String `json:"firstname"`
	LastName  types.String `json:"lastname"`
	Email     types.String `json:"email"`
}

type Registration struct {
	PreRegistration types.Bool `json:"pre_registration"`
}

func (user *Registration) ScanRow(rows database.RowInterface) (interface{}, error) {
	var err error
	// We scan into thw newly allocated object
	err = rows.Scan(
		&user.PreRegistration,
	)
	if err != nil {
		return User{}, err
	}
	return *user, nil
}

type UserStatus struct {
	ID              string       `json:"id"`
	Email           types.String `json:"email"`
	Activated       types.Bool   `json:"activated"`
	PreRegistration types.Bool   `json:"pre_registration"`
	AccountType     types.String `json:"account_type"`
}

func (user *UserStatus) ScanRow(rows database.RowInterface) (interface{}, error) {
	var err error
	// We scan into thw newly allocated object
	err = rows.Scan(
		&user.ID,
		&user.Email,
		&user.Activated,
		&user.PreRegistration,
		&user.AccountType,
	)
	if err != nil {
		return User{}, err
	}
	return *user, nil
}

// ScanRow implements the database.Scanner interface which scans database rows into custom types
func (user *BasicWebUser) ScanRow(rows database.RowInterface) (interface{}, error) {
	var err error
	// We scan into thw newly allocated object
	err = rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
	)
	if err != nil {
		return User{}, err
	}
	return *user, nil
}

// ScanRow implements the database.Scanner interface which scans database rows into custom types
func (user *User) ScanRow(rows database.RowInterface) (interface{}, error) {
	var err error
	// We scan into thw newly allocated object
	err = rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.AccountType,
		&user.Password,
		&user.CreatedAt,
		&user.Activated,
	)
	if err != nil {
		return User{}, err
	}
	return *user, nil
}

// Get User ID
func (user User) GetID() types.UUID {
	return user.ID
}

// GetEmail returns the email of a user
func (user User) GetEmail() string {
	return types.AsString(user.Email)
}

// GetFullName returns the full name of a user
func (user User) GetFullName() string {
	return strings.Title(strings.Trim(types.AsString(user.FirstName), " ")) + " " +
		strings.Title(strings.Trim(types.AsString(user.LastName), " "))
}

func GetUserList(query database.IQuery) ([]interface{}, error) {
	var rows []interface{}
	var err error
	var user User
	// Connect to the database
	// Execute a query
	rows, err = query.QueryMany(queries.GetUsers, &user)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (user *User) ChangePassword(pass crypto.Password) error {
	var query database.IQuery
	var err error
	err = query.Exec(queries.UpdatePassword, pass, user.GetID())
	if err != nil {
		return err
	}
	return nil
}

// SetAuthToken sets the authentication token for later access if needed
func (user User) SetAuthToken(token crypto.Token) {
	user.authToken = token
}

// Create a Session record when a user sign into app
func (user User) CreateSession(transaction database.IQuery) (AuthenticationData, error) {
	var hash []byte
	var err error
	hash, err = crypto.RandomHash()
	if err != nil {
		return AuthenticationData{}, helpers.NewHttpError(http.StatusInternalServerError, err.Error())
	}
	// Close any dangling sessions
	err = user.destroyDanglingSessions(transaction)
	if err != nil {
		return AuthenticationData{}, helpers.NewHttpError(http.StatusInternalServerError, err.Error())
	}
	err = transaction.Exec(
		queries.CreateSession,
		// Parameters go here
		hex.EncodeToString(hash),
		user.ID,
	)
	if err != nil {
		return AuthenticationData{}, helpers.NewHttpError(http.StatusInternalServerError, err.Error())
	}
	return AuthenticationData{Token: crypto.NewToken(hash), User: user, Preregistration: false}, nil
}

func (user User) destroyDanglingSessions(transaction database.IQuery) error {
	var err error
	err = transaction.Exec(queries.InvalidateAllTokens, user.ID)
	if err != nil {

		return err
	}
	return nil
}

func (user *User) Exists(query database.IQuery) (bool, error) {
	var err error
	err = query.QueryOne(queries.GetUserByEmail, user, user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, err
		}
		return false, helpers.NewHttpError(http.StatusInternalServerError, err.Error())
	}
	return true, nil
}

func (user *User) Create(query database.IQuery) error {
	var scanner scanners.UUID
	var err error

	err = query.QueryOne(
		queries.CreateUser,
		&scanner,
		user.FirstName,
		user.LastName,
		user.Password,
		user.Email,
		user.AccountType,
		user.ProfessionalCategory,
		user.Message,
	)
	if err != nil {
		return err
	}
	user.ID = scanner.Value
	return nil
}

func (user *User) Update(query database.IQuery, firstname types.String, lastname types.String, password crypto.Password, accountType AccountType, email types.String, preregistration bool, activated bool) error {

	var err error
	err = query.Exec(
		queries.UpdateUser,
		firstname,
		lastname,
		password,
		accountType,
		preregistration,
		activated,
		email,
	)

	user.FirstName = firstname
	user.LastName = lastname
	user.Activated = types.NewBool(activated)
	user.AccountType = accountType

	if err != nil {
		return err
	}

	return nil
}

func (user User) DestroySession(transaction database.IQuery) error {
	var err error
	if err != nil {
		return err
	}
	defer transaction.Rollback()
	err = transaction.Exec(queries.InvalidateToken, user.ID, user.authToken)
	if err != nil {
		return err
	}
	transaction.Commit()
	return nil
}

// GetUserByToken finds a user using with matching token in the database
func GetUserByToken(token string) (User, error) {
	var query database.IQuery
	var err error
	var user User
	err = query.QueryOne(queries.GetUserByToken, &user, token)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func FindUser(list []interface{}, fn func(*User, interface{}) bool, data interface{}) *User {
	if list == nil {
		return nil
	}
	for i := range list {
		var user = list[i].(User)
		if fn(&user, data) {
			return &user
		}
	}
	return nil
}

func (user User) GetSessionCount() int {
	var err error
	var scanner scanners.Int
	var query database.IQuery
	// Connect to the database
	// Execute a query
	err = query.QueryOne(queries.GetUserSessionCount, &scanner, user.ID)
	if err != nil {
		return 0
	}
	if types.IsNil(scanner.Value) {
		// FIXME: not sure if this is semantically correct
		return 0
	}
	return types.AsInt(scanner.Value)
}
