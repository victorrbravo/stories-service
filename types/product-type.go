package types

import (
	"database/sql/driver"
	"encoding/json"
	"stories/constants"
	"fmt"
)

type ProductType string

const (
	Debt      ProductType = "Debt"
	Equity                = "Equity"
	Mezzanine             = "Mezzanine"
	Services              = "Services"
	Shares                = "Shares or Equity-linked"
	Hedging               = "Hedging & Financial Markets"
	OtherType             = "Other"
)

func (enum ProductType) IsValid() bool {
	switch enum {
	case Debt, Equity, Mezzanine, Services, Shares, Hedging, OtherType:
		return true
	default:
		return false
	}
}

func (pt ProductType) isNil() bool {
	if len(string(pt)) == 0 {
		return true
	}
	return false
}

func (pt ProductType) MarshalJSON() ([]byte, error) {
	var value string
	value = string(pt)
	if len(value) == 0 {
		return json.Marshal(nil)
	}
	return json.Marshal(value)
}

func (pt *ProductType) UnmarshalJSON(src []byte) error {
	var valueValue string
	valueValue = string(src)
	if valueValue != constants.JSONNil {
		var s string
		err := json.Unmarshal(src, &s)
		if err != nil {
			return err
		}
		*pt = ProductType(s)
	}
	return nil
}

func (pt *ProductType) Scan(src interface{}) error {
	switch value := src.(type) {
	case nil:
		pt = nil
	case string:
		*pt = ProductType(value)
	case []byte:
		*pt = ProductType(string(value))
	default:
		return fmt.Errorf(constants.CanntReadAProductTypeFromNonString)
	}
	return nil
}

func (pt ProductType) Value() (driver.Value, error) {
	if len(string(pt)) == 0 {
		return nil, nil
	}
	return string(pt), nil
}
