package types_test

import (
	"testing"

	"stories/types"
)

func TestScanBool(t *testing.T) {
	stringTests := []struct {
		name         string
		value        interface{}
		spectedError bool
		spectedValue interface{}
	}{
		// Valid Scan
		{name: "[1]: Simple Bool", value: true, spectedValue: true},
		{name: "[2]: Simple Bool", value: false, spectedValue: false},
		{name: "[3]: Simple Bool", value: nil, spectedValue: nil},
		// Invalid Scan
		{name: "[4]: Simple Bool", value: "true", spectedError: true},
		{name: "[5]: Simple Bool", value: 2, spectedError: true},
		{name: "[6]: Simple Bool", value: "nil", spectedError: true},
	}
	// using tt.name from the case to use it as the `t.Run` test name
	for _, tt := range stringTests {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run(tt.name, func(t *testing.T) {
			var str = new(types.Bool)
			var err = str.Scan(tt.value)
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

			var vl = value
			if vl != tt.spectedValue {
				t.Errorf("Error value not valid | Name: %s, spected %v, value: %v \n", tt.name, tt.spectedValue, vl)
				return
			}
		})
	}
}

func TestMarshallBool(t *testing.T) {
	stringTests := []struct {
		name    string
		value   types.Bool
		spected string
	}{
		// Valid Marshal
		{name: "[1]: Simple Bool", value: types.NewBool(true), spected: "true"},
		{name: "[2]: Simple Bool", value: types.NewBool(false), spected: "false"},
		{name: "[3]: Simple Bool", value: types.Bool{}, spected: "null"},
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
				t.Errorf("Error value not valid | Name: %s, spected %v, value: %v \n", tt.name, tt.spected, value)
				return
			}
		})
	}
}

func TestUnmarshalJSONBool(t *testing.T) {
	stringTests := []struct {
		name         string
		value        string
		spected      interface{}
		spectedError bool
	}{
		// Valid Bytes Content
		{name: "[1]: Simple Bool", value: "true", spected: true},
		{name: "[2]: Simple Bool", value: "false", spected: false},
		{name: "[3]: Simple Bool", value: "null", spected: nil},
		// Invalid Value
		{name: "[4]: Simple Bool", value: "1", spectedError: true},
		{name: "[5]: Simple Bool", value: `"false"`, spectedError: true},
		{name: "[6]: Simple Bool", value: `"true"`, spectedError: true},
	}

	// using tt.name from the case to use it as the `t.Run` test name

	for _, tt := range stringTests {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run(tt.name, func(t *testing.T) {
			var parseElement = new(types.Bool)
			var err = parseElement.UnmarshalJSON([]byte(tt.value))
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
