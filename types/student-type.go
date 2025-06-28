package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"stories/constants"
)

type StudentType string

var NilStudentType = StudentType("")

const (
	NormalStudent StudentType = "NormalStudent"
	Admin                     = "Admin"
	Distinguished             = "Distinguished"
)

func (enum StudentType) IsValid() bool {
	switch enum {
	case NormalStudent, Admin, Distinguished:
		return true
	default:
		return false
	}
}

func (pt StudentType) isNil() bool {
	if len(string(pt)) == 0 {
		return true
	}
	return false
}

func (pt StudentType) MarshalJSON() ([]byte, error) {
	var value string
	value = string(pt)
	if len(value) == 0 {
		return json.Marshal(nil)
	}
	return json.Marshal(value)
}

func (pt *StudentType) UnmarshalJSON(src []byte) error {
	var valueValue string
	valueValue = string(src)
	if valueValue != constants.JSONNil {
		var s string
		err := json.Unmarshal(src, &s)
		if err != nil {
			return err
		}
		*pt = StudentType(s)
	}
	return nil
}

func (pt *StudentType) Scan(src interface{}) error {
	switch value := src.(type) {
	case nil:
		pt = nil
	case string:
		*pt = StudentType(value)
	case []byte:
		*pt = StudentType(string(value))
	default:
		return fmt.Errorf("Error Student")
	}
	return nil
}

func (pt StudentType) Value() (driver.Value, error) {
	if len(string(pt)) == 0 {
		return nil, nil
	}
	return string(pt), nil
}
