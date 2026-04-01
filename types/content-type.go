package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"

	"stories/constants"
)

type ContentType string

const (
	Cooking ContentType = "Cooking"
	Tourist ContentType = "Tourist"
	Errands ContentType = "Errands"
	Senda   ContentType = "Senda"
)

var NilContentType = NewNilContentType()

func NewNilContentType() ContentType {
	return ContentType("")
}

func (lt ContentType) IsValid() error {
	switch lt {
	case Cooking, Tourist, Errands, Senda:
		return nil
	}
	return fmt.Errorf(constants.FieldMustBeType, "ContentType", "Cooking, Tourist, Errands, Senda")
}

func (pt ContentType) isNil() bool {
	if len(string(pt)) == 0 {
		return true
	}
	return false
}

func (pt ContentType) MarshalJSON() ([]byte, error) {
	var value string
	value = string(pt)
	if len(value) == 0 {
		return json.Marshal(nil)
	}
	return json.Marshal(value)
}

func (pt *ContentType) UnmarshalJSON(src []byte) error {
	var valueValue string
	valueValue = string(src)
	if valueValue != constants.JSONNil {
		var s string
		err := json.Unmarshal(src, &s)
		if err != nil {
			return err
		}
		*pt = ContentType(s)
	}
	return nil
}

func (pt *ContentType) Scan(src interface{}) error {
	switch value := src.(type) {
	case nil:
		pt = nil
	case string:
		*pt = ContentType(value)
	case []byte:
		log.Println("value", value)
		*pt = ContentType(string(value))
	default:
		return fmt.Errorf(constants.CanntReadAContentTypeFromNonString)
	}
	return nil
}

func (pt ContentType) Value() (driver.Value, error) {
	if len(string(pt)) == 0 {
		return nil, nil
	}
	return string(pt), nil
}
