package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"stories/constants"
)

// Bool a boolean type that understands json's `null'
type Bool struct {
	notNil bool
	value  bool
}

var NilBool = Bool{value: false, notNil: false}

// NewBool creates new nullable boolean from a native bool
func NewBool(value bool) Bool {
	return Bool{notNil: true, value: value}
}

// Value implements the sql.Valuer interface
func (boolean Bool) Value() (driver.Value, error) {
	if !boolean.notNil {
		return nil, nil
	}
	return boolean.value, nil
}

// Scan implements the sql.Scanner interface
func (boolean *Bool) Scan(src interface{}) error {
	switch value := src.(type) {
	case nil:
		break
	case bool:
		boolean.notNil = true
		boolean.value = value
		break
	default:
		return fmt.Errorf(constants.InvalidValueForBoolean, src)
	}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (boolean *Bool) UnmarshalJSON(src []byte) error {
	valueValue := string(src)
	if valueValue != constants.JSONNil {
		var err error
		err = json.Unmarshal(src, &boolean.value)
		if err != nil {
			return err
		}
		boolean.notNil = true
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (boolean Bool) MarshalJSON() ([]byte, error) {
	if !boolean.notNil {
		return json.Marshal(nil)
	}
	return json.Marshal(boolean.value)
}

func (boolean Bool) String() string {
	if boolean.notNil {
		return fmt.Sprintf("%v", boolean.value)
	}
	return "nil"
}

func (boolean Bool) isNil() bool {
	return !boolean.notNil
}

func (boolean Bool) get() bool {
	return boolean.value
}

func (boolean Bool) ToBool() bool {
	if boolean.notNil {
		return boolean.value
	}

	return false
}
