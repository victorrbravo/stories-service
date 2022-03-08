package types_test

import (
	"fmt"
	"testing"

	"stories/types"
)

func TestScanFloat(t *testing.T) {
	stringTests := []struct {
		name         string
		value        interface{}
		spectedValue interface{}
		spectedError bool
	}{
		// Valid Scan
		{name: "[1]: Simple Float", value: 1.2, spectedValue: float64(1.2)},
		{name: "[2]: Simple Float", value: []byte("2"), spectedValue: float64(2)},
		//{name: "[3]: Simple Float", value: nil, spectedValue: nil}, // Get nil return 0?

		// Invalid Scan
		{name: "[4]: Simple Float", value: []byte(`"2"`), spectedError: true},
		{name: "[5]: Simple Float", value: "2", spectedError: true},
	}
	// using tt.name from the case to use it as the `t.Run` test name
	for _, tt := range stringTests {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run(tt.name, func(t *testing.T) {
			var str = new(types.Float)
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
			// Verifica Valor
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

func TestMarshallFloat(t *testing.T) {
	stringTests := []struct {
		name    string
		value   types.Float
		spected []byte
	}{
		//Valid Marshal
		{
			name:    "[1]: Simple Float",
			value:   types.NewFloat(1),
			spected: []byte(fmt.Sprintf("%.2f", 1.00)),
		},
		{name: "[2]: Simple Float", value: types.NewFloat(1.2), spected: []byte(fmt.Sprintf("%.2f", 1.2))},
		{name: "[3]: Simple Float", value: types.NewFloat(1.434), spected: []byte(fmt.Sprintf("%.2f", 1.434))},
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

func TestUnmarshalJSONFloat(t *testing.T) {
	stringTests := []struct {
		name         string
		value        []byte
		spected      interface{}
		spectedError bool
	}{
		// Valid Bytes Content
		{name: "[1]: Simple Float", value: []byte(fmt.Sprintf("%f", 1.2)), spected: 1.2},
		{name: "[2]: Simple Float", value: []byte(fmt.Sprintf("%f", 1.3)), spected: 1.3},
	}

	// using tt.name from the case to use it as the `t.Run` test name

	for _, tt := range stringTests {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run(tt.name, func(t *testing.T) {
			var parseElement = new(types.Float)
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
