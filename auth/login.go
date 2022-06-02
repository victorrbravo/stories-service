package auth

import (
	"fmt"
	"time"

	"stories/database"
	"stories/types"
)

type Login struct {
	ID                  types.UUID  `json:"id"`
	LoginSuccess        types.Int   `json:"login_success"`
	LastLoginSuccess    *time.Time  `json:"last_login_success"`
	LoginErrors5Minutes types.Int   `json:"login_errors_5_minutes"`
	FirstLoginError     *time.Time  `json:"first_login_error"`
	LastLoginError      *time.Time  `json:"last_login_error"`
	FirstLoginErrorDiff types.Float `json:"first_login_error_diff"`
	LastLoginErrorDiff  types.Float `json:"last_login_error_diff"`
}

type LoginError struct {
	Blocked           bool       `json:"blocked"`
	LoginErrorsCount  int        `json:"login_errors_count"`
	LoginUnblockTime  *time.Time `json:"login_unlocking_time"`
	UnlockingLeftTime float64    `json:"time_remaining_unlock"`
}

func (login *Login) ScanRow(rows database.RowInterface) (interface{}, error) {
	var err error
	// We scan into thw newly allocated object
	err = rows.Scan(
		&login.ID,
		&login.LoginSuccess,
		&login.LastLoginSuccess,
		&login.LoginErrors5Minutes,
		&login.FirstLoginError,
		&login.LastLoginError,
		&login.FirstLoginErrorDiff,
		&login.LastLoginErrorDiff,
	)
	if err != nil {
		return User{}, err
	}
	return *login, nil
}

const UpdateLoginInformation = `UPDATE "users" SET "%s" = $1 WHERE "id" = $2`

func (login *Login) UpdateLoginBlocked(value bool, userID types.UUID) error {
	var query database.IQuery
	var err error
	var attribute string = "blocked"
	err = query.Exec(fmt.Sprintf(UpdateLoginInformation, attribute), value, userID)
	if err != nil {
		return err
	}

	return nil
}

func (login *Login) UpdateLoginDate(value time.Time, attribute string, userID types.UUID) error {
	var query database.IQuery
	var err error
	err = query.Exec(fmt.Sprintf(UpdateLoginInformation, attribute), value, userID)
	if err != nil {
		return err
	}
	return nil
}

func (login *Login) CleanLoginErrors(userID types.UUID, attribute string, value interface{}) error {
	var query database.IQuery
	var err error
	err = query.Exec(fmt.Sprintf(UpdateLoginInformation, attribute), value, userID)
	if err != nil {
		return err
	}
	return nil
}

func (login *Login) UpdateLastLoginSuccess(userID types.UUID) error {
	var now time.Time
	var err error
	now = time.Now()
	err = login.UpdateLoginDate(now, "last_login_success", userID)
	if err != nil {
		return err
	}
	return nil
}

func (login *Login) UpdateLoginErrorCount(userID types.UUID, value int) error {
	var query database.IQuery
	var err error
	err = query.Exec(fmt.Sprintf(UpdateLoginInformation, "login_errors_5_minutes"), value, userID)
	if err != nil {
		return err
	}
	return nil
}

func (login *Login) SetLoginSuccessStatus(user User) error {
	var err error
	err = login.CleanLoginErrors(user.GetID(), "login_errors_5_minutes", 0)
	if err != nil {
		return err
	}
	err = login.CleanLoginErrors(user.GetID(), "first_login_error", nil)
	if err != nil {
		return err
	}
	err = login.CleanLoginErrors(user.GetID(), "last_login_error", nil)
	if err != nil {
		return err
	}
	return nil
}
