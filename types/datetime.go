package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"stories/constants"
)

// DateTime wrapper around the native `time.Time' with an special Format
type DateTime time.Time

const (
	Iso8601         = "2006-01-02T15:04:05.999Z07:00"
	withoutTimezone = "2006-01-02T15:04:05"
	onlyDate        = "2006-01-02"
	iso8601Quoted   = "\"" + Iso8601 + "\""
	standardFormat  = "2006-01-02 15:04:05"
)

func (date DateTime) Format(format string) string {
	native := time.Time(date)
	if native.IsZero() {
		return ""
	}
	return native.Format(format)
}

// Format Verbose
func (date DateTime) FormatVerbose() string {
	native := time.Time(date)
	return native.Format(time.RFC850)
}

// Round wrapper around the native function
func (date DateTime) Round(duration time.Duration) DateTime {
	native := time.Time(date)
	return DateTime(native.Round(duration))
}

// IsZero wrapper around the native function
func (date DateTime) IsZero() bool {
	native := time.Time(date)
	return native.IsZero()
}

// Value implements the sql.Valuer interface
func (date DateTime) Value() (driver.Value, error) {
	native := time.Time(date)
	if native.IsZero() {
		return nil, nil
	}
	return date.Format(Iso8601), nil
}

// String implements the `Stringer' interface
func (date DateTime) String() string {
	return date.Format("01/02/2006 03:04 PM")
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
	dateValue = time.Time(date)
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
	newDate = time.Time(*date)
	newDate = newDate.In(ServerTZ)
	*date = DateTime(newDate)
	return nil
}

// Scan implements the sql.Scanner interface
func (date *DateTime) Scan(value interface{}) error {
	var err error
	if value == nil {
		*date = DateTime(time.Time{})
	} else {
		switch instance := value.(type) {
		case time.Time:
			*date = DateTime(instance)
			break
		default:
			return fmt.Errorf(constants.InvalidTypeForInputParameter)
		}
	}
	return err
}

// MarshalJSON implements the json.Marshaler interface
func (date DateTime) MarshalJSON() ([]byte, error) {
	native := time.Time(date)
	if native.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(date.Format(Iso8601))
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (date *DateTime) UnmarshalJSON(data []byte) error {
	var err error
	var t time.Time
	if string(data) == constants.JSONNil {
		*date = DateTime(time.Time{})
	} else {
		t, err = time.Parse(iso8601Quoted, string(data))
		if err == nil {
			*date = DateTime(t)
		}
	}
	return err
}

// Before wrapper around the native function
func (date DateTime) Before(other DateTime) bool {
	native := time.Time(date)
	return native.Before(time.Time(other))
}

// After wrapper around the native function
func (date DateTime) After(other DateTime) bool {
	native := time.Time(date)
	return native.After(time.Time(other))
}

// Add wrapper around the native function
func (date DateTime) Add(duration time.Duration) DateTime {
	native := time.Time(date)
	return DateTime(native.Add(duration))
}

// Sub wrapper around the native function
func (date DateTime) Sub(what DateTime) time.Duration {
	native := time.Time(date)
	return native.Sub(time.Time(what))
}

// Now wrapper around the native function
func Now() DateTime {
	return DateTime(time.Now())
}

func NowUTC() DateTime {

	return DateTime(time.Now().UTC())
}
