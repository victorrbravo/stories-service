package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"

	"stories/constants"
)

// Int a boolean type that understands json's `null'
type Int struct {
	notNil bool
	value  int
}

// NilInt Nil for Int values
var NilInt = Int{value: -1, notNil: false}

// Value implements the sql.Valuer interface
func (integer Int) Value() (driver.Value, error) {
	if !integer.notNil {
		return nil, nil
	}
	return int64(integer.value), nil
}

// Scan implements the sql.Scanner interface
func (integer *Int) Scan(src interface{}) error {
	switch value := src.(type) {
	case nil:
		break
	case int:
		integer.notNil = true
		integer.value = value
	case int32:
		integer.notNil = true
		integer.value = int(value)
		break
	case int64:
		integer.value = int(value)
		integer.notNil = true
	default:
		return fmt.Errorf(constants.InvalidValueForInteger, value)
	}
	return nil
}

func (integer *Int) UnmarshalJSON(src []byte) error {
	var strValue = string(src)
	if strValue != constants.JSONNil {
		var value int
		var err error
		value, err = strconv.Atoi(string(src))
		if err != nil {
			return err
		}
		integer.notNil = true
		integer.value = value
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (integer Int) MarshalJSON() ([]byte, error) {
	if !integer.notNil {
		return json.Marshal(nil)
	}
	return json.Marshal(integer.value)
}

// Stringer implementation
func (integer Int) String() string {
	if !integer.notNil {
		return "nil"
	}
	return strconv.Itoa(integer.value)
}

func (integer Int) isNil() bool {
	return !integer.notNil
}

func (integer Int) get() int {
	return integer.value
}

// NewInt builds an `Int' object from a native `int'
func NewInt(value int) Int {
	return Int{value: value, notNil: true}
}
