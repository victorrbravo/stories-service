package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/victorrbravo/stories-service/constants"
)

type ProfessionalCategory string

const (
	FinancialProfessional  ProfessionalCategory = "Financial Professional"
	IndividualInvestor     ProfessionalCategory = "Individual Investor"
	CommercialOrganization ProfessionalCategory = "Commercial Organization"
	AcademicInstitution    ProfessionalCategory = "Academic Institution"
	Other                  ProfessionalCategory = "Other"
)

func (enum ProfessionalCategory) IsValid() error {
	switch enum {
	case FinancialProfessional, IndividualInvestor, CommercialOrganization, AcademicInstitution, Other:
		return nil
	default:
		return fmt.Errorf(fmt.Sprintf(constants.FieldMustBeType, "Professional Category", "Financial Professional, Individual Investor, Commercial Organization, Academic Institution, Other"))
	}
}

func (pt ProfessionalCategory) isNil() bool {
	if len(string(pt)) == 0 {
		return true
	}
	return false
}

func (pt ProfessionalCategory) MarshalJSON() ([]byte, error) {
	var value string
	value = string(pt)
	if len(value) == 0 {
		return json.Marshal(nil)
	}
	return json.Marshal(value)
}

func (pt *ProfessionalCategory) UnmarshalJSON(src []byte) error {
	var valueValue string
	valueValue = string(src)
	if valueValue != constants.JSONNil {
		var s string
		err := json.Unmarshal(src, &s)
		if err != nil {
			return err
		}
		*pt = ProfessionalCategory(s)
	}
	return nil
}

func (pt *ProfessionalCategory) Scan(src interface{}) error {
	switch value := src.(type) {
	case nil:
		pt = nil
	case string:
		*pt = ProfessionalCategory(value)
	case []byte:
		*pt = ProfessionalCategory(string(value))
	default:
		return fmt.Errorf(constants.CanntReadAProductTypeFromNonString)
	}
	return nil
}

func (pt ProfessionalCategory) Value() (driver.Value, error) {
	if len(string(pt)) == 0 {
		return nil, nil
	}
	return string(pt), nil
}
