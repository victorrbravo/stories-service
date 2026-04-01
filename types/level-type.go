package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"stories/constants"
)

type LevelType string

const (
	LevelHigh       LevelType = "High"
	LevelMediumHigh LevelType = "Medium-High"
	LevelMedium     LevelType = "Medium"
	LevelMediumLow  LevelType = "Medium-Low"
	LevelLow        LevelType = "Low"
)

func (lt LevelType) IsValid() error {
	switch lt {
	case LevelHigh, LevelMedium, LevelLow, LevelMediumLow, LevelMediumHigh:
		return nil
	}
	return fmt.Errorf(constants.FieldMustBeType, "LevelType", "High, Medium, Low, Medium-Low, Medium-High")
}

func (pt LevelType) isNil() bool {
	if len(string(pt)) == 0 {
		return true
	}
	return false
}

func (pt LevelType) MarshalJSON() ([]byte, error) {
	var value string
	value = string(pt)
	if len(value) == 0 {
		return json.Marshal(nil)
	}
	return json.Marshal(value)
}

func (pt *LevelType) UnmarshalJSON(src []byte) error {
	var valueValue string
	valueValue = string(src)
	if valueValue != constants.JSONNil {
		var s string
		err := json.Unmarshal(src, &s)
		if err != nil {
			return err
		}
		*pt = LevelType(s)
	}
	return nil
}

func (pt *LevelType) Scan(src interface{}) error {
	switch value := src.(type) {
	case nil:
		pt = nil
	case string:
		*pt = LevelType(value)
	case []byte:
		*pt = LevelType(string(value))
	default:
		return fmt.Errorf(constants.CanntReadALevelTypeFromNonString)
	}
	return nil
}

func (pt LevelType) Value() (driver.Value, error) {
	if len(string(pt)) == 0 {
		return nil, nil
	}
	return string(pt), nil
}
