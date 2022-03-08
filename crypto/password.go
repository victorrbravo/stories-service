package crypto

import (
	"database/sql/driver"
	"encoding/json"
)

// Password an encrypted password type
type Password []byte

// UnmarshalJSON implements the json.Unmarshaler interface
func (password *Password) UnmarshalJSON(input []byte) error {
	var raw string
	var err error
	err = json.Unmarshal(input, &raw)
	if err != nil {
		return err
	}
	*password = []byte(raw)
	return nil
}

// Value implements the sql.Valuer interface
func (password Password) Value() (driver.Value, error) {
	return HashPassword(string(password))
}

// MarshalJSON implements the json.Marshaler interface
// This method is useful to prevent accidentally sending a password hash to someone
func (password Password) MarshalJSON() ([]byte, error) {
	return []byte("\"\""), nil
}

// Scan implements the sql.Scanner interface
func (password *Password) Scan(src interface{}) error {
	*password = []byte(src.(string))
	return nil
}
