package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/victorrbravo/stories-service/constants"
)

// DateTime wrapper around the native `time.Time' with an special Format
type DateTime struct {
	DateTime time.Time
	NotNil   bool
}

var NilDateTime = DateTime{time.Time{}, false}

const (
	Iso8601         = "2006-01-02T15:04:05.999Z07:00"
	Iso8601Seconds  = "2006-01-02T15:04:05Z07:00"
	withoutTimezone = "2006-01-02T15:04:05"
	onlyDate        = "2006-01-02"
	iso8601Quoted   = "\"" + Iso8601 + "\""
	standardFormat  = "2006-01-02 15:04:05"
)

func (date DateTime) Format(format string) string {
	native := time.Time(date.DateTime)
	if native.IsZero() {
		return ""
	}
	return native.Format(format)
}

func FromTime(datetime time.Time) DateTime {
	return DateTime{datetime, true}
}

// Format Verbose
func (date DateTime) FormatVerbose() string {
	native := time.Time(date.DateTime)
	return native.Format(time.RFC850)
}

// Round wrapper around the native function
func (date DateTime) Round(duration time.Duration) DateTime {
	native := time.Time(date.DateTime)
	return DateTime{native.Round(duration), true}
}

// IsZero wrapper around the native function
func (date DateTime) IsZero() bool {
	native := time.Time(date.DateTime)
	return native.IsZero()
}

// Value implements the sql.Valuer interface
func (date DateTime) Value() (driver.Value, error) {
	native := time.Time(date.DateTime)
	if native.IsZero() {
		return nil, nil
	}
	return date.Format(Iso8601), nil
}

// String implements the `Stringer' interface
func (date DateTime) String() string {
	return date.Format(standardFormat)
}

func (date DateTime) ToIso8601Seconds() string {
	return date.Format(Iso8601Seconds)
}

func (date DateTime) ToOnlyDate() string {
	return date.Format(onlyDate)
}

func (date DateTime) WithoutTimezone() string {
	return date.Format(withoutTimezone)
}

// ValidForTimeZone Checks if value is valid for current time Zone (configuration)
func (date DateTime) ValidForTimeZone() error {
	var ServerTZ *time.Location
	var val string
	var now string
	var dateValue time.Time
	var newDate time.Time
	newDate = time.Now() // ServerConfig Current Date
	ServerTZ, _ = time.LoadLocation(constants.CurrentTimeZone)
	if ServerTZ == nil {
		return nil
	}
	dateValue = time.Time(date.DateTime)
	newDate = newDate.In(ServerTZ)
	now = newDate.Format(onlyDate)
	val = dateValue.Format(onlyDate)
	if now != val {
		return fmt.Errorf(constants.NotValidValueForTimeZone)
	}
	return nil
}

func (date *DateTime) AdjustToCurrentTimeZone() error {
	var ServerTZ *time.Location
	var newDate time.Time
	ServerTZ, _ = time.LoadLocation(constants.CurrentTimeZone)
	if ServerTZ == nil {
		return nil
	}
	newDate = time.Time((*date).DateTime)
	newDate = newDate.In(ServerTZ)
	*date = DateTime{newDate, true}
	return nil
}

// Scan implements the sql.Scanner interface
func (date *DateTime) Scan(value interface{}) error {
	var err error
	if value == nil {
		*date = DateTime{time.Time{}, true}
	} else {
		switch instance := value.(type) {
		case time.Time:
			*date = DateTime{instance, true}
		default:
			return fmt.Errorf(constants.InvalidTypeForInputParameter)
		}
	}
	return err
}

// MarshalJSON implements the json.Marshaler interface
func (date DateTime) MarshalJSON() ([]byte, error) {
	native := time.Time(date.DateTime)
	if native.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(date.DateTime.Format(Iso8601))
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (date *DateTime) UnmarshalJSON(data []byte) error {
	var err error
	if string(data) == constants.JSONNil {
		*date = DateTime{time.Time{}, true}
		return nil
	}

	var raw string
	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	var t time.Time
	if raw == "" {
		*date = DateTime{time.Time{}, true}
		return nil
	}

	// try RFC3339 first, then our accepted patterns
	if t, err = time.Parse(time.RFC3339, raw); err == nil {
		*date = DateTime{t, true}
		return nil
	}
	if t, err = time.Parse(Iso8601, raw); err == nil {
		*date = DateTime{t, true}
		return nil
	}
	if t, err = time.Parse(withoutTimezone, raw); err == nil {
		*date = DateTime{t, true}
		return nil
	}
	if t, err = time.Parse(standardFormat, raw); err == nil {
		*date = DateTime{t, true}
		return nil
	}
	return err
}

func (v *DateTime) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	date, err := NewDateTime(value)
	if err != nil {
		return err
	}
	*v = date
	return nil
}

// Before wrapper around the native function
func (date DateTime) Before(other DateTime) bool {
	native := time.Time(date.DateTime)
	return native.Before(time.Time(other.DateTime))
}

// After wrapper around the native function
func (date DateTime) After(other DateTime) bool {
	native := time.Time(date.DateTime)
	return native.After(time.Time(other.DateTime))
}

// Add wrapper around the native function
func (date DateTime) Add(duration time.Duration) DateTime {
	native := time.Time(date.DateTime)
	return DateTime{native.Add(duration), true}
}

// Sub wrapper around the native function
func (date DateTime) Sub(what DateTime) time.Duration {
	native := time.Time(date.DateTime)
	return native.Sub(time.Time(what.DateTime))
}

func NewDateTime(value string) (DateTime, error) {
	if value == "" {
		return NilDateTime, nil
	}

	var t time.Time
	var err error

	if t, err = time.Parse(time.RFC3339, value); err == nil {
		return DateTime{t, true}, nil
	}
	if t, err = time.Parse(Iso8601, value); err == nil {
		return DateTime{t, true}, nil
	}
	if t, err = time.Parse(withoutTimezone, value); err == nil {
		return DateTime{t, true}, nil
	}
	if t, err = time.Parse(standardFormat, value); err == nil {
		return DateTime{t, true}, nil
	}

	return NilDateTime, err
}

// Now wrapper around the native function
func Now() DateTime {
	return DateTime{time.Now(), true}
}

func NowUTC() DateTime {
	return DateTime{time.Now().UTC(), true}
}
