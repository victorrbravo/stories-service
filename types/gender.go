package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/victorrbravo/stories-service/constants"
)

// Gender a simple gender type for easier validation
type Gender rune

// Gender values
const (
	Female          Gender = 'F'
	Male            Gender = 'M'
	unknownGender   Gender = 0
	MaLeLowerCase   Gender = 'm'
	FemaleLowerCase Gender = 'f'
)

func (gender Gender) convertToValue() (string, error) {
	switch gender {
	case unknownGender:
		return "?", nil
	case Male, MaLeLowerCase:
		return "M", nil
	case Female, FemaleLowerCase:
		return "F", nil
	default:
		return "", fmt.Errorf(constants.InvalidGenderValue, gender)
	}
}

func (gender *Gender) convertFromValue(value byte) error {
	switch value {
	case 'M', 'm':
		*gender = Male
	case 'F', 'f':
		*gender = Female
	default:
		return fmt.Errorf(constants.InvalidGenderValue, value)
	}
	return nil
}

// Scan implements the sql.Scanner interface
func (gender *Gender) Scan(src interface{}) error {
	switch value := src.(type) {
	case string:
		if len(value) == 0 {
			return fmt.Errorf(constants.GenderCannotBeEmpty, value)
		}
		return gender.convertFromValue(value[0])
	case []byte:
		return gender.convertFromValue(value[0])
	case byte:
		return gender.convertFromValue(value)
	case rune:
		return gender.convertFromValue(byte(value))
	default:
		return fmt.Errorf(constants.InvalidGenderValue, gender)
	}
}

// Value implements the sql.Valuer interface
func (gender Gender) Value() (driver.Value, error) {
	return gender.convertToValue()
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (gender *Gender) UnmarshalJSON(src []byte) error {
	var value string
	var err error
	err = json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	if len(value) < 1 {
		return fmt.Errorf(constants.InvalidValueForGender, value)
	}
	return gender.convertFromValue(value[0])
}

// MarshalJSON implements the json.Marshaler interface
func (gender Gender) MarshalJSON() ([]byte, error) {
	var value string
	var err error
	value, err = gender.convertToValue()
	if err != nil {
		return nil, err
	}
	return json.Marshal(value)
}
