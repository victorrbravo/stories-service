package types_test

import (
	"testing"
	"time"

	"stories/types"
)

var iso8601 = "2006-01-02T15:04:05.999Z07:00"

func TestScanDateTime(t *testing.T) {
	init := time.Now()
	year, err := time.Parse("2006-01-02", "2020-05-15")
	if err != nil {
		t.Error("DateTime-Scan: Parse init error")
		return
	}
	stringTests := []struct {
		name         string
		value        interface{}
		spectedError bool
		spectedValue interface{}
	}{
		// Valid Scan
		{name: "[1]: Simple DateTime", value: init, spectedValue: init.Format(iso8601)},
		{name: "[2]: Simple DateTime", value: year, spectedValue: year.Format(iso8601)},
		// Invalid Scan
		{name: "[4]: Simple DateTime", value: "true", spectedError: true},
		{name: "[5]: Simple DateTime", value: 2, spectedError: true},
		{name: "[6]: Simple DateTime", value: "nil", spectedError: true},
	}
	// using tt.name from the case to use it as the `t.Run` test name
	for _, tt := range stringTests {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run(tt.name, func(t *testing.T) {
			str := new(types.DateTime)
			err := str.Scan(tt.value)
			// Error quest
			if err != nil {
				if !tt.spectedError {
					t.Errorf("Error Scan | Name: %s, error: %v \n", tt.name, err)
				}
				return // On error next element
			}
			if tt.spectedError {
				t.Errorf("Error Not Found | Name: %s, value: %v \n", tt.name, tt.value)
				return
			}
			// Get Value
			value, err := str.Value()
			if err != nil {
				t.Errorf("Error get value | Name: %s, err %v, value: %v \n", tt.name, err, tt.value)
				return
			}

			vl := value
			if vl != tt.spectedValue {
				t.Errorf("Error value not valid | Name: %s, spected %v, value: %v \n", tt.name, tt.spectedValue, vl)
				return
			}
		})
	}
}

func TestMarshallDateTime(t *testing.T) {
	init := time.Now()
	year, err := time.Parse("2006-01-02", "2020-05-15")
	if err != nil {
		t.Error("DateTime-Marshall: Parse init error")
		return
	}
	stringTests := []struct {
		name    string
		value   types.DateTime
		spected string
	}{
		// Valid Marshal
		{name: "[1]: Simple DateTime", value: types.DateTime(init), spected: init.Format(`"` + iso8601 + `"`)},
		{name: "[2]: Simple DateTime", value: types.DateTime(year), spected: year.Format(`"` + iso8601 + `"`)},
		{name: "[3]: Simple DateTime", value: types.DateTime{}, spected: "null"},
	}
	// using tt.name from the case to use it as the `t.Run` test name
	for _, tt := range stringTests {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run(tt.name, func(t *testing.T) {
			value, err := tt.value.MarshalJSON()
			if err != nil {
				t.Errorf("Error MarshalJSON | Name: %s, error: %v \n", tt.name, err)
				return
			}
			if string(value) != string(tt.spected) {
				t.Errorf("Error value not valid | Name: %s, spected %v, value: %s \n", tt.name, tt.spected, value)
				return
			}
		})
	}
}

func TestUnmarshalJSONDateTime(t *testing.T) {
	init := time.Now()
	year, err := time.Parse("2006-01-02", "2020-05-15")
	if err != nil {
		t.Error("DateTime-Marshall: Parse init error")
		return
	}
	stringTests := []struct {
		name         string
		value        string
		spected      interface{}
		spectedError bool
	}{
		// Valid Bytes Content
		{name: "[1]: Simple DateTime", value: init.Format(`"` + iso8601 + `"`), spected: init.Format(iso8601)},
		{name: "[2]: Simple DateTime", value: year.Format(`"` + iso8601 + `"`), spected: year.Format(iso8601)},
		// Invalid Value
		{name: "[4]: Simple DateTime", value: `""`, spectedError: true},
		{name: "[5]: Simple DateTime", value: `{}`, spectedError: true},
	}

	// using tt.name from the case to use it as the `t.Run` test name

	for _, tt := range stringTests {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run(tt.name, func(t *testing.T) {
			parseElement := new(types.DateTime)
			err := parseElement.UnmarshalJSON([]byte(tt.value))
			// Error quest
			if err != nil {
				if !tt.spectedError {
					t.Errorf("Error UnmarshalJSON | Name: %s, error: %v \n", tt.name, err)
				}
				return // On error next element
			}
			if tt.spectedError {
				t.Errorf("Error Not Found | Name: %s, value: %v \n", tt.name, tt.value)
				return
			}
			// Get Value
			value, err := parseElement.Value()
			if err != nil {
				t.Errorf("Error get Value | Name: %s, error: %v \n", tt.name, err)
				return
			}
			if value != tt.spected {
				t.Errorf("Error value not valid | Name: %s, spected %v, value: %v \n", tt.name, tt.spected, value)
			}
		})
	}
}
