package types

import (
	"database/sql/driver"
	"encoding/json"
	"stories/constants"
	"fmt"
)

type RelationshipType string

const (
	Owner             RelationshipType = "Owner or Investor"
	Client                             = "Client or Customer"
	Investment                         = "Investment"
	Partner                            = "Partner"
	Supplier                           = "Supplier"
	Employer                           = "Employer"
	OtherRelationship                  = "Other"
)

func (rt RelationshipType) MarshalJSON() ([]byte, error) {
	var value string
	value = string(rt)
	if len(value) == 0 {
		return json.Marshal(nil)
	}
	return json.Marshal(value)
}

func (enum RelationshipType) IsValid() bool {
	switch enum {
	case Owner, Client, Investment, Partner, Supplier, Employer, OtherRelationship:
		return true
	default:
		return false
	}
}

func (rt *RelationshipType) Scan(src interface{}) error {
	switch value := src.(type) {
	case nil:
		rt = nil
	case []byte:
		*rt = RelationshipType(string(value))
	default:
		return fmt.Errorf(constants.CanntReadAProductTypeFromNonString)
	}
	return nil
}

func (pt *RelationshipType) UnmarshalJSON(src []byte) error {
	var err error
	var valueValue string
	valueValue = string(src)
	if valueValue != constants.JSONNil {
		var s string
		err = json.Unmarshal(src, &s)
		if err != nil {
			return err
		}
		*pt = RelationshipType(s)
	}
	return nil
}

func (rt RelationshipType) isNil() bool {
	if len(string(rt)) == 0 {
		return true
	}
	return false
}

func (rt RelationshipType) Value() (driver.Value, error) {
	if len(string(rt)) == 0 {
		return nil, nil
	}
	return string(rt), nil
}
