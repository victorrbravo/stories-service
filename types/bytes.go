package types

import (
	"database/sql/driver"
	"fmt"

	"github.com/victorrbravo/stories-service/constants"
)

// Gender a simple gender type for easier validation
type Bytes []byte

var NilBytes = Bytes{}

// Gender values

// Scan implements the sql.Scanner interface
func (bytes *Bytes) Scan(src interface{}) error {
	switch value := src.(type) {
	case []byte:
		*bytes = value
	case string:
		*bytes = []byte(value)
	default:
		return fmt.Errorf(constants.InvalidGenderValue, value)
	}
	return nil
}

func (str Bytes) Value() (driver.Value, error) {
	return []byte(str), nil
}

func (str *Bytes) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	value = string(*str)
	if err := unmarshal(&value); err != nil {
		return err
	}
	*str = []byte(value)

	return nil
}

func (str Bytes) MarshalYAML() (interface{}, error) {
	return string(str), nil
}

// Value implements the sql.Valuer interface
