package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"

	"stories/constants"
)

// Float a floating point value that marshals to 2 decimal places
type Float struct {
	notNil bool
	value  float64
}

// MarshalJSON implements the json.Marshaler interface
func (value Float) MarshalJSON() ([]byte, error) {
	if !value.notNil {
		return json.Marshal(nil)
	}
	return []byte(strconv.FormatFloat(value.value, 'f', 2, 64)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (value *Float) UnmarshalJSON(src []byte) error {
	var valueValue = string(src)
	var number float64
	if valueValue != constants.JSONNil {
		var err error
		number, err = strconv.ParseFloat(string(src), 64)
		if err != nil {
			return err
		}
		value.notNil = true
		value.value = number
	}
	return nil
}

// Value implements the sql.Valuer interface
func (value Float) Value() (driver.Value, error) {
	if IsNil(value) {
		return nil, nil
	}
	return value.value, nil
}

// Scan implements the sql.Scanner interface
func (value *Float) Scan(src interface{}) error {
	var raw float64
	var err error
	switch val := src.(type) {
	case nil:
		*value = NilFloat()
		return nil
	case []byte:
		raw, err = strconv.ParseFloat(string(val), 64)
		if err != nil {
			return err
		}
		break
	case float32:
		raw = float64(val)
		break
	case float64:
		raw = val
		break
	default:
		return fmt.Errorf(constants.InvalidValueForFloat, src)
	}
	*value = Float{true, raw}
	return nil
}

func (value Float) isNil() bool {
	return !value.notNil
}

func (value Float) get() float64 {
	return value.value
}

func AsFloat(value Float) float64 {
	return value.get()
}

func NewFloat(num float64) Float {
	return Float{true, num}
}

func NilFloat() Float {
	return Float{false, 0}
}
