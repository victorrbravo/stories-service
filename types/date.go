package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"stories/constants"
)

// Date wrapper around the native `time.Time' that strips the time part
type Date time.Time

func (date Date) format(format string) string {
	var native = time.Time(date)
	if native.IsZero() {
		return ""
	}
	return native.Format(format)
}

func (date Date) toString() string {
	return date.format("01/02/2006")
}

// Value implements the sql.Valuer interface
func (date Date) Value() (driver.Value, error) {
	var native = time.Time(date)
	if native.IsZero() {
		return nil, nil
	}
	return date.format("2006-01-02"), nil
}

// Scan implements the sql.Scanner interface
func (date *Date) Scan(value interface{}) error {
	var err error
	if value == nil {
		*date = Date(time.Time{})
	} else {
		switch instance := value.(type) {
		case time.Time:
			*date = Date(instance)
			break
		default:
			return fmt.Errorf(constants.InvalidTypeForInputParameter)
		}
	}
	return err
}

// MarshalJSON implements the json.Marshaler interface
func (date Date) MarshalJSON() ([]byte, error) {
	var native time.Time
	native = time.Time(date)
	if native.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(date.toString())
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (date *Date) UnmarshalJSON(data []byte) error {
	var err error
	var t time.Time
	t, err = time.Parse("\"2006-01-02\"", string(data))
	if err == nil {
		*date = Date(t)
	}
	return err
}

// After wrapper around the native `After()'
func (date Date) After(other Date) bool {
	var native = time.Time(date)
	return native.After(time.Time(other))
}

// IsZero wrapper around the native `IsZero()'
func (date Date) IsZero() bool {
	var native = time.Time(date)
	return native.IsZero()
}
