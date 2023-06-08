package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"stories/constants"
)

// String a boolean type that understands jsons `null'
type String struct {
	notNil bool
	value  string
}

// Value implements the sql.Valuer interface
func (str String) Value() (driver.Value, error) {
	if !str.notNil {
		return nil, nil
	}
	return str.value, nil
}

// Scan implements the sql.Scanner interface
func (str *String) Scan(src interface{}) error {
	switch value := src.(type) {
	case nil:
		str.notNil = false
		str.value = ""
		break
	case string:
		str.notNil = true
		str.value = value
		break
	case []byte:
		str.notNil = true
		str.value = string(value)
	default:
		return fmt.Errorf(constants.InvalidValueForString, value)
	}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (str *String) UnmarshalJSON(src []byte) error {
	strValue := string(src)
	if strValue != constants.JSONNil {
		var err error
		err = json.Unmarshal(src, &str.value)
		if err != nil {
			return err
		}
		str.notNil = true
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (str String) MarshalJSON() ([]byte, error) {
	if !str.notNil {
		return json.Marshal(nil)
	}
	return json.Marshal(str.value)
}

func (str String) isNil() bool {
	return !str.notNil
}

func (str String) get() string {
	return str.value
}

// NewString builds a nullable string from a native string
func NewString(str string) String {
	return String{value: str, notNil: true}
}

func NilString() String {
	return String{value: "", notNil: false}
}

func (str String) String() string {
	if str.notNil {
		return str.value
	} else {
		return "nil"
	}
}
